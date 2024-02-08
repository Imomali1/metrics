package repository

import (
	"github.com/Imomali1/metrics/internal/entity"
	store "github.com/Imomali1/metrics/internal/pkg/storage"
)

type metricRepository struct {
	storage *store.Storage
}

func newMetricRepository(storage *store.Storage) *metricRepository {
	return &metricRepository{storage: storage}
}

func (r *metricRepository) UpdateCounter(name string, counter int64) error {
	err := r.storage.Memory.UpdateCounter(name, counter)
	if err != nil {
		return err
	}

	if r.storage.SyncWriteFile {
		err = r.storage.File.WriteCounter(name, counter)
	}
	return err
}

func (r *metricRepository) UpdateGauge(name string, gauge float64) error {
	err := r.storage.Memory.UpdateGauge(name, gauge)
	if err != nil {
		return err
	}

	if r.storage.SyncWriteFile {
		err = r.storage.File.WriteGauge(name, gauge)
	}
	return err
}

func (r *metricRepository) GetCounterValue(name string) (int64, error) {
	value, err := r.storage.GetCounterValue(name)
	return value, err
}

func (r *metricRepository) GetGaugeValue(name string) (float64, error) {
	value, err := r.storage.GetGaugeValue(name)
	return value, err
}

func (r *metricRepository) ListMetrics() (entity.MetricsWithoutPointerList, error) {
	allMetrics, err := r.storage.ListMetrics()
	return allMetrics, err
}
