package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *MetricHandler) PingDB(ctx *gin.Context) {
	c, cancel := context.WithTimeout(ctx, _timeout)
	defer cancel()

	err := h.serviceManager.Ping(c)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.log.Logger.Info().Err(err).Msg("cannot ping database")
		return
	}

	ctx.Status(http.StatusOK)
}
