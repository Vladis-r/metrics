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
	split_url := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(split_url) != 4 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	metric_type := split_url[1]
	if !slices.Contains([]string{models.Counter, models.Gauge}, metric_type) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metric_name := split_url[2]
	metric_value := split_url[3]
	if value, err = checkMetricsType(metric_type, metric_value); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch v := value.(type) {
	case float64:
		models.Storage.SaveFloatMetric(metric_name, metric_type, v)
	case int64:
		models.Storage.SaveIntMetric(metric_name, metric_type, v)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// checkMetricsType - get the metric depending on the type of int64 or float64.s
func checkMetricsType(metric_type, metric_value string) (value interface{}, err error) {
	switch metric_type {
	case models.Counter:
		if value, err = strconv.ParseInt(metric_value, 10, 64); err != nil {
			return value, err
		}
	case models.Gauge:
		if value, err = strconv.ParseFloat(metric_value, 64); err != nil {
			return value, err
		}
	default:
		return value, fmt.Errorf("bad request with metric_type=%v, metric_value=%v", metric_type, metric_value)
	}
	return value, nil
}
