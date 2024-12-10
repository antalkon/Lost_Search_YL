package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gitlab.crja72.ru/gospec/go21/go_final_project/pkg/db/postgres"
)

type Config struct {
	KafkaBroker          string
	RequestsTopic        string
	NotifyResponsesTopic string

	postgres.Config
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	return &Config{
		KafkaBroker:          os.Getenv("KAFKA_BROKER"),
		RequestsTopic:        os.Getenv("KAFKA_REQUESTS_TOPIC"),
		NotifyResponsesTopic: os.Getenv("KAFKA_NOTIFY_RESPONSES_TOPIC"), // Новое поле
	}
}
