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
		metric := entity.Metrics{
			ID:    name,
			MType: entity.Counter,
			Delta: &counter,
		}
		err = r.storage.File.WriteMetrics([]entity.Metrics{metric})
	}
	return err
}

func (r *metricRepository) UpdateGauge(name string, gauge float64) error {
	err := r.storage.Memory.UpdateGauge(name, gauge)
	if err != nil {
		return err
	}

	if r.storage.SyncWriteFile {
		metric := entity.Metrics{
			ID:    name,
			MType: entity.Gauge,
			Value: &gauge,
		}
		err = r.storage.File.WriteMetrics([]entity.Metrics{metric})
	}
	return err
}

func (r *metricRepository) GetCounterValue(name string) (int64, error) {
	value, err := r.storage.Memory.GetCounterValue(name)
	return value, err
}

func (r *metricRepository) GetGaugeValue(name string) (float64, error) {
	value, err := r.storage.Memory.GetGaugeValue(name)
	return value, err
}

func (r *metricRepository) ListMetrics() (entity.MetricsList, error) {
	allMetrics, err := r.storage.Memory.ListMetrics()
	return allMetrics, err
}
