//go:generate go run

package main

import (
	"fmt"
	"net/http"

	"github.com/ksusonic/go-devops-mon/internal/controllers"
	"github.com/ksusonic/go-devops-mon/internal/hash"
	"github.com/ksusonic/go-devops-mon/internal/server"
	"github.com/ksusonic/go-devops-mon/internal/server/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Println(getBuildInfo())

	config, err := server.NewConfig()
	if err != nil {
		panic("error reading config: " + err.Error())
	}

	logger, _ := getLogger(config.Debug)

	router := chi.NewRouter()
	router.Use(middleware.GzipEncoder)
	if config.Debug {
		router.Mount("/debug", chiMiddleware.Profiler())
	}

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

func getLogger(debug bool) (*zap.Logger, error) {
	if debug {
		return zap.NewDevelopment()
	} else {
		return zap.NewProduction()
	}
}

func getBuildInfo() string {
	return fmt.Sprintf(
		"------------------\n"+
			"Build version: %s\n"+
			"Build    date: %s\n"+
			"Build  commit: %s\n"+
			"------------------\n",
		buildVersion,
		buildDate,
		buildCommit,
	)
}
