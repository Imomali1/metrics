package handlers

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/services"
)

func TestMetricHandler_Updates(t *testing.T) {
	type fields struct {
		log            logger.Logger
		serviceManager *services.Services
	}
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &MetricHandler{
				log:            tt.fields.log,
				serviceManager: tt.fields.serviceManager,
			}
			h.Updates(tt.args.ctx)
		})
	}
}
