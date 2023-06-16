package pusher

import (
	"context"
	"net"
	"net/http"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
	pb "github.com/ksusonic/go-devops-mon/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcPusher struct {
	PushURL string
	Client  http.Client
	Addr    net.IP

	client pb.ServerClient
}

func (g *GrpcPusher) Connect() (func(), error) {
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	g.client = pb.NewServerClient(conn)
	return func() {
		g.client = nil
		conn.Close()
	}, nil
}

func (g *GrpcPusher) SendMetric(ctx context.Context, metric *metrics.Metric) error {
	_, err := g.client.UpdateMetric(ctx, &pb.UpdateMetricRequest{Metric: metric.AsProto()})
	if err != nil {
		return err
	}

	return nil
}
