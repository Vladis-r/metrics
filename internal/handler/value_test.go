package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	models "github.com/Vladis-r/metrics.git/internal/model"
	"github.com/gin-gonic/gin"
)

func TestValue(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		method string
		save   bool
		want   int
	}{
		{
			name:   "Test 1. Check statusOK.",
			method: http.MethodGet,
			url:    "/value/testMetricType/testMetricName",
			save:   true,
			want:   http.StatusOK,
		},
		{
			name:   "Test 2. Check statusNotFound.",
			method: http.MethodGet,
			url:    "/value/wrongMetricType/testMetricName",
			want:   http.StatusNotFound,
		},
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.LoadHTMLGlob("../../templates/*.html")

	r.GET("/value/:metricType/:metricName", Value)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.save {
				m := models.Storage
				splitUrl := strings.Split(tt.url, "/")
				m.SaveFloatMetric(strings.ToLower(splitUrl[3]), strings.ToLower(splitUrl[2]), 99)
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
