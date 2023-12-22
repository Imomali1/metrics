package entity

type Metric struct {
	Type         string
	Name         string
	ValueCounter int64
	ValueGauge   float64
}

const (
	// Gauge Counter are metric types
	Gauge   = "gauge"
	Counter = "counter"
)
