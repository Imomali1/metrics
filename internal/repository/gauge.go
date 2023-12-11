package repository

import "github.com/Imomali1/metrics/internal/pkg/storage"

type gaugeRepository struct {
	memStorage storage.GaugeStorage
}

func newGaugeRepository(memStorage storage.GaugeStorage) *gaugeRepository {
	return &gaugeRepository{memStorage: memStorage}
}

func (r *gaugeRepository) UpdateGauge(name string, gauge float64) error {
	r.memStorage.UpdateGaugeValue(name, gauge)
	return nil
}
