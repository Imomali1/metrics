package repository

import (
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/storage"
)

type MetricRepository interface {
	UpdateCounter(name string, counter int64) error
	UpdateGauge(name string, gauge float64) error
	GetCounterValue(name string) (int64, error)
	GetGaugeValue(name string) (float64, error)
	ListMetrics() ([]entity.MetricsWithoutPointer, error)
}

type Repository struct {
	MetricRepository
}

func New(memStorage *storage.Storage) *Repository {
	return &Repository{
		MetricRepository: newMetricRepository(memStorage),
	}
}
