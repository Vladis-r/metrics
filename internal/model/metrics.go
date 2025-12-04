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
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty" doc:"Возвращает значение для типа counter."`
	Value *float64 `json:"value,omitempty" doc:"Возвращает значение для типа gauge."`
	Hash  string   `json:"hash,omitempty"`
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
	if ok := m.ValidateMetric(metric); !ok {
		return fmt.Errorf("func: SaveMetric; bad request. metric: %v", metric)
	}
	metric.MType = strings.ToLower(metric.MType)
	switch metric.MType {
	case Counter:
		if item, ok := m.Store[fmt.Sprintf("counter_%v", metric.ID)]; ok {
			*metric.Delta += *item.Delta
		}
		m.Store[fmt.Sprintf("counter_%v", metric.ID)] = *metric
	case Gauge:
		m.Store[fmt.Sprintf("gauge_%v", metric.ID)] = *metric
	}

	return nil
}

func (m *MemStorage) saveFloatMetric(id string, metricValue float64) {
	m.Store[id] = Metric{
		ID:    id,
		MType: Gauge,
		Value: &metricValue,
	}
}

func (m *MemStorage) saveIntMetric(id string, metricValue int64) {
	if item, ok := m.Store[id]; ok {
		metricValue += *item.Delta
	}
	m.Store[id] = Metric{
		ID:    id,
		MType: Counter,
		Delta: &metricValue,
	}
}

// ValidateMetric - check metric type and value in Metrict struct.
func (m *MemStorage) ValidateMetric(metric *Metric) bool {
	mType := strings.ToLower(metric.MType)
	if mType == Counter && metric.Delta == nil {
		return false
	}
	if mType == Gauge && metric.Value == nil {
		return false
	}
	// Check metric type.
	if !slices.Contains([]string{Counter, Gauge}, mType) {
		return false
	}
	// Limit ID is 40 characters.
	if len(metric.ID) == 0 || len(metric.ID) > 40 {
		return false
	}
	return true
}
