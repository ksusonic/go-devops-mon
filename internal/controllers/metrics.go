package controllers

import (
	"encoding/json"
	"github.com/ksusonic/go-devops-mon/internal/metrics"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Controller struct {
	Storage     metrics.ServerMetricStorage
	HashService metrics.HashService
}

func NewMetricController(storage metrics.ServerMetricStorage, hashService metrics.HashService) *Controller {
	return &Controller{
		Storage:     storage,
		HashService: hashService,
	}
}

func (c *Controller) Router() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/value/{type}/{name}", c.getMetricPathHandler)
	router.Get("/", c.getAllMetricsHandler)

	router.Post("/update/{type}/{name}/{value}", c.updateMetricPathHandler)

	router.Post("/update/", c.updateMetricHandler)
	router.Post("/updates/", c.updatesMetricHandler)
	router.Post("/value/", c.getMetricHandler)

	router.Get("/ping", c.pingHandler)

	return router
}

func (c *Controller) getMetricPathHandler(w http.ResponseWriter, r *http.Request) {
	reqType := chi.URLParam(r, "type")
	reqName := chi.URLParam(r, "name")
	value, err := c.Storage.GetMetric(r.Context(), reqType, reqName)
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
	_, err = w.Write([]byte(stringValue))
	if err != nil {
		log.Println(err)
	}
}

func (c *Controller) getMetricHandler(w http.ResponseWriter, r *http.Request) {
	m := &metrics.Metrics{}
	if err := render.Bind(r, m); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	value, err := c.Storage.GetMetric(r.Context(), m.MType, m.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	c.HashService.SetHash(&value)

	marshal, err := json.Marshal(value)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(marshal)
	if err != nil {
		log.Println(err)
	}
}

func (c *Controller) getAllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	metricsMapping, err := c.Storage.GetMappedByTypeAndNameMetrics(r.Context())
	if err != nil {
		log.Println(err)
		return
	}
	marshall, err := json.Marshal(metricsMapping)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Add("Content-Type", "text/html")
	_, err = w.Write(marshall)
	if err != nil {
		log.Println(err)
	}
}

// updateMetricPathHandler — updates metric by type, name and value
func (c *Controller) updateMetricPathHandler(w http.ResponseWriter, r *http.Request) {
	reqType := chi.URLParam(r, "type")
	reqName := chi.URLParam(r, "name")
	reqRawValue := chi.URLParam(r, "value")

	if reqType == metrics.GaugeMType {
		value, err := strconv.ParseFloat(reqRawValue, 64)
		if err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}
		m := metrics.Metrics{
			ID:    reqName,
			MType: reqType,
			Value: &value,
		}
		if err = c.HashService.SetHash(&m); err != nil {
			render.Render(w, r, ErrInternalError(err))
			return
		}

		if _, err = c.Storage.SetMetric(r.Context(), m); err != nil {
			render.Render(w, r, ErrInternalError(err))
			return
		}
		log.Printf("Updated gauge %s: %f\n", reqName, value)
	} else if reqType == metrics.CounterMType {
		value, err := strconv.ParseInt(reqRawValue, 10, 64)
		if err != nil {
			log.Printf("Incorrect value: %s\n", reqRawValue)
			w.WriteHeader(http.StatusBadRequest)
		}
		m := metrics.Metrics{
			ID:    reqName,
			MType: reqType,
			Delta: &value,
		}
		if err = c.HashService.SetHash(&m); err != nil {
			render.Render(w, r, ErrInternalError(err))
			return
		}

		if _, err = c.Storage.SetMetric(r.Context(), m); err != nil {
			render.Render(w, r, ErrInternalError(err))
			return
		}
		log.Printf("Updated counter %s: %d\n", reqName, value)
	} else {
		log.Println("unexpected metric type!")
		w.WriteHeader(http.StatusNotImplemented)
	}
}

// updateMetricHandler — updates metric by Metrics data in body
func (c *Controller) updateMetricHandler(w http.ResponseWriter, r *http.Request) {
	m := &metrics.Metrics{}
	if err := render.Bind(r, m); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if err := c.HashService.ValidateHash(m); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if m.MType != metrics.GaugeMType && m.MType != metrics.CounterMType {
		w.WriteHeader(http.StatusNotImplemented)
		log.Printf("Unknown metric type %s\n", m.MType)
	} else {
		resultMetric, err := c.Storage.SetMetric(r.Context(), *m)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Updated %s\n", m)

		marshal, err := json.Marshal(resultMetric)
		if err != nil {
			log.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(marshal)
		if err != nil {
			log.Println(err)
		}
	}
}

// updatesMetricHandler — updates metric by []Metrics data in body
func (c *Controller) updatesMetricHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	var metricSlice []metrics.Metrics
	err = json.Unmarshal(body, &metricSlice)
	if err != nil {
		log.Println(err)
		return
	}

	for _, m := range metricSlice {
		if err := c.HashService.ValidateHash(&m); err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}
	}

	err = c.Storage.SetMetrics(r.Context(), &metricSlice)
	if err != nil {
		log.Printf("Error in updatesMetricHandler: %s\n", err)
		render.Render(w, r, ErrInternalError(err))
	}
}
