package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func (m MetricCollector) PushMetrics() error {
	allMetrics := m.Storage.GetAllMetrics()
	for i := range allMetrics {
		err := m.HashService.SetHash(&allMetrics[i])
		if err != nil {
			m.Logger.Error("could not set hash", zap.String("id", allMetrics[i].ID), zap.Error(err))
		}
	}

	if len(allMetrics) == 0 {
		m.Logger.Debug("currently no metrics. push skipped")
		return nil
	}

	marshall, err := json.Marshal(allMetrics)
	if err != nil {
		return fmt.Errorf("error marshalling metric: %v", err)
	}
	r, err := http.NewRequest(http.MethodPost, m.pushURL, bytes.NewReader(marshall))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	r.Header.Add("Content-Type", "application/json")

	response, err := m.client.Do(r)
	if err != nil {
		return fmt.Errorf("error sending push metrics request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		readBody, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("status %s while sending metrics: %v", response.Status, err)
		}
		return fmt.Errorf("status %s while sending metrics on \"updates\" path: %s", response.Status, string(readBody))
	}

	return nil
}
