package agent

import (
	"net/http"
	"runtime"
	"time"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

type MetricCollector struct {
	Storage     metrics.MetricStorage
	CollectChan <-chan time.Time
	PushChan    <-chan time.Time
	ServerHost  string
	ServerPort  int
	Client      http.Client
}

func MakeMetricCollector(
	storage metrics.MetricStorage,
	collectInterval time.Duration,
	pushInterval time.Duration,
	serverHost string,
	serverPort int,
) MetricCollector {
	return MetricCollector{
		Storage:     storage,
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

	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.Alloc,
		Type:  metrics.GaugeType,
		Value: float64(rtm.Alloc),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.BuckHashSys,
		Type:  metrics.GaugeType,
		Value: float64(rtm.BuckHashSys),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.Frees,
		Type:  metrics.GaugeType,
		Value: float64(rtm.Frees),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.GccpuFraction,
		Type:  metrics.GaugeType,
		Value: float64(rtm.GCCPUFraction),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.GcSys,
		Type:  metrics.GaugeType,
		Value: float64(rtm.GCSys),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.HeapAlloc,
		Type:  metrics.GaugeType,
		Value: float64(rtm.HeapAlloc),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.HeapIdle,
		Type:  metrics.GaugeType,
		Value: float64(rtm.HeapIdle),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.HeapInuse,
		Type:  metrics.GaugeType,
		Value: float64(rtm.HeapInuse),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.HeapObjects,
		Type:  metrics.GaugeType,
		Value: float64(rtm.HeapObjects),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.HeapReleased,
		Type:  metrics.GaugeType,
		Value: float64(rtm.HeapReleased),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.HeapSys,
		Type:  metrics.GaugeType,
		Value: float64(rtm.HeapSys),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.LastGC,
		Type:  metrics.GaugeType,
		Value: float64(rtm.LastGC),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.Lookups,
		Type:  metrics.GaugeType,
		Value: float64(rtm.Lookups),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.MCacheInuse,
		Type:  metrics.GaugeType,
		Value: float64(rtm.MCacheInuse),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.MCacheSys,
		Type:  metrics.GaugeType,
		Value: float64(rtm.MCacheSys),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.MSpanInuse,
		Type:  metrics.GaugeType,
		Value: float64(rtm.MSpanInuse),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.MSpanSys,
		Type:  metrics.GaugeType,
		Value: float64(rtm.MSpanSys),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.Mallocs,
		Type:  metrics.GaugeType,
		Value: float64(rtm.Mallocs),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.NextGC,
		Type:  metrics.GaugeType,
		Value: float64(rtm.NextGC),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.NumForcedGC,
		Type:  metrics.GaugeType,
		Value: float64(rtm.NumForcedGC),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.NumGC,
		Type:  metrics.GaugeType,
		Value: float64(rtm.NumGC),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.OtherSys,
		Type:  metrics.GaugeType,
		Value: float64(rtm.OtherSys),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.PauseTotalNs,
		Type:  metrics.GaugeType,
		Value: float64(rtm.PauseTotalNs),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.StackInuse,
		Type:  metrics.GaugeType,
		Value: float64(rtm.StackInuse),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.StackSys,
		Type:  metrics.GaugeType,
		Value: float64(rtm.StackSys),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.Sys,
		Type:  metrics.GaugeType,
		Value: float64(rtm.Sys),
	})
	m.Storage.SetMetric(metrics.AtomicMetric{
		Name:  metrics.TotalAlloc,
		Type:  metrics.GaugeType,
		Value: float64(rtm.TotalAlloc),
	})

	m.Storage.IncPollCount()
	m.Storage.RandomizeRandomValue()
}
