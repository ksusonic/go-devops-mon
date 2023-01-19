package main

import (
	"log"

	"github.com/ksusonic/go-devops-mon/internal/agent"
	"github.com/ksusonic/go-devops-mon/internal/storage"
)

func main() {
	cfg, err := agent.NewConfig()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	memStorage := storage.NewMemStorage()
	collector, err := agent.NewMetricCollector(cfg, memStorage)
	if err != nil {
		log.Fatalf("Error in metric collector: %v", err)
	}

	for {
		select {
		case <-collector.CollectChan:
			log.Println("collected metrics")
			collector.CollectStat()
		case <-collector.PushChan:
			log.Println("pushing metrics...")
			err := collector.PushMetrics()
			if err != nil {
				log.Println(err)
			}
		}
	}
}
