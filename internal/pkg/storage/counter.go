package storage

type counterStorage struct {
	memStorage map[string]int64
}

func newCounterStorage() *counterStorage {
	return &counterStorage{memStorage: make(map[string]int64)}
}

func (s *counterStorage) UpdateCounterValue(key string, value int64) {
	s.memStorage[key] += value
}

func (s *counterStorage) SetCounterValue(key string, value int64) {
	s.memStorage[key] = value
}

func (s *counterStorage) GetCounterValue(key string) int64 {
	return s.memStorage[key]
}
