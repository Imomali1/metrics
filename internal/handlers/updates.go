package handlers

import (
	"context"
	"errors"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
	"io"
	"net/http"
	"time"
)

func (h *MetricHandler) Updates(ctx *gin.Context) {
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

	var batch entity.MetricsList
	err = easyjson.Unmarshal(body, &batch)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		h.log.Logger.Info().Err(err).Msg("cannot unmarshal json")
		return
	}

	for _, metrics := range batch {
		if metrics.MType != entity.Gauge && metrics.MType != entity.Counter {
			err = errors.New("invalid metric type")
			ctx.AbortWithStatus(http.StatusBadRequest)
			h.log.Logger.Info().Err(err).Send()
			return
		}
	}

	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = h.serviceManager.UpdateMetrics(c, batch)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.log.Logger.Info().Err(err).Msg("cannot update batch of metric value")
		return
	}

	ctx.JSON(http.StatusOK, batch)
}
