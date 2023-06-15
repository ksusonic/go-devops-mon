package metrics

import "context"

type DecryptService interface {
	DecryptBytes(b []byte) ([]byte, error)
}

type ServerMetricStorage interface {
	// SetMetric Set value to metric
	SetMetric(ctx context.Context, m *Metric) (*Metric, error)
	SetMetrics(ctx context.Context, m []*Metric) error

	// GetMetric Get metric or error
	GetMetric(ctx context.Context, type_, name string) (*Metric, error)
	// GetAllMetrics Get all metrics as slice
	GetAllMetrics(ctx context.Context) ([]Metric, error)
	// GetMappedByTypeAndNameMetrics Get mapping of type -> name -> value
	GetMappedByTypeAndNameMetrics(ctx context.Context) (map[string]map[string]interface{}, error)

	Close() error
	Ping(ctx context.Context) error
}
