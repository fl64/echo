package app

import (
	"context"
	"echo-http/internal/app/api"
	"echo-http/internal/app/api/middleware"
	"echo-http/internal/app/handlers"
	"echo-http/internal/app/processor"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

type App struct {
	addr    string
	srv     *http.Server
	prom    *prometheus.Registry
	keyFile string
	crtFile string
}

const (
	Namespace = "echo"
)

func NewApp(addr, crt, key string, prom *prometheus.Registry) *App {
	return &App{
		addr:    addr,
		crtFile: crt,
		keyFile: key,
		prom:    prom,
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

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Infof("Working dir: %s", pwd)

	_, errKey := os.Stat(a.keyFile)
	_, errCrt := os.Stat(a.crtFile)

	a.srv = &http.Server{
		Addr:    a.addr,
		Handler: r,
	}
	if errKey == nil && errCrt == nil {
		log.Info("Serving https")
		err = a.srv.ListenAndServeTLS(a.crtFile, a.keyFile)
	} else {
		log.Info("Serving http")
		err = a.srv.ListenAndServe()
	}

	if err != nil && err != http.ErrServerClosed {
		return err
	}
	log.Info("App stopped")
	return nil

}
