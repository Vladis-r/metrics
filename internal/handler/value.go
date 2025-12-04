package handler

import (
	"fmt"
	"net/http"
	"strings"

	models "github.com/Vladis-r/metrics.git/internal/model"
	"github.com/gin-gonic/gin"
)

// Value - POST handler for get metric with JSON in body.
func Value(m *models.MemStorage) gin.HandlerFunc {
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
			existItem, ok = m.Store[fmt.Sprintf("counter_%v", metric.ID)]
		case models.Gauge:
			existItem, ok = m.Store[fmt.Sprintf("gauge_%v", metric.ID)]
		default:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Metric not found"})
		}

		mType := strings.ToLower(metric.MType)
		if !ok || mType != existItem.MType {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Metric not found"})
			return
		}
		c.JSON(http.StatusOK, existItem)
	}
}
