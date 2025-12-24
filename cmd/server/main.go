package main

import (
	"github.com/gin-gonic/gin"

	"github.com/Vladis-r/metrics.git/cmd/config"
	"github.com/Vladis-r/metrics.git/internal/handler"
	"github.com/Vladis-r/metrics.git/internal/middleware"
	models "github.com/Vladis-r/metrics.git/internal/model"
)

func main() {
	c := config.GetConfig() // get config
	middleware.InitLogger() // create logger
	defer middleware.Log.Sync()

	var s = models.NewMemStorage() // Global storage for metrics.

	r := gin.New()                           // Create a new Gin instance
	r.Use(middleware.Logger(middleware.Log)) // Add logger middleware
	r.Use(middleware.Gzip())                 // Add gzip comression and decompression.
	r.LoadHTMLGlob("templates/*.html")       // Load HTML templates
	r.Static("/static", "./static")          // Serve static files from the "static" directory

	// handlers
	r.GET("/", handler.Root(s))

	r.POST("/update", handler.Update(s))
	r.POST("/update/:metricType/:metricName/:metricValue", handler.UpdateTypeNameValue(s))

	r.POST("/value", handler.Value(s))
	r.GET("/value/:metricType/:metricName", handler.ValueTypeName(s))

	r.Run(c.Addr) // Start server localhost:8080 by default
}
