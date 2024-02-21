package handlers

import (
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/logger"
)

type IMetricService interface {
	UpdateCounter(name string, counter int64) error
	UpdateGauge(name string, gauge float64) error
	GetCounterValue(name string) (int64, error)
	GetGaugeValue(name string) (float64, error)
	ListMetrics() (entity.MetricsList, error)
}

type MetricHandler struct {
	log            logger.Logger
	serviceManager IMetricService
}

func NewMetricHandler(log logger.Logger, sm IMetricService) *MetricHandler {
	return &MetricHandler{
		log:            log,
		serviceManager: sm,
	}
}
