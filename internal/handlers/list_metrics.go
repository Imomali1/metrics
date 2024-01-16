package handlers

import (
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/logger"
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
		logger.Log.Info(err)
		return
	}

	metrics := modifyMetrics(allMetrics)

	ctx.HTML(http.StatusOK, "index.html", metrics)
}

func modifyMetrics(metrics []entity.Metric) []ModifiedMetric {
	if len(metrics) == 0 {
		return nil
	}

	ms := make([]ModifiedMetric, len(metrics))
	for i, metric := range metrics {
		var value interface{}
		if metric.Type == entity.Counter {
			value = metric.ValueCounter
		} else if metric.Type == entity.Gauge {
			value = metric.ValueGauge
		}
		ms[i] = ModifiedMetric{
			Type:  metric.Type,
			Name:  metric.Name,
			Value: value,
		}
	}
	return ms
}
