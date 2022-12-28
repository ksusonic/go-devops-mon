package agent

import (
	"log"
	"strconv"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

const contentType = "text/plain"

func (m MetricCollector) makePushURL(metricName string, metricValue string) string {
	return "http://" + m.ServerHost + ":" + strconv.FormatInt(int64(m.ServerPort), 10) +
		"/update/" + metrics.GaugeName + "/" + metricName + "/" + metricValue
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
	for metricName, metricValue := range m.Storage.GaugeMetricStorage {
		stringMetricValue := strconv.FormatFloat(metricValue, 'f', -1, 64)
		var path = m.makePushURL(metricName, stringMetricValue)
		m.sendMetric(path)
	}
	for metricName, metricValue := range m.Storage.CounterMetricStorage {
		stringMetricValue := strconv.FormatInt(metricValue, 10)
		var path = m.makePushURL(metricName, stringMetricValue)
		m.sendMetric(path)
	}
}
