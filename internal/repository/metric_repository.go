package repository

import (
	"context"
	"github.com/Imomali1/metrics/internal/entity"
	store "github.com/Imomali1/metrics/internal/pkg/storage/v2"
)

type metricRepository struct {
	storage store.IStorage
}

func newMetricRepository(storage store.IStorage) *metricRepository {
	return &metricRepository{storage: storage}
}

func (r *metricRepository) UpdateMetrics(ctx context.Context, batch entity.MetricsList) error {
	return r.storage.Update(ctx, batch)
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
