package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getRequest(method, url string, body io.Reader) (rr *httptest.ResponseRecorder, req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil, fmt.Errorf("Error : %v", err)
	}
	rr = httptest.NewRecorder()
	return rr, req, nil
}

func TestMainHandler(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Test 1. Check statusOK.",
			url:  "/",
			want: http.StatusOK,
		},
		{
			name: "Test 2. Check statusNotFound.",
			url:  "/index/",
			want: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr, req, err := getRequest("GET", tt.url, nil)
			if err != nil {
				t.Errorf("Error : %v", err)
			}

			MainHandler(rr, req)

			if rr.Code != tt.want {
				t.Errorf("Wrong status: got: %v, want: %v", rr.Code, http.StatusOK)
			}
		})
	}
}
