package main

import (
	"context"
	app_http "echo/internal/app-http"
	app_tcp "echo/internal/app-tcp"
	"echo/internal/cfg"
	"echo/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var BuildDatetime = "none"
var BuildVer = "devel"

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.Infof("App ver %s, build time %s", BuildVer, BuildDatetime)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	config, err := cfg.GetConfig()
	if err != nil {
		log.Fatalf("Can't get config: %v \n", err)
	}

	prom := prometheus.NewRegistry()
	m := metrics.NewMetricsServer(config.MetricsServerAddr, prom)
	t := app_tcp.NewTCPServer(config.TCPServerAddr, prom)
	a := app_http.NewApp(config.HTTPServerAddr, config.HTTPSServerAddr, config.TLSCrtFile, config.TLSKeyFile, prom)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		s := <-sigChan
		close(sigChan)
		log.Println("Catch signal: ", s)
		cancel()
	}()
	go func() {
		err = m.Run(ctx)
		if err != nil {
			log.Fatalf("Can't run metrics server: %v \n", err)
		}
	}()
	go func() {
		err = t.Run()
		if err != nil {
			log.Fatalf("Can't run TCP server: %v \n", err)
		}
	}()
	err = a.Run(ctx)
	if err != nil {
		log.Fatalf("Can't run app: %v \n", err)
	}
}
