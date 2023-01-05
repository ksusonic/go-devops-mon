package agent

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

func (m MetricCollector) sendMetric(metric metrics.Metrics) {
	marshall, err := json.Marshal(metric)
	if err != nil {
		log.Printf("Error marshalling metric: %v\n", err)
	}
	r, _ := http.NewRequest(http.MethodPost, m.ServerURL+"update/", bytes.NewReader(marshall))
	r.Header.Add("Content-Type", "application/json")

	response, err := m.Client.Do(r)
	if err != nil {
		log.Printf("Error sending push metric request: %v\n", err)
	} else {
		if response.StatusCode != http.StatusOK {
			readBody, err := io.ReadAll(response.Body)
			if err != nil {
				log.Printf("status %s while sending metric\n", response.Status)
			} else {
				log.Printf("status %s while sending metric on \"update/\" path : %s\n", response.Status, string(readBody))
			}
		}
		response.Body.Close()
	}
}

func (m MetricCollector) PushMetrics() {
	for _, metric := range m.Storage.GetAllMetrics() {
		m.sendMetric(metric)
	}
}
