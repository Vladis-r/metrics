package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
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

	db, err := sql.Open("pgx", conf.DatabaseDsn) // Connect to db.
	if err != nil {
		panic(err)
	}
	defer db.Close()

	server.LoadMetricsFromFile(s)
	go server.SaveMetricsToFile(s)

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
	r.GET("/ping", handler.Ping(db))

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

func newServer(conf *config.ConfigServer, r *gin.Engine) *http.Server {
	return &http.Server{
		Addr:    conf.Addr,
		Handler: r,
	}
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

	server.SaveMetricsToFileLogic(s) // Save metric into file.

	if err := srv.Shutdown(ctx); err != nil {
		s.Log.Fatal("Server forced to shutdown", zap.Error(err))
	}
	s.Log.Info("Server stopped")
}
