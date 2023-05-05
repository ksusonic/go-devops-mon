package metrics

import (
	"context"
	"time"
)

type HashService interface {
	// SetHash Calculates hash for metric
	SetHash(m *Metrics) error
	// ValidateHash Returns bool if calculated and metric hash is correct
	ValidateHash(m *Metrics) error
}

type EncryptService interface {
	EncryptBytes(b []byte) ([]byte, error)
}

type DecryptService interface {
	DecryptBytes(b []byte) ([]byte, error)
}

type Repository interface {
	SaveMetrics([]Metrics) error
	ReadCurrentState() []Metrics
	DebugInfo() string
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
