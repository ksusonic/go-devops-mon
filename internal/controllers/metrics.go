package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/ksusonic/go-devops-mon/internal/metrics"
	metricspb "github.com/ksusonic/go-devops-mon/proto/metrics"
	"go.uber.org/zap"
)

type DecodeService interface {
	DecryptBytes(b []byte) ([]byte, error)
}

type HashService interface {
	// SetHash Calculates hash for metric
	SetHash(m *metrics.Metric) error
	// ValidateHash Returns bool if calculated and metric hash is correct
	ValidateHash(m *metrics.Metric) error
	ValidateHashProto(m *metricspb.Metric) error
}

type Controller struct {
	Logger  *zap.Logger
	Storage metrics.ServerMetricStorage

	hashService HashService
	decoder     DecodeService
}

func NewMetricController(logger *zap.Logger, storage metrics.ServerMetricStorage, hashService HashService) *Controller {
	return &Controller{
		Logger:      logger,
		Storage:     storage,
		hashService: hashService,
	}
}

func (c *Controller) Router() *chi.Mux {
	router := chi.NewRouter()

	router.Use()
	router.Get("/ping", c.PingHandler)
	router.Get("/value/{type}/{name}", c.GetMetricPathHandler)
	router.Post("/update/{type}/{name}/{value}", c.UpdateMetricPathHandler)

	// These handlers are grpc-supported. TODO: union handlers
	router.Get("/", c.GetAllMetricsHandler)
	router.Post("/update/", c.UpdateMetricHandler)
	router.Post("/updates/", c.UpdatesMetricHandler)
	router.Post("/value/", c.GetMetricHandler)

	return router
}

// GetMetricPathHandler - get metric by query params: type, name
func (c *Controller) GetMetricPathHandler(w http.ResponseWriter, r *http.Request) {
	reqType := chi.URLParam(r, "type")
	reqName := chi.URLParam(r, "name")
	value, err := c.Storage.GetMetric(r.Context(), reqType, reqName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var stringValue string
	if reqType == metrics.CounterType {
		stringValue = strconv.FormatInt(value.GetDelta(), 10)
	} else {
		stringValue = strconv.FormatFloat(value.GetValue(), 'f', -1, 64)
	}
	_, err = w.Write([]byte(stringValue))
	if err != nil {
		c.Logger.Error("error writing response", zap.Error(err))
	}
}

// GetMetricHandler - get metric by body in post request
func (c *Controller) GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	m := &metrics.Metric{}
	if err := render.Bind(r, m); err != nil {
		render.Render(w, r, ErrBadRequest(err, c.Logger))
		return
	}

	value, err := c.Storage.GetMetric(r.Context(), m.Type, m.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = c.hashService.SetHash(value)
	if err != nil {
		c.Logger.Error("error setting hash", zap.Error(err))
		return
	}

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

// GetAllMetricsHandler - get all current metrics in a list
func (c *Controller) GetAllMetricsHandler(w http.ResponseWriter, r *http.Request) {
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

// UpdateMetricPathHandler — updates metric by type, name and value
func (c *Controller) UpdateMetricPathHandler(w http.ResponseWriter, r *http.Request) {
	reqType := chi.URLParam(r, "type")
	reqName := chi.URLParam(r, "name")
	reqRawValue := chi.URLParam(r, "value")

	if reqType == metrics.GaugeType {
		value, err := strconv.ParseFloat(reqRawValue, 64)
		if err != nil {
			render.Render(w, r, ErrBadRequest(err, c.Logger))
			return
		}
		m := metrics.Metric{
			ID:    reqName,
			Type:  metrics.GaugeType,
			Value: &value,
		}
		if err = c.hashService.SetHash(&m); err != nil {
			render.Render(w, r, ErrInternalError(err, c.Logger))
			return
		}

		if _, err = c.Storage.SetMetric(r.Context(), &m); err != nil {
			render.Render(w, r, ErrInternalError(err, c.Logger))
			return
		}
		c.Logger.Info("Updated gauge", zap.String("id", reqName), zap.Float64("value", value))
	} else if reqType == metrics.CounterType {
		value, err := strconv.ParseInt(reqRawValue, 10, 64)
		if err != nil {
			c.Logger.Error("incorrect value", zap.String("value", reqRawValue))
			w.WriteHeader(http.StatusBadRequest)
		}
		m := metrics.Metric{
			ID:    reqName,
			Type:  metrics.CounterType,
			Delta: &value,
		}
		if err = c.hashService.SetHash(&m); err != nil {
			render.Render(w, r, ErrInternalError(err, c.Logger))
			return
		}

		if _, err = c.Storage.SetMetric(r.Context(), &m); err != nil {
			render.Render(w, r, ErrInternalError(err, c.Logger))
			return
		}
		c.Logger.Info("Updated counter", zap.String("id", reqName), zap.Int64("delta", value))
	} else {
		c.Logger.Error("unexpected metric type", zap.String("type", reqType))
		w.WriteHeader(http.StatusNotImplemented)
	}
}

// UpdateMetricHandler — updates metric by Metric data in body
func (c *Controller) UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	m := &metrics.Metric{}
	if err := render.Bind(r, m); err != nil {
		render.Render(w, r, ErrBadRequest(err, c.Logger))
		return
	}

	if err := c.hashService.ValidateHash(m); err != nil {
		render.Render(w, r, ErrBadRequest(err, c.Logger))
		return
	}

	if m.Type != metrics.GaugeType && m.Type != metrics.CounterType {
		w.WriteHeader(http.StatusNotImplemented)
		c.Logger.Error("unexpected metric type", zap.String("type", m.Type))
	} else {
		resultMetric, err := c.Storage.SetMetric(r.Context(), m)
		if err != nil {
			c.Logger.Error("error while setting metric", zap.Error(err))
			return
		}
		c.Logger.Debug("Updated", zap.String("metric", m.String()))

		marshal, err := json.Marshal(resultMetric)
		if err != nil {
			c.Logger.Fatal("error while marshalling in UpdateMetricHandler", zap.Error(err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(marshal)
		if err != nil {
			c.Logger.Error("error writing response", zap.Error(err))
		}
	}
}

// UpdatesMetricHandler — updates metric by []Metric data in body
func (c *Controller) UpdatesMetricHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.Logger.Error("could not read request body", zap.Error(err))
		return
	}
	var metricSlice []*metrics.Metric
	err = json.Unmarshal(body, &metricSlice)
	if err != nil {
		c.Logger.Error("error unmarshalling metric slice", zap.Error(err))
		return
	}

	for _, m := range metricSlice {
		if err := c.hashService.ValidateHash(m); err != nil {
			render.Render(w, r, ErrBadRequest(err, c.Logger))
			return
		}
	}

	err = c.Storage.SetMetrics(r.Context(), metricSlice)
	if err != nil {
		c.Logger.Error("error setting metric", zap.Error(err))
		render.Render(w, r, ErrInternalError(err, c.Logger))
	}
	w.Header().Set("Content-Type", "application/json")
	marshall, _ := json.Marshal(metricSlice)
	_, err = w.Write(marshall)
	if err != nil {
		c.Logger.Error("error marshalling status", zap.Error(err))
	}
}
