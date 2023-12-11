package storage

type GaugeStorage interface {
	UpdateGaugeValue(key string, value float64)
	SetGaugeValue(key string, value float64)
	GetGaugeValue(key string) float64
}

type CounterStorage interface {
	UpdateCounterValue(key string, value int64)
	SetCounterValue(key string, value int64)
	GetCounterValue(key string) int64
}

type Storage struct {
	GaugeStorage
	CounterStorage
}

func NewStorage() *Storage {
	return &Storage{
		GaugeStorage:   newGaugeStorage(),
		CounterStorage: newCounterStorage(),
	}
}
