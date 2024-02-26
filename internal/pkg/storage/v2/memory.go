package v2

import (
	"context"
	"fmt"
	"github.com/Imomali1/metrics/internal/entity"
	"sync"
)

type memoryStorage struct {
	mu             sync.RWMutex
	counterStorage map[string]int64
	gaugeStorage   map[string]float64
}

func newMemoryStorage() *memoryStorage {
	return &memoryStorage{
		counterStorage: make(map[string]int64),
		gaugeStorage:   make(map[string]float64),
	}
}

func (s *memoryStorage) Update(ctx context.Context, batch entity.MetricsList) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, one := range batch {
		if one.MType == entity.Counter {
			delta := *one.Delta
			s.counterStorage[one.ID] += delta
		} else if one.MType == entity.Gauge {
			value := *one.Value
			s.gaugeStorage[one.ID] = value
		}
	}

	return nil
}

func (s *memoryStorage) GetOne(ctx context.Context, id string, mType string) (entity.Metrics, error) {
	var metric = entity.Metrics{ID: id, MType: mType}

	s.mu.RLock()
	defer s.mu.RUnlock()
	if mType == entity.Counter {
		delta, ok := s.counterStorage[id]
		if !ok {
			return entity.Metrics{}, entity.ErrMetricNotFound
		}
		metric.Delta = &delta
	} else if mType == entity.Gauge {
		value, ok := s.gaugeStorage[id]
		if !ok {
			return entity.Metrics{}, entity.ErrMetricNotFound
		}
		metric.Value = &value
	}
	return metric, nil
}

func (s *memoryStorage) GetAll(ctx context.Context) (entity.MetricsList, error) {
	allMetrics := make(entity.MetricsList, len(s.counterStorage)+len(s.gaugeStorage))
	idx := 0

	s.mu.RLock()
	defer s.mu.RUnlock()

	for name, delta := range s.counterStorage {
		tmp := delta
		allMetrics[idx] = entity.Metrics{
			MType: entity.Counter,
			ID:    name,
			Delta: &tmp,
		}
		idx++
	}

	for name, value := range s.gaugeStorage {
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

func (s *memoryStorage) Ping(ctx context.Context) error {
	return fmt.Errorf("storage instance is not db, it is memory based")
}

func (s *memoryStorage) Close() {
	s.gaugeStorage = nil
	s.counterStorage = nil
}
