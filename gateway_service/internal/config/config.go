package config

import (
	"context"
	"gateway_service/pkg/kafka"
	"gateway_service/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	kafka.Config
	RestServerPort int `env:"REST_PORT" envDefault:"9090"`
}

func New(ctx context.Context) *Config {
	cfg := &Config{}
	log := logger.GetLogger(ctx)
	err := cleanenv.ReadConfig("../../config/local.env", cfg)
	if err != nil {
		log.Error(ctx, err.Error())
	}
	return cfg
}
