package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ksusonic/go-devops-mon/internal/metrics"

	"github.com/go-chi/chi/v5"
)

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
