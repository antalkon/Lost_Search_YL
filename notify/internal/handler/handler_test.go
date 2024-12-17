package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

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

func TestHandleRequest_Success(t *testing.T) {
	writer := &TestKafkaWriter{}
	message := []byte(`{"request_id": "12345", "service": "notify", "action": "test-action", "data": {"key": "value"}}`)

	err := HandleRequest(writer, message)

	assert.NoError(t, err, "HandleRequest should process message successfully")
	assert.Len(t, writer.Messages, 1, "One message should be sent")
	assert.Equal(t, []byte("12345"), writer.Messages[0].Key)
}

func TestHandleRequest_InvalidMessage(t *testing.T) {
	writer := &TestKafkaWriter{}
	invalidMessage := []byte(`invalid-json`)

	err := HandleRequest(writer, invalidMessage)

	assert.Error(t, err, "HandleRequest should return an error for invalid JSON")
}

func TestHandleRequest_FailureToSend(t *testing.T) {
	writer := &TestKafkaWriter{Fail: true}
	message := []byte(`{"request_id": "12345", "service": "notify", "action": "test-action", "data": {"key": "value"}}`)

	err := HandleRequest(writer, message)

	assert.Error(t, err, "HandleRequest should return an error if Kafka write fails")
}
