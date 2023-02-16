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

func (m *AgentStorage) SetMetric(metric metrics.Metrics) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	if metric.MType == metrics.CounterMType {
		var lastValue int64 = 0
		if found := m.typeToNameMapping.getMetric(metric); found != nil {
			lastValue = *found.Delta
		}
		value := lastValue + *metric.Delta
		metric.Delta = &value
	}

	m.typeToNameMapping.safeInsert(metric)

	return nil
}

func (m *AgentStorage) GetAllMetrics() []metrics.Metrics {
	m.mux.RLock()
	defer m.mux.RUnlock()
	var result []metrics.Metrics
	for _, t := range m.typeToNameMapping {
		for _, m := range t {
			result = append(result, m)
		}
	}
	return result
}
