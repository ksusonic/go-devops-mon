package agent

import (
	"net/http"
	"runtime"
	"time"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

type MetricCollector struct {
	Storage             metrics.AgentMetricStorage
	CollectChan         <-chan time.Time
	PushChan            <-chan time.Time
	ServerRequestScheme string
	ServerHost          string
	ServerPort          int
	Client              http.Client
}

func NewMetricCollector(
	storage metrics.AgentMetricStorage,
	collectInterval time.Duration,
	pushInterval time.Duration,
	serverRequestScheme string,
	serverHost string,
	serverPort int,
) *MetricCollector {
	return &MetricCollector{
		Storage:             storage,
		CollectChan:         time.NewTicker(collectInterval).C,
		PushChan:            time.NewTicker(pushInterval).C,
		ServerRequestScheme: serverRequestScheme,
		ServerHost:          serverHost,
		ServerPort:          serverPort,
		Client:              http.Client{},
	}
}

func (m MetricCollector) CollectStat() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	currentMetrics := []metrics.AtomicMetric{
		{
			Name:  "Alloc",
			Type:  metrics.GaugeType,
			Value: float64(rtm.Alloc),
		}, {
			Name:  "BuckHashSys",
			Type:  metrics.GaugeType,
			Value: float64(rtm.BuckHashSys),
		},
		{
			Name:  "Frees",
			Type:  metrics.GaugeType,
			Value: float64(rtm.Frees),
		},
		{
			Name:  "GCCPUFraction",
			Type:  metrics.GaugeType,
			Value: float64(rtm.GCCPUFraction),
		},
		{
			Name:  "GCSys",
			Type:  metrics.GaugeType,
			Value: float64(rtm.GCSys),
		},
		{
			Name:  "HeapAlloc",
			Type:  metrics.GaugeType,
			Value: float64(rtm.HeapAlloc),
		},
		{
			Name:  "HeapIdle",
			Type:  metrics.GaugeType,
			Value: float64(rtm.HeapIdle),
		},
		{
			Name:  "HeapInuse",
			Type:  metrics.GaugeType,
			Value: float64(rtm.HeapInuse),
		},
		{
			Name:  "HeapObjects",
			Type:  metrics.GaugeType,
			Value: float64(rtm.HeapObjects),
		},
		{
			Name:  "HeapReleased",
			Type:  metrics.GaugeType,
			Value: float64(rtm.HeapReleased),
		},
		{
			Name:  "HeapSys",
			Type:  metrics.GaugeType,
			Value: float64(rtm.HeapSys),
		},
		{
			Name:  "LastGC",
			Type:  metrics.GaugeType,
			Value: float64(rtm.LastGC),
		},
		{
			Name:  "Lookups",
			Type:  metrics.GaugeType,
			Value: float64(rtm.Lookups),
		},
		{
			Name:  "MCacheInuse",
			Type:  metrics.GaugeType,
			Value: float64(rtm.MCacheInuse),
		},
		{
			Name:  "MCacheSys",
			Type:  metrics.GaugeType,
			Value: float64(rtm.MCacheSys),
		},
		{
			Name:  "MSpanInuse",
			Type:  metrics.GaugeType,
			Value: float64(rtm.MSpanInuse),
		},
		{
			Name:  "MSpanSys",
			Type:  metrics.GaugeType,
			Value: float64(rtm.MSpanSys),
		},
		{
			Name:  "Mallocs",
			Type:  metrics.GaugeType,
			Value: float64(rtm.Mallocs),
		},
		{
			Name:  "NextGC",
			Type:  metrics.GaugeType,
			Value: float64(rtm.NextGC),
		},
		{
			Name:  "NumForcedGC",
			Type:  metrics.GaugeType,
			Value: float64(rtm.NumForcedGC),
		},
		{
			Name:  "NumGC",
			Type:  metrics.GaugeType,
			Value: float64(rtm.NumGC),
		},
		{
			Name:  "OtherSys",
			Type:  metrics.GaugeType,
			Value: float64(rtm.OtherSys),
		},
		{
			Name:  "PauseTotalNs",
			Type:  metrics.GaugeType,
			Value: float64(rtm.PauseTotalNs),
		},
		{
			Name:  "StackInuse",
			Type:  metrics.GaugeType,
			Value: float64(rtm.StackInuse),
		},
		{
			Name:  "StackSys",
			Type:  metrics.GaugeType,
			Value: float64(rtm.StackSys),
		},
		{
			Name:  "Sys",
			Type:  metrics.GaugeType,
			Value: float64(rtm.Sys),
		},
		{
			Name:  "TotalAlloc",
			Type:  metrics.GaugeType,
			Value: float64(rtm.TotalAlloc),
		},
	}
	m.Storage.AddMetrics(currentMetrics)

	// counters
	m.Storage.IncPollCount()
	m.Storage.RandomizeRandomValue()
}
