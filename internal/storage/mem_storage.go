package storage

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

type TypeToNameToMetric map[string]map[string]metrics.Metrics

type MemStorage struct {
	typeToNameMapping TypeToNameToMetric
}

func NewMemStorage() *MemStorage {
	var typeToNameToMetric = make(TypeToNameToMetric)

	// init map for known types
	typeToNameToMetric[metrics.GaugeMType] = make(map[string]metrics.Metrics)
	typeToNameToMetric[metrics.CounterMType] = make(map[string]metrics.Metrics)

	return &MemStorage{
		typeToNameMapping: typeToNameToMetric,
	}
}

func (m *MemStorage) RepositoryDropRoutine(repository metrics.Repository, duration time.Duration) {
	ticker := time.NewTicker(duration)
	for {
		<-ticker.C
		err := repository.SaveMetrics(m.GetAllMetrics())
		if err != nil {
			log.Println("Error while saving metrics to repository: ", err)
		}
	}
}

func (m *MemStorage) SetMetric(metric metrics.Metrics) metrics.Metrics {
	var result *metrics.Metrics
	if metric.MType == metrics.CounterMType {
		var lastValue int64 = 0
		_, ok := m.typeToNameMapping[metric.MType][metric.ID]
		if ok {
			lastValue = *m.typeToNameMapping[metric.MType][metric.ID].Delta
		}
		value := lastValue + *metric.Delta
		result = &metrics.Metrics{
			ID:    metric.ID,
			MType: metrics.CounterMType,
			Delta: &value,
		}
		m.typeToNameMapping[metric.MType][metric.ID] = *result
	} else {
		result = &metric
		m.typeToNameMapping[metric.MType][metric.ID] = *result
	}
	return *result
}

func (m *MemStorage) AddMetrics(atomicMetrics []metrics.Metrics) {
	for i := range atomicMetrics {
		m.SetMetric(atomicMetrics[i])
	}
}

func (m *MemStorage) GetMetric(type_, name string) (metrics.Metrics, error) {
	_, ok := m.typeToNameMapping[type_]
	if !ok {
		return metrics.Metrics{}, fmt.Errorf("no metric type '%s'", type_)
	}
	value, ok := m.typeToNameMapping[type_][name]
	if ok {
		return value, nil
	} else {
		return metrics.Metrics{}, fmt.Errorf("no metric '%s'", name)
	}
}

func (m *MemStorage) GetAllMetrics() []metrics.Metrics {
	var result []metrics.Metrics
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
			_, ok := res[m.MType]
			if !ok {
				res[m.MType] = make(map[string]interface{})
			}
			if m.MType == metrics.GaugeMType {
				res[m.MType][m.ID] = *m.Value
			} else if m.MType == metrics.CounterMType {
				res[m.MType][m.ID] = *m.Delta
			}
		}
	}
	return res
}

func (m *MemStorage) IncPollCount() {
	metric, ok := m.typeToNameMapping[metrics.CounterMType]["PollCount"]
	var previousValue int64 = 0
	if ok {
		previousValue = *metric.Delta
	}

	value := previousValue + 1
	m.typeToNameMapping[metrics.CounterMType]["PollCount"] = metrics.Metrics{
		ID:    "PollCount",
		MType: metrics.CounterMType,
		Delta: &value,
	}
}
func (m *MemStorage) RandomizeRandomValue() {
	value := rand.Float64()
	m.typeToNameMapping[metrics.GaugeMType]["RandomValue"] = metrics.Metrics{
		ID:    "RandomValue",
		MType: metrics.GaugeMType,
		Value: &value,
	}
}
