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

func Main(c *gin.Context) {
	// getExampleMetrics()
	r := prepareMetrics()

	c.HTML(http.StatusOK, "main.html", gin.H{"Items": r})
}

// getExampleMetrics - example metrics for testing.
func getExampleMetrics() {
	models.Storage.SaveFloatMetric("FirstMetric", "gauge", 0.12345)
	models.Storage.SaveIntMetric("SecondMetric", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric1", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric2", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric3", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric4", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric5", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric6", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric7", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric8", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric9", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric10", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric11", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric12", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric13", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric14", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric15", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric16", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric17", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric18", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric19", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric20", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric21", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric22", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric23", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric24", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric25", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric26", "counter", 6789)
	models.Storage.SaveIntMetric("SecondMetric27", "counter", 6789)
}

func prepareMetrics() []metricsResult {
	var val interface{}

	r := []metricsResult{}
	for _, v := range models.Storage.Store {
		val = v.Value
		if v.Value == nil {
			val = v.Delta
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
