package main

import (
	"context"
	"fmt"
	"gateway_service/internal/config"
	"gateway_service/internal/transport/gateway"
	"gateway_service/internal/transport/rest"
	"gateway_service/pkg/logger"
	"gateway_service/pkg/metrics"
	"os"
	"os/signal"
	"syscall"
)

const (
	serviceName = "gateway_service"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, logger.LoggerKey, logger.New(serviceName))
	cfg := config.New(ctx)
	restSrv, err := rest.New(ctx, cfg.RestServerPort)
	mainLogger := logger.GetLogger(ctx)

	if err != nil || restSrv == nil {
		mainLogger.Error(ctx, "Server creation failed")
		return
	}

	metricsSrv, err := metrics.NewMetricsServer(cfg.PrometheusConfig)

	if err != nil || metricsSrv == nil {
		mainLogger.Error(ctx, "Server creation failed")
		return
	}

	serv := gateway.NewServer(restSrv, metricsSrv)

	go func() {
		if err = serv.Start(ctx); err != nil {
			mainLogger.Error(ctx, "Server start failed")
			return
		}
	}()

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	<-graceCh

	if err = serv.Stop(); err != nil {
		mainLogger.Error(ctx, "Graceful shutdown failed")
	}
	fmt.Println("Server Stopped")
}
