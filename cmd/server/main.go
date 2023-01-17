package main

import (
	"log"
	"net/http"

	"github.com/ksusonic/go-devops-mon/internal/controllers"
	"github.com/ksusonic/go-devops-mon/internal/filemanage"
	"github.com/ksusonic/go-devops-mon/internal/server"
	"github.com/ksusonic/go-devops-mon/internal/storage"

	"github.com/go-chi/chi/v5"
)

func main() {
	config, err := server.NewConfig()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	memStorage := storage.NewMemStorage()

	if config.StoreFile != "" {
		if config.RestoreFile {
			restored, err := filemanage.RestoreMetrics(config.StoreFile)
			if err != nil {
				log.Fatalf("Error restoring metrics: %v", err)
			}
			memStorage.AddMetrics(restored)
			log.Printf("Restored %d metrics\n", len(restored))
		}
		fileProducer, err := filemanage.NewFileProducer(config.StoreFile, config.RestoreFile)
		if err != nil {
			log.Fatalf("Error creating file producer: %v", err)
		}
		defer fileProducer.Close()

		go fileProducer.DropRoutine(memStorage, config.FileDropInterval)
		log.Printf("Enabled drop metrics to %s\n", config.StoreFile)
	}

	router := chi.NewRouter()
	router.Use(server.GzipEncoder)
	metricController := controllers.NewMetricController(memStorage)
	router.Mount("/", metricController.Router())

	log.Printf("Server started on %s\n", config.Address)
	err = http.ListenAndServe(config.Address, router)
	if err != nil {
		log.Fatal(err)
	}
}
