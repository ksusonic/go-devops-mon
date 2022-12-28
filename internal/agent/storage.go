package agent

import (
	"math/rand"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

type GaugeMetrics map[string]float64
type CounterMetrics map[string]int64

type CollectorStorage struct {
	GaugeMetricStorage   GaugeMetrics
	CounterMetricStorage CounterMetrics
}

func NewCollectorStorage() CollectorStorage {
	return CollectorStorage{
		GaugeMetricStorage:   make(map[string]float64),
		CounterMetricStorage: make(map[string]int64),
	}
}

func (c *CounterMetrics) RandomizeRandomValue() {
	(*c)[metrics.RandomValue] = rand.Int63()
}

func (c *CounterMetrics) IncPollCountValue() {
	(*c)[metrics.PollCount]++
}
