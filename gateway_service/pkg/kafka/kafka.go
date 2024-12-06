package kafka

type Config struct {
	port int `env:"KAFKA_PORT" envDefault:"8080"`
}
