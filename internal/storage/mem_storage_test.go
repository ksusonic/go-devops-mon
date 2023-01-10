package storage

import (
	"testing"

	"github.com/ksusonic/go-devops-mon/internal/metrics"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemStorage_IncPollCount(t *testing.T) {
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
			_, err := tt.memStorage.GetMetric("counter", "PollCount")
			require.Error(t, err)

			tt.memStorage.SetMetric(metrics.AtomicMetric{
				Name:  "PollCount",
				Type:  metrics.CounterType,
				Value: int64(1),
			})
			value, err := tt.memStorage.GetMetric("counter", "PollCount")
			require.NoError(t, err)
			require.NotNil(t, value, "value from storage is nil")
			var expected int64 = 1
			require.IsType(t, expected, value.Value)
			assert.Equal(t, expected, value.Value)
		})
	}
}

func TestMemStorage_SetMetric_GetMetric(t *testing.T) {

	tests := []struct {
		name       string
		memStorage *MemStorage
		args       metrics.AtomicMetric
	}{
		{
			name:       "simple test #1",
			memStorage: NewMemStorage(),
			args: metrics.AtomicMetric{
				Name:  "PauseTotalNs",
				Type:  metrics.GaugeType,
				Value: 7.0023,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.memStorage.GetMetric(tt.args.Type, tt.args.Name)
			require.Error(t, err)

			tt.memStorage.SetMetric(tt.args)
			result, err := tt.memStorage.GetMetric(tt.args.Type, tt.args.Name)
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, tt.args, result)
		})
	}
}
