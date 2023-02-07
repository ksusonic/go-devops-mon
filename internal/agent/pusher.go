package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (m MetricCollector) PushMetrics() error {
	allMetrics := m.Storage.GetAllMetrics()
	for i := range allMetrics {
		err := m.HashService.SetHash(&allMetrics[i])
		if err != nil {
			log.Println("Could not set hash!")
		}
	}

	if len(allMetrics) == 0 {
		log.Println("Currently no metrics. Push skipped")
		return nil
	}

	marshall, err := json.Marshal(allMetrics)
	if err != nil {
		return fmt.Errorf("error marshalling metric: %v", err)
	}
	r, err := http.NewRequest(http.MethodPost, m.PushURL, bytes.NewReader(marshall))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	r.Header.Add("Content-Type", "application/json")

	response, err := m.Client.Do(r)
	if err != nil {
		return fmt.Errorf("error sending push metrics request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		readBody, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("status %s while sending metrics: %v", response.Status, err)
		}
		log.Printf("status %s while sending metrics on \"updates\" path : %s\n", response.Status, string(readBody))
	}

	return nil
}
