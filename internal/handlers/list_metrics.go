package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *MetricHandler) ListMetrics(ctx *gin.Context) {
	allMetrics, err := h.serviceManager.ListMetrics()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	ctx.HTML(http.StatusOK, "index.html", allMetrics)
}
