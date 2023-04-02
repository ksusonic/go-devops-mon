package hash

import (
	"fmt"
	"testing"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
	"github.com/stretchr/testify/assert"
)

var (
	HashKey            = "some-secret-key"
	gaugeValue         = 123.45
	counterValue int64 = 12345
)

func TestService_SetHash(t *testing.T) {

	type args struct {
		m *metrics.Metrics
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test gauge",
			args: args{
				m: &metrics.Metrics{
					ID:    "some-gauge-metric",
					MType: metrics.GaugeMType,
					Value: &gaugeValue,
				},
			},
			wantErr: false,
		},
		{
			name: "test counter",
			args: args{
				m: &metrics.Metrics{
					ID:    "some-counter-metric",
					MType: metrics.CounterMType,
					Delta: &counterValue,
				},
			},
			wantErr: false,
		},
		{
			args: args{
				m: &metrics.Metrics{
					ID:    "some-invalid",
					MType: "none",
				},
			},
			wantErr: true,
		},
	}

	s := Service{&HashKey}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.SetHash(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("SetHash() error = %v, wantErr %v", err, tt.wantErr)
				assert.NotEmpty(t, tt.args.m.Hash, "calculated hash is empty")
			}
		})
	}
}

func ExampleService_SetHash() {
	s := Service{&HashKey}
	metric1 := &metrics.Metrics{
		ID:    "some-gauge-metric",
		MType: metrics.GaugeMType,
		Value: &gaugeValue,
	}
	s.SetHash(metric1)
	fmt.Println(metric1.Hash)

	metric2 := &metrics.Metrics{
		ID:    "some-counter-metric",
		MType: metrics.CounterMType,
		Delta: &counterValue,
	}
	s.SetHash(metric2)
	fmt.Println(metric2.Hash)

	// Output:
	// 4d526de041d1d955d86e090e94b1cb2709a029782773674139cafb9901ec44f4
	// 88ca53e61b82bc27dd3825e9318c239d81b54aa424d282a1a97bcb61ac447446
}
