package main

import (
	"github.com/gin-gonic/gin"

	"github.com/Vladis-r/metrics.git/cmd/config"
	"github.com/Vladis-r/metrics.git/internal/handler"
	"github.com/Vladis-r/metrics.git/internal/model/logger"
)

func main() {
	c := config.GetConfig() // get config
	logger.InitLogger()     // create logger
	defer logger.Log.Sync()

	r := gin.New()                       // Create a new Gin instance
	r.Use(logger.Middleware(logger.Log)) // Add logger middleware
	r.LoadHTMLGlob("templates/*.html")   // Load HTML templates
	r.Static("/static", "./static")      // Serve static files from the "static" directory

	// handlers
	r.GET("/", handler.Root)
	r.GET("/value/:metricType/:metricName", handler.Value)
	r.POST("/update/:metricType/:metricName/:metricValue", handler.Update)

	r.Run(c.Addr) // Start server localhost:8080 by default
}
