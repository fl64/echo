package main

import (
	"context"
	app_http "github.com/fl64/echo/internal/app-http"
	app_tcp "github.com/fl64/echo/internal/app-tcp"
	"github.com/fl64/echo/internal/cfg"
	"github.com/fl64/echo/internal/k8s"
	"github.com/fl64/echo/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

var BuildDatetime = "none"
var BuildVer = "devel"

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		s := <-sigChan
		close(sigChan)
		log.Infof("Catch signal: %s", s)
		cancel()
	}()

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.Infof("App ver %s, build time %s", BuildVer, BuildDatetime)

	config, err := cfg.GetConfig()
	if err != nil {
		log.Fatalf("Can't get config: %v \n", err)
	}

	prom := prometheus.NewRegistry()
	metrics := metrics.NewMetricsServer(config.MetricsServerAddr, prom)
	tcpServer := app_tcp.NewTCPServer(config.TCPServerAddr, prom)

	httpResponseStatus := &atomic.Int32{}
	httpServer := app_http.NewApp(config.HTTPServerAddr, config.HTTPSServerAddr, config.TLSCrtFile, config.TLSKeyFile, prom, httpResponseStatus)

	// run annotations checker
	k, err := k8s.NewK8sClient(config.PodNS, config.PodName, config.TickerDuration, httpResponseStatus)
	if err != nil {
		log.Warnf("Can't create k8s client: %+v", err)
	} else {
		go k.Run(ctx)
	}

	// metrics
	go func() {
		err = metrics.Run(ctx)
		if err != nil {
			log.Fatalf("Can't run metrics server: %v \n", err)
		}
	}()

	// tcp
	go func() {
		err = tcpServer.Run(ctx)
		if err != nil {
			log.Fatalf("Can't run TCP server: %v \n", err)
		}
	}()

	// http/https
	err = httpServer.Run(ctx)
	if err != nil {
		log.Fatalf("Can't run app: %v \n", err)
	}
}

// https://stackoverflow.com/questions/43631854/gracefully-shutdown-gorilla-server
