package handler

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	models "github.com/Vladis-r/metrics.git/internal/model"
	"github.com/gin-gonic/gin"
)

func Update(c *gin.Context) {
	var (
		value interface{}
		err   error
	)

	if c.Request.Method != http.MethodPost {
		c.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}
	metricType := strings.ToLower(c.Param("metricType"))
	metricName := strings.ToLower(c.Param("metricName"))
	metricValue := strings.ToLower(c.Param("metricValue"))
	// check possible metricTypes from constants.
	if !slices.Contains([]string{models.Counter, models.Gauge}, metricType) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// check metricValues.
	if value, err = checkMetricsType(metricType, metricValue); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// save metric by type.
	switch v := value.(type) {
	case float64:
		models.Storage.SaveFloatMetric(metricName, metricType, v)
	case int64:
		models.Storage.SaveIntMetric(metricName, metricType, v)
	default:
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}

// checkMetricsType - get the metric depending on the type of int64 or float64.
func checkMetricsType(metricType, metricValue string) (value interface{}, err error) {
	switch metricType {
	case models.Counter:
		if value, err = strconv.ParseInt(metricValue, 10, 64); err != nil {
			return value, fmt.Errorf("func: checkMetricsType. strconv.ParseInt: parsing %v: invalid syntax", metricValue)
		}
	case models.Gauge:
		if value, err = strconv.ParseFloat(metricValue, 64); err != nil {
			return value, fmt.Errorf("func: checkMetricsType. strconv.ParseFloat: parsing %v: invalid syntax", metricValue)
		}
	default:
		return value, fmt.Errorf("bad request with metricType=%v, metricValue=%v", metricValue, metricValue)
	}
	return value, nil
}
