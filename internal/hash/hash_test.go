package hash

import (
	"fmt"
	"testing"

	metricspb "github.com/ksusonic/go-devops-mon/proto/metrics"
	"github.com/stretchr/testify/assert"
)

var (
	HashKey            = "some-secret-key"
	gaugeValue         = 123.45
	counterValue int64 = 12345
)

func TestService_SetHash(t *testing.T) {

	type args struct {
		m *metricspb.Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test gauge",
			args: args{
				m: &metricspb.Metric{
					ID:      "some-gauge-metric",
					Type:    metricspb.MetricType_gauge,
					Payload: &metricspb.Metric_Value{Value: gaugeValue},
				},
			},
			wantErr: false,
		},
		{
			name: "test counter",
			args: args{
				m: &metricspb.Metric{
					ID:      "some-counter-metric",
					Type:    metricspb.MetricType_counter,
					Payload: &metricspb.Metric_Delta{Delta: counterValue},
				},
			},
			wantErr: false,
		},
		{
			args: args{
				m: &metricspb.Metric{
					ID: "some-invalid",
				},
			},
			wantErr: true,
		},
	}

	s := Service{&HashKey}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.SetHashProto(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("SetHashProto() error = %v, wantErr %v", err, tt.wantErr)
				assert.NotEmpty(t, tt.args.m.Hash, "calculated hash is empty")
			}
		})
	}
}

func ExampleService_SetHash() {
	s := Service{&HashKey}
	metric1 := &metricspb.Metric{
		ID:      "some-gauge-metric",
		Type:    metricspb.MetricType_gauge,
		Payload: &metricspb.Metric_Value{Value: gaugeValue},
	}
	err := s.SetHashProto(metric1)
	if err != nil {
		panic(err)
	}
	fmt.Println(metric1.Hash)

	metric2 := &metricspb.Metric{
		ID:      "some-counter-metric",
		Type:    metricspb.MetricType_counter,
		Payload: &metricspb.Metric_Delta{Delta: counterValue},
	}
	err = s.SetHashProto(metric2)
	if err != nil {
		panic(err)
	}
	fmt.Println(metric2.Hash)

	// Output:
	// 4d526de041d1d955d86e090e94b1cb2709a029782773674139cafb9901ec44f4
	// 88ca53e61b82bc27dd3825e9318c239d81b54aa424d282a1a97bcb61ac447446
}
