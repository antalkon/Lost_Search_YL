package config

import (
	"context"
	"gateway_service/pkg/kafka"
	"gateway_service/pkg/logger"
	"gateway_service/pkg/metrics"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	kafka.BrokerConfig
	metrics.PrometheusConfig
	RestServerPort int `env:"REST_PORT" envDefault:"7070"`
}

func New(ctx context.Context) *Config {
	cfg := &Config{}
	log := logger.GetLogger(ctx)
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Error(ctx, err.Error())
	}
	return cfg
}
