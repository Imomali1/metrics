package services

import "github.com/Imomali1/metrics/internal/repository"

type gaugeService struct {
	repo repository.GaugeRepository
}

func newGaugeService(repo repository.GaugeRepository) *gaugeService {
	return &gaugeService{repo: repo}
}

func (s *gaugeService) UpdateGauge(name string, gauge float64) error {
	err := s.repo.UpdateGauge(name, gauge)
	return err
}
