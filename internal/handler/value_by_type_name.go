package handler

import (
	"fmt"
	"net/http"
	"strings"

	models "github.com/Vladis-r/metrics.git/internal/model"
	"github.com/Vladis-r/metrics.git/internal/utils"
	"github.com/gin-gonic/gin"
)

// Value - get metric by type and name.
func ValueTypeName(s *models.MemStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var val interface{}

		mType := strings.ToLower(c.Param("metricType"))
		id := c.Param("metricName")
		// Check metric type.
		if _, err := utils.CheckMetric(mType, "1"); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		// Get metric from storage.
		metric, ok := s.GetMetric(id)
		if !ok {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Metric not found"})
			return
		}
		if metric.Value != nil {
			val = *metric.Value
		} else {
			val = metric.DeltaSum
		}

		c.String(http.StatusOK, fmt.Sprintf("%v", val))
	}
}
