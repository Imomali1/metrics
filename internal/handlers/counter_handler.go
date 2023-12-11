package handlers

import "github.com/Imomali1/metrics/internal/services"

type CounterHandlerOptions struct {
	ServiceManager services.CounterService
}

type CounterHandler struct {
	serviceManager services.CounterService
}

func NewCounterHandler(options CounterHandlerOptions) *CounterHandler {
	return &CounterHandler{
		serviceManager: options.ServiceManager,
	}
}
