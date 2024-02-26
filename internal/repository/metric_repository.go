package repository

import (
	"context"
	"github.com/Imomali1/metrics/internal/entity"
	store "github.com/Imomali1/metrics/internal/pkg/storage/v2"
)

type metricRepository struct {
	storage *store.Storage
}

func newMetricRepository(storage *store.Storage) *metricRepository {
	return &metricRepository{storage: storage}
}

func (r *metricRepository) UpdateMetrics(ctx context.Context, batch entity.MetricsList) error {
	err := r.storage.Update(ctx, batch)
	if err != nil {
		return err
	}
	if r.storage.SyncWriteAllowed {
		var list entity.MetricsList
		list, err = r.storage.GetAll(ctx)
		if err != nil {
			return err
		}

		err = r.storage.Sync.Write(list)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *metricRepository) GetMetrics(ctx context.Context, metric entity.Metrics) (entity.Metrics, error) {
	id, mType := metric.ID, metric.MType
	return r.storage.GetOne(ctx, id, mType)
}

func (r *metricRepository) ListMetrics(ctx context.Context) (entity.MetricsList, error) {
	return r.storage.GetAll(ctx)
}

func (r *metricRepository) Ping(ctx context.Context) error {
	return r.storage.Ping(ctx)
}
