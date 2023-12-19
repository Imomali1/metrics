package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

const (
	gauge   = "gauge"
	counter = "counter"
)

func (h *MetricHandler) UpdateMetric(ctx *gin.Context) {
	metricType := ctx.Param("type")
	if metricType != gauge && metricType != counter {
		err := errors.New("Invalid metric type! ")
		ctx.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
		return
	}

	metricName := ctx.Param("name")
	if metricName == "" {
		err := errors.New("Metric name is empty! ")
		ctx.AbortWithStatus(http.StatusNotFound)
		log.Println(err)
		return
	}

	metricValue := ctx.Param("value")

	switch metricType {
	case gauge:
		gaugeValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			log.Println(err)
			return
		}
		err = h.serviceManager.UpdateGauge(metricName, gaugeValue)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	case counter:
		counterValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			log.Println(err)
			return
		}
		err = h.serviceManager.UpdateCounter(metricName, counterValue)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}

	ctx.Status(http.StatusOK)
}
