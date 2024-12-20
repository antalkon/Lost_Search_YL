package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gitlab.crja72.ru/gospec/go21/go_final_project/pkg/db/postgres"
	"gitlab.crja72.ru/gospec/go21/go_final_project/pkg/db/redis"
)

type Config struct {
	KafkaBroker          string
	RequestsTopic        string
	SearchResponsesTopic string

	postgres.PostgresConfig
	redis.RedisConfig
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	return &Config{
		KafkaBroker:          os.Getenv("KAFKA_BROKER"),
		RequestsTopic:        os.Getenv("KAFKA_REQUESTS_TOPIC"),
		SearchResponsesTopic: os.Getenv("KAFKA_NOTIFY_RESPONSES_TOPIC"), // Новое поле
		PostgresConfig:       postgres.LoadConfig(),
		RedisConfig:          redis.LoadConfig(),
	}
}
