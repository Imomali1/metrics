package handlers

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"

	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/services"
)

func NewTestHandler() (*gin.Engine, *MockServiceManager, MetricHandler) {
	gin.SetMode(gin.TestMode)

	mockServiceManager := new(MockServiceManager)
	log := logger.NewLogger(os.Stdout, "info", "test_metrics")

	handler := MetricHandler{
		serviceManager: &services.Services{
			MetricService: mockServiceManager,
		},
		log: log,
	}

	r := gin.Default()
	return r, mockServiceManager, handler
}

type MockServiceManager struct {
	mock.Mock
}

func (m *MockServiceManager) UpdateMetrics(ctx context.Context, metrics entity.MetricsList) error {
	args := m.Called(ctx, metrics)
	return args.Error(0)
}

func (m *MockServiceManager) GetMetrics(ctx context.Context, metric entity.Metrics) (entity.Metrics, error) {
	args := m.Called(ctx, metric)
	return args.Get(0).(entity.Metrics), args.Error(1)
}

func (m *MockServiceManager) ListMetrics(ctx context.Context) (entity.MetricsList, error) {
	args := m.Called(ctx)
	return args.Get(0).(entity.MetricsList), args.Error(1)
}

func (m *MockServiceManager) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
