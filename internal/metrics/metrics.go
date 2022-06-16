package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type MetricsServer struct {
	addr string
	srv  *http.Server
	prom *prometheus.Registry
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func NewMetricsServer(addr string, prom *prometheus.Registry) *MetricsServer {
	return &MetricsServer{
		addr: addr,
		prom: prom,
	}
}

func (m *MetricsServer) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		ctxShutDown := context.Background()
		ctxShutDown, cancel := context.WithTimeout(ctxShutDown, time.Second*5)
		defer func() {
			cancel()
		}()
		if err := m.srv.Shutdown(ctxShutDown); err != nil {
			log.Fatalf("server Shutdown Failed:%s", err)
		}

	}()
	log.Infof("Starting metrics server on %s", m.addr)
	r := http.NewServeMux()
	r.HandleFunc("/healthz", okHandler)
	r.HandleFunc("/ready", okHandler)
	m.prom.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	handler := promhttp.HandlerFor(m.prom, promhttp.HandlerOpts{})
	r.Handle("/metrics", handler)

	m.srv = &http.Server{
		Addr:    m.addr,
		Handler: r,
	}
	err := m.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	log.Info("Metrics server stopped")
	return nil
}
