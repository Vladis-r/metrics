package utils

import "testing"

func TestCheckMetric(t *testing.T) {
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
			value, err := CheckMetric(tt.metricType, tt.metricValue)
			if err != nil && !tt.checkError {
				t.Errorf("Error: %v", err)
			}
			if value != tt.want {
				t.Errorf("Wrong status: got: %v, want: %v", value, tt.want)
			}
		})
	}
}
