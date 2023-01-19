package storage

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

type MemStorage struct {
	typeToNameMapping TypeToNameToMetric
	repository        metrics.Repository
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		typeToNameMapping: make(TypeToNameToMetric),
		repository:        nil,
	}
}

func NewMemStorageWithRepository(repository metrics.Repository, needToRestoreMetrics bool) *MemStorage {
	var typeToNameToMetric = make(TypeToNameToMetric)

	if needToRestoreMetrics {
		restored := repository.ReadCurrentState()
		if len(restored) > 0 {
			for _, m := range restored {
				typeToNameToMetric.safeInsert(m)
			}
			log.Printf("Restored %d metrics\n", len(restored))
		} else {
			log.Println("No metrics to restore")
		}
	}

	return &MemStorage{
		typeToNameMapping: typeToNameToMetric,
		repository:        repository,
	}
}

func (m *MemStorage) RepositoryDropRoutine(duration time.Duration) {
	if m.repository == nil {
		log.Fatal("Failed to launch RepositoryDropRoutine: repository is nil")
	}

	ticker := time.NewTicker(duration)
	for {
		<-ticker.C
		err := m.repository.SaveMetrics(m.GetAllMetrics())
		if err != nil {
			log.Println("Error while saving metrics to repository: ", err)
		}
	}
}

func (m *MemStorage) SetMetric(metric metrics.Metrics) metrics.Metrics {
	var result metrics.Metrics
	if metric.MType == metrics.CounterMType {
		var lastValue int64 = 0
		if m.typeToNameMapping.hasMetric(metric) {
			lastValue = *m.typeToNameMapping[metric.MType][metric.ID].Delta
		}
		value := lastValue + *metric.Delta
		result = metrics.Metrics{
			ID:    metric.ID,
			MType: metrics.CounterMType,
			Delta: &value,
		}
		m.typeToNameMapping.safeInsert(result)
	} else {
		result = metric
		m.typeToNameMapping.safeInsert(result)
	}
	return result
}

func (m *MemStorage) AddMetrics(atomicMetrics []metrics.Metrics) {
	for i := range atomicMetrics {
		m.SetMetric(atomicMetrics[i])
	}
}

func (m *MemStorage) GetMetric(type_, name string) (metrics.Metrics, error) {
	metric := m.typeToNameMapping.getMetric(metrics.Metrics{
		ID:    name,
		MType: type_,
	})
	if metric == nil {
		return metrics.Metrics{}, fmt.Errorf("metric %s of type %s not found", name, type_)
	}
	return *metric, nil
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
	var previousValue int64 = 0

	if currentMetric := m.typeToNameMapping.getMetric(metrics.Metrics{
		ID:    "PollCount",
		MType: metrics.CounterMType,
	}); currentMetric != nil {
		previousValue = *currentMetric.Delta
	}

	value := previousValue + 1
	m.typeToNameMapping.safeInsert(metrics.Metrics{
		ID:    "PollCount",
		MType: metrics.CounterMType,
		Delta: &value,
	})
}
func (m *MemStorage) RandomizeRandomValue() {
	value := rand.Float64()
	m.typeToNameMapping.safeInsert(metrics.Metrics{
		ID:    "RandomValue",
		MType: metrics.GaugeMType,
		Value: &value,
	})
}

type TypeToNameToMetric map[string]map[string]metrics.Metrics

func (t *TypeToNameToMetric) safeInsert(m metrics.Metrics) {
	_, ok := (*t)[m.MType]
	if !ok {
		(*t)[m.MType] = make(map[string]metrics.Metrics)
	}
	(*t)[m.MType][m.ID] = m
}

func (t *TypeToNameToMetric) hasMetric(m metrics.Metrics) bool {
	_, ok := (*t)[m.MType]
	if !ok {
		return false
	}
	_, ok = (*t)[m.MType][m.ID]
	return ok
}

func (t *TypeToNameToMetric) getMetric(m metrics.Metrics) *metrics.Metrics {
	_, ok := (*t)[m.MType]
	if !ok {
		return nil
	}
	metric, ok := (*t)[m.MType][m.ID]
	if !ok {
		return nil
	}
	return &metric
}
