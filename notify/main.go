package main

import (
	"context"
	"log"
	"notify-service/protos"

	"github.com/golang/protobuf/proto"
	"github.com/segmentio/kafka-go"
)

func consumeKafkaRequests(reader *kafka.Reader) {
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}
		log.Printf("Received message: %s", string(msg.Value))
		handleRequest(msg)
	}
}

func handleRequest(msg kafka.Message) {
	var request protos.KafkaRequest
	err := proto.Unmarshal(msg.Value, &request)
	if err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return
	}

	log.Printf("Processing request: %v", request)

	var response protos.KafkaResponse
	response.RequestId = request.RequestId
	response.Service = "notify"
	response.Status = "success"
	response.Message = "Notification sent successfully"

	if request.Action == "email" {
		log.Println("Sending email...")
		response.Data = map[string]string{"status": "Email sent"}
	} else if request.Action == "push" {
		log.Println("Sending push notification...")
		response.Data = map[string]string{"status": "Push notification sent"}
	} else {
		log.Println("Unknown action")
		response.Status = "error"
		response.Message = "Unknown action"
	}

	sendKafkaResponse(response)
}

func sendKafkaResponse(response protos.KafkaResponse) {
	writer := kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "notify-responses",
		Balancer: &kafka.LeastBytes{},
	}

	data, err := proto.Marshal(&response)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		return
	}

	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: data,
		},
	)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	} else {
		log.Printf("Response sent to Kafka")
	}
}

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "requests",
		GroupID: "notify-service",
	})

	go consumeKafkaRequests(reader)

	select {}
}
