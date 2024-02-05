//go:generate easyjson -no_std_marshalers metrics.go
package entity

//easyjson:json
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type MetricsWithoutPointer struct {
	ID    string
	MType string
	Delta int64
	Value float64
}

const (
	// Gauge Counter are metric types
	Gauge   = "gauge"
	Counter = "counter"
)
