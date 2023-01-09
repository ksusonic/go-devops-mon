package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ksusonic/go-devops-mon/internal/metrics"

	"github.com/go-chi/chi/v5"
)

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
