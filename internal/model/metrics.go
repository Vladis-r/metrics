package models

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/Vladis-r/metrics.git/cmd/config"
	"go.uber.org/zap"
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
	Conf  *config.ConfigServer
	Log   *zap.Logger
	Mu    sync.RWMutex
}

func NewMemStorage(config *config.ConfigServer, logger *zap.Logger) *MemStorage {
	return &MemStorage{
		Store: make(map[string]Metric),
		Conf:  config,
		Log:   logger,
	}
}

func (s *MemStorage) GetMetric(id string) (Metric, bool) {
	metric, ok := s.Store[id]
	return metric, ok
}

func (s *MemStorage) DeleteMetric(id, mType string) {
	delete(s.Store, id)
}

// SaveMetricByTypeValue - save metric by type and value.
func (s *MemStorage) SaveMetricByTypeValue(id, mType string, value interface{}) (err error) {
	switch v := value.(type) {
	case float64:
		s.saveFloatMetric(id, v)
	case int64:
		s.saveIntMetric(id, v)
	default:
		return fmt.Errorf("func: SaveMetricByTypeValue; bad request. id: %v, mType: %v, value: %v . err: %v", id, mType, value, err)
	}
	return nil
}

// SaveMetric - save metric by struct Metric.
func (s *MemStorage) SaveMetric(metric *Metric) error {
	if ok := s.ValidateMetric(metric); !ok {
		return fmt.Errorf("func: SaveMetric; bad request. metric: %v", metric)
	}
	metric.MType = strings.ToLower(metric.MType)
	switch metric.MType {
	case Counter:
		if item, ok := s.Store[metric.ID]; ok {
			*metric.Delta += *item.Delta
		}
		s.Store[metric.ID] = *metric
	case Gauge:
		s.Store[metric.ID] = *metric
	default:
		s.Log.Info("Unexpected type", zap.String("MType", metric.MType))
	}

	return nil
}

func (s *MemStorage) saveFloatMetric(id string, metricValue float64) {
	s.Store[id] = Metric{
		ID:    id,
		MType: Gauge,
		Value: &metricValue,
	}
}

func (s *MemStorage) saveIntMetric(id string, metricValue int64) {
	if item, ok := s.Store[id]; ok {
		metricValue += *item.Delta
	}
	s.Store[id] = Metric{
		ID:    id,
		MType: Counter,
		Delta: &metricValue,
	}
}

// ValidateMetric - check metric type and value in Metrict struct.
func (s *MemStorage) ValidateMetric(metric *Metric) bool {
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
