package agent

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
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
	copyData := m.CopyData()

	jsonData, err := json.Marshal(copyData)
	if err != nil {
		return fmt.Errorf("func: sendMetrics; error while json.Marshal(copyData): %w", err)
	}
	fullURL := fmt.Sprintf("http://%s/%s", c.Addr, "update")
	resp, err := http.Post(fullURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("func: sendMetrics; error while http.Post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("func: sendMetrics; statusCode is not OK: %w", err)
	}

	return nil
}

// GoUpdateMetrics - func for update metrics.
func GoUpdateMetrics(m *models.MetricsMap, cfg *config.Config) {
	defer m.Wg.Done()
	ticker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		m.Mu.Lock()
		updateMetrics(m.Data)
		m.Mu.Unlock()
		fmt.Println("Update metrics. Tick at ", t)
	}
}

// updateMetrics - save metrics in data by key.
func updateMetrics(data map[string]models.Metric) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	upRuntimeMetrics(data, memStats)
	upPollCounterMetric(data)
	upRandomValueMetric(data)
}

func upRuntimeMetrics(data map[string]models.Metric, memStats runtime.MemStats) {
	floatptr := func(v float64) *float64 { return &v }
	intptr := func(v int64) *int64 { return &v }

	data["Alloc"] = models.Metric{ID: "Alloc", MType: "counter", Delta: intptr(int64(memStats.Alloc))}
	data["TotalAlloc"] = models.Metric{ID: "TotalAlloc", MType: "counter", Delta: intptr(int64(memStats.TotalAlloc))}
	data["Sys"] = models.Metric{ID: "Sys", MType: "counter", Delta: intptr(int64(memStats.Sys))}
	data["Lookups"] = models.Metric{ID: "Lookups", MType: "counter", Delta: intptr(int64(memStats.Lookups))}
	data["Mallocs"] = models.Metric{ID: "Mallocs", MType: "counter", Delta: intptr(int64(memStats.Mallocs))}
	data["Frees"] = models.Metric{ID: "Frees", MType: "counter", Delta: intptr(int64(memStats.Frees))}
	data["HeapAlloc"] = models.Metric{ID: "HeapAlloc", MType: "counter", Delta: intptr(int64(memStats.HeapAlloc))}
	data["HeapSys"] = models.Metric{ID: "HeapSys", MType: "counter", Delta: intptr(int64(memStats.HeapSys))}
	data["HeapIdle"] = models.Metric{ID: "HeapIdle", MType: "counter", Delta: intptr(int64(memStats.HeapIdle))}
	data["HeapInuse"] = models.Metric{ID: "HeapInuse", MType: "counter", Delta: intptr(int64(memStats.HeapInuse))}
	data["HeapReleased"] = models.Metric{ID: "HeapReleased", MType: "counter", Delta: intptr(int64(memStats.HeapReleased))}
	data["HeapObjects"] = models.Metric{ID: "HeapObjects", MType: "counter", Delta: intptr(int64(memStats.HeapObjects))}
	data["StackInuse"] = models.Metric{ID: "StackInuse", MType: "counter", Delta: intptr(int64(memStats.StackInuse))}
	data["StackSys"] = models.Metric{ID: "StackSys", MType: "counter", Delta: intptr(int64(memStats.StackSys))}
	data["MSpanInuse"] = models.Metric{ID: "MSpanInuse", MType: "counter", Delta: intptr(int64(memStats.MSpanInuse))}
	data["MSpanSys"] = models.Metric{ID: "MSpanSys", MType: "counter", Delta: intptr(int64(memStats.MSpanSys))}
	data["MCacheInuse"] = models.Metric{ID: "MCacheInuse", MType: "counter", Delta: intptr(int64(memStats.MCacheInuse))}
	data["MCacheSys"] = models.Metric{ID: "MCacheSys", MType: "counter", Delta: intptr(int64(memStats.MCacheSys))}
	data["BuckHashSys"] = models.Metric{ID: "BuckHashSys", MType: "counter", Delta: intptr(int64(memStats.BuckHashSys))}
	data["GCSys"] = models.Metric{ID: "GCSys", MType: "counter", Delta: intptr(int64(memStats.GCSys))}
	data["OtherSys"] = models.Metric{ID: "OtherSys", MType: "counter", Delta: intptr(int64(memStats.OtherSys))}
	data["NextGC"] = models.Metric{ID: "NextGC", MType: "counter", Delta: intptr(int64(memStats.NextGC))}
	data["LastGC"] = models.Metric{ID: "LastGC", MType: "counter", Delta: intptr(int64(memStats.LastGC))}
	data["PauseTotalNs"] = models.Metric{ID: "PauseTotalNs", MType: "counter", Delta: intptr(int64(memStats.PauseTotalNs))}
	data["NumGC"] = models.Metric{ID: "NumGC", MType: "counter", Delta: intptr(int64(memStats.NumGC))}
	data["GCCPUFraction"] = models.Metric{ID: "GCCPUFraction", MType: "gauge", Value: floatptr(float64(memStats.GCCPUFraction))}
}

// getRandomValueMetric - The RandomValue metric. A random number is generated each time.
func upRandomValueMetric(data map[string]models.Metric) {
	var key = "RandomValue"
	randInt, _ := rand.Prime(rand.Reader, 64)
	randFloat, _ := randInt.Float64()
	data[key] = models.Metric{ID: key, MType: "counter", Value: &randFloat}
}

// getPollCounterMetric - The PollCount metric. Counts the number of updates.
func upPollCounterMetric(data map[string]models.Metric) {
	var (
		counter int64
		key     = "PollCount"
	)
	if _, ok := data[key]; ok {
		counter = *data[key].Delta
	}
	counter++
	data[key] = models.Metric{ID: key, MType: "counter", Delta: &counter}
}
