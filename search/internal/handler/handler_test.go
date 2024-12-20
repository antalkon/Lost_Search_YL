package handler

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"gitlab.crja72.ru/gospec/go21/go_final_project/internal/repository"
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

func (t *TestRepo) AddFind(req repository.AddReq) (repository.AddResp, error) {
	t.Table = append(t.Table, "!")
	return repository.AddResp{}, nil
}

func (t *TestRepo) GetFind(req repository.GetReq) (repository.GetResp, error) {
	return repository.GetResp{}, nil
}

func (t *TestRepo) RespondToFind(req repository.RespondReq) (repository.RespondResp, error) {
	return repository.RespondResp{}, nil
}

func TestHandleAdd_Success(t *testing.T) {
	repo := &TestRepo{}
	writer := &TestKafkaWriter{}
	msgData := `{"user_login": "test", "name": "телефон", "description": "маленький", "type": "техника", "location": {"country": "Россия", "city": "Иваново", "district": "Фрунценский"}}`
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
