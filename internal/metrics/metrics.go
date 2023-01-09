package metrics

import (
	"fmt"
	"strconv"
)

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

type AtomicMetric struct {
	Name  string
	Type  string
	Value interface{}
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
