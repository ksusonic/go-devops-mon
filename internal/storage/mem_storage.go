package storage

import (
	"fmt"
	"math/rand"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

type TypeToNameToMetric map[string]map[string]metrics.AtomicMetric

type MemStorage struct {
	typeToNameMapping TypeToNameToMetric
}

func NewMemStorage() *MemStorage {
	var typeToNameToMetric = make(TypeToNameToMetric)

	// init map for known types
	typeToNameToMetric[metrics.GaugeType] = make(map[string]metrics.AtomicMetric)
	typeToNameToMetric[metrics.CounterType] = make(map[string]metrics.AtomicMetric)

	return &MemStorage{
		typeToNameMapping: typeToNameToMetric,
	}
}

func (m *MemStorage) SetMetric(metric metrics.AtomicMetric) {
	if metric.Type == metrics.CounterType {
		var lastValue int64 = 0
		_, ok := (*m).typeToNameMapping[metric.Type][metric.Name]
		if ok {
			lastValue = (*m).typeToNameMapping[metric.Type][metric.Name].Value.(int64)
		}
		(*m).typeToNameMapping[metric.Type][metric.Name] = metrics.AtomicMetric{
			Name:  metric.Name,
			Type:  metrics.CounterType,
			Value: lastValue + metric.Value.(int64),
		}
	} else {
		(*m).typeToNameMapping[metric.Type][metric.Name] = metric
	}
}

func (m *MemStorage) AddMetrics(atomicMetrics []metrics.AtomicMetric) {
	for i := range atomicMetrics {
		m.SetMetric(atomicMetrics[i])
	}
}

func (m *MemStorage) GetMetric(type_, name string) (metrics.AtomicMetric, error) {
	value, ok := m.typeToNameMapping[type_][name]
	if ok {
		return value, nil
	} else {
		return metrics.AtomicMetric{}, fmt.Errorf("no metric '%s'", name)
	}
}

func (m *MemStorage) GetAllMetrics() []metrics.AtomicMetric {
	var result []metrics.AtomicMetric
	for _, t := range m.typeToNameMapping {
		for _, m := range t {
			result = append(result, m)
		}
	}
	return result
}

func (m *MemStorage) GetMappedByTypeAndNameMetrics() map[string]map[string]interface{} {
	res := make(map[string]map[string]interface{})
	for _, t := range m.typeToNameMapping {
		for _, m := range t {
			_, ok := res[m.Type]
			if !ok {
				res[m.Type] = make(map[string]interface{})
			}
			res[m.Type][m.Name] = m.Value
		}
	}
	return res
}

func (m *MemStorage) IncPollCount() {
	metric, ok := (*m).typeToNameMapping[metrics.CounterType]["PollCount"]
	var previousValue int64 = 0
	if ok {
		previousValue = metric.Value.(int64)
	}

	(*m).typeToNameMapping[metrics.CounterType]["PollCount"] = metrics.AtomicMetric{
		Name:  "PollCount",
		Type:  metrics.CounterType,
		Value: previousValue + 1,
	}
}
func (m *MemStorage) RandomizeRandomValue() {
	(*m).typeToNameMapping[metrics.CounterType]["RandomValue"] = metrics.AtomicMetric{
		Name:  "RandomValue",
		Type:  metrics.CounterType,
		Value: rand.Int63(),
	}
}
