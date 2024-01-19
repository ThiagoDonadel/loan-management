package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var statusByMethodCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "funding_http_requests_status_counter",
		Help: "Number of requests with the response status by method",
	},
	[]string{"method", "status"},
)

var requestDurationSumary = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:       "funding_http_requests_duration",
		Help:       "Number of requests made to the funding api",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"route"},
)

var simuulatedAmmountCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "funding_simulated_ammount_sum",
		Help: "The total ammount simuleted in the API",
	},
)

func RegisterMetrics() {
	prometheus.MustRegister(statusByMethodCounter, requestDurationSumary, simuulatedAmmountCounter)
}

func PrometheusMetricsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		start := time.Now()

		ctx.Next()

		requestDurationSumary.With(
			prometheus.Labels{
				"route": ctx.FullPath(),
			},
		).Observe(float64(time.Since(start).Seconds()))

		statusByMethodCounter.With(prometheus.Labels{
			"method": ctx.Request.Method,
			"status": strconv.Itoa(ctx.Writer.Status()),
		}).Inc()
	}
}

func IncrementSimulatedAmmountCounter(amount float64) {
	simuulatedAmmountCounter.Add(amount)
}
