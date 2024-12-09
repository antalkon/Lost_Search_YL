package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(address string, topic string) (*Producer, error) {
	writer := kafka.Writer{
		Addr:  kafka.TCP(address),
		Topic: topic,
	}
	return &Producer{writer: &writer}, nil
}

func (p *Producer) SendMessage(ctx context.Context, msg string) error {
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Value: []byte(msg),
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *Producer) SearchAds(uuid, msg string) error {
	err := p.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(uuid),
		Value: []byte(msg),
	})
	if err != nil {
		return err
	}
	return nil
}
