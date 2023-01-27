package main

import (
	"log"
	"net/http"

	"github.com/ksusonic/go-devops-mon/internal/controllers"
	"github.com/ksusonic/go-devops-mon/internal/filerepository"
	"github.com/ksusonic/go-devops-mon/internal/server"
	"github.com/ksusonic/go-devops-mon/internal/server/middleware"
	"github.com/ksusonic/go-devops-mon/internal/storage"

	"github.com/go-chi/chi/v5"
)

func main() {
	config, err := server.NewConfig()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	var repository *storage.MemStorageRepository
	if config.StoreFile != "" {
		repository = &storage.MemStorageRepository{
			Repository:         nil,
			DropInterval:       config.FileDropIntervalDuration,
			NeedRestoreMetrics: config.RestoreFile,
		}

		repository.Repository, err = filerepository.NewFileRepository(config.StoreFile)
		if err != nil {
			log.Fatalf("Error in repository: %v", err)
		}
		defer repository.Repository.Close()
	}

	memStorage := storage.NewMemStorage(repository)

	router := chi.NewRouter()
	router.Use(middleware.GzipEncoder)
	metricController := controllers.NewMetricController(memStorage)
	router.Mount("/", metricController.Router())

	log.Printf("Server started on %s\n", config.Address)
	err = http.ListenAndServe(config.Address, router)
	if err != nil {
		log.Fatal(err)
	}
}
