package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *MetricHandler) ListMetrics(ctx *gin.Context) {
	allMetrics, err := h.serviceManager.ListMetrics()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.log.Logger.Info().Err(err).Msg("cannot list metrics")
		return
	}
	ctx.HTML(http.StatusOK, "index.html", allMetrics)
}
