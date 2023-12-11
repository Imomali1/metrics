package repository

import "github.com/Imomali1/metrics/internal/pkg/storage"

type GaugeRepository interface {
	UpdateGauge(name string, gauge float64) error
}

type CounterRepository interface {
	UpdateCounter(name string, counter int64) error
}

type Repository struct {
	GaugeRepository
	CounterRepository
}

func New(memStorage *storage.Storage) *Repository {
	return &Repository{
		GaugeRepository:   newGaugeRepository(memStorage),
		CounterRepository: newCounterRepository(memStorage),
	}
}
