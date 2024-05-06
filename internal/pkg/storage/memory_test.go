package storage

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/utils"
)

func Test_memoryStorage_Close(t *testing.T) {
	s := newMemoryStorage()
	s.Close()
}

func Test_memoryStorage_GetAll(t *testing.T) {
	ctx := context.Background()

	s := newMemoryStorage()

	metrics := entity.MetricsList{
		{ID: "gauge1", MType: entity.Gauge, Value: utils.Ptr(123.0)},
		{ID: "counter1", MType: entity.Counter, Delta: utils.Ptr(int64(123))},
	}

	s.Update(ctx, metrics)

	got, err := s.GetAll(ctx)
	require.NoError(t, err)
	require.ElementsMatch(t, metrics, got)
}

func Test_memoryStorage_GetOne(t *testing.T) {
	ctx := context.Background()

	s := newMemoryStorage()

	metrics := entity.MetricsList{
		{ID: "gauge1", MType: entity.Gauge, Value: utils.Ptr(123.0)},
		{ID: "counter1", MType: entity.Counter, Delta: utils.Ptr(int64(123))},
	}

	s.Update(ctx, metrics)

	type args struct {
		id    string
		mType string
	}
	tests := []struct {
		name    string
		idx     int
		args    args
		want    entity.Metrics
		wantErr bool
	}{
		{
			name: "get valid gauge metrics",
			args: args{
				id:    "gauge1",
				mType: entity.Gauge,
			},
			want: entity.Metrics{
				ID:    "gauge1",
				MType: entity.Gauge,
				Value: utils.Ptr(123.0),
			},
		},
		{
			name: "get valid counter metrics",
			args: args{
				id:    "counter1",
				mType: entity.Counter,
			},
			want: entity.Metrics{
				ID:    "counter1",
				MType: entity.Counter,
				Delta: utils.Ptr(int64(123)),
			},
		},
		{
			name: "gauge type but non-existing metrics",
			args: args{
				id:    "non-existing",
				mType: entity.Gauge,
			},
			wantErr: true,
		},
		{
			name: "counter type but non-existing metrics",
			args: args{
				id:    "non-existing",
				mType: entity.Counter,
			},
			wantErr: true,
		},
		{
			name: "invalid type metrics",
			args: args{
				id:    "counter1",
				mType: "invalid",
			},
			want: entity.Metrics{ID: "counter1", MType: "invalid"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetOne(ctx, tt.args.id, tt.args.mType)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, got, tt.want)
		})
	}

}

func Test_memoryStorage_Ping(t *testing.T) {
	s := newMemoryStorage()
	err := s.Ping(context.Background())
	require.Error(t, err)
}

func Test_memoryStorage_Update(t *testing.T) {
	s := newMemoryStorage()

	ctx := context.Background()
	type args struct {
		batch entity.MetricsList
		want  entity.MetricsList
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "proper work of update",
			args: args{
				batch: entity.MetricsList{
					{ID: "gauge1", MType: entity.Gauge, Value: utils.Ptr(123.0), Delta: utils.Ptr(int64(123))},
					{ID: "counter1", MType: entity.Counter, Value: utils.Ptr(123.0), Delta: utils.Ptr(int64(123))},
					{ID: "invalid-type", MType: "invalid", Value: utils.Ptr(123.0), Delta: utils.Ptr(int64(123))},
				},
				want: entity.MetricsList{
					{ID: "gauge1", MType: entity.Gauge, Value: utils.Ptr(123.0)},
					{ID: "counter1", MType: entity.Counter, Delta: utils.Ptr(int64(123))},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.Update(ctx, tt.args.batch)
			require.NoError(t, err)
			batch, _ := s.GetAll(ctx)
			require.ElementsMatch(t, batch, tt.args.want)
		})
	}
}

func Test_newMemoryStorage(t *testing.T) {
	want := &memoryStorage{
		mu:             sync.RWMutex{},
		counterStorage: make(map[string]int64),
		gaugeStorage:   make(map[string]float64),
	}

	got := newMemoryStorage()
	require.Equal(t, got, want)
}

func Test_memoryStorage_DeleteOne(t *testing.T) {
	ctx := context.Background()

	s := newMemoryStorage()

	s.Update(ctx, entity.MetricsList{
		{ID: "gauge1", MType: entity.Gauge, Value: utils.Ptr(123.0)},
		{ID: "counter1", MType: entity.Counter, Delta: utils.Ptr(int64(123))},
	})

	type args struct {
		id    string
		mType string
	}
	tests := []struct {
		name    string
		idx     int
		args    args
		want    entity.Metrics
		wantErr bool
	}{
		{
			name: "delete existing gauge metrics",
			args: args{
				id:    "gauge1",
				mType: entity.Gauge,
			},
		},
		{
			name: "delete non-existing gauge metrics",
			args: args{
				id:    "non-existing",
				mType: entity.Gauge,
			},
		},
		{
			name: "delete valid counter metrics",
			args: args{
				id:    "counter1",
				mType: entity.Counter,
			},
		},
		{
			name: "delete non-existing counter metrics",
			args: args{
				id:    "non-existing",
				mType: entity.Counter,
			},
		},
		{
			name: "delete invalid type metrics",
			args: args{
				id:    "gauge1",
				mType: "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var exists bool
			_, err := s.GetOne(ctx, tt.args.id, tt.args.mType)
			if err != nil {
				exists = true
			}

			err = s.DeleteOne(ctx, tt.args.id, tt.args.mType)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if exists {
				_, err = s.GetOne(ctx, tt.args.id, tt.args.mType)
				require.Equal(t, err, entity.ErrMetricNotFound)
			}
		})
	}

}

func Test_memoryStorage_DeleteAll(t *testing.T) {
	ctx := context.Background()

	want := &memoryStorage{
		mu:             sync.RWMutex{},
		counterStorage: make(map[string]int64),
		gaugeStorage:   make(map[string]float64),
	}

	s := newMemoryStorage()

	s.Update(ctx, entity.MetricsList{
		{ID: "gauge1", MType: entity.Gauge, Value: utils.Ptr(123.0)},
		{ID: "counter1", MType: entity.Counter, Delta: utils.Ptr(int64(123))},
	})

	err := s.DeleteAll(ctx)
	require.NoError(t, err)
	require.Equal(t, s.counterStorage, want.counterStorage)
	require.Equal(t, s.gaugeStorage, want.gaugeStorage)
}
