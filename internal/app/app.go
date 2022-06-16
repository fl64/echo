package app

import (
	"context"
	"echo-http/internal/app/api"
	"echo-http/internal/app/api/middleware"
	"echo-http/internal/app/handlers"
	"echo-http/internal/app/processor"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type App struct {
	addr string
	srv  *http.Server
	prom *prometheus.Registry
}

const (
	Namespace = "echo"
)

func NewApp(addr string, prom *prometheus.Registry) *App {
	return &App{
		addr: addr,
		prom: prom,
	}
}

func (a *App) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		ctxShutDown := context.Background()
		ctxShutDown, cancel := context.WithTimeout(ctxShutDown, time.Second*5)
		defer func() {
			cancel()
		}()
		if err := a.srv.Shutdown(ctxShutDown); err != nil {
			log.Fatalf("server Shutdown Failed:%s", err)
		}

	}()

	log.Infof("Starting app on %s", a.addr)
	p := processor.NewProcessor(a.prom)
	h := handlers.NewHandler(p)
	r := api.CreateRoutes(h)
	m := middleware.NewMiddleware()
	r.Use(m.Logging)
	a.srv = &http.Server{
		Addr:    a.addr,
		Handler: r,
	}
	err := a.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	log.Info("App stopped")
	return nil

}
