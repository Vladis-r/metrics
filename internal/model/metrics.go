package models

import (
	"fmt"
	"slices"
	"strings"
)

const (
	Counter = "counter"
	Gauge   = "gauge"
)

// NOTE: Не усложняем пример, вводя иерархическую вложенность структур.
// Органичиваясь плоской моделью.
// Delta и Value объявлены через указатели,
// что бы отличать значение "0", от не заданного значения
// и соответственно не кодировать в структуру.
type Metric struct {
	ID       string   `json:"id"`
	MType    string   `json:"type"`
	Delta    *int64   `json:"delta,omitempty"`
	Value    *float64 `json:"value,omitempty"`
	Hash     string   `json:"hash,omitempty"`
	ValueSum float64
	DeltaSum int64
}

type MemStorage struct {
	Store map[string]Metric
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Store: make(map[string]Metric),
	}
}

func (m *MemStorage) GetMetric(id string) (Metric, bool) {
	metric, ok := m.Store[id]
	return metric, ok
}

func (m *MemStorage) DeleteMetric(id, mType string) {
	delete(m.Store, id)
}

// SaveMetricByTypeValue - save metric by type and value.
func (m *MemStorage) SaveMetricByTypeValue(id, mType string, value interface{}) (err error) {
	switch v := value.(type) {
	case float64:
		m.saveFloatMetric(id, v)
	case int64:
		m.saveIntMetric(id, v)
	default:
		return fmt.Errorf("func: SaveMetricByTypeValue; bad request. id: %v, mType: %v, value: %v . err: %v", id, mType, value, err)
	}
	return nil
}

// SaveMetric - save metric by struct Metric.
func (m *MemStorage) SaveMetric(metric *Metric) error {
	// Validation metric.
	if ok := m.validateMetric(metric); !ok {
		return fmt.Errorf("func: SaveMetric; bad request. metric: %v", metric)
	}
	// Save metric.
	metric.ID = strings.ToLower(metric.ID)
	metric.MType = strings.ToLower(metric.MType)
	m.Store[metric.ID] = *metric
	return nil
}

func (m *MemStorage) saveFloatMetric(id string, metricValue float64) {
	if metric, ok := m.Store[id]; ok {
		valueSum := metric.ValueSum + metricValue
		m.Store[id] = Metric{
			ID:       id,
			MType:    Gauge,
			Value:    &metricValue,
			ValueSum: valueSum,
		}
	} else {
		m.Store[id] = Metric{
			ID:       id,
			MType:    Gauge,
			Value:    &metricValue,
			ValueSum: metricValue,
		}
	}
}

func (m *MemStorage) saveIntMetric(id string, metricValue int64) {
	if metric, ok := m.Store[id]; ok {
		deltaSum := metric.DeltaSum + metricValue
		m.Store[id] = Metric{
			ID:       id,
			MType:    Counter,
			Delta:    &metricValue,
			DeltaSum: deltaSum,
		}
	} else {
		m.Store[id] = Metric{
			ID:       id,
			MType:    Counter,
			Delta:    &metricValue,
			DeltaSum: metricValue,
		}
	}
}

// ValidateMetric - check metric type and value in Metrict struct.
func (m *MemStorage) validateMetric(metric *Metric) bool {
	// Check metric type.
	if !slices.Contains([]string{Counter, Gauge}, strings.ToLower(metric.MType)) {
		return false
	}
	// Check that value is exist.
	if metric.Value == nil && metric.Delta == nil {
		return false
	}
	// Limit ID is 30 characters.
	if len(metric.ID) > 30 {
		return false
	}
	return true
}
