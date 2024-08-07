package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Imomali1/metrics/internal/entity"
)

func (h *MetricHandler) UpdateMetricValue(ctx *gin.Context) {
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

	var err error
	delta, value := new(int64), new(float64)

	metricValue := ctx.Param("value")

	switch metricType {
	case entity.Gauge:
		*value, err = strconv.ParseFloat(metricValue, 64)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			h.log.Logger.Info().Err(err).Msg("gauge metric value is not float64")
			return
		}
	case entity.Counter:
		*delta, err = strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			h.log.Logger.Info().Err(err).Msg("counter metric value is not int64")
			return
		}
	default:
	}

	metrics := entity.Metrics{
		ID:    metricName,
		MType: metricType,
		Delta: delta,
		Value: value,
	}

	c, cancel := context.WithTimeout(ctx, _timeout)
	defer cancel()

	err = h.uc.UpdateMetrics(c, []entity.Metrics{metrics})
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.log.Logger.Info().Err(err).Msgf("cannot update %s metric value", metrics.MType)
		return
	}

	ctx.Status(http.StatusOK)
}
