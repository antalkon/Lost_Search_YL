package main

import (
	"context"

	"gitlab.crja72.ru/gospec/go21/go_final_project/config"
	"gitlab.crja72.ru/gospec/go21/go_final_project/internal/handler"
	"gitlab.crja72.ru/gospec/go21/go_final_project/internal/kafka"
	"gitlab.crja72.ru/gospec/go21/go_final_project/internal/logger"
	"gitlab.crja72.ru/gospec/go21/go_final_project/internal/repository"
)

func main() {
	log := logger.NewLogger()

	cfg := config.LoadConfig()

	repo, err := repository.NewPostgresRepo(cfg.Config)
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer := kafka.NewKafkaConsumer(cfg.KafkaBroker, cfg.RequestsTopic)
	defer consumer.Close()

	writer := kafka.NewKafkaWriter(cfg.KafkaBroker, cfg.NotifyResponsesTopic)
	defer writer.Close()

	go func() {
		err := kafka.ListenMessages(ctx, consumer, func(b []byte) error {
			log.Infof("received message: %s", string(b))
			if err := handler.HandleRequest(writer, repo, b); err != nil {
				log.Errorf("Error processing message: %v", err)
				// TODO: inc message
				return err
			}
			// TODO: inc error
			return nil
		})
		if err != nil {
			log.Fatalf("Error consuming messages: %v", err)
		}
	}()
}
