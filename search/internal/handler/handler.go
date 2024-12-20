package handler

import (
	"encoding/json"
	"fmt"
	"log"

	"gitlab.crja72.ru/gospec/go21/go_final_project/internal/handler/methods"
	"gitlab.crja72.ru/gospec/go21/go_final_project/internal/kafka"
	"gitlab.crja72.ru/gospec/go21/go_final_project/internal/repository"
)

const (
	ServiceSearch string = "search"
)

const (
	ActionAdd     string = "add"
	ActionRespond string = "respond"
	ActionGet     string = "get"
)

const (
	StatusSuccess string = "success"
	StatusFailed  string = "failed"
)

type Message struct {
	RequestID string         `json:"request_id"`
	Service   string         `json:"service"`
	Action    string         `json:"action"`
	Data      map[string]any `json:"data"`
}

func HandleRequest(writer kafka.KafkaWriter, repo repository.Repository, message []byte) error {
	var (
		req  Message
		resp []byte
		data any
		err  error
	)

	req, err = formRequest(message)
	if err != nil {
		return err
	}

	log.Printf("Processing request: %+v", req)
	switch req.Action {
	case ActionAdd:
		data, err = methods.HandleAdd(req.Data, repo)
	case ActionRespond:
		data, err = methods.HandleRespond(req.Data, repo)
	case ActionGet:
		data, err = methods.HandleGet(req.Data, repo)
	}

	resp, err = formResponse(req.RequestID, data, err)
	if err != nil {
		return err
	}
	return kafka.SendMessage(writer, []byte(req.RequestID), resp)
}

func formRequest(msg []byte) (Message, error) {
	var req Message
	if err := json.Unmarshal(msg, &req); err != nil {
		return Message{}, fmt.Errorf("failed to parse message: %w", err)
	}

	if req.Service != ServiceSearch {
		return Message{}, fmt.Errorf("message was expected for service `%s`, it came for: `%s`", ServiceSearch, req.Service)
	}
	return req, nil
}

func formResponse(reqID string, data any, err error) ([]byte, error) {
	var status string
	if err == nil {
		status = StatusSuccess
	} else {
		status = StatusFailed
	}

	resp := map[string]interface{}{
		"request_id": reqID,
		"service":    ServiceSearch,
		"status":     status,
		"data":       data,
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to marshal response: %w", err)
	}
	return respBytes, nil
}
