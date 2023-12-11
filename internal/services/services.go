package services

import "github.com/Imomali1/metrics/internal/repository"

type GaugeService interface {
	UpdateGauge(name string, gauge float64) error
}

type CounterService interface {
	UpdateCounter(name string, counter int64) error
}

type Services struct {
	GaugeService
	CounterService
}

func New(repo *repository.Repository) *Services {
	return &Services{
		GaugeService:   newGaugeService(repo.GaugeRepository),
		CounterService: newCounterService(repo.CounterRepository),
	}
}
