package main

import (
	"github.com/ksusonic/go-devops-mon/internal/agent"
	"github.com/ksusonic/go-devops-mon/internal/crypt"
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
	encryptService, err := crypt.NewEncrypter(cfg.CryptoKeyPath)
	if err != nil {
		logger.Fatal("error creating encrypter", zap.Error(err))
	}

	collector, err := agent.NewMetricCollector(
		cfg,
		logger,
		memStorage,
		hashService,
		encryptService,
	)
	if err != nil {
		logger.Fatal("error creating collector", zap.Error(err))
	}

	logger.Info("started agent!")
	for {
		select {
		case <-collector.CollectChan:
			logger.Debug("started metrics collector")
			go collector.CollectStat()
			go collector.CollectPsUtil()
		case <-collector.PushChan:
			logger.Debug("started pushing metrics")
			go collector.PushMetrics()
		}
	}
}
