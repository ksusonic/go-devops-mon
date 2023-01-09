package agent

import (
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

type MetricCollector struct {
	Storage     metrics.AgentMetricStorage
	CollectChan <-chan time.Time
	PushChan    <-chan time.Time
	PushUrl     string
	Client      http.Client
}

func NewMetricCollector(
	storage metrics.AgentMetricStorage,
	collectInterval time.Duration,
	pushInterval time.Duration,
	serverRequestScheme string,
	serverAddress string,
) *MetricCollector {
	u := url.URL{
		Scheme: serverRequestScheme,
		Host:   serverAddress,
		Path:   "/update/",
	}
	return &MetricCollector{
		Storage:     storage,
		CollectChan: time.NewTicker(collectInterval).C,
		PushChan:    time.NewTicker(pushInterval).C,
		PushUrl:     u.String(),
		Client:      http.Client{},
	}
}

func (m MetricCollector) CollectStat() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	currentGaugeMetrics := []struct {
		Name  string
		Value float64
	}{
		{
			Name:  "Alloc",
			Value: float64(rtm.Alloc),
		}, {
			Name:  "BuckHashSys",
			Value: float64(rtm.BuckHashSys),
		},
		{
			Name:  "Frees",
			Value: float64(rtm.Frees),
		},
		{
			Name:  "GCCPUFraction",
			Value: rtm.GCCPUFraction,
		},
		{
			Name:  "GCSys",
			Value: float64(rtm.GCSys),
		},
		{
			Name:  "HeapAlloc",
			Value: float64(rtm.HeapAlloc),
		},
		{
			Name:  "HeapIdle",
			Value: float64(rtm.HeapIdle),
		},
		{
			Name:  "HeapInuse",
			Value: float64(rtm.HeapInuse),
		},
		{
			Name:  "HeapObjects",
			Value: float64(rtm.HeapObjects),
		},
		{
			Name:  "HeapReleased",
			Value: float64(rtm.HeapReleased),
		},
		{
			Name:  "HeapSys",
			Value: float64(rtm.HeapSys),
		},
		{
			Name:  "LastGC",
			Value: float64(rtm.LastGC),
		},
		{
			Name:  "Lookups",
			Value: float64(rtm.Lookups),
		},
		{
			Name:  "MCacheInuse",
			Value: float64(rtm.MCacheInuse),
		},
		{
			Name:  "MCacheSys",
			Value: float64(rtm.MCacheSys),
		},
		{
			Name:  "MSpanInuse",
			Value: float64(rtm.MSpanInuse),
		},
		{
			Name:  "MSpanSys",
			Value: float64(rtm.MSpanSys),
		},
		{
			Name:  "Mallocs",
			Value: float64(rtm.Mallocs),
		},
		{
			Name:  "NextGC",
			Value: float64(rtm.NextGC),
		},
		{
			Name:  "NumForcedGC",
			Value: float64(rtm.NumForcedGC),
		},
		{
			Name:  "NumGC",
			Value: float64(rtm.NumGC),
		},
		{
			Name:  "OtherSys",
			Value: float64(rtm.OtherSys),
		},
		{
			Name:  "PauseTotalNs",
			Value: float64(rtm.PauseTotalNs),
		},
		{
			Name:  "StackInuse",
			Value: float64(rtm.StackInuse),
		},
		{
			Name:  "StackSys",
			Value: float64(rtm.StackSys),
		},
		{
			Name:  "Sys",
			Value: float64(rtm.Sys),
		},
		{
			Name:  "TotalAlloc",
			Value: float64(rtm.TotalAlloc),
		},
	}

	for _, metric := range currentGaugeMetrics {
		m.Storage.SetMetric(metrics.Metrics{
			ID:    metric.Name,
			MType: metrics.GaugeMType,
			Delta: nil,
			Value: &metric.Value,
		})
	}

	// counters
	m.Storage.IncPollCount()
	m.Storage.RandomizeRandomValue()
}
