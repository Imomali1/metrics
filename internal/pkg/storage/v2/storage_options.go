package v2

import (
	"bufio"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/mailru/easyjson"
	"os"
)

type dbOptions struct {
	allowed bool
	dsn     string
}

type fileOptions struct {
	allowed bool
	path    string
}

type memoryOptions struct {
	allowed    bool
	syncWrite  bool
	filepath   string
	counterMap map[string]int64
	gaugeMap   map[string]float64
}

type storageOptions struct {
	db     dbOptions
	file   fileOptions
	memory memoryOptions
}

type OptionsStorage func(s *storageOptions) error

func WithDB(dsn string) OptionsStorage {
	return func(s *storageOptions) error {
		s.db.allowed = true
		s.db.dsn = dsn
		return nil
	}
}

func WithFileStorage(path string) OptionsStorage {
	return func(s *storageOptions) error {
		s.file.allowed = true
		s.file.path = path
		return nil
	}
}

func WithMemoryStorage(syncWrite bool, restore bool, filename string) OptionsStorage {
	return func(s *storageOptions) error {
		s.memory.allowed = true

		if syncWrite {
			s.memory.syncWrite = true
			s.memory.filepath = filename
		}

		if !restore {
			s.memory.counterMap = make(map[string]int64)
			s.memory.gaugeMap = make(map[string]float64)
			return nil
		}

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

		s.memory.counterMap = counterMap
		s.memory.gaugeMap = gaugeMap
		return nil
	}
}
