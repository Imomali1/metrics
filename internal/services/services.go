package services

import (
	"context"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/repository"
)

type MetricService interface {
	UpdateMetrics(context.Context, entity.MetricsList) error
	GetMetrics(context.Context, entity.Metrics) (entity.Metrics, error)
	ListMetrics(context.Context) (entity.MetricsList, error)
	Ping(ctx context.Context) error
}

type Services struct {
	MetricService
}

func New(repo *repository.Repository) *Services {
	return &Services{
		MetricService: newMetricService(repo),
	}
}
