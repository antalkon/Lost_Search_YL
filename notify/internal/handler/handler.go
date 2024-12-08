package handler

import (
	"encoding/json"
	"fmt"
	"log"

	"gitlab.crja72.ru/gospec/go21/go_final_project/internal/kafka"
)

type Message struct {
	RequestID string                 `json:"request_id"`
	Service   string                 `json:"service"`
	Action    string                 `json:"action"`
	Data      map[string]interface{} `json:"data"`
}

func HandleRequest(writer kafka.KafkaWriter, message []byte) error {
	var req Message
	if err := json.Unmarshal(message, &req); err != nil {
		return fmt.Errorf("failed to parse message: %w", err)
	}

	log.Printf("Processing request: %+v", req)

	response := map[string]interface{}{
		"request_id": req.RequestID,
		"service":    "notify",
		"status":     "success",
		"data":       fmt.Sprintf("Processed action: %s", req.Action),
	}

	respBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	return kafka.SendMessage(writer, []byte(req.RequestID), respBytes)
}
