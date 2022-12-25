package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/ksusonic/go-devops-mon/internal/metrics"
	"log"
	"net/http"
	"strconv"
)

type updateRequest struct {
	Type     string
	Name     string
	RawValue string
}

// UpdateMetric — обработчик обновления метрики по типу и названию
func (s Server) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	requestData := updateRequest{
		Type:     chi.URLParam(r, "type"),
		Name:     chi.URLParam(r, "name"),
		RawValue: chi.URLParam(r, "value"),
	}
	if requestData.Type == metrics.GaugeName {
		value, err := strconv.ParseFloat(requestData.RawValue, 64)
		if err != nil {
			log.Printf("Incorrect value: %s\n", requestData.RawValue)
			w.WriteHeader(http.StatusBadRequest)
		}
		s.MemStorage.AddGaugeValue(requestData.Name, value)
		log.Printf("Updated gauge %s: %f\n", requestData.Name, value)
	} else if requestData.Type == metrics.CounterName {
		value, err := strconv.ParseInt(requestData.RawValue, 10, 64)
		if err != nil {
			log.Printf("Incorrect value: %s\n", requestData.RawValue)
			w.WriteHeader(http.StatusBadRequest)
		}
		log.Printf("Updated counter %s: %d\n", requestData.Name, value)
	} else {
		log.Println("unexpected metric type!")
		w.WriteHeader(http.StatusNotImplemented)
	}

	w.WriteHeader(http.StatusOK)
}
