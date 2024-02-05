package services

import (
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/repository"
)

type metricService struct {
	repo *repository.Repository
}

func newMetricService(repo *repository.Repository) *metricService {
	return &metricService{repo: repo}
}

func (s *metricService) UpdateCounter(name string, counter int64) error {
	return s.repo.UpdateCounter(name, counter)
}

func (s *metricService) UpdateGauge(name string, gauge float64) error {
	return s.repo.UpdateGauge(name, gauge)
}

func (s *metricService) GetCounterValue(name string) (int64, error) {
	value, err := s.repo.GetCounterValue(name)
	return value, err
}

func (s *metricService) GetGaugeValue(name string) (float64, error) {
	value, err := s.repo.GetGaugeValue(name)
	return value, err
}

func (s *metricService) ListMetrics() ([]entity.MetricsWithoutPointer, error) {
	allMetrics, err := s.repo.ListMetrics()
	return allMetrics, err
}
