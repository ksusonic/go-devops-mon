package main

import (
	"github.com/ksusonic/go-devops-mon/internal/filemanage"
	"log"
	"net/http"

	"github.com/ksusonic/go-devops-mon/internal/controllers"
	"github.com/ksusonic/go-devops-mon/internal/server"
	"github.com/ksusonic/go-devops-mon/internal/storage"

	"github.com/go-chi/chi/v5"
)

func main() {
	config := server.NewConfig()

	memStorage := storage.NewMemStorage()

	if config.StoreFile != "" {
		if config.RestoreFile {
			restored := filemanage.RestoreMetrics(config.StoreFile)
			memStorage.AddMetrics(restored)
			log.Printf("Restored %d metrics\n", len(restored))
		}
		fileProducer := filemanage.NewFileProducer(config.StoreFile, config.RestoreFile)
		defer fileProducer.Close()

		go fileProducer.DropRoutine(memStorage, config.FileDropInterval)
		log.Printf("Enabled drop metrics to %s\n", config.StoreFile)
	}

	router := chi.NewRouter()
	metricController := controllers.NewMetricController(memStorage)
	router.Mount("/", metricController.Router())

	log.Printf("Server started on %s\n", config.Address)
	err := http.ListenAndServe(config.Address, router)
	if err != nil {
		log.Fatal(err)
	}
}
