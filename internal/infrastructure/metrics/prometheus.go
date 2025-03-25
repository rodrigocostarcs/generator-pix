package metrics

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusMiddleware é o middleware para métricas do Prometheus
type PrometheusMiddleware struct {
	// Contador de requisições total
	requestCounter *prometheus.CounterVec
	// Histograma para duração das requisições
	requestDuration *prometheus.HistogramVec
	// Contador de requisições em andamento
	requestsInProgress *prometheus.GaugeVec
	// Contador de respostas por código de status
	responseStatus *prometheus.CounterVec
}

// NewPrometheusMiddleware cria uma nova instância do middleware Prometheus
func NewPrometheusMiddleware() *PrometheusMiddleware {
	// Criar um registry para as métricas
	registry := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = registry
	prometheus.DefaultGatherer = registry

	// Métricas para monitoramento HTTP
	requestCounter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pix_generator_http_requests_total",
			Help: "Total de requisições HTTP recebidas",
		},
		[]string{"method", "endpoint"},
	)

	requestDuration := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "pix_generator_http_request_duration_seconds",
		Help:    "Duração das requisições HTTP em segundos",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "endpoint", "status"})

	requestsInProgress := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pix_generator_http_requests_in_progress",
		Help: "Número de requisições HTTP em andamento",
	}, []string{"method", "endpoint"})

	responseStatus := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pix_generator_http_response_status",
		Help: "Contador por código de status HTTP",
	}, []string{"method", "endpoint", "status"})

	return &PrometheusMiddleware{
		requestCounter:     requestCounter,
		requestDuration:    requestDuration,
		requestsInProgress: requestsInProgress,
		responseStatus:     responseStatus,
	}
}

// Middleware retorna a função de middleware do Gin
func (p *PrometheusMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrair informações da requisição
		path := c.FullPath()
		if path == "" {
			path = "notfound" // Requisições que não correspondem a nenhuma rota
		}
		method := c.Request.Method

		// Registrar início da requisição
		start := time.Now()

		// Incrementar contador de requisições em andamento
		p.requestsInProgress.WithLabelValues(method, path).Inc()

		// Incrementar contador de requisições total
		p.requestCounter.WithLabelValues(method, path).Inc()

		// Prosseguir com a cadeia de middleware
		c.Next()

		// Decrementar contador de requisições em andamento
		p.requestsInProgress.WithLabelValues(method, path).Dec()

		// Calcular duração da requisição
		duration := time.Since(start).Seconds()

		// Registrar código de status
		status := c.Writer.Status()
		statusStr := string(rune(status))

		// Incrementar contador de status
		p.responseStatus.WithLabelValues(method, path, statusStr).Inc()

		// Registrar duração da requisição
		p.requestDuration.WithLabelValues(method, path, statusStr).Observe(duration)
	}
}

// RegisterEndpoint registra o endpoint para exposição das métricas
func (p *PrometheusMiddleware) RegisterEndpoint(r *gin.Engine) {
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
