package services

import (
	"context"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/repository"
)

type MetricService interface {
	UpdateCounter(name string, counter int64) error
	UpdateGauge(name string, gauge float64) error
	GetCounterValue(name string) (int64, error)
	GetGaugeValue(name string) (float64, error)
	ListMetrics() (entity.MetricsList, error)
	PingDB(ctx context.Context) error
}

type Services struct {
	MetricService
}

func New(repo *repository.Repository) *Services {
	return &Services{
		MetricService: newMetricService(repo),
	}
}
