package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"database/sql"

	"github.com/Vladis-r/metrics.git/cmd/config"
	"github.com/Vladis-r/metrics.git/internal/handler"
	"github.com/Vladis-r/metrics.git/internal/middleware"
	models "github.com/Vladis-r/metrics.git/internal/model"
	"github.com/Vladis-r/metrics.git/internal/server"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	logger, err := middleware.InitLogger() // create logger
	if err != nil {
		log.Fatalf("Cant create logger: %v", err)
	}
	defer logger.Sync()
	conf := config.GetConfigServer(logger)     // get config
	var s = models.NewMemStorage(conf, logger) // Global storage for views.
	s.Log.Info("Start server with config", zap.Any("config", conf))

	chooseStorage(conf, s) // Save metrics in db or file

	r := gin.New()                   // Create a new Gin instance
	r.Use(middleware.Logger(logger)) // Add logger middleware
	r.Use(middleware.Gzip())         // Add gzip comression and decompression.

	r.LoadHTMLGlob("templates/*.html") // Load HTML templates
	r.Static("/static", "./static")    // Serve static files from the "static" directory

	// handlers
	r.GET("/", handler.Root(s))
	r.POST("/update", handler.Update(s))
	r.POST("/update/:metricType/:metricName/:metricValue", handler.UpdateTypeNameValue(s))
	r.POST("/value", handler.Value(s))
	r.GET("/value/:metricType/:metricName", handler.ValueTypeName(s))

	// service handlers
	r.GET("/ping", handler.Ping(s.DB))

	srv := newServer(conf, r)
	go startServer(srv, s)

	gracefullShutdown(srv, s)
}

// startServer - starts the HTTP server in a separate goroutine.
func startServer(srv *http.Server, s *models.MemStorage) {
	s.Log.Info("Server is listening", zap.String("addr", s.Conf.Addr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.Log.Fatal("Server failed to start", zap.Error(err))
	}
}

// newServer - creates a new HTTP server.
func newServer(conf *config.ConfigServer, r *gin.Engine) *http.Server {
	return &http.Server{
		Addr:    conf.Addr,
		Handler: r,
	}
}

// chooseStorage - choose storage for metrics.
func chooseStorage(conf *config.ConfigServer, s *models.MemStorage) (err error) {
	switch {
	case conf.DatabaseDsn != "":
		s.DB, err = sql.Open("pgx", conf.DatabaseDsn) // Connect to db.
		if err != nil {
			panic(err)
		}
		// uncomment for up migrations.
		// err = runMigrations(conf.DatabaseDsn, s)
		// if err != nil {
		// 	panic(err)
		// }
		server.LoadMetricsFromDatabase(s)
		go server.SaveMetricToDB(s)
	case conf.FileStoragePath != "":
		server.LoadMetricsFromFile(s)
		go server.SaveMetricsToFile(s)
	default:
		s.Log.Warn("No storage configured")
	}
	return nil
}

// gracefullShutdown - gracefully shutdown the server. Save metric into file.
func gracefullShutdown(srv *http.Server, s *models.MemStorage) {
	// Wait for OS signal to shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c // Block until signal received

	s.Log.Info("Shutting down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.Log.Fatal("Server forced to shutdown", zap.Error(err))
	}
	s.Log.Info("Server stopped")

	switch {
	case s.Conf.DatabaseDsn != "":
		server.SaveMetricToDBLogic(s)
		s.Log.Info("Metrics save in DB")
	case s.Conf.FileStoragePath != "":
		server.SaveMetricsToFileLogic(s)
		s.Log.Info("Metrics save in file")
	default:
		s.Log.Fatal("Metrics not save after shutdown!")
	}
	s.DB.Close() // close database connection
	s.Log.Info("Database connection is closed")
}

// runMigrations - func for run migrations while start server if needed
func runMigrations(dsn string, s *models.MemStorage) error {
	absPath, err := filepath.Abs("./migrations")
	if err != nil {
		s.Log.Error("failed to get absolute path", zap.Error(err))
		return err
	}
	sourceURL := "file://" + filepath.ToSlash(absPath)
	m, err := migrate.New(
		sourceURL,
		dsn,
	)
	if err != nil {
		return err
	}
	defer m.Close()

	_, dirty, _ := m.Version()

	if dirty {
		m.Force(-1)
	}

	err = m.Down()
	if err != nil {
		s.Log.Info("Cant migrations down")
		return nil
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		s.Log.Info("No migration changes")
		return nil
	}
	if err != nil {
		s.Log.Error("failed to run migrations", zap.Error(err))
		return err
	}

	s.Log.Info("Migrations applied successfully")
	return nil
}
