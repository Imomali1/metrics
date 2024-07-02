package usecase

import (
	"context"

	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/repository"
)

type MetricUseCase struct {
	repo repository.Repository
}

func New(repo repository.Repository) UseCase {
	return &MetricUseCase{
		repo: repo,
	}
}

func (uc *MetricUseCase) UpdateMetrics(
	ctx context.Context,
	batch entity.MetricsList,
) error {
	if err := uc.repo.UpdateMetrics(ctx, batch); err != nil {
		return err
	}

	if err := uc.repo.SyncWrite(batch); err != nil {
		return err
	}

	return nil
}

func (uc *MetricUseCase) GetMetrics(
	ctx context.Context,
	metric entity.Metrics,
) (entity.Metrics, error) {
	return uc.repo.GetMetrics(ctx, metric)
}

func (uc *MetricUseCase) ListMetrics(
	ctx context.Context,
) (entity.MetricsList, error) {
	return uc.repo.ListMetrics(ctx)
}

func (uc *MetricUseCase) Ping(
	ctx context.Context,
) error {
	return uc.repo.Ping(ctx)
}
