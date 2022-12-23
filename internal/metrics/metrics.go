package metrics

import "math/rand"

type GaugeStorage map[string][]float64
type CounterStorage map[string][]int64

func HasMetricName(name string) bool {
	for _, metric := range MetricTypeToName[GaugeName] {
		if metric == name {
			return true
		}
	}
	for _, metric := range MetricTypeToName[CounterName] {
		if metric == name {
			return true
		}
	}
	return false
}

func (s *GaugeStorage) Alloc() []float64 {
	return (*s)[Alloc]
}
func (s *GaugeStorage) AddAlloc(v float64) {
	(*s)[Alloc] = append((*s)[Alloc], v)
}
func (s *GaugeStorage) BuckHashSys() []float64 {
	return (*s)[BuckHashSys]
}
func (s *GaugeStorage) AddBuckHashSys(v float64) {
	(*s)[BuckHashSys] = append((*s)[BuckHashSys], v)
}
func (s *GaugeStorage) Frees() []float64 {
	return (*s)[Frees]
}
func (s *GaugeStorage) AddFrees(v float64) {
	(*s)[Frees] = append((*s)[Frees], v)
}
func (s *GaugeStorage) GCCPUFraction() []float64 {
	return (*s)[GccpuFraction]
}
func (s *GaugeStorage) AddGCCPUFraction(v float64) {
	(*s)[GccpuFraction] = append((*s)[GccpuFraction], v)
}
func (s *GaugeStorage) GCSys() []float64 {
	return (*s)[GcSys]
}
func (s *GaugeStorage) AddGCSys(v float64) {
	(*s)[GcSys] = append((*s)[GcSys], v)
}
func (s *GaugeStorage) HeapAlloc() []float64 {
	return (*s)[HeapAlloc]
}
func (s *GaugeStorage) AddHeapAlloc(v float64) {
	(*s)[HeapAlloc] = append((*s)[HeapAlloc], v)
}
func (s *GaugeStorage) HeapIdle() []float64 {
	return (*s)[HeapIdle]
}
func (s *GaugeStorage) AddHeapIdle(v float64) {
	(*s)[HeapIdle] = append((*s)[HeapIdle], v)
}
func (s *GaugeStorage) HeapInuse() []float64 {
	return (*s)[HeapInuse]
}
func (s *GaugeStorage) AddHeapInuse(v float64) {
	(*s)[HeapInuse] = append((*s)[HeapInuse], v)
}
func (s *GaugeStorage) HeapObjects() []float64 {
	return (*s)[HeapObjects]
}
func (s *GaugeStorage) AddHeapObjects(v float64) {
	(*s)[HeapObjects] = append((*s)[HeapObjects], v)
}
func (s *GaugeStorage) HeapReleased() []float64 {
	return (*s)[HeapReleased]
}
func (s *GaugeStorage) AddHeapReleased(v float64) {
	(*s)[HeapReleased] = append((*s)[HeapReleased], v)
}
func (s *GaugeStorage) HeapSys() []float64 {
	return (*s)[HeapSys]
}
func (s *GaugeStorage) AddHeapSys(v float64) {
	(*s)[HeapSys] = append((*s)[HeapSys], v)
}
func (s *GaugeStorage) LastGC() []float64 {
	return (*s)[LastGC]
}
func (s *GaugeStorage) AddLastGC(v float64) {
	(*s)[LastGC] = append((*s)[LastGC], v)
}
func (s *GaugeStorage) Lookups() []float64 {
	return (*s)[Lookups]
}
func (s *GaugeStorage) AddLookups(v float64) {
	(*s)[Lookups] = append((*s)[Lookups], v)
}
func (s *GaugeStorage) MCacheInuse() []float64 {
	return (*s)[MCacheInuse]
}
func (s *GaugeStorage) AddMCacheInuse(v float64) {
	(*s)[MCacheInuse] = append((*s)[MCacheInuse], v)
}
func (s *GaugeStorage) MCacheSys() []float64 {
	return (*s)[MCacheSys]
}
func (s *GaugeStorage) AddMCacheSys(v float64) {
	(*s)[MCacheSys] = append((*s)[MCacheSys], v)
}
func (s *GaugeStorage) MSpanInuse() []float64 {
	return (*s)[MSpanInuse]
}
func (s *GaugeStorage) AddMSpanInuse(v float64) {
	(*s)[MSpanInuse] = append((*s)[MSpanInuse], v)
}
func (s *GaugeStorage) MSpanSys() []float64 {
	return (*s)[MSpanSys]
}
func (s *GaugeStorage) AddMSpanSys(v float64) {
	(*s)[MSpanSys] = append((*s)[MSpanSys], v)
}
func (s *GaugeStorage) Mallocs() []float64 {
	return (*s)[Mallocs]
}
func (s *GaugeStorage) AddMallocs(v float64) {
	(*s)[Mallocs] = append((*s)[Mallocs], v)
}
func (s *GaugeStorage) NextGC() []float64 {
	return (*s)[NextGC]
}
func (s *GaugeStorage) AddNextGC(v float64) {
	(*s)[NextGC] = append((*s)[NextGC], v)
}
func (s *GaugeStorage) NumForcedGC() []float64 {
	return (*s)[NumForcedGC]
}
func (s *GaugeStorage) AddNumForcedGC(v float64) {
	(*s)[NumForcedGC] = append((*s)[NumForcedGC], v)
}
func (s *GaugeStorage) NumGC() []float64 {
	return (*s)[NumGC]
}
func (s *GaugeStorage) AddNumGC(v float64) {
	(*s)[NumGC] = append((*s)[NumGC], v)
}
func (s *GaugeStorage) OtherSys() []float64 {
	return (*s)[OtherSys]
}
func (s *GaugeStorage) AddOtherSys(v float64) {
	(*s)[OtherSys] = append((*s)[OtherSys], v)
}
func (s *GaugeStorage) PauAddotalNs() []float64 {
	return (*s)[PauseTotalNs]
}
func (s *GaugeStorage) AddPauAddotalNs(v float64) {
	(*s)[PauseTotalNs] = append((*s)[PauseTotalNs], v)
}
func (s *GaugeStorage) StackInuse() []float64 {
	return (*s)[StackInuse]
}
func (s *GaugeStorage) AddStackInuse(v float64) {
	(*s)[StackInuse] = append((*s)[StackInuse], v)
}
func (s *GaugeStorage) StackSys() []float64 {
	return (*s)[StackSys]
}
func (s *GaugeStorage) AddStackSys(v float64) {
	(*s)[StackSys] = append((*s)[StackSys], v)
}
func (s *GaugeStorage) Sys() []float64 {
	return (*s)[Sys]
}
func (s *GaugeStorage) AddSys(v float64) {
	(*s)[Sys] = append((*s)[Sys], v)
}
func (s *GaugeStorage) TotalAlloc() []float64 {
	return (*s)[TotalAlloc]
}
func (s *GaugeStorage) AddTotalAlloc(v float64) {
	(*s)[TotalAlloc] = append((*s)[TotalAlloc], v)
}

func (s *CounterStorage) PollCount() []int64 {
	return (*s)[PollCount]
}
func (s *CounterStorage) IncPollCount() {
	if len((*s)[PollCount]) == 0 {
		(*s)[PollCount] = append((*s)[PollCount], 1)
	} else {
		var lastElement = (*s)[PollCount][len((*s)[PollCount])-1]
		(*s)[PollCount] = append((*s)[PollCount], lastElement+1)
	}
}
func (s *CounterStorage) RandomValue() []int64 {
	return (*s)[RandomValue]
}
func (s *CounterStorage) RandomizeRandomValue() {
	(*s)[RandomValue] = append((*s)[RandomValue], rand.Int63())
}
