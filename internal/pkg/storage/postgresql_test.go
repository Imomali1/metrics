package storage_test

import (
	"context"
	"fmt"
	"github.com/Imomali1/metrics/internal/pkg/storage"
	"log"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/utils"
)

var (
	DSN           string
	TestContainer testcontainers.Container
)

func init() {
	startTestContainer()
}

func startTestContainer() {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:alpine",
		ExposedPorts: []string{"5432/tcp"},
		Hostname:     "test-postgres",
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
		},
		WaitingFor: wait.
			ForListeningPort("5432/tcp").
			WithStartupTimeout(2 * time.Minute),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to start container: %w", err))
	}

	ipAddress, err := postgresContainer.Host(ctx)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get container host: %w", err))
	}

	mappedPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get mapped port: %w", err))
	}

	DSN = fmt.Sprintf("postgres://user:password@%s:%s/testdb", ipAddress, mappedPort.Port())

	TestContainer = postgresContainer
}

func stopTestContainer() {
	ctx := context.Background()
	if err := TestContainer.Terminate(ctx); err != nil {
		log.Fatal(fmt.Errorf("failed to terminate container: %w", err))
	}
	TestContainer = nil
	DSN = ""
}

func restartTestContainer() {
	stopTestContainer()
	startTestContainer()
}

func createTestPool() (*pgxpool.Pool, error) {
	restartTestContainer()

	ctx := context.Background()

	config, err := pgxpool.ParseConfig(DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	var pool *pgxpool.Pool
	err = utils.DoWithRetries(func() error {
		pool, err = pgxpool.NewWithConfig(ctx, config)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to Pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping Pool: %w", err)
	}

	return pool, nil
}

func Test_newDBStorage(t *testing.T) {
	tests := []struct {
		name    string
		dsn     string
		wantErr bool
	}{
		{
			name:    "valid dsn",
			dsn:     DSN,
			wantErr: false,
		},
		{
			name:    "empty dsn",
			dsn:     "",
			wantErr: true,
		},
		{
			name:    "not existing database",
			dsn:     DSN + "non-existing",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := storage.NewDB(context.Background(), tt.dsn)
			if tt.wantErr {
				require.Error(t, err)
				t.Log(err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, db)

			err = db.Ping(context.Background())
			require.NoError(t, err)
		})
	}
}

func Test_createTable(t *testing.T) {
	pool, err := createTestPool()
	require.NoError(t, err)
	require.NotNil(t, pool)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "valid",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = storage.CreateTable(context.Background(), pool)
			require.NoError(t, err)
		})
	}
}

func Test_dbStorage_Update(t *testing.T) {
	ctx := context.Background()
	db, err := storage.NewDB(ctx, DSN)
	require.NoError(t, err)

	tests := []struct {
		name    string
		batch   entity.MetricsList
		wanted  entity.MetricsList
		wantErr bool
	}{
		{
			name:    "zero length metrics",
			batch:   entity.MetricsList{},
			wanted:  entity.MetricsList{},
			wantErr: false,
		},
		{
			name: "valid gauge metrics",
			batch: entity.MetricsList{
				{
					ID:    "gauge1",
					MType: entity.Gauge,
					Value: utils.Ptr(123.0),
				},
			},
			wanted: entity.MetricsList{
				{
					ID:    "gauge1",
					MType: entity.Gauge,
					Value: utils.Ptr(123.0),
				},
			},
			wantErr: false,
		},
		{
			name: "valid counter metrics",
			batch: entity.MetricsList{
				{
					ID:    "counter1",
					MType: entity.Counter,
					Delta: utils.Ptr(int64(123)),
				},
			},
			wanted: entity.MetricsList{
				{
					ID:    "counter1",
					MType: entity.Counter,
					Delta: utils.Ptr(int64(123)),
				},
			},
			wantErr: false,
		},
		{
			name: "valid both metrics",
			batch: entity.MetricsList{
				{
					ID:    "gauge1",
					MType: entity.Gauge,
					Value: utils.Ptr(123.0),
				},
				{
					ID:    "counter1",
					MType: entity.Counter,
					Delta: utils.Ptr(int64(123)),
				},
			},
			wanted: entity.MetricsList{
				{
					ID:    "gauge1",
					MType: entity.Gauge,
					Value: utils.Ptr(123.0),
				},
				{
					ID:    "counter1",
					MType: entity.Counter,
					Delta: utils.Ptr(int64(123)),
				},
			},
			wantErr: false,
		},
		{
			name: "one invalid metrics",
			batch: entity.MetricsList{
				{
					ID:    "invalid",
					MType: "invalid",
					Value: utils.Ptr(123.0),
					Delta: utils.Ptr(int64(123)),
				},
				{
					ID:    "counter1",
					MType: entity.Counter,
					Delta: utils.Ptr(int64(123)),
				},
			},
			wanted: entity.MetricsList{
				{
					ID:    "counter1",
					MType: entity.Counter,
					Delta: utils.Ptr(int64(123)),
				},
			},
			wantErr: false,
		},
		{
			name: "both invalid metrics",
			batch: entity.MetricsList{
				{
					ID:    "invalid",
					MType: "invalid",
					Value: utils.Ptr(123.0),
					Delta: utils.Ptr(int64(123)),
				},
				{
					ID:    "test",
					MType: "test",
					Value: utils.Ptr(123.0),
					Delta: utils.Ptr(int64(123)),
				},
			},
			wanted:  entity.MetricsList{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean storage for test
			err = db.DeleteAll(ctx)
			require.NoError(t, err)

			err = db.Update(ctx, tt.batch)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			var got entity.MetricsList
			got, err = db.GetAll(ctx)
			require.NoError(t, err)
			require.ElementsMatch(t, tt.wanted, got)
		})
	}
}

