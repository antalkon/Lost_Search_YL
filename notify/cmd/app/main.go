package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.crja72.ru/gospec/go21/go_final_project/config"
	"gitlab.crja72.ru/gospec/go21/go_final_project/internal/handler"
	"gitlab.crja72.ru/gospec/go21/go_final_project/internal/kafka"
	"gitlab.crja72.ru/gospec/go21/go_final_project/internal/logger"
)

func main() {
	log := logger.NewLogger()

	kafka.RegisterMetrics()

	cfg := config.LoadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	consumer := kafka.NewKafkaConsumer(cfg.KafkaBroker, cfg.RequestsTopic)
	defer consumer.Close()

	writer := kafka.NewKafkaWriter(cfg.KafkaBroker, cfg.NotifyResponsesTopic)
	defer writer.Close()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Info("Starting metrics server at :8080/metrics")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Metrics server failed: %v", err)
		}
	}()

	go func() {
		err := kafka.ListenMessages(ctx, consumer, func(msg []byte) error {
			log.Infof("Received message: %s", string(msg))
			if err := handler.HandleRequest(writer, msg); err != nil {
				log.Errorf("Error processing message: %v", err)
				kafka.ErrorMessages.Inc()
				return err
			}
			kafka.ProcessedMessages.Inc()
			return nil
		})
		if err != nil {
			log.Fatalf("Error consuming messages: %v", err)
		}
	}()

	select {
	case <-signalChan:
		log.Warn("Received termination signal. Shutting down...")
		cancel()
	}

	log.Info("Notify Service stopped gracefully.")
}
