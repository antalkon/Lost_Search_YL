package kafka

type BrokerConfig struct {
	port int `env:"KAFKA_PORT" envDefault:"8080"`
}
