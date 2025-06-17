package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// KafkaEvents — счётчик успешно отправленных событий в Kafka,
// разбитый по типу события.
//
// Метрика: kafka_events_total
// Labels:
// - event_type: тип события Kafka (например, user.registered, user.logged_in)
var KafkaEvents = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "kafka_events_total",
		Help: "Количество успешно отправленных событий Kafka",
	},
	[]string{"event_type"},
)

// KafkaErrors — счётчик ошибок при отправке событий в Kafka,
// разбитый по типу события.
//
// Метрика: kafka_errors_total
// Labels:
// - event_type: тип события Kafka (например, user.registered, user.logged_in)
var KafkaErrors = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "kafka_errors_total",
		Help: "Количество ошибок при отправке событий Kafka",
	},
	[]string{"event_type"},
)

func init() {
	// Регистрируем метрики в Prometheus, чтобы они были доступны для сбора
	prometheus.MustRegister(KafkaEvents)
	prometheus.MustRegister(KafkaErrors)
}
