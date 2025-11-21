package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/Vladis-r/metrics.git/internal/model"
	"github.com/gin-gonic/gin"
)

func TestUpdateTypeNameValue(t *testing.T) {
	tests := []struct {
		storage *models.MemStorage
		name    string
		method  string
		url     string
		want    int
	}{
		{
			name:   "Test 1. Check wrong method",
			method: "GET",
			url:    "/update/wrongType/testMetric/101",
			want:   http.StatusNotFound,
		},
		{
			name:   "Test 3. Check wrong type",
			method: http.MethodPost,
			url:    "/update/wrongType/testMetric/101",
			want:   http.StatusBadRequest,
		},
		{
			name:   "Test 4. Check wrong value",
			method: http.MethodPost,
			url:    "/update/gauge/testMetric/errorValue",
			want:   http.StatusBadRequest,
		},
		{
			name:   "Test 5. Check status ok",
			method: http.MethodPost,
			url:    "/update/counter/testMetric/101",
			want:   http.StatusOK,
		},
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	storage := &models.MemStorage{Store: map[string]models.Metric{}}
	r.POST("/update/:metricType/:metricName/:metricValue", UpdateTypeNameValue(storage))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storage = storage

			req, _ := http.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.want {
				t.Errorf("Wrong status: got: %v, want: %v", w.Code, tt.want)
			}
		})
	}
}
