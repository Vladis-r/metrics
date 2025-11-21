package handler

import (
	"net/http"
	"slices"
	"strings"

	models "github.com/Vladis-r/metrics.git/internal/model"
	"github.com/Vladis-r/metrics.git/internal/utils"
	"github.com/gin-gonic/gin"
)

func UpdateTypeNameValue(s *models.MemStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			value interface{}
			err   error
		)

		metricType := strings.ToLower(c.Param("metricType"))
		metricName := strings.ToLower(c.Param("metricName"))
		metricValue := strings.ToLower(c.Param("metricValue"))
		// check possible metricTypes from constants.
		if !slices.Contains([]string{models.Counter, models.Gauge}, metricType) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		// check metricValues.
		if value, err = utils.CheckMetric(metricType, metricValue); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err = s.SaveMetricByTypeValue(metricName, metricType, value); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.Status(http.StatusOK)
	}
}
