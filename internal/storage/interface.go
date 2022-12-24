package storage

type Storage interface {
	AddGaugeValue(name string, value float64)
	AddCounterValue(name string, value int64)

	Alloc() []float64
	AddAlloc(v float64)

	BuckHashSys() []float64
	AddBuckHashSys(v float64)

	Frees() []float64
	AddFrees(v float64)

	GCCPUFraction() []float64
	AddGCCPUFraction(v float64)

	GCSys() []float64
	AddGCSys(v float64)

	HeapAlloc() []float64
	AddHeapAlloc(v float64)

	HeapIdle() []float64
	AddHeapIdle(v float64)

	HeapInuse() []float64
	AddHeapInuse(v float64)

	HeapObjects() []float64
	AddHeapObjects(v float64)

	HeapReleased() []float64
	AddHeapReleased(v float64)

	HeapSys() []float64
	AddHeapSys(v float64)

	LastGC() []float64
	AddLastGC(v float64)

	Lookups() []float64
	AddLookups(v float64)

	MCacheInuse() []float64
	AddMCacheInuse(v float64)

	MCacheSys() []float64
	AddMCacheSys(v float64)

	MSpanInuse() []float64
	AddMSpanInuse(v float64)

	MSpanSys() []float64
	AddMSpanSys(v float64)

	Mallocs() []float64
	AddMallocs(v float64)

	NextGC() []float64
	AddNextGC(v float64)

	NumForcedGC() []float64
	AddNumForcedGC(v float64)

	NumGC() []float64
	AddNumGC(v float64)

	OtherSys() []float64
	AddOtherSys(v float64)

	PauAddotalNs() []float64
	AddPauAddotalNs(v float64)

	StackInuse() []float64
	AddStackInuse(v float64)

	StackSys() []float64
	AddStackSys(v float64)

	Sys() []float64
	AddSys(v float64)

	TotalAlloc() []float64
	AddTotalAlloc(v float64)

	PollCount() []int64
	IncPollCount()

	RandomValue() []int64
	RandomizeRandomValue()
}
