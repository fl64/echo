package app_tcp

import (
	"bufio"
	"context"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
)

type TcpEchoServer struct {
	addr string
	prom *prometheus.Registry
}

func NewTCPServer(addr string, prom *prometheus.Registry) *TcpEchoServer {
	return &TcpEchoServer{
		addr: addr,
		prom: prom,
	}
}

func (t *TcpEchoServer) Run(ctx context.Context) error {
	log.Infof("Starting tcp server on %s", t.addr)
	listener, err := net.Listen("tcp", t.addr)
	if err != nil {
		return err
	}
	for {
		con, err := listener.Accept()
		if err != nil {
			log.Warningf("Failed to accept connection: %+v", err)
			continue
		}
		tcp_con_handle(ctx, con)
	}
}

func tcp_con_handle(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		select {
		case <-ctx.Done():
			log.Info("TCP server stopped")
			break
		default:
			bytes, err := reader.ReadBytes(byte('\n'))
			if err != nil {
				if err != io.EOF {
					log.Errorf("Failed to read data: %+v", err)
				}
				return
			}
			log.Infof("Request: '%s'", bytes)
			_, err = conn.Write(bytes)
			if err != nil {
				log.Errorf("Failed to write data: %+v", err)
			}
		}

	}
}
