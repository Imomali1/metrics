package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (h *GaugeHandler) UpdateGauge(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		err := errors.New("Название метрики не указано ")
		ctx.AbortWithStatus(http.StatusNotFound)
		log.Println(err)
		return
	}

	value := ctx.Param("value")
	gauge, err := strconv.ParseFloat(value, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
		return
	}

	err = h.serviceManager.UpdateGauge(name, gauge)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	ctx.Status(http.StatusOK)
}
