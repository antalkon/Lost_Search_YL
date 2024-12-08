package gateway

import (
	"context"
	"golang.org/x/sync/errgroup"
)

type RestServer interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type MetricsServer interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	NewRequest()
}

type Server struct {
	restServer    RestServer
	metricsServer MetricsServer
}

func NewServer(restSrv RestServer, metricSrv MetricsServer) *Server {
	return &Server{restServer: restSrv, metricsServer: metricSrv}
}

func (s *Server) Start(ctx context.Context) error {
	eg := errgroup.Group{}
	eg.Go(func() error {
		return s.restServer.Start(ctx)
	})
	eg.Go(func() error {
		return s.metricsServer.Start(ctx)
	})

	return eg.Wait()
}

func (s *Server) Stop() error {
	err := s.restServer.Stop(nil)
	if err != nil {
		return err
	}
	return s.metricsServer.Stop(nil)
}
