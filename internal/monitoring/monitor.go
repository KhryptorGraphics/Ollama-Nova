package monitoring

import (
"net/http"
"time"

"github.com/prometheus/client_golang/prometheus"
"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Monitor struct {
requestsTotal prometheus.Counter
requestDuration prometheus.Histogram
activePeers   prometheus.Gauge
}

func NewMonitor() *Monitor {
m := &Monitor{
requestsTotal: prometheus.NewCounter(prometheus.CounterOpts{
Name: "ollama_nova_requests_total",
Help: "Total number of requests",
}),
requestDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
Name: "ollama_nova_request_duration_seconds",
Help: "Request duration in seconds",
}),
activePeers: prometheus.NewGauge(prometheus.GaugeOpts{
Name: "ollama_nova_active_peers",
Help: "Number of active peers",
}),
}

prometheus.MustRegister(m.requestsTotal, m.requestDuration, m.activePeers)
return m
}

func (m *Monitor) StartMetricsServer(port int) {
http.Handle("/metrics", promhttp.Handler())
http.ListenAndServe(":"+string(port), nil)
}
