package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	KafkaEvents = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_events_total",
			Help: "Количество успешно отправленных событий Kafka",
		},
		[]string{"event_type"},
	)

	KafkaErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_errors_total",
			Help: "Количество ошибок при отправке событий Kafka",
		},
		[]string{"event_type"},
	)
)

func init() {
	prometheus.MustRegister(KafkaEvents)
	prometheus.MustRegister(KafkaErrors)
}
