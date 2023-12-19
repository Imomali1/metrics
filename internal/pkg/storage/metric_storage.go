package storage

import (
	"errors"
	"github.com/Imomali1/metrics/internal/entity"
)

type metricStorage struct {
	counterStorage map[string]int64
	gaugeStorage   map[string]float64
}

func newMetricStorage() *metricStorage {
	return &metricStorage{
		counterStorage: make(map[string]int64),
		gaugeStorage:   make(map[string]float64),
	}
}

func (s *metricStorage) UpdateCounter(name string, counter int64) error {
	s.counterStorage[name] += counter
	return nil
}

func (s *metricStorage) UpdateGauge(name string, gauge float64) error {
	s.gaugeStorage[name] = gauge
	return nil
}

func (s *metricStorage) GetCounterValue(name string) (int64, error) {
	value, ok := s.counterStorage[name]
	if !ok {
		return 0, errors.New("not found")
	}
	return value, nil
}

func (s *metricStorage) GetGaugeValue(name string) (float64, error) {
	value, ok := s.gaugeStorage[name]
	if !ok {
		return 0, errors.New("not found")
	}
	return value, nil
}

func (s *metricStorage) ListMetrics() ([]entity.Metric, error) {
	allMetrics := make([]entity.Metric, 0)
	for name, value := range s.counterStorage {
		allMetrics = append(allMetrics, entity.Metric{
			Type:  "counter",
			Name:  name,
			Value: value,
		})
	}

	for name, value := range s.gaugeStorage {
		allMetrics = append(allMetrics, entity.Metric{
			Type:  "counter",
			Name:  name,
			Value: value,
		})
	}

	return allMetrics, nil
}
