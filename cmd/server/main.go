package main

import (
	"flag"

	"github.com/gin-gonic/gin"

	"github.com/Vladis-r/metrics.git/cmd/server/flags"
	"github.com/Vladis-r/metrics.git/internal/handler"
)

func main() {
	flag.Parse() // Parse command-line arguments

	r := gin.Default()                 // Create a new Gin instance
	r.LoadHTMLGlob("templates/*.html") // Load HTML templates
	r.Static("/static", "./static")    // Serve static files from the "static" directory

	// handlers
	r.GET("/", handler.Root)
	r.GET("/value/:metricType/:metricName", handler.Value)
	r.POST("/update/:metricType/:metricName/:metricValue", handler.Update)

	r.Run(flags.Addr) // Start server localhost:8080 by default
}
