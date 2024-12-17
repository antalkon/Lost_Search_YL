package kafka

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ProcessedMessages = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "notify_processed_messages_total",
			Help: "Total number of processed messages",
		},
	)
	ErrorMessages = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "notify_error_messages_total",
			Help: "Total number of messages that failed to process",
		},
	)
)

// RegisterMetrics регистрирует метрики в Prometheus
func RegisterMetrics() {
	prometheus.MustRegister(ProcessedMessages)
	prometheus.MustRegister(ErrorMessages)
}
