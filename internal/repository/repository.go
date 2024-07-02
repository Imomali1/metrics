package repository

import (
	"context"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/file"
	"github.com/Imomali1/metrics/internal/pkg/storage"
)

type MetricsRepo struct {
	store storage.Storage
	syncFileWriter file.SyncFileWriter
}

func New(
	store storage.Storage,
	syncFileWriter file.SyncFileWriter,
) Repository {
	return &MetricsRepo{
		store: store,
		syncFileWriter: syncFileWriter,
	}
}

func (r *MetricsRepo) UpdateMetrics(ctx context.Context, batch entity.MetricsList) error {
	err := r.store.Update(ctx, batch)
	if err != nil {
		return err
	}

	return nil
}

func (r *MetricsRepo) GetMetrics(ctx context.Context, metric entity.Metrics) (entity.Metrics, error) {
	id, mType := metric.ID, metric.MType
	return r.store.GetOne(ctx, id, mType)
}

func (r *MetricsRepo) ListMetrics(ctx context.Context) (entity.MetricsList, error) {
	return r.store.GetAll(ctx)
}

func (r *MetricsRepo) Ping(ctx context.Context) error {
	return r.store.Ping(ctx)
}

func (r *MetricsRepo) SyncWrite(list entity.MetricsList) error {
	if r.syncFileWriter != nil {
		return r.syncFileWriter.Write(list)
	}

	return nil
}