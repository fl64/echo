package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Middleware struct {
	prom    *prometheus.Registry
	metrics MiddlewareMetrics
}

type MiddlewareMetrics struct {
	requestDuration *prometheus.HistogramVec
	requestCount    *prometheus.CounterVec
}

func NewMiddleware(prom *prometheus.Registry) *Middleware {
	m := Middleware{
		prom: prom,
	}
	m.initMetrics()
	return &m
}

func (p *Middleware) initMetrics() {
	p.metrics.requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   "echo",
		Subsystem:   "request",
		Name:        "duration",
		Help:        "A histogram of latencies for request duration",
		ConstLabels: nil,
		Buckets:     prometheus.DefBuckets,
	}, []string{"method"})
	p.metrics.requestCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "echo",
		Subsystem:   "request",
		Name:        "count_total",
		Help:        "A requests counter",
		ConstLabels: nil,
	}, []string{"method", "path"})
	if p.prom != nil {
		p.prom.MustRegister(p.metrics.requestCount)
		p.prom.MustRegister(p.metrics.requestDuration)
	}
}

func (m *Middleware) Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		//log.Println(r.RequestURI)
		startTime := time.Now()
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
		//log.Infof("PIU")

		m.metrics.requestCount.With(prometheus.Labels{
			"method": r.Method,
			"path":   r.RequestURI,
		}).Inc()
		m.metrics.requestDuration.With(prometheus.Labels{
			"method": r.Method,
		}).Observe(time.Since(startTime).Seconds())
	})
}

func (m *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		//log.Println(r.RequestURI)
		t1 := time.Now()
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.WithFields(log.Fields{
			"url":            r.URL.String(),
			"host":           r.Host,
			"user-agent":     r.Header["User-Agent"],
			"method":         r.Method,
			"remote-addr":    r.RemoteAddr,
			"content-length": r.ContentLength,
			"duration":       t2.Sub(t1).String(),
		}).Info()
	})
}
