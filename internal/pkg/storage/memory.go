package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/Imomali1/metrics/internal/entity"
)

type Memory struct {
	mu             sync.RWMutex
	CounterStorage map[string]int64
	GaugeStorage   map[string]float64
}

func NewMemory() (Storage, error) {
	return &Memory{
		CounterStorage: make(map[string]int64),
		GaugeStorage:   make(map[string]float64),
	}, nil
}

func (s *Memory) DeleteOne(_ context.Context, id string, mType string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	switch mType {
	case entity.Gauge:
		delete(s.GaugeStorage, id)
	case entity.Counter:
		delete(s.CounterStorage, id)
	default:
		return fmt.Errorf("unknown type: %s", mType)
	}
	return nil
}

func (s *Memory) DeleteAll(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.CounterStorage = make(map[string]int64)
	s.GaugeStorage = make(map[string]float64)
	return nil
}

func (s *Memory) Update(_ context.Context, batch entity.MetricsList) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, one := range batch {
		if one.MType == entity.Counter {
			delta := *one.Delta
			s.CounterStorage[one.ID] += delta
		} else if one.MType == entity.Gauge {
			value := *one.Value
			s.GaugeStorage[one.ID] = value
		}
	}

	return nil
}

func (s *Memory) GetOne(_ context.Context, id string, mType string) (entity.Metrics, error) {
	var metric = entity.Metrics{ID: id, MType: mType}

	s.mu.RLock()
	defer s.mu.RUnlock()
	if mType == entity.Counter {
		delta, ok := s.CounterStorage[id]
		if !ok {
			return entity.Metrics{}, entity.ErrMetricNotFound
		}
		metric.Delta = &delta
	} else if mType == entity.Gauge {
		value, ok := s.GaugeStorage[id]
		if !ok {
			return entity.Metrics{}, entity.ErrMetricNotFound
		}
		metric.Value = &value
	}
	return metric, nil
}

func (s *Memory) GetAll(_ context.Context) (entity.MetricsList, error) {
	allMetrics := make(entity.MetricsList, len(s.CounterStorage)+len(s.GaugeStorage))
	idx := 0

	s.mu.RLock()
	defer s.mu.RUnlock()

	for name, delta := range s.CounterStorage {
		tmp := delta
		allMetrics[idx] = entity.Metrics{
			MType: entity.Counter,
			ID:    name,
			Delta: &tmp,
		}
		idx++
	}

	for name, value := range s.GaugeStorage {
		tmp := value
		allMetrics[idx] = entity.Metrics{
			MType: entity.Gauge,
			ID:    name,
			Value: &tmp,
		}
		idx++
	}

	return allMetrics, nil
}

func (s *Memory) Ping(_ context.Context) error {
	return fmt.Errorf("storage instance is not database, it is memory based")
}

func (s *Memory) Close() {
	s.GaugeStorage = nil
	s.CounterStorage = nil
}
