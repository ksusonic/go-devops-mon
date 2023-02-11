package metrics

import (
	"context"
	"time"
)

type HashService interface {
	SetHash(m *Metrics) error
	ValidateHash(m *Metrics) error
}

type Repository interface {
	SaveMetrics([]Metrics) error
	ReadCurrentState() []Metrics
	Info() string
	Close() error
	DropRoutine(ctx context.Context, getMetricsFunc func(context.Context) ([]Metrics, error), duration time.Duration)
}

type ServerMetricStorage interface {
	// SetMetric Set value to metric
	SetMetric(ctx context.Context, m Metrics) (Metrics, error)
	SetMetrics(ctx context.Context, m *[]Metrics) error

	// GetMetric Get metric or error
	GetMetric(ctx context.Context, type_, name string) (Metrics, error)
	// GetAllMetrics Get all metrics as slice
	GetAllMetrics(ctx context.Context) ([]Metrics, error)
	// GetMappedByTypeAndNameMetrics Get mapping of type -> name -> value
	GetMappedByTypeAndNameMetrics(ctx context.Context) (map[string]map[string]interface{}, error)

	Close() error
	Ping(ctx context.Context) error
}

type AgentMetricStorage interface {
	// SetMetric Set value to metric
	SetMetric(m Metrics) error
	// GetAllMetrics Get all metrics as slice
	GetAllMetrics() []Metrics
}
