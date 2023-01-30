package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

const contentType = "application/json"

func (m MetricCollector) sendMetric(metric metrics.Metrics) error {
	marshall, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("error marshalling metric: %v", err)
	}
	r, _ := http.NewRequest(http.MethodPost, m.PushURL, bytes.NewReader(marshall))
	r.Header.Add("Content-Type", contentType)

	response, err := m.Client.Do(r)
	if err != nil {
		return fmt.Errorf("error sending push metric %s request: %v", metric.ID, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		readBody, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("status %s while sending metric %s", response.Status, metric.ID)
		}
		log.Printf("status %s while sending metric %s on \"update\" path : %s\n", response.Status, metric.ID, string(readBody))
	}
	return nil
}

func (m MetricCollector) PushMetrics() error {
	var accumulatedErrs string
	for _, metric := range m.Storage.GetAllMetrics() {
		err := m.sendMetric(metric)
		if err != nil {
			accumulatedErrs += err.Error() + "\n"
		}
	}

	if accumulatedErrs != "" {
		return fmt.Errorf(accumulatedErrs)
	}
	return nil
}
