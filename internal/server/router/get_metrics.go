package router

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/ksusonic/go-devops-mon/internal/metrics"

	"github.com/go-chi/chi/v5"
)

func init() {
	registerHandler("POST", "/value/", getMetricHandler)
	registerHandler("GET", "/value/{type}/{name}", getOneMetricHandler)
	registerHandler("GET", "/", getAllMetricsHandler)
}

var getMetricHandler = func(w http.ResponseWriter, r *http.Request, c context) {
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

	value, err := (*c.storage).GetMetric(m.ID)
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

var getOneMetricHandler = func(w http.ResponseWriter, r *http.Request, c context) {
	reqName := chi.URLParam(r, "name")
	value, err := (*c.storage).GetMetric(reqType, reqName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var stringValue string
	if value.MType == metrics.CounterMType {
		stringValue = strconv.FormatInt(*value.Delta, 10)
	} else {
		stringValue = strconv.FormatFloat(*value.Value, 'f', -1, 64)
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
