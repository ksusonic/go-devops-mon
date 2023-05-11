package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sync"

	"go.uber.org/zap"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

func (m MetricCollector) PushMetrics() {
	metricCh := make(chan metrics.Metrics)
	wg := sync.WaitGroup{}

	for i := 0; i < m.RateLimit; i++ {
		go func() {
			for metric := range metricCh {
				wg.Done()
				err := m.sendMetric(metric)
				if err != nil {
					m.Logger.Error("error pushing metric", zap.Error(err))
				}
			}
		}()
	}

	allMetrics := m.Storage.GetAllMetrics()
	for _, metric := range allMetrics {
		wg.Add(1)
		metricCh <- metric
	}
	wg.Wait()
}

func (m MetricCollector) sendMetric(metric metrics.Metrics) error {
	err := m.hashService.SetHash(&metric)
	if err != nil {
		return fmt.Errorf("could not set hash for metric %s: %v", metric.ID, err)
	}
	marshall, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("could not marshall %s: %v", metric.ID, err)
	}
	if m.encryptService != nil && !reflect.ValueOf(m.encryptService).IsNil() {
		marshall, err = m.encryptService.EncryptBytes(marshall)
		if err != nil {
			return fmt.Errorf("could not encrypt request %s: %v", metric.ID, err)
		}
	}

	r, err := http.NewRequest(http.MethodPost, m.pushURL, bytes.NewReader(marshall))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("X-Real-IP", m.currentIP.String())

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
