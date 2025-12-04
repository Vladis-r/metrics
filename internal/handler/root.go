package handler

import (
	"net/http"
	"sort"

	models "github.com/Vladis-r/metrics.git/internal/model"
	"github.com/gin-gonic/gin"
)

type metricsResult struct {
	ID    string
	MType string
	Value interface{}
}

func Root(s *models.MemStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := prepareMetrics(s)

		c.HTML(http.StatusOK, "main.html", gin.H{"Items": r})
	}
}

func prepareMetrics(s *models.MemStorage) []metricsResult {
	var val interface{}

	r := []metricsResult{}
	for _, v := range s.Store {
		val = *v.Value
		if v.Value == nil {
			val = *v.Delta
		}
		add := metricsResult{
			ID:    v.ID,
			MType: v.MType,
			Value: val,
		}
		r = append(r, add)
	}

	sort.Slice(r, func(i, j int) bool {
		return r[i].ID < r[j].ID
	})

	return r
}
