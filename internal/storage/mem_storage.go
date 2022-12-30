package storage

import (
	"fmt"
	"math/rand"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

type MemStorage struct {
	nameMapping map[string]metrics.AtomicMetric
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		nameMapping: make(map[string]metrics.AtomicMetric),
	}
}

func (m *MemStorage) SetMetric(metric metrics.AtomicMetric) {
	if metric.Type == metrics.CounterType {
		var lastValue int64 = 0
		_, ok := (*m).nameMapping[metric.Name]
		if ok {
			lastValue = (*m).nameMapping[metric.Name].Value.(int64)
		}
		(*m).nameMapping[metric.Name] = metrics.AtomicMetric{
			Name:  metric.Name,
			Type:  metrics.CounterType,
			Value: lastValue + metric.Value.(int64),
		}
	} else {
		(*m).nameMapping[metric.Name] = metric
	}
}

func (m *MemStorage) AddMetrics(atomicMetrics []metrics.AtomicMetric) {
	for i := range atomicMetrics {
		m.SetMetric(atomicMetrics[i])
	}
}

func (m *MemStorage) GetMetric(name string) (metrics.AtomicMetric, error) {
	value, ok := m.nameMapping[name]
	if ok {
		return value, nil
	} else {
		return metrics.AtomicMetric{}, fmt.Errorf("no metric '%s'", name)
	}
}

func (m *MemStorage) GetAllMetrics() []metrics.AtomicMetric {
	var result []metrics.AtomicMetric
	for _, m := range m.nameMapping {
		result = append(result, m)
	}
	return result
}

func (m *MemStorage) GetAllTypedMetrics(type_ string) []metrics.AtomicMetric {
	var result []metrics.AtomicMetric
	for _, m := range m.nameMapping {
		if m.Type == type_ {
			result = append(result, m)
		}
	}
	return result
}

func (m *MemStorage) GetMappedByTypeAndNameMetrics() map[string]map[string]interface{} {
	res := make(map[string]map[string]interface{})
	for _, m := range m.nameMapping {
		_, ok := res[m.Type]
		if !ok {
			res[m.Type] = make(map[string]interface{})
		}
		res[m.Type][m.Name] = m.Value
	}
	return res
}

func (m *MemStorage) IncPollCount() {
	metric, ok := (*m).nameMapping["PollCount"]
	var previousValue int64 = 0
	if ok {
		previousValue = metric.Value.(int64)
	}

	(*m).nameMapping["PollCount"] = metrics.AtomicMetric{
		Name:  "PollCount",
		Type:  metrics.CounterType,
		Value: previousValue + 1,
	}
}
func (m *MemStorage) RandomizeRandomValue() {
	(*m).nameMapping["RandomValue"] = metrics.AtomicMetric{
		Name:  "RandomValue",
		Type:  metrics.CounterType,
		Value: rand.Int63(),
	}
}
