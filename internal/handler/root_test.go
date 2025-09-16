package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRoot(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		method string
		want   int
	}{
		{
			name:   "Test 1. Check statusOK.",
			method: http.MethodGet,
			url:    "/",
			want:   http.StatusOK,
		},
		{
			name:   "Test 2. Check statusNotFound.",
			method: http.MethodGet,
			url:    "/index/",
			want:   http.StatusNotFound,
		},
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.LoadHTMLGlob("../../templates/*.html")

	r.GET("/", Root)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.want {
				t.Errorf("Wrong status: got: %v, want: %v", w.Code, tt.want)
			}
		})
	}
}
