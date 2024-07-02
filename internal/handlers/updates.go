package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"

	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/utils"
)

func (h *MetricHandler) Updates(ctx *gin.Context) {
	ct := ctx.GetHeader("Content-Type")
	if ct != "application/json" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		h.log.Logger.Info().Msg("content-type is not application/json")
		return
	}

	body, err := utils.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.log.Logger.Info().Err(err).Msg("cannot read request body")
		return
	}

	var batch entity.MetricsList
	err = easyjson.Unmarshal(body, &batch)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		h.log.Logger.Info().Err(err).Msg("cannot unmarshal json")
		return
	}

	for i, metrics := range batch {
		switch metrics.MType {
		case entity.Counter:
			h.log.Logger.Info().Msgf("#%d counter %s %d", i+1, metrics.ID, *metrics.Delta)
		case entity.Gauge:
			h.log.Logger.Info().Msgf("#%d gauge %s %f", i+1, metrics.ID, *metrics.Value)
		}

		if metrics.MType != entity.Gauge && metrics.MType != entity.Counter {
			err = errors.New("invalid metric type")
			ctx.AbortWithStatus(http.StatusBadRequest)
			h.log.Logger.Info().Err(err).Send()
			return
		}

	}

	c, cancel := context.WithTimeout(ctx, _timeout)
	defer cancel()

	err = h.uc.UpdateMetrics(c, batch)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.log.Logger.Info().Err(err).Msg("cannot update batch of metric value")
		return
	}

	ctx.JSON(http.StatusOK, batch)
}
