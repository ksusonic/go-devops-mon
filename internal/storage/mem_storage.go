package storage

import (
	"context"
	"fmt"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

type MemStorage struct {
	typeToNameMapping TypeToNameToMetric
	repository        metrics.Repository
}

func (m *MemStorage) Ping(context.Context) error {
	return fmt.Errorf("in-memory storage does not support ping")
}

func (m *MemStorage) Close() error {
	if m.repository != nil {
		return m.repository.Close()
	}
	// no additional actions needed
	return nil
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		typeToNameMapping: make(TypeToNameToMetric),
	}
}

func (m *MemStorage) SetMetric(_ context.Context, metric metrics.Metrics) (metrics.Metrics, error) {
	if metric.MType == metrics.CounterMType {
		var lastValue int64 = 0
		if found := m.typeToNameMapping.getMetric(metric); found != nil {
			lastValue = *found.Delta
		}
		value := lastValue + *metric.Delta
		metric.Delta = &value
	}

	m.typeToNameMapping.safeInsert(metric)

	return metric, nil
}

func (m *MemStorage) SetMetrics(ctx context.Context, metrics *[]metrics.Metrics) error {
	for _, metric := range *metrics {
		_, err := m.SetMetric(ctx, metric)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MemStorage) GetMetric(_ context.Context, type_, name string) (metrics.Metrics, error) {
	metric := m.typeToNameMapping.getMetric(metrics.Metrics{
		ID:    name,
		MType: type_,
	})
	if metric == nil {
		return metrics.Metrics{}, fmt.Errorf("metric %s of type %s not found", name, type_)
	}
	return *metric, nil
}

func (m *MemStorage) GetAllMetrics(context.Context) ([]metrics.Metrics, error) {
	var result []metrics.Metrics
	for _, t := range m.typeToNameMapping {
		for _, m := range t {
			result = append(result, m)
		}
	}
	return result, nil
}

func (m *MemStorage) GetMappedByTypeAndNameMetrics(context.Context) (map[string]map[string]interface{}, error) {
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
	return res, nil
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
