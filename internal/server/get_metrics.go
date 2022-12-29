package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ksusonic/go-devops-mon/internal/metrics"

	"github.com/go-chi/chi/v5"
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
	value, err := s.Storage.GetMetric(requestData.Name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var stringValue string
	if requestData.Type == metrics.CounterType {
		stringValue = strconv.FormatInt(value.Value.(int64), 10)
	} else {
		stringValue = strconv.FormatFloat(value.Value.(float64), 'f', -1, 64)
	}
	w.Write([]byte(stringValue))
}

func (s Server) GetAllMetrics(w http.ResponseWriter, _ *http.Request) {
	result := s.Storage.GetMappedByTypeAndNameMetrics()
	marshall, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(marshall)
}
