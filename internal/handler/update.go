package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	models "github.com/Vladis-r/metrics.git/internal/model"
	"github.com/gin-gonic/gin"
)

// Update - handler for update metric with POST request and JSON in body.
func Update(s *models.MemStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		fmt.Println(string(body))
		if err != nil || len(body) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Parsed slice of metrics.
		var metrics []models.Metric
		if err := json.Unmarshal(body, &metrics); err == nil {
			errModels := []models.Metric{} // save only corrected metrics.
			for _, item := range metrics {
				if err := s.SaveMetric(&item); err != nil {
					errModels = append(errModels, item)
				}
			}
			// Bad request if one metric is wrong.
			if len(errModels) > 0 {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Wrong metrics": errModels, "count": len(errModels)})
				return
			}
			// all metrics is ok.
			c.JSON(http.StatusOK, gin.H{"data": metrics, "count": len(metrics)})
			return
		}
		// Parsed single metric.
		var metric models.Metric
		if err := json.Unmarshal(body, &metric); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Save metric to storage.
		if err := s.SaveMetric(&metric); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": metric})
	}
}
