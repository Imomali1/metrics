package storage

type gaugeStorage struct {
	memStorage map[string]float64
}

func newGaugeStorage() *gaugeStorage {
	return &gaugeStorage{memStorage: make(map[string]float64)}
}

func (s *gaugeStorage) UpdateGaugeValue(key string, value float64) {
	s.memStorage[key] = value
}

func (s *gaugeStorage) SetGaugeValue(key string, value float64) {
	s.memStorage[key] = value
}

func (s *gaugeStorage) GetGaugeValue(key string) float64 {
	return s.memStorage[key]
}
