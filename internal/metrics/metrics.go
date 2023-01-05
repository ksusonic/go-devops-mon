package metrics

import "fmt"

const (
	GaugeMType   = "gauge"
	CounterMType = "counter"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // Значение метрики в случае передачи gauge
}

func (m Metrics) String() string {
	switch m.MType {
	case CounterMType:
		return fmt.Sprintf("metric %s of type %s with value %d", m.ID, m.MType, *m.Delta)
	case GaugeMType:
		return fmt.Sprintf("metric %s of type %s with value %f", m.ID, m.MType, *m.Value)
	default:
		return ""
	}
}

type MetricStorage interface {
	// SetMetric Set value to metric
	SetMetric(Metrics)
	AddMetrics([]Metrics)

	// GetMetric Get metric or error
	GetMetric(type_, name string) (Metrics, error)
	// GetAllMetrics Get all metrics as slice
	GetAllMetrics() []Metrics
	// GetMappedByTypeAndNameMetrics Get mapping of type -> name -> value
	GetMappedByTypeAndNameMetrics() map[string]map[string]interface{}

	// IncPollCount Increases field PollCount by 1
	IncPollCount()
	// RandomizeRandomValue Set RandomValue to random number
	RandomizeRandomValue()
}
