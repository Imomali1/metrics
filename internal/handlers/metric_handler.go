package handlers

import (
	"time"

	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/usecase"
)

const _timeout = 1 * time.Second

type MetricHandler struct {
	log logger.Logger
	uc  usecase.UseCase
}

func NewMetricHandler(log logger.Logger, uc usecase.UseCase) *MetricHandler {
	return &MetricHandler{
		log: log,
		uc:  uc,
	}
}
