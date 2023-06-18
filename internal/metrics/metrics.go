package metrics

import (
	"fmt"
	"net/http"
	"strings"

	metricspb "github.com/ksusonic/go-devops-mon/proto/metrics"
)

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

type GenericMetric interface {
	*Metric | *metricspb.Metric
	GetType() metricspb.MetricType
	GetID() string
	GetDelta() int64
	GetValue() float64
}

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	Type  string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // Значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

func (m Metric) GetID() string {
	return m.ID
}

func (m Metric) GetType() metricspb.MetricType {
	return m.TypeAsProtoType()
}

func (m Metric) GetDelta() int64 {
	if m.Delta != nil {
		return *m.Delta
	}
	return 0
}

func (m Metric) GetValue() float64 {
	if m.Value != nil {
		return *m.Value
	}
	return 0
}

func (m Metric) TypeAsProtoType() metricspb.MetricType {
	switch m.Type {
	case GaugeType:
		return metricspb.MetricType_gauge
	case CounterType:
		return metricspb.MetricType_counter
	default:
		return metricspb.MetricType_unknown
	}
}

func (m Metric) AsProto() *metricspb.Metric {
	res := metricspb.Metric{
		ID:   m.ID,
		Type: m.TypeAsProtoType(),
		Hash: m.Hash,
	}
	if m.Value != nil {
		res.Payload = &metricspb.Metric_Value{Value: *m.Value}
	} else if m.Delta != nil {
		res.Payload = &metricspb.Metric_Delta{Delta: *m.Delta}
	}
	return &res
}

func FromProto(m *metricspb.Metric) Metric {
	metric := Metric{
		ID:   m.GetID(),
		Type: m.GetType().String(),
		Hash: m.GetHash(),
	}
	if delta := m.GetDelta(); delta != 0 {
		metric.Delta = &delta
	}
	if value := m.GetValue(); value != 0 {
		metric.Value = &value
	}
	return metric
}

func (m Metric) Bind(*http.Request) error {
	if m.ID == "" {
		return fmt.Errorf("missing ID of metric")
	} else if m.Type == "" {
		return fmt.Errorf("missing type of metric")
	}

	return nil
}

func (m Metric) String() string {
	builder := strings.Builder{}
	switch m.Type {
	case CounterType:
		builder.WriteString(fmt.Sprintf("metric %s of type %s with value %d", m.ID, m.Type, *m.Delta))
	case GaugeType:
		builder.WriteString(fmt.Sprintf("metric %s of type %s with value %f", m.ID, m.Type, *m.Value))
	default:
		return ""
	}
	if m.Hash != "" {
		builder.WriteString(" and hash: ")
		builder.WriteString(m.Hash)
	}
	return builder.String()
}
