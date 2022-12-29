package metrics

import (
	"fmt"
	"strconv"
)

type AtomicMetric struct {
	Name  string
	Type  string
	Value interface{}
}

type MetricStorage interface {
	// SetMetric Set value to metric
	SetMetric(metric AtomicMetric)

	// GetMetric Get metric or error
	GetMetric(name string) (AtomicMetric, error)
	// GetAllMetrics Get all metrics as slice
	GetAllMetrics() []AtomicMetric
	// GetMappedByTypeAndNameMetrics Get mapping of type -> name -> value
	GetMappedByTypeAndNameMetrics() map[string]map[string]interface{}

	// IncPollCount Increases field PollCount by 1
	IncPollCount()
	// RandomizeRandomValue Set RandomValue to random number
	RandomizeRandomValue()
}

func (am AtomicMetric) GetStringValue() (string, error) {
	var stringValue string
	switch am.Value.(type) {
	case float64:
		stringValue = strconv.FormatFloat(am.Value.(float64), 'f', -1, 64)
	case int64:
		stringValue = strconv.FormatInt(am.Value.(int64), 10)
	case string:
		stringValue = am.Value.(string)
	default:
		err := fmt.Errorf("unknown metric type: %s", am.Type)
		return "", err
	}
	return stringValue, nil
}
