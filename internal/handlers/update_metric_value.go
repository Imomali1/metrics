package handlers

import (
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *MetricHandler) UpdateMetricValue(ctx *gin.Context) {
	metricType := ctx.Param("type")
	if metricType != entity.Gauge && metricType != entity.Counter {
		//err := errors.New("invalid metric type ")
		ctx.AbortWithStatus(http.StatusBadRequest)
		//logger.Logger.Info(err)
		return
	}

	metricName := ctx.Param("name")
	if metricName == "" {
		//err := errors.New("metric name is empty ")
		ctx.AbortWithStatus(http.StatusNotFound)
		//logger.Logger.Info(err)
		return
	}

	metricValue := ctx.Param("value")

	switch metricType {
	case entity.Gauge:
		gaugeValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			//logger.Logger.Info(err)
			return
		}
		err = h.serviceManager.UpdateGauge(metricName, gaugeValue)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			//logger.Logger.Info(err)
			return
		}
	case entity.Counter:
		counterValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			//logger.Logger.Info(err)
			return
		}
		err = h.serviceManager.UpdateCounter(metricName, counterValue)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			//logger.Logger.Info(err)
			return
		}
	}

	ctx.Status(http.StatusOK)
}
