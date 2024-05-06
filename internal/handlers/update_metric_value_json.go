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

func (h *MetricHandler) UpdateMetricValueJSON(ctx *gin.Context) {
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

	var metrics entity.Metrics
	err = easyjson.Unmarshal(body, &metrics)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		h.log.Logger.Info().Err(err).Msg("cannot unmarshal json")
		return
	}

	if metrics.MType != entity.Gauge && metrics.MType != entity.Counter {
		err = errors.New("invalid metric type")
		ctx.AbortWithStatus(http.StatusBadRequest)
		h.log.Logger.Info().Err(err).Send()
		return
	}

	c, cancel := context.WithTimeout(ctx, _timeout)
	defer cancel()

	err = h.serviceManager.UpdateMetrics(c, []entity.Metrics{metrics})
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.log.Logger.Info().Err(err).Msgf("cannot update %s metric value", metrics.MType)
		return
	}

	ctx.JSON(http.StatusOK, metrics)
}
