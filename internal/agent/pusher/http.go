package pusher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"reflect"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

type HTTPPusher struct {
	PushURL        string
	Client         http.Client
	encryptService EncryptService
	Addr           net.IP
}

type EncryptService interface {
	EncryptBytes(b []byte) ([]byte, error)
}

func (h *HTTPPusher) SendMetric(metric *metrics.Metric) error {
	marshall, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("could not marshall %s: %v", metric.ID, err)
	}
	if h.encryptService != nil && !reflect.ValueOf(h.encryptService).IsNil() {
		marshall, err = h.encryptService.EncryptBytes(marshall)
		if err != nil {
			return fmt.Errorf("could not encrypt request %s: %v", metric.ID, err)
		}
	}

	r, err := http.NewRequest(http.MethodPost, h.PushURL, bytes.NewReader(marshall))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("X-Real-IP", h.Addr.String())

	response, err := h.Client.Do(r)
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
