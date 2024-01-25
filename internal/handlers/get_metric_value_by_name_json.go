package handlers

import (
	"errors"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
	"io"
	"net/http"
)

func (h *MetricHandler) GetMetricValueByNameJSON(ctx *gin.Context) {
	ct := ctx.GetHeader("Content-Type")
	if ct != "application/json" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		h.log.Logger.Info().Msg("content-type is not application/json")
		return
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.log.Logger.Info().Err(err).Msg("cannot read request body")
		return
	}

	var metrics entity.Metrics
	err = easyjson.Unmarshal(body, &metrics)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		h.log.Logger.Info().Err(err).Msg("cannot unmarshal json")
		return
	}

	switch metrics.MType {
	case gauge:
		var value float64
		value, err = h.serviceManager.GetGaugeValue(metrics.ID)
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
		metrics.Value = &value
	case counter:
		var delta int64
		delta, err = h.serviceManager.GetCounterValue(metrics.ID)
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
		metrics.Delta = &delta
	}

	ctx.JSON(http.StatusOK, metrics)
}
