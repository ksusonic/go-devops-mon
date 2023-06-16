//go:generate go run

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/ksusonic/go-devops-mon/internal/controllers"
	"github.com/ksusonic/go-devops-mon/internal/crypt"
	"github.com/ksusonic/go-devops-mon/internal/hash"
	"github.com/ksusonic/go-devops-mon/internal/server"
	"github.com/ksusonic/go-devops-mon/internal/server/middleware"
	"github.com/ksusonic/go-devops-mon/internal/trust"
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
	router.Use(chiMiddleware.Logger)
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
	decryptService, err := crypt.NewDecrypter(config.CryptoKeyPath, logger.Named("decrypter"))
	if config.CryptoKeyPath != "" && err != nil {
		logger.Fatal("error creating decrypter", zap.Error(err))
	}

	if decryptService != nil {
		logger.Info("using decrypt middleware")
		router.Use(decryptService.Middleware)
	}
	if len(config.TrustedSubnet) > 0 {
		trustService, err := trust.NewNetTrustService(config.TrustedSubnet, logger.Named("trust"))
		if err != nil {
			logger.Fatal("incorrect CIDR subnet from config", zap.String("subnet", config.TrustedSubnet))
		}
		router.Use(trustService.Middleware)
	}

	metricController := controllers.NewMetricController(
		logger.Named("MetricController"),
		metricsStorage,
		hashService,
	)
	router.Mount("/", metricController.Router())

	var srv = http.Server{Addr: config.Address, Handler: router}
	var grpcSrv = server.NewServer(metricsStorage, hashService, logger)

	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() { // graceful shutdown watcher
		<-sigint

		grpcSrv.GrpcServer.GracefulStop()
		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Info("HTTP server Shutdown", zap.Error(err))
		}
		close(idleConnsClosed)
	}()

	go grpcSrv.Start(3200)
	go func() {
		logger.Info("Listening http", zap.String("address", config.Address))
		if err := http.ListenAndServe(config.Address, router); err != http.ErrServerClosed {
			logger.Fatal("shutdown", zap.Error(err))
		}
		logger.Info("http server down")
	}()

	<-idleConnsClosed
	logger.Info("Server Shutdown gracefully")
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
