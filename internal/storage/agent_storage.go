package storage

import (
	"sync"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

type AgentStorage struct {
	typeToNameMapping TypeToNameToMetric
	mux               sync.RWMutex
}

func NewAgentStorage() *AgentStorage {
	return &AgentStorage{
		typeToNameMapping: make(TypeToNameToMetric),
	}
}

func (m *AgentStorage) SetMetric(metric metrics.Metric) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	if metric.Type == metrics.CounterType {
		var lastValue int64 = 0
		if found := m.typeToNameMapping.getMetric(metric.ID, metric.Type); found != nil {
			lastValue = *found.Delta
		}
		value := lastValue + *metric.Delta
		metric.Delta = &value
	}

	m.typeToNameMapping.safeInsert(metric)

	return nil
}

func (m *AgentStorage) GetAllMetrics() []metrics.Metric {
	m.mux.RLock()
	defer m.mux.RUnlock()
	var result []metrics.Metric
	for _, t := range m.typeToNameMapping {
		for _, m := range t {
			result = append(result, m)
		}
	}
	return result
}
