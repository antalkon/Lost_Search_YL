package kafka

import (
	"context"
	"gateway_service/pkg/logger"
	"gateway_service/pkg/syncmap"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type BrokerConfig struct {
	Host string `env:"KAFKA_HOST" envDefault:"kafka"`
	Port int    `env:"KAFKA_PORT" envDefault:"8080"`
}

type Consumer struct {
	reader   *kafka.Reader
	requests *syncmap.SyncMap
}

func NewConsumer(address string, topic string, requests *syncmap.SyncMap) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{address},
		Topic:   topic,
	})

	return &Consumer{reader: r, requests: requests}
}

func (c *Consumer) Consume(ctx context.Context) error {
	log := logger.GetLogger(ctx)
	log.Info(ctx, "Started consume")
	for {
		select {
		case <-ctx.Done():
			log.Info(ctx, "consumer stopped")
			return nil
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Error(ctx, "error", zap.String("Logging error", err.Error()))
			} else {
				uuid := string(msg.Key)
				ch, ok := c.requests.Read(uuid)
				if !ok {
					continue
				}
				c.requests.Delete(uuid)
				ch <- msg.Value
			}
		}
	}
}
