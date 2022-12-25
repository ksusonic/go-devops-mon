package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/ksusonic/go-devops-mon/internal/metrics"
	"log"
	"net/http"
	"strconv"
)

type getRequest struct {
	Type string
	Name string
}

func (s Server) GetMetric(w http.ResponseWriter, r *http.Request) {
	requestData := getRequest{
		Type: chi.URLParam(r, "type"),
		Name: chi.URLParam(r, "name"),
	}
	if requestData.Type == metrics.GaugeName {
		value, err := s.MemStorage.GetCurrentGaugeMetric(requestData.Name)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write([]byte(strconv.FormatFloat(value, 'f', -1, 64)))
	} else if requestData.Type == metrics.CounterName {
		value, err := s.MemStorage.GetCurrentCounterMetric(requestData.Name)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write([]byte(strconv.FormatInt(value, 10)))
	} else {
		log.Println("unexpected metric type!")
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (s Server) GetAllMetrics(w http.ResponseWriter, _ *http.Request) {
	result := make(map[string]interface{})
	result["gauge"] = s.MemStorage.GetAllCurrentGaugeMetrics()
	result["counter"] = s.MemStorage.GetAllCurrentCounterMetrics()
	marshall, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(marshall)
}
