package metrics

import (
	"fmt"
)

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
