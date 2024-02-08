package storage

import (
	"github.com/Imomali1/metrics/internal/entity"
)

type MemoryStorage interface {
	UpdateCounter(name string, counter int64) error
	UpdateGauge(name string, gauge float64) error
	GetCounterValue(name string) (int64, error)
	GetGaugeValue(name string) (float64, error)
	ListMetrics() (entity.MetricsWithoutPointerList, error)
}

type memoryStorage struct {
	counterStorage map[string]int64
	gaugeStorage   map[string]float64
}

type OptionsMemoryStorage func(m *memoryStorage)

func newMemoryStorage(opts ...OptionsMemoryStorage) *memoryStorage {
	m := &memoryStorage{
		counterStorage: make(map[string]int64),
		gaugeStorage:   make(map[string]float64),
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func WithCounterMap(counterMap map[string]int64) OptionsMemoryStorage {
	return func(m *memoryStorage) {
		m.counterStorage = counterMap
	}
}

func WithGaugeMap(gaugeMap map[string]float64) OptionsMemoryStorage {
	return func(m *memoryStorage) {
		m.gaugeStorage = gaugeMap
	}
}

func (s *memoryStorage) UpdateCounter(name string, counter int64) error {
	s.counterStorage[name] += counter
	return nil
}

func (s *memoryStorage) UpdateGauge(name string, gauge float64) error {
	s.gaugeStorage[name] = gauge
	return nil
}

func (s *memoryStorage) GetCounterValue(name string) (int64, error) {
	value, ok := s.counterStorage[name]
	if !ok {
		return 0, entity.ErrMetricNotFound
	}
	return value, nil
}

func (s *memoryStorage) GetGaugeValue(name string) (float64, error) {
	value, ok := s.gaugeStorage[name]
	if !ok {
		return 0, entity.ErrMetricNotFound
	}
	return value, nil
}

func (s *memoryStorage) ListMetrics() (entity.MetricsWithoutPointerList, error) {
	allMetrics := make(entity.MetricsWithoutPointerList, 0)
	for name, delta := range s.counterStorage {
		allMetrics = append(allMetrics, entity.MetricsWithoutPointer{
			MType: entity.Counter,
			ID:    name,
			Delta: delta,
		})
	}

	for name, value := range s.gaugeStorage {
		allMetrics = append(allMetrics, entity.MetricsWithoutPointer{
			MType: entity.Gauge,
			ID:    name,
			Value: value,
		})
	}

	return allMetrics, nil
}
