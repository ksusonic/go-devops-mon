package pusher

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
	"google.golang.org/protobuf/proto"
)

type GrpcPusher struct {
	PushURL string
	Client  http.Client
	Addr    net.IP
}

func (g *GrpcPusher) SendMetric(metric *metrics.Metric) error {
	marshall, err := proto.Marshal(metric.AsProto())
	if err != nil {
		return fmt.Errorf("could not marshall %s: %v", metric.ID, err)
	}

	r, err := http.NewRequest(http.MethodPost, g.PushURL, bytes.NewReader(marshall))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	r.Header.Add("Content-Type", "application/protobuf")
	r.Header.Add("X-Real-IP", g.Addr.String())

	response, err := g.Client.Do(r)
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
