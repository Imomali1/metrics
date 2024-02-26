package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *MetricHandler) PingDB(ctx *gin.Context) {
	err := h.serviceManager.PingDB()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.log.Logger.Info().Err(err).Msg("cannot ping database")
		return
	}
	ctx.Status(http.StatusOK)
}
