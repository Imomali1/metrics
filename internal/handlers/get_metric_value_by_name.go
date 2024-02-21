package handlers

import (
	"errors"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *MetricHandler) GetMetricValueByName(ctx *gin.Context) {
	metricType := ctx.Param("type")
	if metricType != entity.Gauge && metricType != entity.Counter {
		err := errors.New("invalid metric type")
		ctx.AbortWithStatus(http.StatusBadRequest)
		h.log.Logger.Info().Err(err).Send()
		return
	}

	metricName := ctx.Param("name")
	if metricName == "" {
		err := errors.New("metric name is empty")
		ctx.AbortWithStatus(http.StatusNotFound)
		h.log.Logger.Info().Err(err).Send()
		return
	}

	var metricValue string
	switch metricType {
	case gauge:
		value, err := h.serviceManager.GetGaugeValue(metricName)
		if err != nil {
			if errors.Is(err, entity.ErrMetricNotFound) {
				ctx.AbortWithStatus(http.StatusNotFound)
				h.log.Logger.Info().Err(err).Send()
				return
			}
			ctx.AbortWithStatus(http.StatusInternalServerError)
			h.log.Logger.Info().Err(err).Msg("cannot get gauge metric value")
			return
		}
		metricValue = strconv.FormatFloat(value, 'f', -1, 64)
	case counter:
		value, err := h.serviceManager.GetCounterValue(metricName)
		if err != nil {
			if errors.Is(err, entity.ErrMetricNotFound) {
				ctx.AbortWithStatus(http.StatusNotFound)
				h.log.Logger.Info().Err(err).Send()
				return
			}
			ctx.AbortWithStatus(http.StatusInternalServerError)
			h.log.Logger.Info().Err(err).Msg("cannot get counter metric value")
			return
		}
		metricValue = strconv.FormatInt(value, 10)
	}

	ctx.String(http.StatusOK, metricValue)
}
