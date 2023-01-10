package agent

import (
	"log"
	"net/url"
	"strconv"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

const contentType = "text/plain"

func (m MetricCollector) makePushURL(metricName, metricType, metricValue string) string {
	u := url.URL{
		Scheme: m.ServerRequestScheme,
		Host:   m.ServerHost + ":" + strconv.FormatInt(int64(m.ServerPort), 10),
		Path:   "/update/" + metricType + "/" + metricName + "/" + metricValue,
	}
	return u.String()
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
			log.Fatalf("Unknown metric type: %s\n", metric.Type)
		}
		var path = m.makePushURL(metric.Name, metric.Type, stringMetricValue)
		m.sendMetric(path)
	}
}
