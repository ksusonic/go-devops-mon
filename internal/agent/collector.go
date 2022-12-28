package agent

import (
	"net/http"
	"runtime"
	"time"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

type MetricCollector struct {
	Storage     CollectorStorage
	CollectChan <-chan time.Time
	PushChan    <-chan time.Time
	ServerHost  string
	ServerPort  int
	Client      http.Client
}

func MakeMetricCollector(
	collectInterval time.Duration,
	pushInterval time.Duration,
	serverHost string,
	serverPort int,
) MetricCollector {
	return MetricCollector{
		Storage:     NewCollectorStorage(),
		CollectChan: time.NewTicker(collectInterval).C,
		PushChan:    time.NewTicker(pushInterval).C,
		ServerHost:  serverHost,
		ServerPort:  serverPort,
		Client:      http.Client{},
	}
}

func (m MetricCollector) CollectStat() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	m.Storage.GaugeMetricStorage[metrics.Alloc] = float64(rtm.Alloc)
	m.Storage.GaugeMetricStorage[metrics.BuckHashSys] = float64(rtm.BuckHashSys)
	m.Storage.GaugeMetricStorage[metrics.Frees] = float64(rtm.Frees)
	m.Storage.GaugeMetricStorage[metrics.GccpuFraction] = rtm.GCCPUFraction
	m.Storage.GaugeMetricStorage[metrics.GcSys] = float64(rtm.GCSys)
	m.Storage.GaugeMetricStorage[metrics.HeapAlloc] = float64(rtm.HeapAlloc)
	m.Storage.GaugeMetricStorage[metrics.HeapIdle] = float64(rtm.HeapIdle)
	m.Storage.GaugeMetricStorage[metrics.HeapInuse] = float64(rtm.HeapInuse)
	m.Storage.GaugeMetricStorage[metrics.HeapObjects] = float64(rtm.HeapObjects)
	m.Storage.GaugeMetricStorage[metrics.HeapReleased] = float64(rtm.HeapReleased)
	m.Storage.GaugeMetricStorage[metrics.HeapSys] = float64(rtm.HeapSys)
	m.Storage.GaugeMetricStorage[metrics.LastGC] = float64(rtm.LastGC)
	m.Storage.GaugeMetricStorage[metrics.Lookups] = float64(rtm.Lookups)
	m.Storage.GaugeMetricStorage[metrics.MCacheInuse] = float64(rtm.MCacheInuse)
	m.Storage.GaugeMetricStorage[metrics.MCacheSys] = float64(rtm.MCacheSys)
	m.Storage.GaugeMetricStorage[metrics.MSpanInuse] = float64(rtm.MSpanInuse)
	m.Storage.GaugeMetricStorage[metrics.MSpanSys] = float64(rtm.MSpanSys)
	m.Storage.GaugeMetricStorage[metrics.Mallocs] = float64(rtm.Mallocs)
	m.Storage.GaugeMetricStorage[metrics.NextGC] = float64(rtm.NextGC)
	m.Storage.GaugeMetricStorage[metrics.NumForcedGC] = float64(rtm.NumForcedGC)
	m.Storage.GaugeMetricStorage[metrics.NumGC] = float64(rtm.NumGC)
	m.Storage.GaugeMetricStorage[metrics.OtherSys] = float64(rtm.OtherSys)
	m.Storage.GaugeMetricStorage[metrics.PauseTotalNs] = float64(rtm.PauseTotalNs)
	m.Storage.GaugeMetricStorage[metrics.StackInuse] = float64(rtm.StackInuse)
	m.Storage.GaugeMetricStorage[metrics.StackSys] = float64(rtm.StackSys)
	m.Storage.GaugeMetricStorage[metrics.Sys] = float64(rtm.Sys)
	m.Storage.GaugeMetricStorage[metrics.TotalAlloc] = float64(rtm.TotalAlloc)

	m.Storage.CounterMetricStorage.IncPollCountValue()
	m.Storage.CounterMetricStorage.RandomizeRandomValue()
}
