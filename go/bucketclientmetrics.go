package bucketclient

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type BucketClientMetrics struct {
	RequestsTotal              *prometheus.CounterVec
	RequestDurationSeconds     *prometheus.SummaryVec
	RequestBytesSentTotal      *prometheus.CounterVec
	ResponseBytesReceivedTotal *prometheus.CounterVec
}

var metricsLabels = []string{
	"endpoint",
	"method",
	"action",
	"code",
}

var globalMetrics *BucketClientMetrics
var globalMetricsLock sync.Mutex

// EnableMetrics enables Prometheus metrics gathering for the provided client and registers
// them in the provided registerer.
//
// Metrics implemented:
//   - `s3_metadata_bucketclient_requests_total`:
//     Number of requests processed (counter)
//   - `s3_metadata_bucketclient_request_duration_seconds`:
//     Time elapsed processing requests to bucketd, in seconds (summary)
//   - `s3_metadata_bucketclient_request_bytes_sent_total`:
//     Number of request body bytes sent to bucketd (counter)
//   - `s3_metadata_bucketclient_response_bytes_received_total`:
//     Number of response body bytes received from bucketd (counter)
//
// Metrics have the following labels attached:
//   - `endpoint`:
//     bucketd endpoint such as `http://localhost:9000`
//   - `method`:
//     HTTP method
//   - `action`:
//     name of the API action, such as `CreateBucket`. Admin actions are prefixed with `Admin`.
//   - `code`:
//     HTTP status code returned, or "0" for generic network or protocol errors
func (client *BucketClient) EnableMetrics(registerer prometheus.Registerer) {
	globalMetricsLock.Lock()
	defer globalMetricsLock.Unlock()

	if globalMetrics == nil {
		globalMetrics = &BucketClientMetrics{
			RequestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: MetricsNamespace,
				Name:      "requests_total",
				Help:      "Number of requests processed",
			}, metricsLabels),

			RequestDurationSeconds: prometheus.NewSummaryVec(prometheus.SummaryOpts{
				Namespace:  MetricsNamespace,
				Name:       "request_duration_seconds",
				Help:       "Time elapsed processing requests to bucketd, in seconds",
				Objectives: MetricsSummaryDefaultObjectives,
			}, metricsLabels),

			RequestBytesSentTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: MetricsNamespace,
				Name:      "request_bytes_sent_total",
				Help:      "Number of request body bytes sent to bucketd",
			}, metricsLabels),

			ResponseBytesReceivedTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: MetricsNamespace,
				Name:      "response_bytes_received_total",
				Help:      "Number of response body bytes received from bucketd",
			}, metricsLabels),
		}
		registerer.MustRegister(
			globalMetrics.RequestsTotal,
			globalMetrics.RequestDurationSeconds,
			globalMetrics.RequestBytesSentTotal,
			globalMetrics.ResponseBytesReceivedTotal,
		)
	}
	client.Metrics = globalMetrics
}
