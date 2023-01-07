package main

import (
	"log"
	"time"

	"github.com/ksusonic/go-devops-mon/internal/agent"
	"github.com/ksusonic/go-devops-mon/internal/storage"
)

const (
	CollectInterval     = time.Second * 2
	PushInterval        = time.Second * 10
	ServerRequestMethod = "http"
	ServerHost          = "localhost"
	ServerPort          = 8080
)

func main() {
	memStorage := storage.NewMemStorage()
	collector := agent.NewMetricCollector(memStorage, CollectInterval, PushInterval, ServerRequestMethod, ServerHost, ServerPort)
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
