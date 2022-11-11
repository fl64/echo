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
	addrTLS string
	srv     *http.Server
	srvTLS  *http.Server
	prom    *prometheus.Registry
	keyFile string
	crtFile string
}

const (
	Namespace = "echo"
)

func NewApp(addr, addrTLS, crt, key string, prom *prometheus.Registry) *App {
	return &App{
		addr:    addr,
		addrTLS: addrTLS,
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
			log.Fatalf("http server Shutdown Failed:%s", err)
		}
		if a.srvTLS != nil {
			if err := a.srvTLS.Shutdown(ctxShutDown); err != nil {
				log.Fatalf("https server Shutdown Failed:%s", err)
			}
		}

	}()

	log.Info("Starting app ...")
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
		a.srvTLS = &http.Server{
			Addr:    a.addrTLS,
			Handler: r,
		}
		log.Infof("Starting https on %s", a.addrTLS)
		go func() {
			err := a.srvTLS.ListenAndServeTLS(a.crtFile, a.keyFile)
			if err != nil {
				log.Errorf("Can't serve https: %+v", err)
			}
		}()
	}

	log.Infof("Starting http on %s", a.addr)
	err = a.srv.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return err
	}
	log.Info("App stopped")
	return nil

}
