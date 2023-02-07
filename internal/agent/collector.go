package agent

import (
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

var pollCounterDelta int64 = 1

type MetricCollector struct {
	Storage     metrics.AgentMetricStorage
	CollectChan <-chan time.Time
	PushChan    <-chan time.Time
	PushURL     string
	Client      http.Client
	HashService metrics.HashService
}

func NewMetricCollector(
	cfg *Config,
	storage metrics.AgentMetricStorage,
	service metrics.HashService,
) (*MetricCollector, error) {
	u := url.URL{
		Scheme: cfg.AddressScheme,
		Host:   cfg.Address,
		Path:   "/update/",
	}
	pollInterval, err := time.ParseDuration(cfg.PollInterval)
	if err != nil {
		return nil, err
	}
	reportInterval, err := time.ParseDuration(cfg.ReportInterval)
	if err != nil {
		return nil, err
	}

	return &MetricCollector{
		Storage:     storage,
		CollectChan: time.NewTicker(pollInterval).C,
		PushChan:    time.NewTicker(reportInterval).C,
		PushURL:     u.String(),
		Client:      http.Client{},
		HashService: service,
	}, nil
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
		{
			Name:  "RandomValue",
			Value: rand.Float64(),
		},
	}

	// gauge
	for i := range currentGaugeMetrics {
		err := m.Storage.SetMetric(metrics.Metrics{
			ID:    currentGaugeMetrics[i].Name,
			MType: metrics.GaugeMType,
			Value: &currentGaugeMetrics[i].Value,
		})
		if err != nil {
			log.Println("Could not set hash:", err)
		}
	}

	// counters
	err := m.Storage.SetMetric(metrics.Metrics{
		ID:    "PollCount",
		MType: metrics.CounterMType,
		Delta: &pollCounterDelta,
	})
	if err != nil {
		log.Println("Could not set hash:", err)
	}
}
