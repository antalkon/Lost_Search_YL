package kafka

import (
	"context"
	"gateway_service/pkg/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type BrokerConfig struct {
	port int `env:"KAFKA_PORT" envDefault:"8080"`
}

type Consumer struct {
	reader   *kafka.Reader
	requests *map[string]chan string
}

func NewConsumer(address string, groupid string, topic string, requests *map[string]chan string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{address},
		GroupID: groupid,
		Topic:   topic,
	})
	return &Consumer{reader: r, requests: requests}
}

func (c *Consumer) Consume(ctx context.Context) error {
	log := logger.GetLogger(ctx)
	select {
	case <-ctx.Done():
		log.Info(ctx, "consumer stopped")
	default:
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Error(ctx, "error", zap.String("Logging error", err.Error()))
		} else {
			// TODO get uuid from msg
			uuid := string(msg.Value)
			(*c.requests)[uuid] <- string(msg.Value)
		}
	}
	return nil
}
