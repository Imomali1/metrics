package services

import (
	"context"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/repository"
)

type metricService struct {
	repo *repository.Repository
}

func newMetricService(repo *repository.Repository) *metricService {
	return &metricService{repo: repo}
}

func (s *metricService) UpdateMetrics(ctx context.Context, batch entity.MetricsList) error {
	return s.repo.UpdateMetrics(ctx, batch)
}

func (s *metricService) GetMetrics(ctx context.Context, metric entity.Metrics) (entity.Metrics, error) {
	return s.repo.GetMetrics(ctx, metric)
}

func (s *metricService) ListMetrics(ctx context.Context) (entity.MetricsList, error) {
	return s.repo.ListMetrics(ctx)
}

func (s *metricService) Ping(ctx context.Context) error {
	return s.repo.Ping(ctx)
}
