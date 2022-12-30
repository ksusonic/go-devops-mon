package router

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ksusonic/go-devops-mon/internal/metrics"

	"github.com/go-chi/chi/v5"
)

func init() {
	registerHandler("GET", "/value/{type}/{name}", getMetricHandler)
	registerHandler("GET", "/", getAllMetricsHandler)
}

var getMetricHandler = func(w http.ResponseWriter, r *http.Request, c context) {
	reqType := chi.URLParam(r, "type")
	reqName := chi.URLParam(r, "name")
	value, err := (*c.storage).GetMetric(reqName)
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

var getAllMetricsHandler = func(w http.ResponseWriter, _ *http.Request, c context) {
	marshall, err := json.Marshal((*c.storage).GetMappedByTypeAndNameMetrics())
	if err != nil {
		log.Fatal(err)
	}
	_, _ = w.Write(marshall)
}
