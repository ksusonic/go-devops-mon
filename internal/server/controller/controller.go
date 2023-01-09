package controller

import (
	"github.com/ksusonic/go-devops-mon/internal/metrics"

	"github.com/go-chi/chi/v5"
)

type Controller struct {
	Storage metrics.MetricStorage
}

func (c *Controller) Register(router *chi.Mux) {
	router.Get("/value/{type}/{name}", c.getMetricHandler)
	router.Get("/", c.getAllMetricsHandler)

	router.Post("/update/{type}/{name}/{value}", c.updateMetricHandler)
}

func NewController(storage metrics.MetricStorage) *Controller {
	return &Controller{
		Storage: storage,
	}
}
