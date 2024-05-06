package handlers

import (
	"time"

	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/services"
)

const _timeout = 1 * time.Second

type MetricHandler struct {
	log            logger.Logger
	serviceManager *services.Services
}

func NewMetricHandler(log logger.Logger, sm *services.Services) *MetricHandler {
	return &MetricHandler{
		log:            log,
		serviceManager: sm,
	}
}
