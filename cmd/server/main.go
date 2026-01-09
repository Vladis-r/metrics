package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Vladis-r/metrics.git/cmd/config"
	"github.com/Vladis-r/metrics.git/internal/handler"
	"github.com/Vladis-r/metrics.git/internal/middleware"
	models "github.com/Vladis-r/metrics.git/internal/model"
	"github.com/Vladis-r/metrics.git/internal/server"
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

	server.LoadMetricsFromFile(s)
	go server.SaveMetricsToFile(s) // Save metrics to file.

	r := gin.New()                     // Create a new Gin instance
	r.Use(middleware.Logger(logger))   // Add logger middleware
	r.Use(middleware.Gzip())           // Add gzip comression and decompression.
	r.LoadHTMLGlob("templates/*.html") // Load HTML templates
	r.Static("/static", "./static")    // Serve static files from the "static" directory

	// handlers
	r.GET("/", handler.Root(s))

	r.POST("/update", handler.Update(s))
	r.POST("/update/:metricType/:metricName/:metricValue", handler.UpdateTypeNameValue(s))

	r.POST("/value", handler.Value(s))
	r.GET("/value/:metricType/:metricName", handler.ValueTypeName(s))

	r.Run(conf.Addr) // Start server localhost:8080 by default
}
