package v2

import (
	"context"
	"github.com/Imomali1/metrics/internal/entity"
)

type IStorage interface {
	Update(ctx context.Context, batch entity.MetricsList) error
	GetOne(ctx context.Context, id string, mType string) (entity.Metrics, error)
	GetAll(ctx context.Context) (entity.MetricsList, error)
	Ping(ctx context.Context) error
	Close()
}

func NewStorage(opts ...OptionsStorage) (IStorage, error) {
	var s storageOptions
	for _, opt := range opts {
		if err := opt(&s); err != nil {
			return nil, err
		}
	}

	if s.db.allowed {
		return newDBStorage(context.Background(), s.db.dsn)
	}

	if s.file.allowed {
		return newFileStorage(s.file.path)
	}

	return newMemoryStorage(
		s.memory.counterMap,
		s.memory.gaugeMap,
		s.memory.syncWrite,
		s.memory.filepath,
	)
}
