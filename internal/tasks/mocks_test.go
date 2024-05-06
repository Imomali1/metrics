package tasks

import (
	"context"
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/Imomali1/metrics/internal/entity"
)

// MockStorage implements IStorage for testing purposes
type MockStorage struct {
	mock.Mock
}

func NewMockStorage(opts ...MockStorageOption) *MockStorage {
	mockStorage := new(MockStorage)
	for _, opt := range opts {
		opt(mockStorage)
	}
	return mockStorage
}

type MockStorageOption func(*MockStorage)

func WithMetrics(metrics entity.MetricsList) MockStorageOption {
	return func(mockStorage *MockStorage) {
		mockStorage.
			On("GetAll", mock.Anything).
			Return(metrics, nil)
	}
}

func WithError() MockStorageOption {
	return func(mockStorage *MockStorage) {
		mockStorage.
			On("GetAll", mock.Anything).
			Return(entity.MetricsList{}, errors.New("storage error"))
	}
}

func (ms *MockStorage) Update(ctx context.Context, batch entity.MetricsList) error {
	args := ms.Called(ctx, batch)
	return args.Error(0)
}

func (ms *MockStorage) GetOne(ctx context.Context, id string, mType string) (entity.Metrics, error) {
	args := ms.Called(ctx, id, mType)
	return args.Get(0).(entity.Metrics), args.Error(1)
}

func (ms *MockStorage) GetAll(ctx context.Context) (entity.MetricsList, error) {
	args := ms.Called(ctx)
	return args.Get(0).(entity.MetricsList), args.Error(1)
}

func (ms *MockStorage) DeleteOne(ctx context.Context, id, mType string) error {
	args := ms.Called(ctx, id, mType)
	return args.Error(0)
}

func (ms *MockStorage) DeleteAll(ctx context.Context) error {
	args := ms.Called(ctx)
	return args.Error(0)
}

func (ms *MockStorage) Ping(ctx context.Context) error {
	args := ms.Called(ctx)
	return args.Error(0)
}

func (ms *MockStorage) Close() {
	ms.Called()
}
