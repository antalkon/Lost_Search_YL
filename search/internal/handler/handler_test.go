package handler

import (
	"context"
	"errors"
	"fmt"
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

type TestRepo struct {
	Table []string
}

func (t *TestRepo) CreateFind() {
	t.Table = append(t.Table, "!")
}

func (t *TestRepo) GetFind() {

}

func TestHandleAdd_Success(t *testing.T) {
	repo := &TestRepo{}
	writer := &TestKafkaWriter{}
	msgData := `{"name": "телефон", "description": "маленький", "type": "техника", "location": ["Россия", "Иваново", "пр.Фрунцзе"]}}`
	message := []byte(messageMarkup(ActionAdd, msgData))

	err := HandleRequest(writer, repo, message)

	assert.NoError(t, err, "HandleRequest should process message successfully")
	assert.Len(t, writer.Messages, 1, "One message should be sent")
}

func TestHandleRespond_Success(t *testing.T) {
	repo := &TestRepo{}
	writer := &TestKafkaWriter{}
	msgData := `{"find_uuid":"fa832e04-4022-4ac9-bf2d-f38c5cb5ccb8"}`
	message := []byte(messageMarkup(ActionRespond, msgData))

	err := HandleRequest(writer, repo, message)

	assert.NoError(t, err, "HandleRequest should process message successfully")
	assert.Len(t, writer.Messages, 1, "One message should be sent")
}

func messageMarkup(action string, data string) string {
	return fmt.Sprintf(`{"request_id": "12345", "service": "search", "action": "%s", "data": %s}`, action, data)
}
