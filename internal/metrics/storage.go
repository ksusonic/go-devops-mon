package metrics

type ServerMetricStorage interface {
	// SetMetric Set value to metric
	SetMetric(AtomicMetric)

	// GetMetric Get metric or error
	GetMetric(type_, name string) (AtomicMetric, error)
	// GetMappedByTypeAndNameMetrics Get mapping of type -> name -> value
	GetMappedByTypeAndNameMetrics() map[string]map[string]interface{}
}

type AgentMetricStorage interface {
	// SetMetric Set value to metric
	SetMetric(AtomicMetric)
	AddMetrics([]AtomicMetric)

	// GetAllMetrics Get all metrics as slice
	GetAllMetrics() []AtomicMetric

	// IncPollCount Increases field PollCount by 1
	IncPollCount()
	// RandomizeRandomValue Set RandomValue to random number
	RandomizeRandomValue()
}
