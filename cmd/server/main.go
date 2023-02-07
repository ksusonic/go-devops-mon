package main

import (
	"net/http"

	"github.com/ksusonic/go-devops-mon/internal/controllers"
	"github.com/ksusonic/go-devops-mon/internal/db"
	"github.com/ksusonic/go-devops-mon/internal/filerepository"
	"github.com/ksusonic/go-devops-mon/internal/hash"
	"github.com/ksusonic/go-devops-mon/internal/metrics"
	"github.com/ksusonic/go-devops-mon/internal/server"
	"github.com/ksusonic/go-devops-mon/internal/server/middleware"
	"github.com/ksusonic/go-devops-mon/internal/storage"

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

	var metricsStorage metrics.ServerMetricStorage

	if config.DatabaseURL != "" {
		logger.Info("using DB-storage")
		metricsStorage, err = db.NewDB(config.DatabaseURL)
		if err != nil {
			logger.Fatal("Error in database: %v", zap.Error(err))
		}
		logger.Debug("successfully initialized db")
	} else {
		logger.Info("using in-memory storage")
		var repository metrics.Repository
		if config.StoreFile != "" {
			repository, err = filerepository.NewFileRepository(config.StoreFile)
			if err != nil {
				logger.Fatal("Error in repository: %v", zap.Error(err))
			}
		}

		metricsStorage = storage.NewMemStorage(
			&storage.MemStorageRepository{
				Repository:         repository,
				DropInterval:       config.FileDropIntervalDuration,
				NeedRestoreMetrics: config.RestoreFile,
			},
		)
	}
	defer metricsStorage.Close()

	hashService := hash.NewService(config.SecretKey)
	metricController := controllers.NewMetricController(logger, metricsStorage, hashService)
	router.Mount("/", metricController.Router())

	logger.Info("Server started!", zap.String("address", config.Address))
	err = http.ListenAndServe(config.Address, router)
	if err != nil {
		logger.Fatal("shutdown", zap.Error(err))
	}
}
