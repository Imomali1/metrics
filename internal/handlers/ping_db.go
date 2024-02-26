package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (h *MetricHandler) PingDB(ctx *gin.Context) {
	c, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	err := h.serviceManager.Ping(c)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.log.Logger.Info().Err(err).Msg("cannot ping database")
		return
	}

	ctx.Status(http.StatusOK)
}
