package handlers

import (
	"reflect"
	"testing"

	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/services"
)

func TestNewMetricHandler(t *testing.T) {
	type args struct {
		log logger.Logger
		sm  *services.Services
	}
	tests := []struct {
		name string
		args args
		want *MetricHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMetricHandler(tt.args.log, tt.args.sm); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetricHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
