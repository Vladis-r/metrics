package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/Vladis-r/metrics.git/internal/model"
	"github.com/gin-gonic/gin"
)

func TestUpdate(t *testing.T) {
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
			method:   http.MethodPost,
			url:      "/update",
			want:     http.StatusOK,
		},
		{
			name:     "Test 2. Check statusOK. Slice metrics.",
			jsonData: []byte(`[{"id": "LastGC","type": "gauge","value": 1744184459}, {"id": "testMetric","type": "Counter","value": 100}]`),
			method:   http.MethodPost,
			url:      "/update",
			want:     http.StatusBadRequest,
		},
		{
			name:     "Test 3. Check StatusBadRequest. Wrong type.",
			jsonData: []byte(`{"id": "LastGC","type": "wrongType","value": 1744184459}`),
			method:   http.MethodPost,
			url:      "/update",
			want:     http.StatusBadRequest,
		},
		{
			name:     "Test 4. Check StatusBadRequest. Wrong value.",
			jsonData: []byte(`{"id": "LastGC","type": "gauge","value": "err"}`),
			method:   http.MethodPost,
			url:      "/update",
			want:     http.StatusBadRequest,
		},
		{
			name:     "Test 5. Check StatusBadRequest. Wrong json.",
			jsonData: []byte(`{"id": 'LastGC',"type": "gauge","value": "err"}`),
			method:   http.MethodPost,
			url:      "/update",
			want:     http.StatusBadRequest,
		},
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	storage := &models.MemStorage{Store: map[string]models.Metric{}}
	r.POST("/update", Update(storage))

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
