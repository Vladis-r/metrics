package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vladis-r/metrics.git/internal/middleware"
	models "github.com/Vladis-r/metrics.git/internal/model"
	"github.com/gin-gonic/gin"
)

func TestValue(t *testing.T) {
	tests := []struct {
		storage  *models.MemStorage
		name     string
		url      string
		method   string
		jsonData []byte
		want     int
	}{
		{
			name:     "Test 1. Check statusOK. Single metric.",
			jsonData: []byte(`{"id": "LastGC","type": "gauge","value": 1744184459}`),
			storage:  &models.MemStorage{},
			method:   http.MethodPost,
			url:      "/value",
			want:     http.StatusOK,
		},
		{
			name:     "Test 2. Check StatusBadRequest. Noname ID.",
			jsonData: []byte(`{"id": "noName","type": "gauge","value": 1744184459}`),
			method:   http.MethodPost,
			url:      "/value",
			want:     http.StatusNotFound,
		},
		{
			name:     "Test 3. Check unexpected type. Noname ID.",
			jsonData: []byte(`{"id": "noName","type": "unexpected","value": 1744184459}`),
			method:   http.MethodPost,
			url:      "/value",
			want:     http.StatusNotFound,
		},
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	var val1 = 1744184459.0
	var val2 = 101.0
	logger, err := middleware.InitLogger()
	if err != nil {
		t.Errorf("Error init logger: %v", err)
		return
	}
	storage := &models.MemStorage{
		Store: map[string]models.Metric{
			"LastGC":     models.Metric{ID: "LastGC", MType: "gauge", Value: &val1},
			"testMetric": models.Metric{ID: "testMetric", MType: "gauge", Value: &val2}},
		Log: logger,
	}
	r.POST("/value", Value(storage))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storage = storage

			req, _ := http.NewRequest(tt.method, tt.url, bytes.NewReader(tt.jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.want {
				t.Errorf("Wrong status: got: %v, want: %v", w.Code, tt.want)
			}
		})
	}
}
