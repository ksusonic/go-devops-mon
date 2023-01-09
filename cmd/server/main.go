package main

import (
	"log"

	"github.com/ksusonic/go-devops-mon/internal/controllers/metric"
	"github.com/ksusonic/go-devops-mon/internal/server"
	"github.com/ksusonic/go-devops-mon/internal/storage"

	"github.com/go-chi/chi/v5"
)

const DefaultAddress = "127.0.0.1:8080"

func main() {
	config := server.Config{
		ServeAddress: DefaultAddress,
	}

	memStorage := storage.NewMemStorage()
	router := chi.NewRouter()
	metricController := metric.NewMetricController(memStorage)
	router.Mount("/", metricController.Router())

	s := server.NewServer(&config, router)
	err := s.Start()
	if err != nil {
		log.Fatal(err)
	}
}
