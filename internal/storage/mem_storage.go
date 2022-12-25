package storage

import (
	"fmt"
	"github.com/ksusonic/go-devops-mon/internal/metrics"
	"math/rand"
)

type GaugeStorage map[string][]float64
type CounterStorage map[string][]int64

type MemStorage struct {
	GaugeStorage
	CounterStorage
}

func (m MemStorage) GetAllCurrentGaugeMetrics() map[string]float64 {
	result := make(map[string]float64)
	for name, value := range m.GaugeStorage {
		if len(value) > 0 {
			result[name] = value[len(value)-1]
		} else {
			result[name] = -1
		}
	}
	return result
}

func (m MemStorage) GetAllCurrentCounterMetrics() map[string]int64 {
	result := make(map[string]int64)
	for name, value := range m.CounterStorage {
		if len(value) > 0 {
			result[name] = value[len(value)-1]
		} else {
			result[name] = -1
		}
	}
	return result
}

func (s *GaugeStorage) GetCurrentGaugeMetric(name string) (float64, *error) {
	val, ok := (*s)[name]
	if !ok {
		err := fmt.Errorf("no such metric")
		return 0, &err
	}
	if len(val) > 0 {
		current := val[len(val)-1]
		return current, nil
	} else {
		err := fmt.Errorf("no values in metric %s", name)
		return 0, &err
	}
}

func (s *CounterStorage) GetCurrentCounterMetric(name string) (int64, *error) {
	val, ok := (*s)[name]
	if !ok {
		err := fmt.Errorf("no such metric")
		return 0, &err
	}
	if len(val) > 0 {
		current := val[len(val)-1]
		return current, nil
	} else {
		err := fmt.Errorf("no values in metric %s", name)
		return 0, &err
	}
}

func (s *GaugeStorage) AddGaugeValue(name string, value float64) {
	(*s)[name] = append((*s)[name], value)
}
func (s *CounterStorage) AddToCounterValue(name string, value int64) {
	if len((*s)[name]) > 0 {
		lastValue := (*s)[name][len((*s)[name])-1]
		(*s)[name] = append((*s)[name], lastValue+value)
	} else {
		(*s)[name] = []int64{value}
	}
}

