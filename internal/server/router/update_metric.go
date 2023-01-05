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
	registerHandler("POST", "/update/", updateMetricHandler)
	registerHandler("POST", "/update/{type}/{name}/{value}", updateOneMetricHandler)
}

type updateRequest struct {
	Type     string
	Name     string
	RawValue string
}

// updateMetricHandler — updates metric by Metrics data in body
var updateMetricHandler = func(w http.ResponseWriter, r *http.Request, c context) {
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

	if m.MType != metrics.GaugeMType && m.MType != metrics.CounterMType {
		w.WriteHeader(http.StatusNotImplemented)
		log.Printf("Unknown metric type %s\n", m.MType)
	} else {
		resultMetric := (*c.storage).SetMetric(m)
		log.Printf("Updated %s\n", m.String())

		marshal, err := json.Marshal(resultMetric)
		if err != nil {
			log.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(marshal)
	}
}

// updateOneMetricHandler — updates metric by type, name and value
var updateOneMetricHandler = func(w http.ResponseWriter, r *http.Request, c context) {
	requestData := updateRequest{
		Type:     chi.URLParam(r, "type"),
		Name:     chi.URLParam(r, "name"),
		RawValue: chi.URLParam(r, "value"),
	}
	if requestData.Type == metrics.GaugeMType {
		value, err := strconv.ParseFloat(requestData.RawValue, 64)
		if err != nil {
			log.Printf("Incorrect value: %s\n", requestData.RawValue)
			w.WriteHeader(http.StatusBadRequest)
		}
		(*c.storage).SetMetric(metrics.Metrics{
			ID:    requestData.Name,
			MType: requestData.Type,
			Value: &value,
		})
		log.Printf("Updated gauge %s: %f\n", requestData.Name, value)
	} else if requestData.Type == metrics.CounterMType {
		value, err := strconv.ParseInt(requestData.RawValue, 10, 64)
		if err != nil {
			log.Printf("Incorrect value: %s\n", requestData.RawValue)
			w.WriteHeader(http.StatusBadRequest)
		}
		(*c.storage).SetMetric(metrics.Metrics{
			ID:    requestData.Name,
			MType: requestData.Type,
			Delta: &value,
		})
		log.Printf("Updated counter %s: %d\n", requestData.Name, value)
	} else {
		log.Println("unexpected metric type!")
		w.WriteHeader(http.StatusNotImplemented)
	}
}
