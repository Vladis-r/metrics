package agent

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/Vladis-r/metrics.git/cmd/config"
	models "github.com/Vladis-r/metrics.git/internal/model"
	ut "github.com/Vladis-r/metrics.git/internal/utils"
)

// GoReportMetrics - func for send metrics to server.
func GoReportMetrics(m *models.MetricsMap, c *config.ConfigAgent) {
	defer m.Wg.Done()
	ticker := time.NewTicker(time.Duration(c.ReportInterval) * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		go sendGzipMetrics(m, c)
		fmt.Println("Send metrics. Tick at ", t)
	}
}

// sendGzipMetrics - send metric with gzip compress.
func sendGzipMetrics(m *models.MetricsMap, c *config.ConfigAgent) (err error) {
	copyData := m.CopyData()

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if err := json.NewEncoder(gz).Encode(copyData); err != nil {
		return fmt.Errorf("func: sendMetrics; error while Encode(copyData): %w", err)
	}
	if err := gz.Close(); err != nil {
		return fmt.Errorf("func: sendMetrics; error while gz.Close(): %w", err)
	}

	fullURL := fmt.Sprintf("http://%s/%s", c.Addr, "update")
	req, err := http.NewRequest("POST", fullURL, &buf)
	if err != nil {
		return fmt.Errorf("func: sendMetrics; error while http.NewRequest: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := m.Client.Do(req)
	if err != nil {
		return fmt.Errorf("func: sendMetrics; error while Client.Do: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("func: sendMetrics; statusCode is not OK: %w", err)
	}

	return nil
}

// GoUpdateMetrics - func for update metrics.
func GoUpdateMetrics(m *models.MetricsMap, cfg *config.ConfigAgent) {
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
	upCounterMetric(data, "PollCount")
	upRandomValueMetric(data)
}

func upRuntimeMetrics(data map[string]models.Metric, memStats runtime.MemStats) {
	data["Alloc"] = models.Metric{ID: "Alloc", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.Alloc))}
	data["TotalAlloc"] = models.Metric{ID: "TotalAlloc", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.TotalAlloc))}
	data["Sys"] = models.Metric{ID: "Sys", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.Sys))}
	data["Lookups"] = models.Metric{ID: "Lookups", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.Lookups))}
	data["Mallocs"] = models.Metric{ID: "Mallocs", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.Mallocs))}
	data["Frees"] = models.Metric{ID: "Frees", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.Frees))}
	data["HeapAlloc"] = models.Metric{ID: "HeapAlloc", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.HeapAlloc))}
	data["HeapSys"] = models.Metric{ID: "HeapSys", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.HeapSys))}
	data["HeapIdle"] = models.Metric{ID: "HeapIdle", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.HeapIdle))}
	data["HeapInuse"] = models.Metric{ID: "HeapInuse", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.HeapInuse))}
	data["HeapReleased"] = models.Metric{ID: "HeapReleased", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.HeapReleased))}
	data["HeapObjects"] = models.Metric{ID: "HeapObjects", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.HeapObjects))}
	data["StackInuse"] = models.Metric{ID: "StackInuse", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.StackInuse))}
	data["StackSys"] = models.Metric{ID: "StackSys", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.StackSys))}
	data["MSpanInuse"] = models.Metric{ID: "MSpanInuse", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.MSpanInuse))}
	data["MSpanSys"] = models.Metric{ID: "MSpanSys", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.MSpanSys))}
	data["MCacheInuse"] = models.Metric{ID: "MCacheInuse", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.MCacheInuse))}
	data["MCacheSys"] = models.Metric{ID: "MCacheSys", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.MCacheSys))}
	data["BuckHashSys"] = models.Metric{ID: "BuckHashSys", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.BuckHashSys))}
	data["GCSys"] = models.Metric{ID: "GCSys", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.GCSys))}
	data["OtherSys"] = models.Metric{ID: "OtherSys", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.OtherSys))}
	data["NextGC"] = models.Metric{ID: "NextGC", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.NextGC))}
	data["LastGC"] = models.Metric{ID: "LastGC", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.LastGC))}
	data["PauseTotalNs"] = models.Metric{ID: "PauseTotalNs", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.PauseTotalNs))}
	data["NumGC"] = models.Metric{ID: "NumGC", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.NumGC))}
	data["GCCPUFraction"] = models.Metric{ID: "GCCPUFraction", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.GCCPUFraction))}
	data["NumForcedGC"] = models.Metric{ID: "NumForcedGC", MType: "gauge", Value: ut.Float64Ptr(float64(memStats.NumForcedGC))}
}

// getRandomValueMetric - The RandomValue metric. A random number is generated each time.
func upRandomValueMetric(data map[string]models.Metric) {
	var key = "RandomValue"
	randInt, _ := rand.Prime(rand.Reader, 64)
	randFloat, _ := randInt.Float64()
	data[key] = models.Metric{ID: key, MType: "gauge", Value: &randFloat}
}

// getPollCounterMetric - The PollCount metric. Counts the number of updates.
func upCounterMetric(data map[string]models.Metric, key string) {
	var counter int64
	if _, ok := data[key]; ok {
		counter = *data[key].Delta
	}
	counter++
	data[key] = models.Metric{ID: key, MType: "counter", Delta: &counter}
}
