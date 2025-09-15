package handler

import (
	"net/http"
	"strings"

	models "github.com/Vladis-r/metrics.git/internal/model"
	"github.com/gin-gonic/gin"
)

// Value - get metric by type and name.
func Value(c *gin.Context) {
	metricType := strings.ToLower(c.Param("metricType"))
	metricName := strings.ToLower(c.Param("metricName"))
	metric, ok := models.Storage.GetMetric(metricName, metricType)
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Metric not found",
		})
		return
	}

	c.JSON(http.StatusOK, metric)
}
