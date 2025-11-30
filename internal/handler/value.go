package handler

import (
	"fmt"
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
		fmt.Printf("/value - OK: %v, get value: %v\n", ok, existItem)
		if !ok {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Metric not found"})
			return
		}

		if existItem.Value != nil {
			c.JSON(http.StatusOK, gin.H{"id": existItem.ID, "type": existItem.MType, "value": *existItem.Value})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": existItem.ID, "type": existItem.MType, "value": *existItem.Delta})
	}
}
