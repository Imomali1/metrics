package storage

import (
	"bufio"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/mailru/easyjson"
	"os"
)

type Storage struct {
	SyncWriteFile bool
	Memory        MemoryStorage
	File          FileStorage
}

func NewStorage(opts ...OptionsStorage) (*Storage, error) {
	s := &Storage{}
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
		s.File, err = newFileStorage(path, s)
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

		scanner := bufio.NewScanner(file)
		if !scanner.Scan() {
			return scanner.Err()
		}

		data := scanner.Bytes()

		var metrics entity.MetricsWithoutPointerList
		err = easyjson.Unmarshal(data, &metrics)
		if err != nil {
			return err
		}

		gaugeMap := make(map[string]float64)
		counterMap := make(map[string]int64)

		for _, m := range metrics {
			if m.MType == entity.Gauge {
				gaugeMap[m.ID] = m.Value
			} else if m.MType == entity.Counter {
				counterMap[m.ID] = m.Delta
			}
		}

		var options []OptionsMemoryStorage
		if gaugeMap != nil {
			options = append(options, WithGaugeMap(gaugeMap))
		}

		if counterMap != nil {
			options = append(options, WithCounterMap(counterMap))
		}

		s.Memory = newMemoryStorage(options...)

		return nil
	}
}
