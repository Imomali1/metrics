package repository

import (
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/storage"
)

type metricRepository struct {
	memStorage *storage.Storage
}

func newMetricRepository(memStorage *storage.Storage) *metricRepository {
	return &metricRepository{memStorage: memStorage}
}

func (r *metricRepository) UpdateCounter(name string, counter int64) error {
	err := r.memStorage.UpdateCounter(name, counter)
	return err
}

func (r *metricRepository) UpdateGauge(name string, gauge float64) error {
	err := r.memStorage.UpdateGauge(name, gauge)
	return err
}

func (r *metricRepository) GetCounterValue(name string) (int64, error) {
	value, err := r.memStorage.GetCounterValue(name)
	return value, err
}

func (r *metricRepository) GetGaugeValue(name string) (float64, error) {
	value, err := r.memStorage.GetGaugeValue(name)
	return value, err
}

func (r *metricRepository) ListMetrics() ([]entity.Metrics, error) {
	allMetrics, err := r.memStorage.ListMetrics()
	return allMetrics, err
}
