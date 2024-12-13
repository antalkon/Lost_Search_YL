package handler

import (
	"auth/internal/kafka"
	"context"
	"encoding/json"
	"fmt"
	"log"
)

const ServiceAuth = "auth"

const (
	ActionCreateUser     = "create_user"
	ActionLoginUser      = "login_user"
	ActionValidateToken  = "validate_token"
	ActionGetUserData    = "get_user_data"
	ActionUpdateUserData = "update_user_data"
)

const (
	StatusSuccess = "success"
	StatusFailed  = "failed"
)

type Message struct {
	RequestID string         `json:"request_id"`
	Service   string         `json:"service"`
	Action    string         `json:"action"`
	Data      map[string]any `json:"data"`
}

type AuthServiceInterface interface {
	CreateUser(ctx context.Context,
		login string, password string, data string) (string, error)
	LoginUser(ctx context.Context, login string, password string) (string, error)
	IsTokenValid(tokenString string) bool
	GetUserData(ctx context.Context, login string) (string, error)
	UpdateUserData(ctx context.Context, login string, data string) (string, error)
}

func HandleRequest(ctx context.Context, writer kafka.KafkaWriter,
	authService AuthServiceInterface, message []byte) error {
	var (
		req  Message
		resp []byte
		data map[string]any
		err  error
	)

	req, err = formRequest(message)
	if err != nil {
		return err
	}

	log.Printf("Processing request: %+v", req)
	switch req.Action {
	case ActionCreateUser:
		data, err = createUserHandler(ctx, req.Data, authService)
	case ActionLoginUser:
		data, err = loginUserHandler(ctx, req.Data, authService)
	case ActionValidateToken:
		data, err = validateTokenHandler(req.Data, authService)
	case ActionGetUserData:
		data, err = getUserDataHandler(ctx, req.Data, authService)
	case ActionUpdateUserData:
		data, err = updateUserDataHandler(ctx, req.Data, authService)
	default:
		err = fmt.Errorf("unknown action: %s", req.Action)
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

	if req.Service != ServiceAuth {
		return Message{}, fmt.Errorf("message was expected for service `%s`, it came for: `%s`", ServiceAuth, req.Service)
	}
	return req, nil
}

func formResponse(reqID string, data map[string]any, err error) ([]byte, error) {
	var status string
	if err == nil {
		status = StatusSuccess
	} else {
		data = map[string]any{"error": err.Error()}
		status = StatusFailed
	}

	resp := map[string]interface{}{
		"request_id": reqID,
		"service":    ServiceAuth,
		"status":     status,
		"data":       data,
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to marshal response: %w", err)
	}
	return respBytes, nil
}
