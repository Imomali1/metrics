package storage

import (
	"context"
	"fmt"
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

type testDB struct {
	dsn       string
	pool      *pgxpool.Pool
	container testcontainers.Container
}

func newTestDB(opts ...testDBOption) (db testDB, err error) {
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
		return db, fmt.Errorf("failed to start container: %w", err)
	}

	ipAddress, err := postgresContainer.Host(ctx)
	if err != nil {
		return db, fmt.Errorf("failed to get container host: %w", err)
	}

	mappedPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return db, fmt.Errorf("failed to get mapped port: %w", err)
	}

	dsn := fmt.Sprintf("postgres://user:password@%s:%s/testdb", ipAddress, mappedPort.Port())

	db = testDB{
		dsn:       dsn,
		container: postgresContainer,
	}

	for _, opt := range opts {
		err = opt(&db)
		if err != nil {
			return testDB{}, err
		}
	}

	return db, err
}

type testDBOption func(db *testDB) error

func WithPool() testDBOption {
	return func(db *testDB) error {
		ctx := context.Background()

		config, err := pgxpool.ParseConfig(db.dsn)
		if err != nil {
			return fmt.Errorf("failed to parse config: %w", err)
		}

		var pool *pgxpool.Pool
		err = utils.DoWithRetries(func() error {
			pool, err = pgxpool.NewWithConfig(ctx, config)
			return err
		})

		if err != nil {
			return fmt.Errorf("failed to connect to pool: %w", err)
		}

		if err = pool.Ping(ctx); err != nil {
			return fmt.Errorf("failed to ping pool: %w", err)
		}

		db.pool = pool

		return nil
	}
}

func (db *testDB) close() {
	if db == nil {
		return
	}

	if db.pool != nil {
		db.pool.Close()
	}

	if db.container == nil {
		return
	}

	err := db.container.Terminate(context.Background())
	if err != nil {
		log.Printf("failed to terminate container: %v", err)
	}
}

func Test_newDBStorage(t *testing.T) {
	db, err := newTestDB()
	defer db.close()
	require.NoError(t, err)

	tests := []struct {
		name    string
		dsn     string
		wantErr bool
	}{
		{
			name:    "valid dsn",
			dsn:     db.dsn,
			wantErr: false,
		},
		{
			name:    "invalid dsn",
			dsn:     "",
			wantErr: true,
		},
		{
			name:    "not existing database",
			dsn:     db.dsn + "non-existing",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			var dbs *dbStorage
			dbs, err = newDBStorage(ctx, tt.dsn)
			if tt.wantErr {
				require.Error(t, err)
				t.Log(err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, dbs)
		})
	}
}

func Test_createTable(t *testing.T) {
	db, err := newTestDB(WithPool())
	defer db.close()
	require.NoError(t, err)
	require.NotNil(t, db.pool)

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
			err = createTable(context.Background(), db.pool)
			require.NoError(t, err)
		})
	}
}

func Test_dbStorage_Update(t *testing.T) {
	db, err := newTestDB()
	defer db.close()
	require.NoError(t, err)

	ctx := context.Background()
	s, err := newDBStorage(ctx, db.dsn)
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
			err = s.DeleteAll(ctx)
			require.NoError(t, err)

			err = s.Update(ctx, tt.batch)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			var got entity.MetricsList
			got, err = s.GetAll(ctx)
			require.NoError(t, err)
			require.ElementsMatch(t, tt.wanted, got)
		})
	}
}

func Test_dbStorage_GetAll(t *testing.T) {
	db, err := newTestDB()
	defer db.close()
	require.NoError(t, err)

	ctx := context.Background()

	s, err := newDBStorage(ctx, db.dsn)
	require.NoError(t, err)
	require.NotNil(t, s)

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
				_ = s.DeleteAll(ctx)
				_ = s.Update(ctx, metrics)
			},
			want:    metrics,
			wantErr: false,
		},
		{
			name: "empty",
			updateFunc: func() {
				_ = s.DeleteAll(ctx)
			},
			want:    entity.MetricsList{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.updateFunc()

			got, err := s.GetAll(ctx)
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
	db, err := newTestDB()
	defer db.close()
	require.NoError(t, err)

	ctx := context.Background()

	s, err := newDBStorage(ctx, db.dsn)
	require.NoError(t, err)
	require.NotNil(t, s)

	_ = s.Update(ctx, entity.MetricsList{
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

func Test_dbStorage_DeleteOne(t *testing.T) {
	db, err := newTestDB()
	defer db.close()
	require.NoError(t, err)

	ctx := context.Background()

	s, err := newDBStorage(ctx, db.dsn)
	require.NoError(t, err)
	require.NotNil(t, s)

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
			updateFunc: func() { _ = s.Update(ctx, metrics) },
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
			updateFunc: func() { _ = s.DeleteAll(ctx) },
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.updateFunc()
			err = s.DeleteOne(ctx, tt.args.id, tt.args.mType)
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
	db, err := newTestDB(WithPool())
	defer db.close()
	require.NoError(t, err)
	require.NotNil(t, db.pool)

	s := &dbStorage{
		pool: db.pool,
	}

	err = s.Ping(context.Background())
	require.NoError(t, err)
}

func Test_dbStorage_Close(t *testing.T) {
	db, err := newTestDB(WithPool())
	defer db.close()
	require.NoError(t, err)
	require.NotNil(t, db.pool)

	s := &dbStorage{
		pool: db.pool,
	}

	s.Close()
}
