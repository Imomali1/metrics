package v1

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type postgresClient struct {
	pool *pgxpool.Pool
}

func newPostgresClient(ctx context.Context, dsn string) (*postgresClient, error) {
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

	return &postgresClient{pool: pool}, nil
}

func createTable(ctx context.Context, pool *pgxpool.Pool) error {
	const schema = `
	CREATE TABLE IF NOT EXISTS metrics (
    	id VARCHAR(63) PRIMARY KEY,
    	type VARCHAR(15) NOT NULL,
    	delta INTEGER,
    	value DOUBLE PRECISION
	)`

	_, err := pool.Exec(ctx, schema)
	return err
}

func (c *postgresClient) Ping(ctx context.Context) error {
	return c.pool.Ping(ctx)
}

func (c *postgresClient) Close() {
	c.pool.Close()
}
