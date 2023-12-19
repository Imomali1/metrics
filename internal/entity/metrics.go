package entity

type MetricType string

type Metric struct {
	Type  MetricType
	Name  string
	Value interface{}
}
