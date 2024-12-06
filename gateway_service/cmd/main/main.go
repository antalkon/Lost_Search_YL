package main

import (
	"context"
	"gateway_service/internal/config"
	"gateway_service/internal/transport/rest"
	"gateway_service/pkg/logger"
)

const (
	serviceName = "gateway_service"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, logger.LoggerKey, logger.New(serviceName))
	cfg := config.New(ctx)
	serv, err := rest.New(ctx, cfg.RestServerPort)
	log := logger.GetLogger(ctx)
	if err != nil || serv == nil {
		log.Error(ctx, "Server creation failed")
		return
	}
	if err = serv.Start(ctx); err != nil {
		log.Error(ctx, "Server start failed")
		return
	}
	
}
