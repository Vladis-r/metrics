package agent

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Vladis-r/metrics.git/cmd/config"
	models "github.com/Vladis-r/metrics.git/internal/model"
)

// GoReportMetrics - func for send metrics to server.
func GoReportMetrics(m *models.MetricsMap, c *config.Config) {
	defer m.Wg.Done()
	ticker := time.NewTicker(time.Duration(c.ReportInterval) * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		go sendMetrics(m, c)
		fmt.Println("Send metrics. Tick at ", t)
	}
}

// sendMetrics - create url and send metrics to server.
func sendMetrics(m *models.MetricsMap, c *config.Config) (err error) {
	client := m.Client
	copyData := m.CopyData()
	for key, metricValue := range copyData {
		splittedString := strings.Split(key, "-")
		metricName, metricType := splittedString[0], splittedString[1]
		fullURL := fmt.Sprintf("http://%s%s/%s/%s/%s", c.Addr, "/update", metricType, metricName, metricValue)
		req, err := http.NewRequest("POST", fullURL, nil)
		if err != nil {
			return fmt.Errorf("func: sendMetrics; error while NewRequest: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("func: sendMetrics; error while client.Do(req): %w", err)
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("func: sendMetrics; statusCode is not OK: %w", err)
		}
	}
	return nil
}

// GoUpdateMetrics - func for update metrics.
func GoUpdateMetrics(m *models.MetricsMap, c *config.Config) {
	defer m.Wg.Done()
	ticker := time.NewTicker(time.Duration(c.PollInterval) * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		m.Mu.Lock()
		updateMetrics(m.Data)
		m.Mu.Unlock()
		fmt.Println("Update metrics. Tick at ", t)
	}
}

// updateMetrics - save metrics in data by key.
func updateMetrics(data map[string]string) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	for key := range data {
		switch key {
		// custom keys
		case "PollCount-counter":
			data[key] = getPollCounterMetric(data[key])
		case "RandomValue-gauge":
			data[key] = getRandomValueMetric()
		// runtime keys
		default:
			data[key] = getRunTimeMetrics(key, memStats)
		}
	}
}

// getRandomValueMetric - return random value in string.
func getRandomValueMetric() string {
	randInt, _ := rand.Prime(rand.Reader, 64)
	randFloat, _ := randInt.Float64()
	return fmt.Sprintf("%f", randFloat)
}

// getPollCounterMetric - count update metrics.
func getPollCounterMetric(counter string) string {
	v, err := strconv.Atoi(counter)
	if err != nil {
		v = -1
	}
	return fmt.Sprintf("%d", v+1)
}

// getRunTimeMetrics - find metric by name and return in string.
func getRunTimeMetrics(key string, memStats runtime.MemStats) string {
	name := strings.Split(key, "-")[0]
	value := fmt.Sprintf("%v", reflect.ValueOf(memStats).FieldByName(name))
	return value
}