func Test_dbStorage_GetAll(t *testing.T) {
	ctx := context.Background()

	db, err := storage.NewDB(ctx, DSN)
	require.NoError(t, err)
	require.NotNil(t, db)

	metrics := entity.MetricsList{
		{
			ID:    "gauge1",
			MType: entity.Gauge,
			Value: utils.Ptr(123.0),
		},
		{
			ID:    "counter1",
			MType: entity.Counter,
			Delta: utils.Ptr(int64(123)),
		},
	}

	tests := []struct {
		name       string
		updateFunc func()
		want       entity.MetricsList
		wantErr    bool
	}{
		{
			name: "valid",
			updateFunc: func() {
				_ = db.DeleteAll(ctx)
				_ = db.Update(ctx, metrics)
			},
			want:    metrics,
			wantErr: false,
		},
		{
			name: "empty",
			updateFunc: func() {
				_ = db.DeleteAll(ctx)
			},
			want:    entity.MetricsList{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.updateFunc()

			got, err := db.GetAll(ctx)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			require.ElementsMatch(t, got, tt.want)
		})
	}
}

func Test_dbStorage_GetOne(t *testing.T) {
	ctx := context.Background()

	db, err := storage.NewDB(ctx, DSN)
	require.NoError(t, err)
	require.NotNil(t, db)

	_ = db.Update(ctx, entity.MetricsList{
		{
			ID:    "gauge1",
			MType: entity.Gauge,
			Value: utils.Ptr(123.0),
		},
		{
			ID:    "counter1",
			MType: entity.Counter,
			Delta: utils.Ptr(int64(123)),
		},
	})

	type args struct {
		id    string
		mType string
	}
	tests := []struct {
		name    string
		args    args
		want    entity.Metrics
		wantErr bool
	}{
		{
			name: "valid-gauge",
			args: args{
				id:    "gauge1",
				mType: entity.Gauge,
			},
			want: entity.Metrics{
				ID:    "gauge1",
				MType: entity.Gauge,
				Value: utils.Ptr(123.0),
			},
			wantErr: false,
		},
		{
			name: "valid-counter",
			args: args{
				id:    "counter1",
				mType: entity.Counter,
			},
			want: entity.Metrics{
				ID:    "counter1",
				MType: entity.Counter,
				Delta: utils.Ptr(int64(123)),
			},
			wantErr: false,
		},
		{
			name: "non-existing-gauge",
			args: args{
				id:    "gauge-non-existing",
				mType: entity.Gauge,
			},
			want: entity.Metrics{
				ID:    "gauge1",
				MType: entity.Gauge,
				Value: utils.Ptr(123.0),
			},
			wantErr: true,
		},
		{
			name: "non-existing-counter",
			args: args{
				id:    "counter-non-existing",
				mType: entity.Counter,
			},
			want: entity.Metrics{
				ID:    "counter1",
				MType: entity.Counter,
				Delta: utils.Ptr(int64(123)),
			},
			wantErr: true,
		},
		{
			name: "non-existing-type",
			args: args{
				id:    "non-existing-type",
				mType: "non-existing",
			},
			want: entity.Metrics{
				ID:    "non-existing-type",
				MType: "non-existing",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.GetOne(ctx, tt.args.id, tt.args.mType)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			require.Equal(t, got, tt.want)
		})
	}
}

func Test_dbStorage_DeleteOne(t *testing.T) {
	ctx := context.Background()

	db, err := storage.NewDB(ctx, DSN)
	require.NoError(t, err)
	require.NotNil(t, db)

	metrics := entity.MetricsList{
		{
			ID:    "gauge1",
			MType: entity.Gauge,
			Value: utils.Ptr(123.0),
		},
		{
			ID:    "counter1",
			MType: entity.Counter,
			Delta: utils.Ptr(int64(123)),
		},
	}

	type args struct {
		id    string
		mType string
	}
	tests := []struct {
		name       string
		args       args
		updateFunc func()
		wantErr    bool
	}{
		{
			name: "valid metrics deletion",
			args: args{
				id:    "gauge1",
				mType: entity.Gauge,
			},
			updateFunc: func() { _ = db.Update(ctx, metrics) },
			wantErr:    false,
		},
		{
			name: "invalid typed metrics deletion",
			args: args{
				id:    "gauge1",
				mType: "invalid",
			},
			updateFunc: func() {},
			wantErr:    true,
		},
		{
			name: "non-existing metrics deletion",
			args: args{
				id:    "non-existing",
				mType: entity.Counter,
			},
			updateFunc: func() {},
			wantErr:    false,
		},
		{
			name: "valid metrics deletion after removing all metrics",
			args: args{
				id:    "gauge1",
				mType: entity.Gauge,
			},
			updateFunc: func() { _ = db.DeleteAll(ctx) },
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.updateFunc()
			err = db.DeleteOne(ctx, tt.args.id, tt.args.mType)
			if tt.wantErr {
				t.Log(err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func Test_dbStorage_Ping(t *testing.T) {
	pool, err := createTestPool()
	require.NoError(t, err)
	require.NotNil(t, pool)

	s := &storage.DB{
		Pool: pool,
	}

	err = s.Ping(context.Background())
	require.NoError(t, err)
}

func Test_dbStorage_Close(t *testing.T) {
	pool, err := createTestPool()
	require.NoError(t, err)
	require.NotNil(t, pool)

	s := &storage.DB{
		Pool: pool,
	}

	s.Close()
}
