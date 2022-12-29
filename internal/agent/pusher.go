package agent

import (
	"log"
	"strconv"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

const contentType = "text/plain"

func (m MetricCollector) makePushURL(metricName string, metricValue string) string {
	return "http://" + m.ServerHost + ":" + strconv.FormatInt(int64(m.ServerPort), 10) +
		"/update/" + metrics.GaugeType + "/" + metricName + "/" + metricValue
}

func (m MetricCollector) sendMetric(path string) {
	response, err := m.Client.Post(path, contentType, nil)
	if err != nil {
		log.Printf("Error sending push metric request: %v\n", err)
	} else {
		response.Body.Close()
	}
}

func (m MetricCollector) PushMetrics() {
	for _, metric := range m.Storage.GetAllMetrics() {
		var stringMetricValue string
		if metric.Type == metrics.GaugeType {
			stringMetricValue = strconv.FormatFloat(metric.Value.(float64), 'f', -1, 64)
		} else if metric.Type == metrics.CounterType {
			stringMetricValue = strconv.FormatInt(metric.Value.(int64), 10)
		} else {
			log.Printf("Unknown metric type: %s\n", metric.Type)
			// using float as default value
			stringMetricValue = strconv.FormatFloat(metric.Value.(float64), 'f', -1, 64)
		}
		var path = m.makePushURL(metric.Name, stringMetricValue)
		m.sendMetric(path)
	}
}
