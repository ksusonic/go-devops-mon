package metric

import (
	"github.com/ksusonic/go-devops-mon/internal/metrics"

	"github.com/go-chi/chi/v5"
)

type Controller struct {
	Storage metrics.ServerMetricStorage
}

func (c *Controller) Router() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/value/{type}/{name}", c.getMetricHandler)
	router.Get("/", c.getAllMetricsHandler)

	router.Post("/update/{type}/{name}/{value}", c.updateMetricHandler)

	return router
}

func NewMetricController(storage metrics.ServerMetricStorage) *Controller {
	return &Controller{
		Storage: storage,
	}
}
