package handlers

import "github.com/Imomali1/metrics/internal/services"

type GaugeHandlerOptions struct {
	ServiceManager services.GaugeService
}

type GaugeHandler struct {
	serviceManager services.GaugeService
}

func NewGaugeHandler(options GaugeHandlerOptions) *GaugeHandler {
	return &GaugeHandler{
		serviceManager: options.ServiceManager,
	}
}
