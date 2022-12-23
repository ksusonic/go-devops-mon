package agent

import (
	"github.com/ksusonic/go-devops-mon/internal/metrics"
	"runtime"
	"time"
)

type MetricCollector struct {
	GaugeMetrics   metrics.GaugeStorage
	CounterMetrics metrics.CounterStorage
	CollectChan    <-chan time.Time
	PushChan       <-chan time.Time
	ServerHost     string
	ServerPort     int
}

func MakeMetricCollector(
	collectInterval time.Duration,
	pushInterval time.Duration,
	serverHost string,
	serverPort int,
) MetricCollector {
	return MetricCollector{
		GaugeMetrics:   metrics.GaugeStorage{},
		CounterMetrics: metrics.CounterStorage{},
		CollectChan:    time.NewTicker(collectInterval).C,
		PushChan:       time.NewTicker(pushInterval).C,
		ServerHost:     serverHost,
		ServerPort:     serverPort,
	}
}

func (m MetricCollector) CollectStat() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	m.GaugeMetrics.AddAlloc(float64(rtm.Alloc))
	m.GaugeMetrics.AddBuckHashSys(float64(rtm.BuckHashSys))
	m.GaugeMetrics.AddFrees(float64(rtm.Frees))
	m.GaugeMetrics.AddGCCPUFraction(rtm.GCCPUFraction)
	m.GaugeMetrics.AddGCSys(float64(rtm.GCSys))
	m.GaugeMetrics.AddHeapAlloc(float64(rtm.HeapAlloc))
	m.GaugeMetrics.AddHeapIdle(float64(rtm.HeapIdle))
	m.GaugeMetrics.AddHeapInuse(float64(rtm.HeapInuse))
	m.GaugeMetrics.AddHeapObjects(float64(rtm.HeapObjects))
	m.GaugeMetrics.AddHeapReleased(float64(rtm.HeapReleased))
	m.GaugeMetrics.AddHeapSys(float64(rtm.HeapSys))
	m.GaugeMetrics.AddLastGC(float64(rtm.LastGC))
	m.GaugeMetrics.AddLookups(float64(rtm.Lookups))
	m.GaugeMetrics.AddMCacheInuse(float64(rtm.MCacheInuse))
	m.GaugeMetrics.AddMCacheSys(float64(rtm.MCacheSys))
	m.GaugeMetrics.AddMSpanInuse(float64(rtm.MSpanInuse))
	m.GaugeMetrics.AddMSpanSys(float64(rtm.MSpanSys))
	m.GaugeMetrics.AddMallocs(float64(rtm.Mallocs))
	m.GaugeMetrics.AddNextGC(float64(rtm.NextGC))
	m.GaugeMetrics.AddNumForcedGC(float64(rtm.NumForcedGC))
	m.GaugeMetrics.AddNumGC(float64(rtm.NumGC))
	m.GaugeMetrics.AddOtherSys(float64(rtm.OtherSys))
	m.GaugeMetrics.AddPauAddotalNs(float64(rtm.PauseTotalNs))
	m.GaugeMetrics.AddStackInuse(float64(rtm.StackInuse))
	m.GaugeMetrics.AddStackSys(float64(rtm.StackSys))
	m.GaugeMetrics.AddSys(float64(rtm.Sys))
	m.GaugeMetrics.AddTotalAlloc(float64(rtm.TotalAlloc))

	m.CounterMetrics.IncPollCount()
	m.CounterMetrics.RandomizeRandomValue()
}
