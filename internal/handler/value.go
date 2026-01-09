package handler

import (
	"net/http"
	"strings"

	models "github.com/Vladis-r/metrics.git/internal/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

		mType := strings.ToLower(metric.MType)
		switch mType {
		case models.Counter:
			existItem, ok = s.Store[metric.ID]
		case models.Gauge:
			existItem, ok = s.Store[metric.ID]
		default:
			s.Log.Info("Unexpected type", zap.String("MType", metric.MType))
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Metric not found"})
		}

		if !ok || mType != strings.ToLower(existItem.MType) {
			s.Log.Info("Not found metric", zap.Any("Metric", metric))
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Metric not found"})
			return
		}
		c.JSON(http.StatusOK, existItem)
	}
}
