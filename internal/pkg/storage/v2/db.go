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
		access = `GRANT ALL PRIVILEGES ON DATABASE metrics TO metrics;`

		schema = `
	CREATE TABLE IF NOT EXISTS metrics (
    	id SERIAL PRIMARY KEY,
    	name TEXT NOT NULL,
    	type TEXT NOT NULL,
    	delta INTEGER,
    	value DOUBLE PRECISION
	)`
	)

	statements := []string{access, schema}

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
	// If metric exists just update, otherwise insert
	statusMap := make(map[string]bool)
	for _, one := range batch {
		metric, err := s.GetOne(ctx, one.ID, one.MType)
		if err != nil {
			return err
		}
		// Set true, if metric exists
		statusMap[one.ID] = metric.ID != ""
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}

	queryInsert := `INSERT INTO metrics(id, type, delta, value) VALUES ($1, $2, $3, $4)`
	queryUpdateDelta := `UPDATE metrics SET delta = $1 WHERE id = $2`
	queryUpdateValue := `UPDATE metrics SET value = $1 WHERE id = $2`

	for _, one := range batch {
		exists := statusMap[one.ID]
		if exists {
			_, err = tx.Exec(ctx, queryInsert, one)
		} else {
			if one.MType == entity.Counter {
				_, err = tx.Exec(ctx, queryUpdateDelta, *one.Delta, one.ID)
			} else if one.MType == entity.Gauge {
				_, err = tx.Exec(ctx, queryUpdateValue, *one.Value, one.ID)
			}
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
		query := `SELECT delta FROM metrics WHERE id = $1 LIMIT 1`
		delta := new(int64)
		if err := s.pool.QueryRow(ctx, query, id).Scan(&delta); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return entity.Metrics{}, nil
			}
			return entity.Metrics{}, err
		}
		metric.Delta = delta
	case entity.Gauge:
		query := `SELECT value FROM metrics WHERE id = $1 LIMIT 1`
		value := new(float64)
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
	query := `SELECT id, type, delta, value FROM metrics`
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list entity.MetricsList
	for rows.Next() {
		id, mType := new(string), new(string)
		delta, value := new(int64), new(float64)

		if err = rows.Scan(&id, &mType, &delta, &value); err != nil {
			return nil, err
		}

		list = append(list, entity.Metrics{
			ID:    *id,
			MType: *mType,
			Delta: delta,
			Value: value,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (s *dbStorage) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

func (s *dbStorage) Close() {
	s.pool.Close()
}
