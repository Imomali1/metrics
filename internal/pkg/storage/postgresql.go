package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB(ctx context.Context, dsn string) (Storage, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	var pool *pgxpool.Pool
	err = utils.DoWithRetries(func() error {
		pool, err = pgxpool.NewWithConfig(ctx, config)
		return err
	})

	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	if err = CreateTable(ctx, pool); err != nil {
		return nil, err
	}

	return &DB{Pool: pool}, nil
}

func CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	var (
		counterTable = `
		CREATE TABLE IF NOT EXISTS counter (
    		id SERIAL PRIMARY KEY,
    		name TEXT NOT NULL UNIQUE,
    		delta BIGINT
		)`

		gaugeTable = `
		CREATE TABLE IF NOT EXISTS gauge (
    		id SERIAL PRIMARY KEY,
    		name TEXT NOT NULL UNIQUE,
    		value DOUBLE PRECISION
		)`
	)

	statements := []string{counterTable, gaugeTable}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	for _, statement := range statements {
		_, err = tx.Exec(ctx, statement)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *DB) Update(ctx context.Context, batch entity.MetricsList) error {
	if len(batch) == 0 {
		return nil
	}

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	queryInsertLayout := `INSERT INTO %[1]s (name, %[2]s)
						VALUES ($1, $2)
						ON CONFLICT (name)
						DO UPDATE
						SET %[2]s = EXCLUDED.%[2]s;`

	counterValMap := make(map[string]int64)

	for _, one := range batch {
		var query string

		if one.MType == entity.Counter {
			_, exists := counterValMap[one.ID]
			if !exists {
				var oldValue int64
				var oldCounter entity.Metrics
				oldCounter, err = s.GetOne(ctx, one.ID, one.MType)
				switch {
				case err == nil:
					oldValue = *oldCounter.Delta
				case errors.Is(err, entity.ErrMetricNotFound):
					oldValue = 0
				default:
					if errRollBack := tx.Rollback(ctx); errRollBack != nil {
						return fmt.Errorf("get_one error: %w; rollback error: %w", err, errRollBack)
					}
					return err
				}
				counterValMap[one.ID] = oldValue
			}

			counterValMap[one.ID] += *one.Delta

			query = fmt.Sprintf(queryInsertLayout, "counter", "delta")
			_, err = tx.Exec(ctx, query, one.ID, counterValMap[one.ID])
		} else if one.MType == entity.Gauge {
			query = fmt.Sprintf(queryInsertLayout, "gauge", "value")
			_, err = tx.Exec(ctx, query, one.ID, *one.Value)
		}

		if err != nil {
			if errRollBack := tx.Rollback(ctx); errRollBack != nil {
				return fmt.Errorf("exec error: %w; rollback error: %w", err, errRollBack)
			}
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *DB) GetOne(ctx context.Context, id, mType string) (entity.Metrics, error) {
	var metric = entity.Metrics{ID: id, MType: mType}
	switch mType {
	case entity.Counter:
		query := `SELECT delta FROM counter WHERE name = $1 LIMIT 1`
		var delta *int64
		if err := s.Pool.QueryRow(ctx, query, id).Scan(&delta); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return entity.Metrics{}, entity.ErrMetricNotFound
			}
			return entity.Metrics{}, err
		}
		metric.Delta = delta
	case entity.Gauge:
		query := `SELECT value FROM gauge WHERE name = $1 LIMIT 1`
		var value *float64
		if err := s.Pool.QueryRow(ctx, query, id).Scan(&value); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return entity.Metrics{}, entity.ErrMetricNotFound
			}
			return entity.Metrics{}, err
		}
		metric.Value = value
	}

	return metric, nil
}

func (s *DB) GetAll(ctx context.Context) (entity.MetricsList, error) {
	querySelectLayout := `SELECT name, %s FROM %s`

	colTables := [][]string{{"delta", entity.Counter}, {"value", entity.Gauge}}

	var list entity.MetricsList
	for _, colTable := range colTables {
		metrics, err := s.fetchMetrics(ctx, querySelectLayout, colTable[0], colTable[1])
		if err != nil {
			return entity.MetricsList{}, err
		}
		list = append(list, metrics...)
	}

	return list, nil
}

func (s *DB) fetchMetrics(
	ctx context.Context,
	queryLayout, colName, tableName string,
) (entity.MetricsList, error) {

	query := fmt.Sprintf(queryLayout, colName, tableName)
	rows, err := s.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var list entity.MetricsList
	for rows.Next() {
		var (
			name, mType string
			delta       *int64
			value       *float64
		)

		if tableName == entity.Counter {
			err = rows.Scan(&name, &delta)
			mType = entity.Counter
		} else {
			err = rows.Scan(&name, &value)
			mType = entity.Gauge
		}

		if err != nil {
			return nil, err
		}

		list = append(list, entity.Metrics{
			ID:    name,
			MType: mType,
			Delta: delta,
			Value: value,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (s *DB) DeleteOne(ctx context.Context, id, mType string) error {
	if mType != entity.Counter && mType != entity.Gauge {
		return entity.ErrInvalidMetricType
	}
	query := fmt.Sprintf(`DELETE FROM %s WHERE name = $1`, mType)
	_, err := s.Pool.Exec(ctx, query, id)
	return err
}

func (s *DB) DeleteAll(ctx context.Context) error {
	_, err := s.Pool.Exec(ctx, "DELETE FROM counter")
	if err != nil {
		return err
	}

	_, err = s.Pool.Exec(ctx, "DELETE FROM gauge")
	return err
}

func (s *DB) Ping(ctx context.Context) error {
	return s.Pool.Ping(ctx)
}

func (s *DB) Close() {
	s.Pool.Close()
}
