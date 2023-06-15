package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ksusonic/go-devops-mon/internal/controllers"
	"github.com/ksusonic/go-devops-mon/internal/metrics"
	protopb "github.com/ksusonic/go-devops-mon/proto"
	metricspb "github.com/ksusonic/go-devops-mon/proto/metrics"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	protopb.UnimplementedServerServer

	GrpcServer  *grpc.Server
	logger      *zap.Logger
	storage     metrics.ServerMetricStorage
	hashService controllers.HashService
}

func NewServer(storage metrics.ServerMetricStorage, hashService controllers.HashService, logger *zap.Logger) *GrpcServer {
	grpcServer := grpc.NewServer()
	s := &GrpcServer{
		GrpcServer:  grpcServer,
		logger:      logger,
		storage:     storage,
		hashService: hashService,
	}
	protopb.RegisterServerServer(grpcServer, s)
	return s
}

func (s *GrpcServer) Start(port int) {
	addr := net.TCPAddr{Port: port}
	listen, err := net.Listen("tcp", addr.String())
	if err != nil {
		log.Fatal(err)
	}

	s.logger.Info("Listening grpc", zap.Int("port", port))
	if err := s.GrpcServer.Serve(listen); err != nil {
		s.logger.Fatal("error serving grpc", zap.Error(err))
	}
}

func (s *GrpcServer) UpdateMetric(ctx context.Context, in *protopb.UpdateMetricRequest) (*protopb.UpdateMetricResponse, error) {
	m := in.GetMetric()
	if err := s.hashService.ValidateHashProto(m); err != nil {
		return nil, fmt.Errorf("invalid metric hash")
	}

	if m.GetType() != metricspb.MetricType_gauge && m.GetType() != metricspb.MetricType_counter {
		return nil, fmt.Errorf("unexpected metric type: %s", m.GetType().String())
	} else {
		metric := metrics.FromProto(m)
		resultMetric, err := s.storage.SetMetric(ctx, &metric)
		if err != nil {
			s.logger.Error("error saving metric", zap.Error(err))
			return nil, fmt.Errorf("unexpected error")
		}
		s.logger.Debug("Updated", zap.String("metric", m.String()))

		return &protopb.UpdateMetricResponse{Metric: resultMetric.AsProto()}, nil
	}
}

func (s *GrpcServer) MetricUpdates(ctx context.Context, in *protopb.MetricUpdatesRequest) (*protopb.MetricUpdatesResponse, error) {

	for _, m := range in.GetMetrics() {
		if err := s.hashService.ValidateHashProto(m); err != nil {
			return nil, fmt.Errorf("invalid metric hash")
		}
	}
	var allmetrics []*metrics.Metric
	for _, m := range in.GetMetrics() {
		metric := metrics.FromProto(m)
		allmetrics = append(allmetrics, &metric)
	}
	err := s.storage.SetMetrics(ctx, allmetrics)
	if err != nil {
		s.logger.Error("error setting metric", zap.Error(err))
		return nil, fmt.Errorf("unexpected error")
	}
	return &protopb.MetricUpdatesResponse{}, nil
}

func (s *GrpcServer) GetValue(ctx context.Context, in *protopb.GetValueRequest) (*protopb.GetValueResponse, error) {
	value, err := s.storage.GetMetric(ctx, in.Type.String(), in.GetID())
	if err != nil {
		return &protopb.GetValueResponse{
			Response: &protopb.GetValueResponse_Error{Error: "not found metric"},
		}, nil
	}

	err = s.hashService.SetHash(value)
	if err != nil {
		s.logger.Error("error setting hash", zap.Error(err))
		return nil, fmt.Errorf("error setting hash")
	}

	return &protopb.GetValueResponse{
		Response: &protopb.GetValueResponse_Metric{Metric: value.AsProto()},
	}, nil
}

func (s *GrpcServer) GetAllMetrics(ctx context.Context, _ *protopb.GetAllMetricsRequest) (*protopb.GetAllMetricsResponse, error) {
	allMetrics, err := s.storage.GetAllMetrics(ctx)
	if err != nil {
		s.logger.Error("could not GetMappedByTypeAndNameMetrics", zap.Error(err))
		return &protopb.GetAllMetricsResponse{}, nil
	}

	result := protopb.GetAllMetricsResponse{}
	for i := range allMetrics {
		metric := allMetrics[i].AsProto()
		result.Mapping[metric.Type.String()].Value[metric.ID] = metric
	}

	return &result, nil
}
