package main

import (
	"github.com/gin-gonic/gin"

	"github.com/Vladis-r/metrics.git/cmd/config"
	"github.com/Vladis-r/metrics.git/internal/handler"
)

func main() {
	c := config.GetConfig()

	r := gin.Default()                 // Create a new Gin instance
	r.LoadHTMLGlob("templates/*.html") // Load HTML templates
	r.Static("/static", "./static")    // Serve static files from the "static" directory

	// handlers
	r.GET("/", handler.Root)
	r.GET("/value/:metricType/:metricName", handler.Value)
	r.POST("/update/:metricType/:metricName/:metricValue", handler.Update)

	r.Run(c.Addr) // Start server localhost:8080 by default
}
