package gateway

import (
	"context"
	"gateway_service/internal/transport/rest"
	"gateway_service/pkg/metrics"
	"github.com/labstack/echo/v4"
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
	restServer    rest.Server
	metricsServer metrics.PrometheusServer
}

func NewServer(restSrv rest.Server, metricSrv metrics.PrometheusServer) *Server {
	return &Server{restServer: restSrv, metricsServer: metricSrv}
}

func (s *Server) Start(ctx context.Context) error {
	eg := errgroup.Group{}
	f := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			s.metricsServer.NewRequest()
			return next(c)
		}
	}
	eg.Go(func() error {
		return s.restServer.Start(ctx, f)
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
