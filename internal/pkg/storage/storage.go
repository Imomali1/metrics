package storage

import (
	"bufio"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/mailru/easyjson"
	"os"
)

type MemoryStorage interface {
	UpdateCounter(name string, counter int64) error
	UpdateGauge(name string, gauge float64) error
	GetCounterValue(name string) (int64, error)
	GetGaugeValue(name string) (float64, error)
	ListMetrics() (entity.MetricsList, error)
}

type FileStorage interface {
	WriteMetrics(metrics []entity.Metrics) error
}

type Storage struct {
	SyncWriteFile bool
	Memory        MemoryStorage
	File          FileStorage
}

func NewStorage(opts ...OptionsStorage) (*Storage, error) {
	s := &Storage{
		Memory: newMemoryStorage(),
	}
	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

type OptionsStorage func(s *Storage) error

func WithFileStorage(path string) OptionsStorage {
	return func(s *Storage) error {
		var err error
		s.File, err = newFileStorage(path)
		if err != nil {
			return err
		}
		s.SyncWriteFile = true
		return nil
	}
}

func RestoreFile(filename string) OptionsStorage {
	return func(s *Storage) error {
		file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)

		var metrics entity.MetricsList
		for scanner.Scan() {
			line := scanner.Bytes()
			var metric entity.Metrics
			if err = easyjson.Unmarshal(line, &metric); err != nil {
				return err
			}
			metrics = append(metrics, metric)
		}

		if err = scanner.Err(); err != nil {
			return err
		}

		gaugeMap := make(map[string]float64)
		counterMap := make(map[string]int64)

		for _, m := range metrics {
			if m.MType == entity.Gauge {
				gaugeMap[m.ID] = *m.Value
			} else if m.MType == entity.Counter {
				counterMap[m.ID] = *m.Delta
			}
		}

		var options []OptionsMemoryStorage
		if len(gaugeMap) != 0 {
			options = append(options, WithGaugeMap(gaugeMap))
		}

		if len(counterMap) != 0 {
			options = append(options, WithCounterMap(counterMap))
		}

		s.Memory = newMemoryStorage(options...)

		return nil
	}
}
