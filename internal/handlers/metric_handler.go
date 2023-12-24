package handlers

import (
	"github.com/Imomali1/metrics/internal/entity"
)

type IMetricService interface {
	UpdateCounter(name string, counter int64) error
	UpdateGauge(name string, gauge float64) error
	GetCounterValue(name string) (int64, error)
	GetGaugeValue(name string) (float64, error)
	ListMetrics() ([]entity.Metric, error)
}

type MetricHandler struct {
	serviceManager IMetricService
}

func NewMetricHandler(sm IMetricService) *MetricHandler {
	return &MetricHandler{
		serviceManager: sm,
	}
}
