package repository

import (
	"context"

	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/file"
	"github.com/Imomali1/metrics/internal/pkg/storage"
)

// MetricsRepo implements all methods of Repository interface.
type MetricsRepo struct {
	store          storage.Storage
	syncFileWriter file.SyncFileWriter
}

// New creates new instance of Repository.
func New(
	store storage.Storage,
	syncFileWriter file.SyncFileWriter,
) Repository {
	return &MetricsRepo{
		store:          store,
		syncFileWriter: syncFileWriter,
	}
}

// UpdateMetrics sets values if metrics do not exist,
// if exist then updates values.
func (r *MetricsRepo) UpdateMetrics(ctx context.Context, batch entity.MetricsList) error {
	err := r.store.Update(ctx, batch)
	if err != nil {
		return err
	}

	return nil
}

// GetMetrics fetches metrics by name and type.
func (r *MetricsRepo) GetMetrics(ctx context.Context, metric entity.Metrics) (entity.Metrics, error) {
	id, mType := metric.ID, metric.MType
	return r.store.GetOne(ctx, id, mType)
}

// ListMetrics fetches all metrics stored in storage.
func (r *MetricsRepo) ListMetrics(ctx context.Context) (entity.MetricsList, error) {
	return r.store.GetAll(ctx)
}

// Ping checks whether database is alive or not.
func (r *MetricsRepo) Ping(ctx context.Context) error {
	return r.store.Ping(ctx)
}

// SyncWrite writes synchronously all metrics from storage.
func (r *MetricsRepo) SyncWrite(list entity.MetricsList) error {
	if r.syncFileWriter != nil {
		return r.syncFileWriter.Write(list)
	}

	return nil
}
