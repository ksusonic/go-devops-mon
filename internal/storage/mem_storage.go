package storage

import (
	"context"
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

type MemStorageRepository struct {
	Repository         metrics.Repository
	DropInterval       time.Duration
	NeedRestoreMetrics bool
}

func NewMemStorage(repository *MemStorageRepository) *MemStorage {
	memStorage := MemStorage{
		typeToNameMapping: make(TypeToNameToMetric),
		repository:        nil,
	}

	if repository != nil {
		memStorage.repository = repository.Repository
		if repository.NeedRestoreMetrics {
			restored := memStorage.repository.ReadCurrentState()
			if len(restored) > 0 {
				for _, m := range restored {
					memStorage.typeToNameMapping.safeInsert(m)
				}
				log.Printf("Restored %d metrics\n", len(restored))
			} else {
				log.Println("No metrics to restore")
			}
		}

		go memStorage.RepositoryDropRoutine(context.Background(), repository.DropInterval)
	}

	return &memStorage
}

func (m *MemStorage) RepositoryDropRoutine(ctx context.Context, duration time.Duration) {
	log.Printf("Started repository drop routine to %s with interval %s\n", m.repository.Info(), duration)
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			err := m.repository.SaveMetrics(m.GetAllMetrics())
			if err != nil {
				log.Println("Error while saving metrics to repository: ", err)
			}
		case <-ctx.Done():
			log.Println("Finished repository routine")
			return
		}
	}
}

func (m *MemStorage) SetMetric(metric metrics.Metrics) metrics.Metrics {
	if metric.MType == metrics.CounterMType {
		var lastValue int64 = 0
		if found := m.typeToNameMapping.getMetric(metric); found != nil {
			lastValue = *found.Delta
		}
		value := lastValue + *metric.Delta
		metric.Delta = &value
	}
	m.typeToNameMapping.safeInsert(metric)

	return metric
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

func (m *MemStorage) IncPollCount(secretKey string) {
	metric := m.typeToNameMapping.getMetric(metrics.Metrics{
		ID:    "PollCount",
		MType: metrics.CounterMType,
	})

	if metric != nil {
		*metric.Delta++
		if secretKey != "" {
			metric.Hash = metric.CalcHash(secretKey)
		}
	} else {
		var startValue int64 = 0
		res := metrics.Metrics{
			ID:    "PollCount",
			MType: metrics.CounterMType,
			Delta: &startValue,
		}
		if secretKey != "" {
			res.Hash = res.CalcHash(secretKey)
		}
		m.typeToNameMapping.safeInsert(res)
	}
}
func (m *MemStorage) RandomizeRandomValue(secretKey string) {
	value := rand.Float64()
	res := metrics.Metrics{
		ID:    "RandomValue",
		MType: metrics.GaugeMType,
		Value: &value,
	}
	if secretKey != "" {
		res.Hash = res.CalcHash(secretKey)
	}
	m.typeToNameMapping.safeInsert(res)
}

type TypeToNameToMetric map[string]map[string]metrics.Metrics

func (t *TypeToNameToMetric) safeInsert(m metrics.Metrics) {
	_, ok := (*t)[m.MType]
	if !ok {
		(*t)[m.MType] = make(map[string]metrics.Metrics)
	}
	(*t)[m.MType][m.ID] = m
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
