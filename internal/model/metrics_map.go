package models

import (
	"net/http"
	"sync"
)

// MetricsMap - struct for agent with save metrics and send them to server.
type MetricsMap struct {
	Data   map[string]Metric
	Mu     sync.RWMutex
	Wg     sync.WaitGroup
	Client *http.Client
}

func NewMetricsMap() *MetricsMap {
	return &MetricsMap{
		Data:   map[string]Metric{},
		Client: &http.Client{},
	}
}

// CopyData - return copy of data map for send on server.
func (m *MetricsMap) CopyData() []Metric {
	m.Mu.RLock()
	sl := make([]Metric, len(m.Data))
	idx := 0
	for _, v := range m.Data {
		sl[idx] = v
		idx++
	}
	m.Mu.RUnlock()
	return sl
}
