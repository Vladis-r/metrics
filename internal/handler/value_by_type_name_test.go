package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	models "github.com/Vladis-r/metrics.git/internal/model"
	"github.com/gin-gonic/gin"
)

func TestValueTypeName(t *testing.T) {
	tests := []struct {
		storage *models.MemStorage
		name    string
		url     string
		method  string
		save    bool
		want    int
	}{
		{
			name:   "Test 1. Check statusOK.",
			method: http.MethodGet,
			url:    "/value/Counter/testMetricName",
			save:   true,
			want:   http.StatusOK,
		},
		{
			name:   "Test 2. Check StatusBadRequest.",
			method: http.MethodGet,
			url:    "/value/wrongMetricType/testMetricName",
			want:   http.StatusBadRequest,
		},
		{
			name:   "Test 3. Check statusNotFound.",
			method: http.MethodGet,
			url:    "/value/gauge/testMetricName/12.5",
			want:   http.StatusNotFound,
		},
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	storage := &models.MemStorage{Store: map[string]models.Metric{}}
	r.GET("/value/:metricType/:metricName", ValueTypeName(storage))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storage = storage
			if tt.save {
				splitURL := strings.Split(tt.url, "/")
				tt.storage.SaveMetricByTypeValue(splitURL[3], splitURL[2], int64(99))
			}

			req, _ := http.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.want {
				t.Errorf("Wrong status: got: %v, want: %v", w.Code, tt.want)
			}
		})
	}
}
