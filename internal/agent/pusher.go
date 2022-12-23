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
	for name, metric := range m.GaugeMetrics {
		var lastValue = metric[len(metric)-1]

		var path = "http://" + m.ServerHost + ":" + strconv.FormatInt(int64(m.ServerPort), 10) + "/update/" +
			metrics.GaugeName + "/" +
			name + "/" +
			strconv.FormatFloat(lastValue, 'E', -1, 64)
		// TODO: sent only last value of metric

		response, err := client.Post(path, contentType, nil)
		if err != nil {
			log.Printf("Error sending post request: %v\n", err)
		}
		response.Body.Close()
	}
	for name, metric := range m.CounterMetrics {
		var lastValue = metric[len(metric)-1]
		var path = "http://" + m.ServerHost + ":" + strconv.FormatInt(int64(m.ServerPort), 10) + "/update/" +
			metrics.CounterName + "/" +
			name + "/" +
			strconv.FormatInt(lastValue, 10)
		// TODO: sent only last value of metric

		response, err := client.Post(path, contentType, nil)
		if err != nil {
			log.Printf("Error sending post request: %v\n", err)
		}
		response.Body.Close()
	}
}
