package metrics

const (
	GaugeName   = "gauge"
	CounterName = "counter"
)

func MetricExists(metricName string) bool {
	for _, name := range MetricTypeToName[GaugeName] {
		if metricName == name {
			return true
		}
	}
	for _, name := range MetricTypeToName[CounterName] {
		if metricName == name {
			return true
		}
	}
	return false
}

const (
	Alloc         = "Alloc"
	BuckHashSys   = "BuckHashSys"
	Frees         = "Frees"
	GccpuFraction = "GCCPUFraction"
	GcSys         = "GCSys"
	HeapAlloc     = "HeapAlloc"
	HeapIdle      = "HeapIdle"
	HeapInuse     = "HeapInuse"
	HeapObjects   = "HeapObjects"
	HeapReleased  = "HeapReleased"
	HeapSys       = "HeapSys"
	LastGC        = "LastGC"
	Lookups       = "Lookups"
	MCacheInuse   = "MCacheInuse"
	MCacheSys     = "MCacheSys"
	MSpanInuse    = "MSpanInuse"
	MSpanSys      = "MSpanSys"
	Mallocs       = "Mallocs"
	NextGC        = "NextGC"
	NumForcedGC   = "NumForcedGC"
	NumGC         = "NumGC"
	OtherSys      = "OtherSys"
	PauseTotalNs  = "PauseTotalNs"
	StackInuse    = "StackInuse"
	StackSys      = "StackSys"
	Sys           = "Sys"
	TotalAlloc    = "TotalAlloc"

	PollCount   = "PollCount"
	RandomValue = "RandomValue"
)

var MetricTypeToName = map[string][]string{
	GaugeName: {
		Alloc,
		BuckHashSys,
		Frees,
		GccpuFraction,
		GcSys,
		HeapAlloc,
		HeapIdle,
		HeapInuse,
		HeapObjects,
		HeapReleased,
		HeapSys,
		LastGC,
		Lookups,
		MCacheInuse,
		MCacheSys,
		MSpanInuse,
		MSpanSys,
		Mallocs,
		NextGC,
		NumForcedGC,
		NumGC,
		OtherSys,
		PauseTotalNs,
		StackInuse,
		StackSys,
		Sys,
		TotalAlloc,
	},
	CounterName: {
		PollCount,
		RandomValue,
	},
}
