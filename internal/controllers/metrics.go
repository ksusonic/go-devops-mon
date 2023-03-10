package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/ksusonic/go-devops-mon/internal/metrics"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type Controller struct {
	Logger      *zap.Logger
	Storage     metrics.ServerMetricStorage
	HashService metrics.HashService
}

func NewMetricController(logger *zap.Logger, storage metrics.ServerMetricStorage, hashService metrics.HashService) *Controller {
	return &Controller{
		Logger:      logger,
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
		c.Logger.Error("error writing response", zap.Error(err))
	}
}

func (c *Controller) getMetricHandler(w http.ResponseWriter, r *http.Request) {
	m := &metrics.Metrics{}
	if err := render.Bind(r, m); err != nil {
		render.Render(w, r, ErrBadRequest(err, c.Logger))
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
		c.Logger.Error("error unmarshalling", zap.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(marshal)
	if err != nil {
		c.Logger.Error("error writing response", zap.Error(err))
	}
}

func (c *Controller) getAllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	metricsMapping, err := c.Storage.GetMappedByTypeAndNameMetrics(r.Context())
	if err != nil {
		c.Logger.Error("could not GetMappedByTypeAndNameMetrics", zap.Error(err))
		return
	}
	marshall, err := json.Marshal(metricsMapping)
	if err != nil {
		c.Logger.Fatal("error while marshalling mappedByTypeAndNameMetrics", zap.Error(err))
	}
	w.Header().Add("Content-Type", "text/html")
	_, err = w.Write(marshall)
	if err != nil {
		c.Logger.Error("error writing response", zap.Error(err))
	}
}

// updateMetricPathHandler ??? updates metric by type, name and value
func (c *Controller) updateMetricPathHandler(w http.ResponseWriter, r *http.Request) {
	reqType := chi.URLParam(r, "type")
	reqName := chi.URLParam(r, "name")
	reqRawValue := chi.URLParam(r, "value")

	if reqType == metrics.GaugeMType {
		value, err := strconv.ParseFloat(reqRawValue, 64)
		if err != nil {
			render.Render(w, r, ErrBadRequest(err, c.Logger))
			return
		}
		m := metrics.Metrics{
			ID:    reqName,
			MType: reqType,
			Value: &value,
		}
		if err = c.HashService.SetHash(&m); err != nil {
			render.Render(w, r, ErrInternalError(err, c.Logger))
			return
		}

		if _, err = c.Storage.SetMetric(r.Context(), m); err != nil {
			render.Render(w, r, ErrInternalError(err, c.Logger))
			return
		}
		c.Logger.Info("Updated gauge", zap.String("id", reqName), zap.Float64("value", value))
	} else if reqType == metrics.CounterMType {
		value, err := strconv.ParseInt(reqRawValue, 10, 64)
		if err != nil {
			c.Logger.Error("incorrect value", zap.String("value", reqRawValue))
			w.WriteHeader(http.StatusBadRequest)
		}
		m := metrics.Metrics{
			ID:    reqName,
			MType: reqType,
			Delta: &value,
		}
		if err = c.HashService.SetHash(&m); err != nil {
			render.Render(w, r, ErrInternalError(err, c.Logger))
			return
		}

		if _, err = c.Storage.SetMetric(r.Context(), m); err != nil {
			render.Render(w, r, ErrInternalError(err, c.Logger))
			return
		}
		c.Logger.Info("Updated counter", zap.String("id", reqName), zap.Int64("delta", value))
	} else {
		c.Logger.Error("unexpected metric type", zap.String("type", reqType))
		w.WriteHeader(http.StatusNotImplemented)
	}
}

// updateMetricHandler ??? updates metric by Metrics data in body
func (c *Controller) updateMetricHandler(w http.ResponseWriter, r *http.Request) {
	m := &metrics.Metrics{}
	if err := render.Bind(r, m); err != nil {
		render.Render(w, r, ErrBadRequest(err, c.Logger))
		return
	}

	if err := c.HashService.ValidateHash(m); err != nil {
		render.Render(w, r, ErrBadRequest(err, c.Logger))
		return
	}

	if m.MType != metrics.GaugeMType && m.MType != metrics.CounterMType {
		w.WriteHeader(http.StatusNotImplemented)
		c.Logger.Error("unexpected metric type", zap.String("type", m.MType))
	} else {
		resultMetric, err := c.Storage.SetMetric(r.Context(), *m)
		if err != nil {
			c.Logger.Error("error while setting metric", zap.Error(err))
			return
		}
		c.Logger.Debug("Updated", zap.String("metric", m.String()))

		marshal, err := json.Marshal(resultMetric)
		if err != nil {
			c.Logger.Fatal("error while marshalling in updateMetricHandler", zap.Error(err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(marshal)
		if err != nil {
			c.Logger.Error("error writing response", zap.Error(err))
		}
	}
}

// updatesMetricHandler ??? updates metric by []Metrics data in body
func (c *Controller) updatesMetricHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.Logger.Error("could not read request body", zap.Error(err))
		return
	}
	var metricSlice []metrics.Metrics
	err = json.Unmarshal(body, &metricSlice)
	if err != nil {
		c.Logger.Error("error unmarshalling metric slice", zap.Error(err))
		return
	}

	for _, m := range metricSlice {
		if err := c.HashService.ValidateHash(&m); err != nil {
			render.Render(w, r, ErrBadRequest(err, c.Logger))
			return
		}
	}

	err = c.Storage.SetMetrics(r.Context(), &metricSlice)
	if err != nil {
		c.Logger.Error("error setting metric", zap.Error(err))
		render.Render(w, r, ErrInternalError(err, c.Logger))
	}
}
