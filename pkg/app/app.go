package app

import (
	"context"
	"echo-http/api"
	"echo-http/api/middleware"
	"echo-http/pkg/app/handlers"
	"echo-http/pkg/app/processor"
	"echo-http/pkg/cfg"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type App struct {
	cfg *cfg.Cfg
	srv *http.Server
}

func NewApp(c *cfg.Cfg) *App {
	return &App{
		cfg: c,
	}
}

func (a *App) Run(ctx context.Context) error {
	go func() {
		select {
		case <-ctx.Done():
			ctxShutDown := context.Background()
			ctxShutDown, cancel := context.WithTimeout(ctxShutDown, time.Second*5)
			defer func() {
				cancel()
			}()
			if err := a.srv.Shutdown(ctxShutDown); err != nil {
				log.Fatalf("server Shutdown Failed:%s", err)
			}
		}
	}()
	log.SetFormatter(&log.JSONFormatter{})
	log.Println("Starting app")
	p := processor.NewProcessor()
	h := handlers.NewHandler(p)
	r := api.CreateRoutes(h)
	m := middleware.NewMiddleware()
	r.Use(m.Logging)
	a.srv = &http.Server{
		Addr:    a.cfg.Addr,
		Handler: r,
	}
	err := a.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal("Can't run http server:", err)
	}
	log.Println("App stopped")
	return nil

}
