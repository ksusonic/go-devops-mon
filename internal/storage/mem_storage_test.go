package storage

import (
	"context"
	"math/rand"
	"testing"

	"github.com/ksusonic/go-devops-mon/internal/metrics"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemStorage_IncPollCount(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name       string
		memStorage *MemStorage
	}{
		{
			name:       "add to empty test",
			memStorage: NewMemStorage(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.memStorage.GetMetric(ctx, metrics.CounterType, "PollCount")
			require.Error(t, err)

			var number int64 = 1
			tt.memStorage.SetMetric(ctx, &metrics.Metric{
				ID:    "PollCount",
				Type:  metrics.CounterType,
				Delta: &number,
			})
			value, err := tt.memStorage.GetMetric(ctx, metrics.CounterType, "PollCount")
			require.NoError(t, err)
			require.NotNil(t, value, "value from storage is nil")
			var expected int64 = 1
			require.IsType(t, expected, *value.Delta)
			assert.Equal(t, expected, *value.Delta)
		})
	}
}

func TestMemStorage_SetMetric_GetMetric(t *testing.T) {
	ctx := context.Background()
	var value = 7.0023

	tests := []struct {
		name       string
		memStorage *MemStorage
		args       metrics.Metric
	}{
		{
			name:       "simple test #1",
			memStorage: NewMemStorage(),
			args: metrics.Metric{
				ID:    "PauseTotalNs",
				Type:  metrics.GaugeType,
				Value: &value,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.memStorage.GetMetric(ctx, tt.args.Type, tt.args.ID)
			require.Error(t, err)

			tt.memStorage.SetMetric(ctx, &tt.args)
			result, err := tt.memStorage.GetMetric(ctx, tt.args.Type, tt.args.ID)
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, tt.args, *result)
		})
	}
}

func BenchmarkMemStorage(b *testing.B) {
	GenerateMetricPool := func(size int) []*metrics.Metric {
		var metricPool = make([]*metrics.Metric, size)
		for i := 0; i < size; i++ {
			value := rand.Float64()
			metricPool[i] = &metrics.Metric{
				ID:    "metric" + string(rune(rand.Intn(3))),
				Type:  metrics.GaugeType,
				Value: &value,
			}
		}
		return metricPool
	}

	const defaultPoolSize = 5
	memStorage := NewMemStorage()

	b.Run("SetMetrics", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			metricPool := GenerateMetricPool(defaultPoolSize)
			b.StartTimer()

			err := memStorage.SetMetrics(context.Background(), metricPool)
			if err != nil {
				b.Error(err)
			}
		}
	})
	b.Run("SetMetric", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			metricPool := GenerateMetricPool(defaultPoolSize)
			b.StartTimer()

			for j := 0; j < len(metricPool); j++ {
				_, err := memStorage.SetMetric(context.Background(), metricPool[j])
				if err != nil {
					b.Error(err)
				}
			}
		}
	})
	b.Run("GetMetric", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			metricPool := GenerateMetricPool(1)
			b.StartTimer()

			metric := &metricPool[0]
			_, err := memStorage.GetMetric(context.Background(), (*metric).Type, (*metric).ID)
			if err != nil {
				b.Error(err)
			}
		}
	})

}
