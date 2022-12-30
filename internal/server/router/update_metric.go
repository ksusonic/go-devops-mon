package router

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ksusonic/go-devops-mon/internal/metrics"

	"github.com/go-chi/chi/v5"
)

func init() {
	registerHandler("POST", "/update/{type}/{name}/{value}", updateMetricHandler)
}

type updateRequest struct {
	Type     string
	Name     string
	RawValue string
}

// updateMetricHandler — обработчик обновления метрики по типу и названию
var updateMetricHandler = func(w http.ResponseWriter, r *http.Request, c context) {
	requestData := updateRequest{
		Type:     chi.URLParam(r, "type"),
		Name:     chi.URLParam(r, "name"),
		RawValue: chi.URLParam(r, "value"),
	}
	if requestData.Type == metrics.GaugeType {
		value, err := strconv.ParseFloat(requestData.RawValue, 64)
		if err != nil {
			log.Printf("Incorrect value: %s\n", requestData.RawValue)
			w.WriteHeader(http.StatusBadRequest)
		}
		(*c.storage).SetMetric(metrics.AtomicMetric{
			Name:  requestData.Name,
			Type:  requestData.Type,
			Value: value,
		})
		log.Printf("Updated gauge %s: %f\n", requestData.Name, value)
	} else if requestData.Type == metrics.CounterType {
		value, err := strconv.ParseInt(requestData.RawValue, 10, 64)
		if err != nil {
			log.Printf("Incorrect value: %s\n", requestData.RawValue)
			w.WriteHeader(http.StatusBadRequest)
		}
		(*c.storage).SetMetric(metrics.AtomicMetric{
			Name:  requestData.Name,
			Type:  requestData.Type,
			Value: value,
		})
		log.Printf("Updated counter %s: %d\n", requestData.Name, value)
	} else {
		log.Println("unexpected metric type!")
		w.WriteHeader(http.StatusNotImplemented)
	}
}
