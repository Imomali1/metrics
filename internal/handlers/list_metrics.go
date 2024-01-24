package handlers

import (
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ModifiedMetric struct {
	Type  string
	Name  string
	Value interface{}
}

func (h *MetricHandler) ListMetrics(ctx *gin.Context) {
	allMetrics, err := h.serviceManager.ListMetrics()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		//logger.Log.Info(err)
		return
	}

	metrics := modifyMetrics(allMetrics)

	ctx.HTML(http.StatusOK, "index.html", metrics)
}

func modifyMetrics(metrics []entity.Metrics) []ModifiedMetric {
	if len(metrics) == 0 {
		return nil
	}

	ms := make([]ModifiedMetric, len(metrics))
	for i, metric := range metrics {
		var value interface{}
		if metric.MType == entity.Counter {
			value = metric.Delta
		} else if metric.MType == entity.Gauge {
			value = metric.Value
		}
		ms[i] = ModifiedMetric{
			Type:  metric.MType,
			Name:  metric.ID,
			Value: value,
		}
	}
	return ms
}