func (s *GaugeStorage) Alloc() []float64 {
	return (*s)[metrics.Alloc]
}
func (s *GaugeStorage) AddAlloc(v float64) {
	(*s)[metrics.Alloc] = append((*s)[metrics.Alloc], v)
}
func (s *GaugeStorage) BuckHashSys() []float64 {
	return (*s)[metrics.BuckHashSys]
}
func (s *GaugeStorage) AddBuckHashSys(v float64) {
	(*s)[metrics.BuckHashSys] = append((*s)[metrics.BuckHashSys], v)
}
func (s *GaugeStorage) Frees() []float64 {
	return (*s)[metrics.Frees]
}
func (s *GaugeStorage) AddFrees(v float64) {
	(*s)[metrics.Frees] = append((*s)[metrics.Frees], v)
}
func (s *GaugeStorage) GCCPUFraction() []float64 {
	return (*s)[metrics.GccpuFraction]
}
func (s *GaugeStorage) AddGCCPUFraction(v float64) {
	(*s)[metrics.GccpuFraction] = append((*s)[metrics.GccpuFraction], v)
}
func (s *GaugeStorage) GCSys() []float64 {
	return (*s)[metrics.GcSys]
}
func (s *GaugeStorage) AddGCSys(v float64) {
	(*s)[metrics.GcSys] = append((*s)[metrics.GcSys], v)
}
func (s *GaugeStorage) HeapAlloc() []float64 {
	return (*s)[metrics.HeapAlloc]
}
func (s *GaugeStorage) AddHeapAlloc(v float64) {
	(*s)[metrics.HeapAlloc] = append((*s)[metrics.HeapAlloc], v)
}
func (s *GaugeStorage) HeapIdle() []float64 {
	return (*s)[metrics.HeapIdle]
}
func (s *GaugeStorage) AddHeapIdle(v float64) {
	(*s)[metrics.HeapIdle] = append((*s)[metrics.HeapIdle], v)
}
func (s *GaugeStorage) HeapInuse() []float64 {
	return (*s)[metrics.HeapInuse]
}
func (s *GaugeStorage) AddHeapInuse(v float64) {
	(*s)[metrics.HeapInuse] = append((*s)[metrics.HeapInuse], v)
}
func (s *GaugeStorage) HeapObjects() []float64 {
	return (*s)[metrics.HeapObjects]
}
func (s *GaugeStorage) AddHeapObjects(v float64) {
	(*s)[metrics.HeapObjects] = append((*s)[metrics.HeapObjects], v)
}
func (s *GaugeStorage) HeapReleased() []float64 {
	return (*s)[metrics.HeapReleased]
}
func (s *GaugeStorage) AddHeapReleased(v float64) {
	(*s)[metrics.HeapReleased] = append((*s)[metrics.HeapReleased], v)
}
func (s *GaugeStorage) HeapSys() []float64 {
	return (*s)[metrics.HeapSys]
}
func (s *GaugeStorage) AddHeapSys(v float64) {
	(*s)[metrics.HeapSys] = append((*s)[metrics.HeapSys], v)
}
func (s *GaugeStorage) LastGC() []float64 {
	return (*s)[metrics.LastGC]
}
func (s *GaugeStorage) AddLastGC(v float64) {
	(*s)[metrics.LastGC] = append((*s)[metrics.LastGC], v)
}
func (s *GaugeStorage) Lookups() []float64 {
	return (*s)[metrics.Lookups]
}
func (s *GaugeStorage) AddLookups(v float64) {
	(*s)[metrics.Lookups] = append((*s)[metrics.Lookups], v)
}
func (s *GaugeStorage) MCacheInuse() []float64 {
	return (*s)[metrics.MCacheInuse]
}
func (s *GaugeStorage) AddMCacheInuse(v float64) {
	(*s)[metrics.MCacheInuse] = append((*s)[metrics.MCacheInuse], v)
}
func (s *GaugeStorage) MCacheSys() []float64 {
	return (*s)[metrics.MCacheSys]
}
func (s *GaugeStorage) AddMCacheSys(v float64) {
	(*s)[metrics.MCacheSys] = append((*s)[metrics.MCacheSys], v)
}
func (s *GaugeStorage) MSpanInuse() []float64 {
	return (*s)[metrics.MSpanInuse]
}
func (s *GaugeStorage) AddMSpanInuse(v float64) {
	(*s)[metrics.MSpanInuse] = append((*s)[metrics.MSpanInuse], v)
}
func (s *GaugeStorage) MSpanSys() []float64 {
	return (*s)[metrics.MSpanSys]
}
func (s *GaugeStorage) AddMSpanSys(v float64) {
	(*s)[metrics.MSpanSys] = append((*s)[metrics.MSpanSys], v)
}
func (s *GaugeStorage) Mallocs() []float64 {
	return (*s)[metrics.Mallocs]
}
func (s *GaugeStorage) AddMallocs(v float64) {
	(*s)[metrics.Mallocs] = append((*s)[metrics.Mallocs], v)
}
func (s *GaugeStorage) NextGC() []float64 {
	return (*s)[metrics.NextGC]
}
func (s *GaugeStorage) AddNextGC(v float64) {
	(*s)[metrics.NextGC] = append((*s)[metrics.NextGC], v)
}
func (s *GaugeStorage) NumForcedGC() []float64 {
	return (*s)[metrics.NumForcedGC]
}
func (s *GaugeStorage) AddNumForcedGC(v float64) {
	(*s)[metrics.NumForcedGC] = append((*s)[metrics.NumForcedGC], v)
}
func (s *GaugeStorage) NumGC() []float64 {
	return (*s)[metrics.NumGC]
}
func (s *GaugeStorage) AddNumGC(v float64) {
	(*s)[metrics.NumGC] = append((*s)[metrics.NumGC], v)
}
func (s *GaugeStorage) OtherSys() []float64 {
	return (*s)[metrics.OtherSys]
}
func (s *GaugeStorage) AddOtherSys(v float64) {
	(*s)[metrics.OtherSys] = append((*s)[metrics.OtherSys], v)
}
func (s *GaugeStorage) PauAddotalNs() []float64 {
	return (*s)[metrics.PauseTotalNs]
}
func (s *GaugeStorage) AddPauAddotalNs(v float64) {
	(*s)[metrics.PauseTotalNs] = append((*s)[metrics.PauseTotalNs], v)
}
func (s *GaugeStorage) StackInuse() []float64 {
	return (*s)[metrics.StackInuse]
}
func (s *GaugeStorage) AddStackInuse(v float64) {
	(*s)[metrics.StackInuse] = append((*s)[metrics.StackInuse], v)
}
func (s *GaugeStorage) StackSys() []float64 {
	return (*s)[metrics.StackSys]
}
func (s *GaugeStorage) AddStackSys(v float64) {
	(*s)[metrics.StackSys] = append((*s)[metrics.StackSys], v)
}
func (s *GaugeStorage) Sys() []float64 {
	return (*s)[metrics.Sys]
}
func (s *GaugeStorage) AddSys(v float64) {
	(*s)[metrics.Sys] = append((*s)[metrics.Sys], v)
}
func (s *GaugeStorage) TotalAlloc() []float64 {
	return (*s)[metrics.TotalAlloc]
}
func (s *GaugeStorage) AddTotalAlloc(v float64) {
	(*s)[metrics.TotalAlloc] = append((*s)[metrics.TotalAlloc], v)
}

func (s *CounterStorage) PollCount() []int64 {
	return (*s)[metrics.PollCount]
}
func (s *CounterStorage) IncPollCount() {
	if len((*s)[metrics.PollCount]) == 0 {
		(*s)[metrics.PollCount] = append((*s)[metrics.PollCount], 1)
	} else {
		var lastElement = (*s)[metrics.PollCount][len((*s)[metrics.PollCount])-1]
		(*s)[metrics.PollCount] = append((*s)[metrics.PollCount], lastElement+1)
	}
}
func (s *CounterStorage) RandomValue() []int64 {
	return (*s)[metrics.RandomValue]
}
func (s *CounterStorage) RandomizeRandomValue() {
	(*s)[metrics.RandomValue] = append((*s)[metrics.RandomValue], rand.Int63())
}
