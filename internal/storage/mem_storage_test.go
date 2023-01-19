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
			_, err := tt.memStorage.GetMetric(metrics.CounterMType, "PollCount")
			require.Error(t, err)

			var number int64 = 1
			tt.memStorage.SetMetric(metrics.Metrics{
				ID:    "PollCount",
				MType: metrics.CounterMType,
				Delta: &number,
			})
			value, err := tt.memStorage.GetMetric(metrics.CounterMType, "PollCount")
			require.NoError(t, err)
			require.NotNil(t, value, "value from storage is nil")
			var expected int64 = 1
			require.IsType(t, expected, *value.Delta)
			assert.Equal(t, expected, *value.Delta)
		})
	}
}

func TestMemStorage_SetMetric_GetMetric(t *testing.T) {
	var value = 7.0023

	tests := []struct {
		name       string
		memStorage *MemStorage
		args       metrics.Metrics
	}{
		{
			name:       "simple test #1",
			memStorage: NewMemStorage(),
			args: metrics.Metrics{
				ID:    "PauseTotalNs",
				MType: metrics.GaugeMType,
				Value: &value,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.memStorage.GetMetric(tt.args.MType, tt.args.ID)
			require.Error(t, err)

			tt.memStorage.SetMetric(tt.args)
			result, err := tt.memStorage.GetMetric(tt.args.MType, tt.args.ID)
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, tt.args, result)
		})
	}
}
