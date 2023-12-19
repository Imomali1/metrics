package handlers

import "github.com/Imomali1/metrics/internal/services"

type MetricHandlerOptions struct {
	ServiceManager *services.Services
}

type MetricHandler struct {
	serviceManager *services.Services
}

func NewMetricHandler(options MetricHandlerOptions) *MetricHandler {
	return &MetricHandler{
		serviceManager: options.ServiceManager,
	}
}
