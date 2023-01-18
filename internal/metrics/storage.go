package metrics

import "time"

type Repository interface {
	SaveMetrics([]Metrics) error
	ReadCurrentState() []Metrics
	Close() error
}

type ServerMetricStorage interface {
	// SetMetric Set value to metric
	SetMetric(Metrics) Metrics

	// GetMetric Get metric or error
	GetMetric(type_, name string) (Metrics, error)
	// GetAllMetrics Get all metrics as slice
	GetAllMetrics() []Metrics
	// GetMappedByTypeAndNameMetrics Get mapping of type -> name -> value
	GetMappedByTypeAndNameMetrics() map[string]map[string]interface{}

	RepositoryDropRoutine(time.Duration)
}

type AgentMetricStorage interface {
	// SetMetric Set value to metric
	SetMetric(Metrics) Metrics
	AddMetrics([]Metrics)

	// GetAllMetrics Get all metrics as slice
	GetAllMetrics() []Metrics

	// IncPollCount Increases field PollCount by 1
	IncPollCount()
	// RandomizeRandomValue Set RandomValue to random number
	RandomizeRandomValue()
}
