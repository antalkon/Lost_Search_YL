package kafka

import (
	"context"
	"errors"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

// Тестовый продюсер для моков
type TestKafkaWriter struct {
	Messages []kafka.Message
	Fail     bool
}

func (t *TestKafkaWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	if t.Fail {
		return errors.New("failed to write message")
	}
	t.Messages = append(t.Messages, msgs...)
	return nil
}

func (t *TestKafkaWriter) Close() error {
	return nil
}

// Тест на успешную отправку сообщения
func TestSendMessage_Success(t *testing.T) {
	writer := &TestKafkaWriter{}
	err := SendMessage(writer, []byte("key"), []byte("test message"))

	assert.NoError(t, err, "SendMessage should not return an error")
	assert.Len(t, writer.Messages, 1, "One message should be sent")
	assert.Equal(t, []byte("key"), writer.Messages[0].Key)
	assert.Equal(t, []byte("test message"), writer.Messages[0].Value)
}

// Тест на ошибку при отправке сообщения
func TestSendMessage_Failure(t *testing.T) {
	writer := &TestKafkaWriter{Fail: true}
	err := SendMessage(writer, []byte("key"), []byte("test message"))

	assert.Error(t, err, "SendMessage should return an error")
}
