package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// ============ 统一的业务操作指标 ============
var (
	// operation_result_total{module="grammar", operation="validate", result="success"}
	operationTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "operation_result_total",
			Help: "Total number of business operations by module, operation type and result.",
		},
		[]string{"module", "operation", "result"},
	)

	// operation_duration_seconds{module="fsm", operation="minimize"}
	operationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "operation_duration_seconds",
			Help:    "Duration of business operations in seconds.",
			Buckets: []float64{0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"module", "operation"},
	)
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"method", "path"},
	)
)

func PrometheusRegister() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration, operationTotal, operationDuration)
}

func IncOperation(module, operation, result string) {
	operationTotal.WithLabelValues(module, operation, result).Inc()
}

func ObserveOperationDuration(module, operation string, seconds float64) {
	operationDuration.WithLabelValues(module, operation).Observe(seconds)
}

func HTTPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		
		c.Next()

		status := strconv.Itoa(c.Writer.Status())

		httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(time.Since(start).Seconds())
	}
}
