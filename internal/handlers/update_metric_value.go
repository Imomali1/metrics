package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const (
	gauge   = "gauge"
	counter = "counter"
)

func (h *MetricHandler) UpdateMetricValue(ctx *gin.Context) {
	metricType := ctx.Param("type")
	if metricType != gauge && metricType != counter {
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

	metricValue := ctx.Param("value")

	switch metricType {
	case gauge:
		gaugeValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			h.log.Logger.Info().Err(err).Msg("gauge metric value is not float64")
			return
		}
		err = h.serviceManager.UpdateGauge(metricName, gaugeValue)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			h.log.Logger.Info().Err(err).Msg("cannot update gauge metric value")
			return
		}
	case counter:
		counterValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			h.log.Logger.Info().Err(err).Msg("counter metric value is not int64")
			return
		}
		err = h.serviceManager.UpdateCounter(metricName, counterValue)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			h.log.Logger.Info().Err(err).Msg("cannot update counter metric value")
			return
		}
	}

	ctx.Status(http.StatusOK)
}
