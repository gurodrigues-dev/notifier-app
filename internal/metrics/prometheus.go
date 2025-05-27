package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusImpl struct {
	requestCounter *prometheus.CounterVec
	errorCounter   *prometheus.CounterVec
	durationHist   *prometheus.HistogramVec
}

func NewPrometheusImpl() *PrometheusImpl {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total de requisições HTTP por método e rota",
		},
		[]string{"method", "path"},
	)

	errorCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total de erros por método, rota e tipo de erro",
		},
		[]string{"method", "path", "error_type"},
	)

	durationHist := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duração das requisições HTTP por método e rota",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	prometheus.MustRegister(requestCounter, errorCounter, durationHist)

	return &PrometheusImpl{
		requestCounter: requestCounter,
		errorCounter:   errorCounter,
		durationHist:   durationHist,
	}
}

func (p *PrometheusImpl) IncRequest(method, path string) {
	if path == "" {
		path = "undefined"
	}
	p.requestCounter.WithLabelValues(method, path).Inc()
}

func (p *PrometheusImpl) IncError(method, path string, errType string) {
	if path == "" {
		path = "undefined"
	}
	if errType == "" {
		errType = "unknown"
	}
	p.errorCounter.WithLabelValues(method, path, errType).Inc()
}

func (p *PrometheusImpl) ObserveDuration(method, path string, durationSeconds float64) {
	if path == "" {
		path = "undefined"
	}
	p.durationHist.WithLabelValues(method, path).Observe(durationSeconds)
}

func (p *PrometheusImpl) ExposeHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}
