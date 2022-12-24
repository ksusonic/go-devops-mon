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
	}
	tests := []struct {
		name string
		s    CounterStorage
		args args
	}{
		{
			name: "simple test #1",
			s: map[string][]int64{
				metrics.RandomValue: {1, 2, 3},
			},
			args: args{
				name:       metrics.RandomValue,
				value:      4,
				checkIndex: 3,
			},
		},
		{
			name: "add to empty test #2",
			s: map[string][]int64{
				metrics.RandomValue: {},
			},
			args: args{
				name:       metrics.RandomValue,
				value:      1,
				checkIndex: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.AddCounterValue(tt.args.name, tt.args.value)
			assert.Equal(t, tt.args.value, tt.s[tt.args.name][tt.args.checkIndex])
		})
	}
}
