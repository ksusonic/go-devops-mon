package controllers

import (
	"encoding/json"
	"github.com/ksusonic/go-devops-mon/internal/metrics"
	"log"
	"net/http"
	"strconv"

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

func (c *Controller) getMetricHandler(w http.ResponseWriter, r *http.Request) {
	reqType := chi.URLParam(r, "type")
	reqName := chi.URLParam(r, "name")
	value, err := c.Storage.GetMetric(reqType, reqName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var stringValue string
	if reqType == metrics.CounterType {
		stringValue = strconv.FormatInt(value.Value.(int64), 10)
	} else {
		stringValue = strconv.FormatFloat(value.Value.(float64), 'f', -1, 64)
	}
	_, _ = w.Write([]byte(stringValue))
}

func (c *Controller) getAllMetricsHandler(w http.ResponseWriter, _ *http.Request) {
	marshall, err := json.Marshal(c.Storage.GetMappedByTypeAndNameMetrics())
	if err != nil {
		log.Fatal(err)
	}
	_, _ = w.Write(marshall)
}

// updateMetricHandler — обработчик обновления метрики по типу и названию
func (c *Controller) updateMetricHandler(w http.ResponseWriter, r *http.Request) {
	reqType := chi.URLParam(r, "type")
	reqName := chi.URLParam(r, "name")
	reqRawValue := chi.URLParam(r, "value")

	if reqType == metrics.GaugeType {
		value, err := strconv.ParseFloat(reqRawValue, 64)
		if err != nil {
			log.Printf("Incorrect value: %s\n", reqRawValue)
			w.WriteHeader(http.StatusBadRequest)
		}
		c.Storage.SetMetric(metrics.AtomicMetric{
			Name:  reqName,
			Type:  reqType,
			Value: value,
		})
		log.Printf("Updated gauge %s: %f\n", reqName, value)
	} else if reqType == metrics.CounterType {
		value, err := strconv.ParseInt(reqRawValue, 10, 64)
		if err != nil {
			log.Printf("Incorrect value: %s\n", reqRawValue)
			w.WriteHeader(http.StatusBadRequest)
		}
		c.Storage.SetMetric(metrics.AtomicMetric{
			Name:  reqName,
			Type:  reqType,
			Value: value,
		})
		log.Printf("Updated counter %s: %d\n", reqName, value)
	} else {
		log.Println("unexpected metric type!")
		w.WriteHeader(http.StatusNotImplemented)
	}
}
