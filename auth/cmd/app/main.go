package main

import (
	"auth/internal/authservice"
	"auth/internal/config"
	"auth/internal/handler"
	"auth/internal/kafka"
	"auth/internal/userrepo"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//TODO: логи логи логи логи

	cfg, err := config.ReadFromFile("config/auth.env")

	if err != nil {
		panic(err)
	}

	repo, err := userrepo.NewUserRepo(cfg.AuthRepoConfig, nil) //TODO: редиска
	if err != nil {
		panic(err)
	}

	as := authservice.NewAuthService(cfg.AuthServiceConfig, repo)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer := kafka.NewKafkaConsumer(cfg.KafkaBroker, cfg.RequestsTopic)
	defer consumer.Close()

	writer := kafka.NewKafkaWriter(cfg.KafkaBroker, cfg.ResponsesTopic)
	defer writer.Close()

	go func() {
		err := kafka.ListenMessages(ctx, consumer, func(b []byte) error {
			log.Printf("received message: %s", string(b))
			if err := handler.HandleRequest(ctx, writer, as, b); err != nil {
				log.Printf("Error processing message: %v", err)
				return err
			}
			return nil
		})
		if err != nil {
			log.Fatalf("Error consuming messages: %v", err)
		}
	}()

	<-signalChan
	log.Print("Received termination signal. Shutting down...")
	cancel()

	log.Println("Auth Service stopped.")

}
