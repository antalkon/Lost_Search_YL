package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	KafkaBroker          string
	RequestsTopic        string
	NotifyResponsesTopic string
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
