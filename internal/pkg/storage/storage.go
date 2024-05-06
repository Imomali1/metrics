package storage

import (
	"bufio"
	"context"
	"os"

	"github.com/mailru/easyjson"

	"github.com/Imomali1/metrics/internal/entity"
)

type IStorage interface {
	Update(ctx context.Context, batch entity.MetricsList) error
	GetOne(ctx context.Context, id string, mType string) (entity.Metrics, error)
	GetAll(ctx context.Context) (entity.MetricsList, error)
	DeleteOne(ctx context.Context, id string, mType string) error
	DeleteAll(ctx context.Context) error
	Ping(ctx context.Context) error
	Close()
}

type Sync interface {
	Write(batch entity.MetricsList) error
}

type Storage struct {
	SyncWriteAllowed bool
	Sync             Sync
	IStorage
}

func NewStorage(opts ...OptionsStorage) (*Storage, error) {
	s := &Storage{}
	// By default, we use memory storage
	s.IStorage = newMemoryStorage()
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

type OptionsStorage func(s *Storage) error

func WithDB(ctx context.Context, dsn string) OptionsStorage {
	return func(s *Storage) error {
		var err error
		s.IStorage, err = newDBStorage(ctx, dsn)
		return err
	}
}

func WithSyncWrite(filename string) OptionsStorage {
	return func(s *Storage) error {
		var err error
		s.Sync, err = newFileWriter(filename)
		if err != nil {
			return err
		}
		s.SyncWriteAllowed = true
		return nil
	}
}

func RestoreFile(ctx context.Context, filename string) OptionsStorage {
	return func(s *Storage) error {
		file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)

		var i int
		m := make(map[string]int)
		var metrics entity.MetricsList
		for scanner.Scan() {
			line := scanner.Bytes()
			var metric entity.Metrics
			if err = easyjson.Unmarshal(line, &metric); err != nil {
				return err
			}
			idx, ok := m[metric.ID]
			if !ok {
				metrics = append(metrics, metric)
				m[metric.ID] = i
				i++
			} else {
				metrics[idx] = metric
			}
		}

		if err = scanner.Err(); err != nil {
			return err
		}

		if s.IStorage != nil {
			err = s.Update(ctx, metrics)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
