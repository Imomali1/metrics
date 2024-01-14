package handlers

import (
	"errors"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *MetricHandler) GetMetricValueByName(ctx *gin.Context) {
	metricType := ctx.Param("type")
	if metricType != entity.Gauge && metricType != entity.Counter {
		err := errors.New("invalid metric type ")
		ctx.AbortWithStatus(http.StatusBadRequest)
		logger.Log.Info(err)
		return
	}

	metricName := ctx.Param("name")
	if metricName == "" {
		err := errors.New("metric name is empty ")
		ctx.AbortWithStatus(http.StatusNotFound)
		logger.Log.Info(err)
		return
	}

	var metricValue string
	switch metricType {
	case entity.Gauge:
		value, err := h.serviceManager.GetGaugeValue(metricName)
		if err != nil {
			if errors.Is(err, entity.ErrMetricNotFound) {
				ctx.AbortWithStatus(http.StatusNotFound)
				logger.Log.Info(err)
				return
			}
			ctx.AbortWithStatus(http.StatusInternalServerError)
			logger.Log.Info(err)
			return
		}
		metricValue = strconv.FormatFloat(value, 'f', -1, 64)
	case entity.Counter:
		value, err := h.serviceManager.GetCounterValue(metricName)
		if err != nil {
			if errors.Is(err, entity.ErrMetricNotFound) {
				ctx.AbortWithStatus(http.StatusNotFound)
				logger.Log.Info(err)
				return
			}
			ctx.AbortWithStatus(http.StatusInternalServerError)
			logger.Log.Info(err)
			return
		}
		metricValue = strconv.FormatInt(value, 10)
	}

	ctx.String(http.StatusOK, metricValue)
}
