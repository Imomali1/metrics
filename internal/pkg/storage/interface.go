package storage

import (
	"context"
	"github.com/Imomali1/metrics/internal/entity"
)

type Storage interface {
	Update(ctx context.Context, batch entity.MetricsList) error
	GetOne(ctx context.Context, id string, mType string) (entity.Metrics, error)
	GetAll(ctx context.Context) (entity.MetricsList, error)
	DeleteOne(ctx context.Context, id string, mType string) error
	DeleteAll(ctx context.Context) error
	Ping(ctx context.Context) error
	Close()
}
