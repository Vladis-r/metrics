package handler

import (
	"fmt"
	"net/http"
	"strings"

	models "github.com/Vladis-r/metrics.git/internal/model"
	"github.com/gin-gonic/gin"
)

// Value - POST handler for get metric with JSON in body.
func Value(s *models.MemStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			metric    models.Metric
			ok        bool
			existItem models.Metric
		)
		if err := c.ShouldBindJSON(&metric); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		switch metric.MType {
		case models.Counter:
			existItem, ok = s.Store[fmt.Sprintf("counter_%v", metric.ID)]
		case models.Gauge:
			existItem, ok = s.Store[fmt.Sprintf("gauge_%v", metric.ID)]
		}

		mType := strings.ToLower(metric.MType)
		if !ok || mType != existItem.MType {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Metric not found"})
			return
		}

		switch mType {
		case "gauge":
			c.JSON(http.StatusOK, gin.H{"id": existItem.ID, "type": existItem.MType, "value": *existItem.Value})
		case "counter":
			c.JSON(http.StatusOK, gin.H{"id": existItem.ID, "type": existItem.MType, "value": *existItem.Delta})
		default:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Metric not found"})
		}
	}
}
