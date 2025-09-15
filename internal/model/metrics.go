package models

const (
	Counter = "counter"
	Gauge   = "gauge"
)

var Storage = NewMemStorage()

// NOTE: Не усложняем пример, вводя иерархическую вложенность структур.
// Органичиваясь плоской моделью.
// Delta и Value объявлены через указатели,
// что бы отличать значение "0", от не заданного значения
// и соответственно не кодировать в структуру.
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

type MemStorage struct {
	Store map[string]Metrics
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Store: make(map[string]Metrics),
	}
}

func (m *MemStorage) SaveFloatMetric(metricName, metricType string, metricValue float64) {
	m.Store[metricType+"_"+metricName] = Metrics{
		ID:    metricName,
		MType: Gauge,
		Value: &metricValue,
	}
}

func (m *MemStorage) SaveIntMetric(metricName, metricType string, metricValue int64) {
	m.Store[metricType+"_"+metricName] = Metrics{
		ID:    metricName,
		MType: Counter,
		Delta: &metricValue,
	}
}

func (m *MemStorage) GetMetric(metricName, metricType string) (Metrics, bool) {
	key := metricType + "_" + metricName
	metric, ok := m.Store[key]
	return metric, ok
}

func (m *MemStorage) DeleteMetric(metricName, metricType string) {
	delete(m.Store, metricType+"_"+metricName)
}
