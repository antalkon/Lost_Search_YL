package config

import (
	"auth/internal/authservice"
	"auth/internal/userrepo"

	cleanenv "github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	KafkaBroker    string `env:"KAFKA_BROKER" env-default:"localhost:9092"`
	RequestsTopic  string `env:"REQUESTS_TOPIC" env-default:"requests"`
	ResponsesTopic string `env:"RESPONSES_TOPIC" env-default:"responses"`

	authservice.AuthServiceConfig
	userrepo.AuthRepoConfig
}

func ReadFromFile(path string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
