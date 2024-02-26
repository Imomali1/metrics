package v2

import (
	"bufio"
	"context"
	"fmt"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/mailru/easyjson"
	"os"
	"sync"
)

type memFileWriter struct {
	file   *os.File
	writer *bufio.Writer
}

type memoryStorage struct {
	mu             sync.RWMutex
	counterStorage map[string]int64
	gaugeStorage   map[string]float64
	syncWrite      bool
	fw             memFileWriter
}

func newMemoryStorage(counterMap map[string]int64, gaugeMap map[string]float64, syncWrite bool, filepath string) (*memoryStorage, error) {
	fileWriter, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &memoryStorage{
		counterStorage: counterMap,
		gaugeStorage:   gaugeMap,
		syncWrite:      syncWrite,
		fw: memFileWriter{
			file:   fileWriter,
			writer: bufio.NewWriter(fileWriter),
		},
	}, nil
}

func (s *memoryStorage) Update(ctx context.Context, batch entity.MetricsList) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, one := range batch {
		if one.MType == entity.Counter {
			s.counterStorage[one.ID] += *one.Delta
		} else if one.MType == entity.Gauge {
			s.gaugeStorage[one.ID] = *one.Value
		}
	}

	if s.syncWrite {
		for _, one := range batch {
			data, err := easyjson.Marshal(one)
			if err != nil {
				return err
			}

			_, err = s.fw.writer.Write(data)
			if err != nil {
				return err
			}

			if err = s.fw.writer.WriteByte('\n'); err != nil {
				return err
			}
		}

		return s.fw.writer.Flush()
	}

	return nil
}

func (s *memoryStorage) GetOne(ctx context.Context, id string, mType string) (entity.Metrics, error) {
	var ok bool
	delta, value := new(int64), new(float64)

	s.mu.RLock()
	defer s.mu.RUnlock()
	if mType == entity.Counter {
		*delta, ok = s.counterStorage[id]
	} else if mType == entity.Gauge {
		*value, ok = s.gaugeStorage[id]
	}

	if !ok {
		return entity.Metrics{}, entity.ErrMetricNotFound
	}

	return entity.Metrics{
		ID:    id,
		MType: mType,
		Delta: delta,
		Value: value,
	}, nil
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
