package metrics

import (
	"context"
	"fmt"
	"gateway_service/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

type PrometheusConfig struct {
	port int `env:"PROMETEUS_PORT" envDefault:"9090"`
}

type PrometheusServer struct {
	httpServer *http.Server
}

var (
	totalrequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "gateway_processed_req_total",
		Help: "The total number of processed requests",
	})
)

func NewMetricsServer(cfg PrometheusConfig) (*PrometheusServer, error) {
	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "hash_duration_seconds",
		Help: "Time taken to create hashes",
	}, []string{"code"})

	r := http.NewServeMux()
	r.Handle("/metrics", promhttp.Handler())

	err := prometheus.Register(histogram)
	if err != nil {
		return nil, err
	}

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.port),
		ReadTimeout:    8 * time.Second,
		WriteTimeout:   8 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        r,
	}
	return &PrometheusServer{httpServer: s}, nil
}

func (s *PrometheusServer) NewRequest() {
	totalrequests.Inc()
}

func (s *PrometheusServer) Start(ctx context.Context) error {
	log := logger.GetLogger(ctx)
	if err := s.httpServer.ListenAndServe(); err != nil {
		log.Error(ctx, "Error starting metrics server")
		return err
	}
	log.Info(ctx, "Metrics server started")
	return nil
}

func (s *PrometheusServer) Stop(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
