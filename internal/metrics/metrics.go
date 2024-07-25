package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

type Metrics interface {
	ObserveRequestDuration(status string, duration time.Duration)
	IncrementIssuedCounter()
	IncrementMethodCallCounter(methodName string)
}

var (
	methodCallCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "cli",
		Name:      "method_calls_total",
		Help:      "Total number of method calls",
	}, []string{"method"})

	issuedOrdersCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "cli",
		Subsystem: "http",
		Name:      "issued_orders_total",
		Help:      "Total number of issued orders",
	})

	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "cli",
		Name:      "request_duration_seconds",
		Help:      "Request latency histogram",
		Buckets:   prometheus.DefBuckets,
	}, []string{"status"})
)

func IncrementMethodCallCounter(methodName string) {
	methodCallCounter.WithLabelValues(methodName).Inc()
}

func IncrementIssuedOrdersCounter() {
	issuedOrdersCounter.Inc()
}

func ObserveRequestDuration(status string, duration time.Duration) {
	requestDuration.WithLabelValues(status).Observe(duration.Seconds())
}