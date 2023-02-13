package main

import (
	"github.com/ksusonic/go-devops-mon/internal/agent"
	"github.com/ksusonic/go-devops-mon/internal/hash"
	"github.com/ksusonic/go-devops-mon/internal/storage"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()

	cfg, err := agent.NewConfig()
	if err != nil {
		logger.Fatal("error reading config", zap.Error(err))
	}
	memStorage := storage.NewAgentStorage()
	hashService := hash.NewService(cfg.SecretKey)

	collector, err := agent.NewMetricCollector(cfg, logger, memStorage, hashService)
	if err != nil {
		logger.Fatal("error in metric collector", zap.Error(err))
	}

	logger.Info("started agent!")
	for {
		select {
		case <-collector.CollectChan:
			logger.Debug("collected metrics")
			collector.CollectStat()
		case <-collector.PushChan:
			logger.Debug("pushing metrics...")
			err := collector.PushMetrics()
			if err != nil {
				logger.Error("error pushing metrics", zap.Error(err))
			}
		}
	}
}
