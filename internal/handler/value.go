package handler

import (
	"net/http"

	models "github.com/Vladis-r/metrics.git/internal/model"
	"github.com/gin-gonic/gin"
)

// Value - POST handler for get metric with JSON in body.
func Value(s *models.MemStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var item models.Metric
		if err := c.ShouldBindJSON(&item); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		existItem, ok := s.Store[item.ID]
		if !ok {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Metric not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": existItem})
	}
}
