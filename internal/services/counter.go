package services

import "github.com/Imomali1/metrics/internal/repository"

type counterService struct {
	repo repository.CounterRepository
}

func newCounterService(repo repository.CounterRepository) *counterService {
	return &counterService{repo: repo}
}

func (s *counterService) UpdateCounter(name string, counter int64) error {
	err := s.repo.UpdateCounter(name, counter)
	return err
}
