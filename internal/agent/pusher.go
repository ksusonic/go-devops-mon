package agent

import (
	"log"
	"strconv"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

const contentType = "text/plain"

func (m MetricCollector) makePushURL(metricName, metricType, metricValue string) string {
	return "http://" + m.ServerHost + ":" + strconv.FormatInt(int64(m.ServerPort), 10) +
		"/update/" + metricType + "/" + metricName + "/" + metricValue
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
		if metric.MType == metrics.GaugeMType {
			stringMetricValue = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
		} else if metric.MType == metrics.CounterMType {
			stringMetricValue = strconv.FormatInt(*metric.Delta, 10)
		} else {
			log.Printf("Unknown metric type: %s\n", metric.MType)
		}
		var path = m.makePushURL(metric.ID, metric.MType, stringMetricValue)
		m.sendMetric(path)
	}
}
