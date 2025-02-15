package kafka

import (
	"context"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

// NewKafkaConsumer создает новый Kafka Consumer
func NewKafkaConsumer(broker, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
	})
}

// ListenMessages читает сообщения из Kafka и передает их обработчику
func ListenMessages(ctx context.Context, consumer *kafka.Reader, handler func([]byte) error) error {
	log.Println("Consume started")
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Kafka Consumer...")
			return nil
		default:
			msg, err := consumer.ReadMessage(ctx)
			log.Println("Consume message ", string(msg.Value))
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			if err := handler(msg.Value); err != nil {
				log.Printf("Handler error: %v", err)
			}
		}
	}

}
