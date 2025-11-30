package handler

import (
	"net/http"

	models "github.com/Vladis-r/metrics.git/internal/model"
	"github.com/gin-gonic/gin"
)

// Value - POST handler for get metric with JSON in body.
func Value(s *models.MemStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var metric models.Metric
		if err := c.ShouldBindJSON(&metric); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		existItem, ok := s.Store[metric.ID]
		if !ok {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Metric not found"})
			return
		}

		switch metric.MType {
		case "gauge":
			c.JSON(http.StatusOK, gin.H{"id": existItem.ID, "type": existItem.MType, "value": *existItem.Value})
		case "counter":
			c.JSON(http.StatusOK, gin.H{"id": existItem.ID, "type": existItem.MType, "value": *existItem.Delta})
		default:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Metric not found"})
		}
	}
}
