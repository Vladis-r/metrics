package handler

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	models "github.com/Vladis-r/metrics.git/internal/model"
)

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var (
		value interface{}
		err   error
	)

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	splitURL := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(splitURL) != 4 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	metricType := splitURL[1]
	if !slices.Contains([]string{models.Counter, models.Gauge}, metricType) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metricName := splitURL[2]
	metricValue := splitURL[3]
	if value, err = checkMetricsType(metricType, metricValue); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch v := value.(type) {
	case float64:
		models.Storage.SaveFloatMetric(metricName, metricType, v)
	case int64:
		models.Storage.SaveIntMetric(metricName, metricType, v)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
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
