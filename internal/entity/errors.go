package entity

import (
	"errors"
)

var (
	ErrMetricNotFound    = errors.New("metric not found")
	ErrInvalidMetricType = errors.New("invalid metric type")
)
