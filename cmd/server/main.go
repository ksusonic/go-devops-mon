package main

import (
	"net/http"

	"github.com/ksusonic/go-devops-mon/internal/controllers"
	"github.com/ksusonic/go-devops-mon/internal/hash"
	"github.com/ksusonic/go-devops-mon/internal/server"
	"github.com/ksusonic/go-devops-mon/internal/server/middleware"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()

	config, err := server.NewConfig()
	if err != nil {
		logger.Fatal("error reading config", zap.Error(err))
	}

	router := chi.NewRouter()
	router.Use(middleware.GzipEncoder)

	metricsStorage, err := server.NewStorage(config, logger)
	if err != nil {
		logger.Fatal("error creating storage", zap.Error(err))
	}
	defer metricsStorage.Close()

	hashService := hash.NewService(config.SecretKey)
	metricController := controllers.NewMetricController(
		logger.Named("MetricController"),
		metricsStorage,
		hashService,
	)
	router.Mount("/", metricController.Router())

	logger.Info("Server started!", zap.String("address", config.Address))
	err = http.ListenAndServe(config.Address, router)
	if err != nil {
		logger.Fatal("shutdown", zap.Error(err))
	}
}
