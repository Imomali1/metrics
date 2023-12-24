package handlers

import (
	"errors"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (h *MetricHandler) GetMetricValueByName(ctx *gin.Context) {
	metricType := ctx.Param("type")
	if metricType != entity.Gauge && metricType != entity.Counter {
		//err := errors.New("Invalid metric type! ")
		ctx.AbortWithStatus(http.StatusBadRequest)
		//log.Println(err)
		return
	}

	metricName := ctx.Param("name")
	if metricName == "" {
		//err := errors.New("Metric name is empty! ")
		ctx.AbortWithStatus(http.StatusNotFound)
		//log.Println(err)
		return
	}

	var metricValue string
	switch metricType {
	case gauge:
		value, err := h.serviceManager.GetGaugeValue(metricName)
		if err != nil {
			if errors.Is(err, entity.ErrMetricNotFound) {
				ctx.AbortWithStatus(http.StatusNotFound)
				log.Println(err)
				return
			}
			ctx.AbortWithStatus(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		metricValue = strconv.FormatFloat(value, 'f', -1, 64)
	case counter:
		value, err := h.serviceManager.GetCounterValue(metricName)
		if err != nil {
			if errors.Is(err, entity.ErrMetricNotFound) {
				ctx.AbortWithStatus(http.StatusNotFound)
				log.Println(err)
				return
			}
			ctx.AbortWithStatus(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		metricValue = strconv.FormatInt(value, 10)
	}

	ctx.String(http.StatusOK, metricValue)
}
