package main

import (
	"context"
	"echo-http/pkg/app"
	"echo-http/pkg/cfg"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config, err := cfg.GetConfig()
	if err != nil {
		log.Fatalf("Can't get config: %v \n", err)
	}
	a := app.NewApp(config)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
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
	err = a.Run(ctx)
	if err != nil {
		log.Fatalf("Can't run app: %v \n", err)
	}
}
