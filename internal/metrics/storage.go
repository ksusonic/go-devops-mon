package metrics

import (
	"context"
	"time"
)

type Repository interface {
	SaveMetrics([]Metrics) error
	ReadCurrentState() []Metrics
	Info() string
	Close() error
}

type ServerMetricStorage interface {
	// SetMetric Set value to metric
	SetMetric(m Metrics, secretKey *string) Metrics

	// GetMetric Get metric or error
	GetMetric(type_, name string) (Metrics, error)
	// GetAllMetrics Get all metrics as slice
	GetAllMetrics() []Metrics
	// GetMappedByTypeAndNameMetrics Get mapping of type -> name -> value
	GetMappedByTypeAndNameMetrics() map[string]map[string]interface{}

	RepositoryDropRoutine(context.Context, time.Duration)
}

type AgentMetricStorage interface {
	// SetMetric Set value to metric
	SetMetric(m Metrics, secretKey *string) Metrics
	// GetAllMetrics Get all metrics as slice
	GetAllMetrics() []Metrics
}
