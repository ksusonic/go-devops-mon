package main

import (
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/ksusonic/go-devops-mon/internal/agent"
	"github.com/ksusonic/go-devops-mon/internal/agent/pusher"
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
	if cfg.CryptoKeyPath != "" && err != nil {
		logger.Fatal("error creating encrypter", zap.Error(err))
	}

	u := url.URL{
		Scheme: cfg.AddressScheme,
		Host:   cfg.Address,
		Path:   "/update/",
	}
	ip, _ := agent.GetFirstIPOfMachine()
	var metricPusher agent.Pusher
	if cfg.UseGrpc {
		metricPusher = &pusher.GrpcPusher{
			PushURL: u.String(),
			Addr:    ip,
		}
	} else {
		metricPusher = &pusher.HTTPPusher{
			PushURL: u.String(),
			Addr:    ip,
		}
	}

	collector, err := agent.NewMetricCollector(
		cfg,
		logger,
		memStorage,
		hashService,
		encryptService,
		metricPusher,
	)
	if err != nil {
		logger.Fatal("error creating collector", zap.Error(err))
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

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
		case <-sigint:
			logger.Info("Caught interrupt signal")
			return
		}
	}
}
