package app_tcp

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/fl64/echo/internal/app-http/models"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
)

type TcpEchoServer struct {
	addr string
	prom *prometheus.Registry
	msg  string
}

func NewTCPServer(addr string, prom *prometheus.Registry, msg string) *TcpEchoServer {
	return &TcpEchoServer{
		addr: addr,
		prom: prom,
		msg:  msg,
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
		t.tcp_con_handle(ctx, con)
	}
}

func (t *TcpEchoServer) tcp_con_handle(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		select {
		case <-ctx.Done():
			log.Info("TCP server stopped")
			break
		default:
			var bytes []byte
			var err error
			log.Infof("Remote connection from %s", conn.RemoteAddr().String())
			switch t.msg {
			case "":
				bytes, err = reader.ReadBytes(byte('\n'))
				if err != nil {
					if err != io.EOF {
						log.Errorf("Failed to read data: %+v", err)
					}
					return
				}
				log.Infof("Request: '%s'", bytes)
			default:
				bytes, _ = json.Marshal(models.Msg{Msg: os.ExpandEnv(t.msg)})
			}

			_, err = conn.Write(bytes)
			if err != nil {
				log.Errorf("Failed to write data: %+v", err)
			}
			return
		}
	}
}
