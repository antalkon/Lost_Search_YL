package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"gitlab.crja72.ru/gospec/go21/go_final_project/config"
	"gitlab.crja72.ru/gospec/go21/go_final_project/internal/handler"
	"gitlab.crja72.ru/gospec/go21/go_final_project/internal/kafka"
	"gitlab.crja72.ru/gospec/go21/go_final_project/internal/logger"
	"gitlab.crja72.ru/gospec/go21/go_final_project/internal/repository"
)

func main() {
	log := logger.NewLogger()

	cfg := config.LoadConfig()

	log.Printf("%+v\n", cfg.RedisConfig)
	repo, err := repository.NewCachedRepo(cfg.PostgresConfig, cfg.RedisConfig)
	if err != nil {
		log.Fatalln(err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer := kafka.NewKafkaConsumer(cfg.KafkaBroker, cfg.RequestsTopic)
	defer consumer.Close()

	writer := kafka.NewKafkaWriter(cfg.KafkaBroker, cfg.SearchResponsesTopic)
	defer writer.Close()

	// if _, err := repo.AddFind(repository.AddReqMok()); err != nil {
	// 	log.Fatalln(err)
	// }
	if r, err := repo.RespondToFind(repository.RespondReq{FindUUID: "66fafbcd-d402-42d4-9301-194b3142e092"}); err != nil {
		log.Fatalln(err)
	} else {
		log.Println(r)
	}
	t := "тел"
	if r, err := repo.GetFind(repository.GetReq{Name: &t}); err != nil {
		log.Fatalln(err)
	} else {
		for _, f := range r.Finds {
			log.Println(f.Name, f.Type, f.Description, f.Location)
		}
	}
	if r, err := repo.GetFind(repository.GetReq{Name: &t}); err != nil {
		log.Fatalln(err)
	} else {
		for _, f := range r.Finds {
			log.Println(f.Name, f.Type, f.Description, f.Location)
		}
	}

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

	<-signalChan
	log.Warn("Received termination signal. Shutting down...")
	cancel()

	log.Info("Search Service stopped gracefully.")
}
