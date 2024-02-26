package repository

import (
	"context"
	"github.com/Imomali1/metrics/internal/entity"
	store "github.com/Imomali1/metrics/internal/pkg/storage/v2"
)

type MetricRepository interface {
	UpdateMetrics(context.Context, entity.MetricsList) error
	GetMetrics(context.Context, entity.Metrics) (entity.Metrics, error)
	ListMetrics(context.Context) (entity.MetricsList, error)
	Ping(ctx context.Context) error
}

type Repository struct {
	MetricRepository
}

func New(memStorage store.IStorage) *Repository {
	return &Repository{
		MetricRepository: newMetricRepository(memStorage),
	}
}
