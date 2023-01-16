package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/ksusonic/go-devops-mon/internal/metrics"

	"github.com/go-chi/chi/v5"
)

type Controller struct {
	Storage metrics.ServerMetricStorage
}

func (c *Controller) Router() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/value/{type}/{name}", c.getMetricPathHandler)
	router.Get("/", c.getAllMetricsHandler)

	router.Post("/update/{type}/{name}/{value}", c.updateMetricPathHandler)

	router.Post("/update/", c.updateMetricHandler)
	router.Post("/value/", c.getMetricHandler)

	return router
}

func NewMetricController(storage metrics.ServerMetricStorage) *Controller {
	return &Controller{
		Storage: storage,
	}
}

func (c *Controller) getMetricPathHandler(w http.ResponseWriter, r *http.Request) {
	reqType := chi.URLParam(r, "type")
	reqName := chi.URLParam(r, "name")
	value, err := c.Storage.GetMetric(reqType, reqName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var stringValue string
	if reqType == metrics.CounterMType {
		stringValue = strconv.FormatInt(*value.Delta, 10)
	} else {
		stringValue = strconv.FormatFloat(*value.Value, 'f', -1, 64)
	}
	_, _ = w.Write([]byte(stringValue))
}

func (c *Controller) getMetricHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var m metrics.Metrics
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Println(err)
		return
	}

	value, err := c.Storage.GetMetric(m.MType, m.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	marshal, err := json.Marshal(value)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(marshal)
}

func (c *Controller) getAllMetricsHandler(w http.ResponseWriter, _ *http.Request) {
	marshall, err := json.Marshal(c.Storage.GetMappedByTypeAndNameMetrics())
	if err != nil {
		log.Fatal(err)
	}
	_, _ = w.Write(marshall)
}

// updateMetricPathHandler — updates metric by type, name and value
func (c *Controller) updateMetricPathHandler(w http.ResponseWriter, r *http.Request) {
	reqType := chi.URLParam(r, "type")
	reqName := chi.URLParam(r, "name")
	reqRawValue := chi.URLParam(r, "value")

	if reqType == metrics.GaugeMType {
		value, err := strconv.ParseFloat(reqRawValue, 64)
		if err != nil {
			log.Printf("Incorrect value: %s\n", reqRawValue)
			w.WriteHeader(http.StatusBadRequest)
		}
		c.Storage.SetMetric(metrics.Metrics{
			ID:    reqName,
			MType: reqType,
			Value: &value,
		})
		log.Printf("Updated gauge %s: %f\n", reqName, value)
	} else if reqType == metrics.CounterMType {
		value, err := strconv.ParseInt(reqRawValue, 10, 64)
		if err != nil {
			log.Printf("Incorrect value: %s\n", reqRawValue)
			w.WriteHeader(http.StatusBadRequest)
		}
		c.Storage.SetMetric(metrics.Metrics{
			ID:    reqName,
			MType: reqType,
			Delta: &value,
		})
		log.Printf("Updated counter %s: %d\n", reqName, value)
	} else {
		log.Println("unexpected metric type!")
		w.WriteHeader(http.StatusNotImplemented)
	}
}

// updateMetricHandler — updates metric by Metrics data in body
func (c *Controller) updateMetricHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var m metrics.Metrics
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Println(err)
		return
	}

	if m.MType != metrics.GaugeMType && m.MType != metrics.CounterMType {
		w.WriteHeader(http.StatusNotImplemented)
		log.Printf("Unknown metric type %s\n", m.MType)
	} else {
		resultMetric := c.Storage.SetMetric(m)
		log.Printf("Updated %s\n", m.String())

		marshal, err := json.Marshal(resultMetric)
		if err != nil {
			log.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(marshal)
	}
}
