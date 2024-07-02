package repository

import (
	"context"
	"github.com/Imomali1/metrics/internal/entity"
)

type Repository interface {
	UpdateMetrics(context.Context, entity.MetricsList) error
	GetMetrics(context.Context, entity.Metrics) (entity.Metrics, error)
	ListMetrics(context.Context) (entity.MetricsList, error)
	Ping(ctx context.Context) error
	SyncWrite(list entity.MetricsList) error
}
