package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// KafkaWriter интерфейс для отправки сообщений в Kafka
type KafkaWriter interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

// NewKafkaWriter создает новый Kafka Writer
func NewKafkaWriter(broker, topic string) KafkaWriter {
	return &kafka.Writer{
		Addr:  kafka.TCP(broker),
		Topic: topic,
	}
}

// SendMessage отправляет сообщение в Kafka
func SendMessage(writer KafkaWriter, key, message []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: message,
	}
	return writer.WriteMessages(context.Background(), msg)
}
