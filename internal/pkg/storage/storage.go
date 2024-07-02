package storage

import (
	"context"
)

func New(ctx context.Context, dsn string) (Storage, error) {
	if dsn != "" {
		return NewDB(ctx, dsn)
	}

	return NewMemory()
}
