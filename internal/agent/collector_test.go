package agent

import (
	"testing"

	"go.uber.org/zap"

	"github.com/ksusonic/go-devops-mon/internal/storage"
)

var logger, _ = zap.NewDevelopment()

func BenchmarkMetricCollector_CollectStat(b *testing.B) {
	collector := MetricCollector{
		Storage: storage.NewAgentStorage(),
		Logger:  logger,
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		collector.CollectStat()
	}
}

func BenchmarkMetricCollector_CollectPsUtil(b *testing.B) {
	collector := MetricCollector{
		Storage: storage.NewAgentStorage(),
		Logger:  logger,
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		collector.CollectPsUtil()
	}
}
