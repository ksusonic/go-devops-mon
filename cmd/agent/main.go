package main

import (
	"log"

	"github.com/ksusonic/go-devops-mon/internal/agent"
	"github.com/ksusonic/go-devops-mon/internal/storage"
)

func main() {
	cfg := agent.NewConfig()
	memStorage := storage.NewMemStorage()
	collector := agent.NewMetricCollector(cfg, memStorage)
	for {
		select {
		case <-collector.CollectChan:
			log.Println("collected metrics")
			collector.CollectStat()
		case <-collector.PushChan:
			log.Println("pushing metrics...")
			collector.PushMetrics()
		}
	}
}
