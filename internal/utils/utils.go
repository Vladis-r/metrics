package utils

import (
	"fmt"
	"strconv"

	models "github.com/Vladis-r/metrics.git/internal/model"
)

// CheckMetric - check possible metricType and metricValue.
func CheckMetric(mType, mValue string) (value interface{}, err error) {
	switch mType {
	case models.Counter:
		if value, err = strconv.ParseInt(mValue, 10, 64); err != nil {
			return value, fmt.Errorf("func: CheckMetric. strconv.ParseInt: parsing %v: invalid syntax", mValue)
		}
	case models.Gauge:
		if value, err = strconv.ParseFloat(mValue, 64); err != nil {
			return value, fmt.Errorf("func: CheckMetric. strconv.ParseFloat: parsing %v: invalid syntax", mValue)
		}
	default:
		return value, fmt.Errorf("bad request with metricType=%v, metricValue=%v", mType, mValue)
	}
	return value, nil
}
