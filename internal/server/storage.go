package server

import (
	"context"
	"fmt"
	"github.com/ksusonic/go-devops-mon/internal/db"
	"github.com/ksusonic/go-devops-mon/internal/filerepository"
	"github.com/ksusonic/go-devops-mon/internal/metrics"
	"github.com/ksusonic/go-devops-mon/internal/storage"

	"go.uber.org/zap"
)

func NewStorage(config *Config, logger *zap.Logger) (metrics.ServerMetricStorage, error) {
	if config.DatabaseURL != "" {
		logger.Info("using DB-storage")
		metricsStorage, err := db.NewDB(config.DatabaseURL)
		if err != nil {
			return nil, fmt.Errorf("error in database: %v", err)
		}
		return metricsStorage, nil
	} else {
		logger.Info("using in-memory storage")
		memStorage := storage.NewMemStorage()

		if config.StoreFile != "" {
			repository, err := filerepository.NewFileRepository(config.StoreFile, logger.Named("FileRepository"))
			if err != nil {
				logger.Fatal("Error in repository: %v", zap.Error(err))
			}
			if config.RestoreFile {
				restoredMetrics := repository.ReadCurrentState()
				if err := memStorage.SetMetrics(context.Background(), &restoredMetrics); err != nil {
					return nil, fmt.Errorf("error setting current metrics from repository: %v", err)
				}
				logger.Info("Restored metrics", zap.Int("amount", len(restoredMetrics)))
			}
			go repository.DropRoutine(context.Background(), memStorage.GetAllMetrics, config.FileDropIntervalDuration)
		}
		return memStorage, nil
	}
}
