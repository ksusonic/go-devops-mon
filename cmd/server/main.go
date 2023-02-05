package main

import (
	"log"
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
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	config, err := server.NewConfig()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.GzipEncoder)

	var metricsStorage metrics.ServerMetricStorage

	if config.DatabaseURL != "" {
		log.Println("Using DB-storage")
		metricsStorage, err = db.NewDB(config.DatabaseURL)
		if err != nil {
			log.Fatalf("Error in database: %v", err)
		}
		log.Println("Initted db")
	} else {
		log.Println("Using in-memory storage")
		var repository metrics.Repository
		if config.StoreFile != "" {
			repository, err = filerepository.NewFileRepository(config.StoreFile)
			if err != nil {
				log.Fatalf("Error in repository: %v", err)
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
	metricController := controllers.NewMetricController(metricsStorage, hashService)
	router.Mount("/", metricController.Router())

	log.Printf("Server started on %s\n", config.Address)
	err = http.ListenAndServe(config.Address, router)
	if err != nil {
		log.Fatal(err)
	}
}
