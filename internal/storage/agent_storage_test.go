package storage

import (
	"github.com/ksusonic/go-devops-mon/internal/hash"
	"testing"

	"github.com/ksusonic/go-devops-mon/internal/metrics"

	"github.com/stretchr/testify/assert"
)

func TestAgentStorage_SetMetric_GetAllMetrics(t *testing.T) {
	var constFloatValue = 123.43
	var constIntValue int64 = 123456

	allMetrics := []metrics.Metrics{
		{
			ID:    "testMetric1",
			MType: "gauge",
			Value: &constFloatValue,
		},
		{
			ID:    "testMetric2",
			MType: "gauge",
			Value: &constFloatValue,
		},
		{
			ID:    "testMetric1",
			MType: "counter",
			Delta: &constIntValue,
		},
		{
			ID:    "testMetric2",
			MType: "counter",
			Delta: &constIntValue,
		},
	}

	t.Run(t.Name(), func(t *testing.T) {
		agentStorage := NewAgentStorage()
		for _, m := range allMetrics {
			agentStorage.SetMetric(m, hash.NewService(""))
		}

		actual := agentStorage.GetAllMetrics()
		assert.Equal(t, 4, len(actual))

		for i := range actual {
			assert.Equal(t, allMetrics[i].String(), actual[i].String())
		}
	})
}
