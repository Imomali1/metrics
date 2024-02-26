package v2

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type dbStorage struct {
	pool *pgxpool.Pool
}

func newDBStorage(ctx context.Context, dsn string) (*dbStorage, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err = createTable(ctx, pool); err != nil {
		return nil, err
	}

	log.Println("DB successfully initialized")

	return &dbStorage{pool: pool}, nil
}

func createTable(ctx context.Context, pool *pgxpool.Pool) error {
	var (
		counterTable = `
		CREATE TABLE IF NOT EXISTS counter (
    		id SERIAL PRIMARY KEY,
    		name TEXT NOT NULL UNIQUE,
    		delta INTEGER
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
			if errRollBack := tx.Rollback(ctx); errRollBack != nil {
				return fmt.Errorf("exec error: %w; rollback error: %w", err, errRollBack)
			}
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *dbStorage) Update(ctx context.Context, batch entity.MetricsList) error {
	if len(batch) == 0 {
		return nil
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}

	queryInsertLayout := `INSERT INTO %[1]s (name, %[2]s)
						VALUES ($1, $2)
						ON CONFLICT (name)
						DO UPDATE
						SET %[2]s = EXCLUDED.%[2]s;`

	for _, one := range batch {
		var query string
		if one.MType == entity.Counter {
			query = fmt.Sprintf(queryInsertLayout, "counter", "delta")
			_, err = tx.Exec(ctx, query, one.ID, *one.Delta)
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

func (s *dbStorage) GetOne(ctx context.Context, id, mType string) (entity.Metrics, error) {
	var metric = entity.Metrics{ID: id, MType: mType}
	switch mType {
	case entity.Counter:
		query := `SELECT delta FROM counter WHERE name = $1 LIMIT 1`
		var delta *int64
		if err := s.pool.QueryRow(ctx, query, id).Scan(&delta); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return entity.Metrics{}, nil
			}
			return entity.Metrics{}, err
		}
		metric.Delta = delta
	case entity.Gauge:
		query := `SELECT value FROM gauge WHERE name = $1 LIMIT 1`
		var value *float64
		if err := s.pool.QueryRow(ctx, query, id).Scan(&value); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return entity.Metrics{}, nil
			}
			return entity.Metrics{}, err
		}
		metric.Value = value
	}

	return metric, nil
}

func (s *dbStorage) GetAll(ctx context.Context) (entity.MetricsList, error) {
	querySelectLayout := `SELECT name, %s FROM %s`

	colTables := [][]string{{"delta", "counter"}, {"value", "gauge"}}

	var list entity.MetricsList
	for i, colTable := range colTables {
		query := fmt.Sprintf(querySelectLayout, colTable[0], colTable[1])
		rows, err := s.pool.Query(ctx, query)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var (
				name, mType string
				delta       *int64
				value       *float64
			)
			if i == 0 {
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
	}

	return list, nil
}

func (s *dbStorage) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

func (s *dbStorage) Close() {
	s.pool.Close()
}
