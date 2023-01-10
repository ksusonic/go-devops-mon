package main

import (
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
	router := chi.NewRouter()
	metricController := controllers.NewMetricController(memStorage)
	router.Mount("/", metricController.Router())

	log.Println("Server started")
	err := http.ListenAndServe(config.Address, router)
	if err != nil {
		log.Fatal(err)
	}
}
