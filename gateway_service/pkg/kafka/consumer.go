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
	reader *kafka.Reader
}

func NewBroker(address string, groupid string, topic string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{address},
		GroupID: groupid,
		Topic:   topic,
	})
	return &Consumer{reader: r}
}

func (c *Consumer) Consume(ctx context.Context, handler func(message kafka.Message)) error {
	log := logger.GetLogger(ctx)
	select {
	case <-ctx.Done():
		log.Info(ctx, "consumer stopped")
	default:
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Error(ctx, "error", zap.String("Logging error", err.Error()))
		} else {
			handler(msg)
		}
	}
	return nil
}
