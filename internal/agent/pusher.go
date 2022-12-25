package agent

import (
	"github.com/ksusonic/go-devops-mon/internal/metrics"
	"log"
	"net/http"
	"strconv"
)

const contentType = "text/plain"

var client = http.Client{}

func (m MetricCollector) PushMetrics() {
	for metricName, metricValue := range m.Storage.GaugeMetricStorage {
		var path = "http://" + m.ServerHost + ":" + strconv.FormatInt(int64(m.ServerPort), 10) + "/update/" +
			metrics.GaugeName + "/" +
			metricName + "/" +
			strconv.FormatFloat(metricValue, 'f', -1, 64)

		response, err := client.Post(path, contentType, nil)
		if err != nil {
			log.Printf("Error sending post request: %v\n", err)
			continue
		}
		response.Body.Close()
	}
	for metricName, metricValue := range m.Storage.CounterMetricStorage {
		var path = "http://" + m.ServerHost + ":" + strconv.FormatInt(int64(m.ServerPort), 10) + "/update/" +
			metrics.CounterName + "/" +
			metricName + "/" +
			strconv.FormatInt(metricValue, 10)

		response, err := client.Post(path, contentType, nil)
		if err != nil {
			log.Printf("Error sending post request: %v\n", err)
			continue
		}
		response.Body.Close()
	}
}
