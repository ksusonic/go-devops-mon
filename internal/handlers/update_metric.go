package handlers

import (
	"fmt"
	"github.com/ksusonic/go-devops-mon/internal/metrics"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const UpdateHandlerName = "/update/"

var GaugeStorage = make(metrics.GaugeStorage)
var CounterStorage = make(metrics.CounterStorage)

type updateRequest struct {
	Type  string
	Name  string
	Value interface{}
}

func parseUpdateURL(url string) (*updateRequest, *error) {
	var args = strings.Split(strings.TrimPrefix(url, UpdateHandlerName), "/")
	if len(args) < 3 {
		err := fmt.Errorf("incorrect url: expected \"/update/<type>/<name>/<value>\"")
		return nil, &err
	}
	result := updateRequest{}

	if args[0] == metrics.GaugeName {
		result.Type = metrics.GaugeName
		parsed, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			err := fmt.Errorf("error parsing value: %s", err)
			return nil, &err
		}
		result.Value = parsed
	} else if args[0] == metrics.CounterName {
		result.Type = metrics.CounterName
		parsed, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			err := fmt.Errorf("error parsing value: %s", err)
			return nil, &err
		}
		result.Value = parsed
	} else {
		err := fmt.Errorf("no such metric type")
		return nil, &err
	}
	if !metrics.HasMetricName(args[1]) {
		err := fmt.Errorf("no such metric name")
		return nil, &err
	}
	result.Name = args[1]

	return &result, nil
}

// UpdateMetric — обработчик обновления метрики по типу и названию
func UpdateMetric(w http.ResponseWriter, r *http.Request) {
	requestData, err := parseUpdateURL(r.URL.Path)
	if err != nil {
		log.Printf("%v", *err)
		return
	}
	if requestData.Type == metrics.GaugeName {
		GaugeStorage[requestData.Name] = append(GaugeStorage[requestData.Name], requestData.Value.(float64))
		log.Printf("Updated gauge metrics: %f\n", requestData.Value.(float64))
	} else if requestData.Type == metrics.CounterName {
		CounterStorage[requestData.Name] = append(CounterStorage[requestData.Name], requestData.Value.(int64))
		log.Printf("Updated counter metrics: %d\n", requestData.Value.(int64))
	} else {
		log.Println("Unknown metric name!")
		w.WriteHeader(404)
	}

}
