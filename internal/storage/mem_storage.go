package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

type MemStorage struct {
	typeToNameMapping TypeToNameToMetric
	repository        Repository
}

type Repository interface {
	SaveMetrics([]metrics.Metric) error
	ReadCurrentState() []metrics.Metric
	DebugInfo() string
	Close() error
	DropRoutine(ctx context.Context, getMetricsFunc func(context.Context) ([]metrics.Metric, error), duration time.Duration)
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

func (m *MemStorage) SetMetric(_ context.Context, metric *metrics.Metric) (*metrics.Metric, error) {
	if metric.Type == metrics.CounterType {
		var lastValue int64 = 0
		if found := m.typeToNameMapping.getMetric(metric.ID, metric.Type); found != nil {
			lastValue = *found.Delta
		}
		value := lastValue + *metric.Delta
		metric.Delta = &value
	}

	m.typeToNameMapping.safeInsert(*metric)

	return metric, nil
}

func (m *MemStorage) SetMetrics(ctx context.Context, metrics []*metrics.Metric) error {
	for _, metric := range metrics {
		_, err := m.SetMetric(ctx, metric)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MemStorage) GetMetric(_ context.Context, type_, name string) (*metrics.Metric, error) {
	metric := m.typeToNameMapping.getMetric(name, type_)
	if metric == nil {
		return nil, fmt.Errorf("metric %s of type %s not found", name, type_)
	}
	return metric, nil
}

func (m *MemStorage) GetAllMetrics(context.Context) ([]metrics.Metric, error) {
	var result []metrics.Metric
	for _, t := range m.typeToNameMapping {
		for _, m := range t {
			result = append(result, m)
		}
	}
	return result, nil
}

func (m *MemStorage) GetMappedByTypeAndNameMetrics(_ context.Context) (map[string]map[string]interface{}, error) {
	res := make(map[string]map[string]interface{})
	for _, t := range m.typeToNameMapping {
		for _, m := range t {
			_, ok := res[m.Type]
			if !ok {
				res[m.Type] = make(map[string]interface{})
			}
			if m.Type == metrics.GaugeType {
				res[m.Type][m.ID] = *m.Value
			} else if m.Type == metrics.CounterType {
				res[m.Type][m.ID] = *m.Delta
			}
		}
	}
	return res, nil
}

type TypeToNameToMetric map[string]map[string]metrics.Metric

func (t *TypeToNameToMetric) safeInsert(m metrics.Metric) {
	_, ok := (*t)[m.Type]
	if !ok {
		(*t)[m.Type] = make(map[string]metrics.Metric)
	}
	(*t)[m.Type][m.ID] = m
}

func (t *TypeToNameToMetric) getMetric(id, mtype string) *metrics.Metric {
	_, ok := (*t)[mtype]
	if !ok {
		return nil
	}
	metric, ok := (*t)[mtype][id]
	if !ok {
		return nil
	}
	return &metric
}
