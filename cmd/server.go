package main

import (
	"echo-http/pkg/echoserver"
	"os"
)

func main() {

	ListenAddr, exists := os.LookupEnv("LISTEN_ADDR")
	if !exists {
		ListenAddr = "0.0.0.0:8000"
	}
	Server := &echoserver.EchoServer{
		ListenAddr: ListenAddr,
	}
	Server.Run()
}
