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

type serverMetrics struct {
	requestDuration *prometheus.HistogramVec
	counterVec      *prometheus.CounterVec
}

func NewServerMetrics(reg prometheus.Registerer) Metrics {
	rd := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "cli",
		Subsystem: "http",
		Name:      "cli_duration",
		Help:      "Request latency histogram",
		Buckets:   prometheus.DefBuckets,
	}, []string{"status"})

	cv := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "cli",
		Subsystem: "http",
		Name:      "cli_total",
		Help:      "Total number of method calls and issued orders",
	}, []string{"type"})
	reg.MustRegister(rd)
	reg.MustRegister(cv)
	return &serverMetrics{
		requestDuration: rd,
		counterVec:      cv,
	}
}

func (sm *serverMetrics) ObserveRequestDuration(status string, duration time.Duration) {
	sm.requestDuration.WithLabelValues(status).Observe(duration.Seconds())
}

func (sm *serverMetrics) IncrementIssuedCounter() {
	sm.counterVec.WithLabelValues("Issued").Inc()
}

func (sm *serverMetrics) IncrementMethodCallCounter(methodName string) {
	sm.counterVec.WithLabelValues(methodName).Inc()
}