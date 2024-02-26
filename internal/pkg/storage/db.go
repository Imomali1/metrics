package storage

import (
	_ "github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
)

type postgresClient struct {
	db *sqlx.DB
}

func newPostgresClient(dsn string) (*postgresClient, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &postgresClient{db: db}, nil
}

func (c *postgresClient) Ping() error {
	return c.db.Ping()
}
