package handler

import (
	"net/http"
	"testing"
)

func TestUpdateHandler(t *testing.T) {
	tests := []struct {
		name   string
		method string
		url    string
		want   int
	}{
		{
			name:   "Test 1. Check wrong method",
			method: "GET",
			url:    "/update/",
			want:   http.StatusMethodNotAllowed,
		},
		{
			name:   "Test 2. Check wrong url",
			method: http.MethodPost,
			url:    "/update/",
			want:   http.StatusNotFound,
		},
		{
			name:   "Test 3. Check wrong type",
			method: http.MethodPost,
			url:    "/update/wrongType/testMetric/101",
			want:   http.StatusBadRequest,
		},
		{
			name:   "Test 4. Check status ok",
			method: http.MethodPost,
			url:    "/update/counter/testMetric/101",
			want:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr, req, err := getRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Errorf("Error : %v", err)
			}
			UpdateHandler(rr, req)

			if rr.Code != tt.want {
				t.Errorf("Wrong status: got: %v, want: %v", rr.Code, tt.want)
			}
		})
	}
}

func TestCheckMetricsType(t *testing.T) {
	tests := []struct {
		name        string
		metricType  string
		metricValue string
		want        interface{}
		checkError  bool
	}{
		{
			name:        "Test 1. Check wrong metricType",
			metricType:  "wrong type",
			metricValue: "101",
			want:        nil,
			checkError:  true,
		},
		{
			name:        "Test 2. Check int64 metricType",
			metricType:  "counter",
			metricValue: "101",
			want:        int64(101),
		},
		{
			name:        "Test 3. Check float64 metricType",
			metricType:  "gauge",
			metricValue: "5.201",
			want:        5.201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := checkMetricsType(tt.metricType, tt.metricValue)
			if err != nil && !tt.checkError {
				t.Errorf("Error: %v", err)
			}
			if value != tt.want {
				t.Errorf("Wrong status: got: %v, want: %v", value, tt.want)
			}
		})
	}
}
