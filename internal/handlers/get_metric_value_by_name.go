package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Imomali1/metrics/internal/entity"
)

func (h *MetricHandler) GetMetricValueByName(ctx *gin.Context) {
	metricType := ctx.Param("type")
	if metricType != entity.Gauge && metricType != entity.Counter {
		err := errors.New("invalid metric type")
		ctx.AbortWithStatus(http.StatusBadRequest)
		h.log.Info().Err(err).Send()
		return
	}

	metricName := ctx.Param("name")

	metrics := entity.Metrics{
		ID:    metricName,
		MType: metricType,
	}

	c, cancel := context.WithTimeout(ctx, _timeout)
	defer cancel()

	result, err := h.uc.GetMetrics(c, metrics)
	if err != nil {
		if errors.Is(err, entity.ErrMetricNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			h.log.Info().Err(err).Send()
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.log.Info().Err(err).Msgf("cannot get %s metric value", metrics.MType)
		return
	}

	var metricValue string
	switch result.MType {
	case entity.Counter:
		metricValue = strconv.FormatInt(*result.Delta, 10)
	case entity.Gauge:
		metricValue = strconv.FormatFloat(*result.Value, 'f', -1, 64)
	}

	ctx.String(http.StatusOK, metricValue)
}
