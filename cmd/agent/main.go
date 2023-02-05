package main

import (
	"log"

	"github.com/ksusonic/go-devops-mon/internal/agent"
	"github.com/ksusonic/go-devops-mon/internal/hash"
	"github.com/ksusonic/go-devops-mon/internal/storage"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	cfg, err := agent.NewConfig()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	memStorage := storage.NewAgentStorage()
	hashService := hash.NewService(cfg.SecretKey)

	collector, err := agent.NewMetricCollector(cfg, memStorage, hashService)
	if err != nil {
		log.Fatalf("Error in metric collector: %v", err)
	}

	log.Println("Started agent!")
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
