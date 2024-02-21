package handlers

import (
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
	"io"
	"net/http"
)

func (h *MetricHandler) UpdateMetricValueJSON(ctx *gin.Context) {
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
		err = h.serviceManager.UpdateGauge(metrics.ID, *metrics.Value)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			h.log.Logger.Info().Err(err).Msg("cannot update gauge metric value")
			return
		}
	case counter:
		err = h.serviceManager.UpdateCounter(metrics.ID, *metrics.Delta)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			h.log.Logger.Info().Err(err).Msg("cannot update counter metric value")
			return
		}
	}

	ctx.JSON(http.StatusOK, metrics)
}
