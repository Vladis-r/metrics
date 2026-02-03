package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/Vladis-r/metrics.git/internal/model"
	"github.com/gin-gonic/gin"
)

func TestRoot(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		method  string
		storage *models.MemStorage
		want    int
	}{
		{
			name:    "Test 1. Check statusOK.",
			storage: &models.MemStorage{},
			method:  http.MethodGet,
			url:     "/",
			want:    http.StatusOK,
		},
		{
			name:    "Test 2. Check statusNotFound.",
			storage: &models.MemStorage{},
			method:  http.MethodGet,
			url:     "/index/",
			want:    http.StatusNotFound,
		},
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.LoadHTMLGlob("../../templates/*.html")
	storage := &models.MemStorage{Store: map[string]models.Metric{}}
	r.GET("/", Root(storage))

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
