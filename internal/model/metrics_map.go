package models

import (
	"net/http"
	"sync"
)

// MetricsMap - struct for agent with save metrics and send them to server.
type MetricsMap struct {
	Data   map[string]string
	Mu     sync.RWMutex
	Wg     sync.WaitGroup
	Client *http.Client
}

func NewMetricsMap() *MetricsMap {
	return &MetricsMap{
		Data: map[string]string{
			// Runtime
			"Alloc-counter":        "value",
			"BuckHashSys-counter":  "value",
			"Frees-counter":        "value",
			"GCCPUFraction-gauge":  "value",
			"GCSys-counter":        "value",
			"HeapAlloc-counter":    "value",
			"HeapIdle-counter":     "value",
			"HeapInuse-counter":    "value",
			"HeapObjects-counter":  "value",
			"HeapReleased-counter": "value",
			"HeapSys-counter":      "value",
			"LastGC-counter":       "value",
			"Lookups-counter":      "value",
			"MCacheInuse-counter":  "value",
			"MCacheSys-counter":    "value",
			"MSpanInuse-counter":   "value",
			"MSpanSys-counter":     "value",
			"Mallocs-counter":      "value",
			"NextGC-counter":       "value",
			"NumForcedGC-counter":  "value",
			"NumGC-counter":        "value",
			"OtherSys-counter":     "value",
			"PauseTotalNs-counter": "value",
			"StackInuse-counter":   "value",
			"StackSys-counter":     "value",
			"Sys-counter":          "value",
			"TotalAlloc-counter":   "value",
			// Custom
			"PollCount-counter": "value",
			"RandomValue-gauge": "value",
		},
		Client: &http.Client{},
	}
}

// CopyData - return copy of data map for send on server.
func (m *MetricsMap) CopyData() map[string]string {
	m.Mu.RLock()
	newMap := make(map[string]string, len(m.Data))
	for k, v := range m.Data {
		newMap[k] = v
	}
	m.Mu.RUnlock()
	return newMap
}
