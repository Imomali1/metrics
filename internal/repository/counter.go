package repository

import "github.com/Imomali1/metrics/internal/pkg/storage"

type counterRepository struct {
	memStorage storage.CounterStorage
}

func newCounterRepository(memStorage storage.CounterStorage) *counterRepository {
	return &counterRepository{memStorage: memStorage}
}

func (r *counterRepository) UpdateCounter(name string, counter int64) error {
	r.memStorage.UpdateCounterValue(name, counter)
	return nil
}
