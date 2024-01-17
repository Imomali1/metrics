package storage

import "github.com/Imomali1/metrics/internal/entity"

type MetricStorage interface {
	UpdateCounter(name string, counter int64) error
	UpdateGauge(name string, gauge float64) error
	GetCounterValue(name string) (int64, error)
	GetGaugeValue(name string) (float64, error)
	ListMetrics() ([]entity.Metrics, error)
}

type Storage struct {
	MetricStorage
}

func New() *Storage {
	return &Storage{
		MetricStorage: newMetricStorage(),
	}
}
