package config

import (
	"auth/internal/authservice"
	"auth/internal/userrepo"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	KafkaBroker    string `env:"KAFKA_BROKER" env-default:"localhost:9092"`
	RequestsTopic  string `env:"REQUESTS_TOPIC" env-default:"request"`
	ResponsesTopic string `env:"RESPONSES_TOPIC" env-default:"responses"`

	authservice.AuthServiceConfig
	userrepo.AuthRepoConfig
}

func ReadFromFile() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
