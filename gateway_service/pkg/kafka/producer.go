package kafka

import (
	"context"
	"fmt"
	"gateway_service/pkg/logger"
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

func (p *Producer) SendMessage(ctx context.Context, uuid, msg string) error {
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(uuid),
		Value: []byte(msg),
	})
	logger.GetLogger(ctx).Info(ctx, fmt.Sprintf("send msg key: %v  \n value: %v", uuid, msg))
	if err != nil {
		return err
	}
	return nil
}
