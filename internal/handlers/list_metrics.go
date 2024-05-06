package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *MetricHandler) ListMetrics(ctx *gin.Context) {
	c, cancel := context.WithTimeout(ctx, _timeout)
	defer cancel()

	allMetrics, err := h.serviceManager.ListMetrics(c)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.log.Logger.Info().Err(err).Msg("cannot list metrics")
		return
	}

	ctx.HTML(http.StatusOK, "index.html", allMetrics)
}
