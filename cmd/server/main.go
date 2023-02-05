package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/ksusonic/go-devops-mon/internal/controllers"
	"github.com/ksusonic/go-devops-mon/internal/filerepository"
	"github.com/ksusonic/go-devops-mon/internal/hash"
	"github.com/ksusonic/go-devops-mon/internal/server"
	"github.com/ksusonic/go-devops-mon/internal/server/middleware"
	"github.com/ksusonic/go-devops-mon/internal/storage"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	config, err := server.NewConfig()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	var repository *storage.MemStorageRepository = nil
	if config.StoreFile != "" {
		rep, err := filerepository.NewFileRepository(config.StoreFile)
		if err != nil {
			log.Fatalf("Error in repository: %v", err)
		}
		defer rep.Close()

		repository = &storage.MemStorageRepository{
			Repository:         rep,
			DropInterval:       config.FileDropIntervalDuration,
			NeedRestoreMetrics: config.RestoreFile,
		}
	}

	memStorage := storage.NewMemStorage(repository)

	var db *sql.DB
	if config.DatabaseURL != "" {
		db, err = sql.Open("pgx", config.DatabaseURL)
		if err != nil {
			panic(err)
		}
		log.Println("Successfully connected to db")
		defer db.Close()
	}

	router := chi.NewRouter()
	router.Use(middleware.GzipEncoder)
	hashService := hash.NewService(config.SecretKey)
	metricController := controllers.NewMetricController(memStorage, hashService)
	router.Mount("/", metricController.Router())

	log.Printf("Server started on %s\n", config.Address)
	err = http.ListenAndServe(config.Address, router)
	if err != nil {
		log.Fatal(err)
	}
}
