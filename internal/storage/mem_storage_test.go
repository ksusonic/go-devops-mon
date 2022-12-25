package storage

import (
	"github.com/ksusonic/go-devops-mon/internal/metrics"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCounterStorage_AddCounterValue(t *testing.T) {
	type args struct {
		name       string
		value      int64
		checkIndex int
		expected   int64
	}
	tests := []struct {
		name string
		s    CounterStorage
		args args
	}{
		{
			name: "simple test #1",
			s: CounterStorage{
				metrics.RandomValue: {1, 2, 3},
			},
			args: args{
				name:       metrics.RandomValue,
				value:      7,
				checkIndex: 3,
				expected:   10,
			},
		},
		{
			name: "add to empty test #2",
			s: CounterStorage{
				metrics.RandomValue: {},
			},
			args: args{
				name:       metrics.RandomValue,
				value:      1,
				checkIndex: 0,
				expected:   1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.AddToCounterValue(tt.args.name, tt.args.value)
			assert.Equal(t, tt.args.expected, tt.s[tt.args.name][tt.args.checkIndex])
		})
	}
}

func TestGaugeStorage_AddGaugeValue(t *testing.T) {
	type args struct {
		name       string
		value      float64
		checkIndex int
	}
	tests := []struct {
		name string
		s    GaugeStorage
		args args
	}{
		{
			name: "simple test #1",
			s: GaugeStorage{
				metrics.PauseTotalNs: {1.001, 2.0003, 3.123},
			},
			args: args{
				name:       metrics.PauseTotalNs,
				value:      7.0003,
				checkIndex: 3,
			},
		},
		{
			name: "add to empty test #2",
			s: GaugeStorage{
				metrics.PauseTotalNs: {},
			},
			args: args{
				name:       metrics.PauseTotalNs,
				value:      1.123456,
				checkIndex: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.AddGaugeValue(tt.args.name, tt.args.value)
			assert.Equal(t, tt.args.value, tt.s[tt.args.name][tt.args.checkIndex])
		})
	}
}
