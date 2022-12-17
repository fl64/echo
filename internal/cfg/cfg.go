package cfg

import (
	"github.com/caarlos0/env/v6"
	"time"
)

type Cfg struct {
	Message           string        `env:"MESSAGE" envDefault:""`
	HTTPServerAddr    string        `env:"HTTP_ADDR" envDefault:":8000"`
	HTTPSServerAddr   string        `env:"HTTPS_ADDR" envDefault:":8443"`
	TCPServerAddr     string        `env:"TCP_ADDR" envDefault:":1234"`
	MetricsServerAddr string        `env:"METRICS_ADDR" envDefault:":8001"`
	TLSKeyFile        string        `env:"TLS_KEY_FILE" envDefault:"tls.key"`
	TLSCrtFile        string        `env:"TLS_CRT_FILE" envDefault:"tls.crt"`
	PodName           string        `env:"POD_NAME"`
	PodNS             string        `env:"POD_NAMESPACE"`
	TickerDuration    time.Duration `env:"K8S_TICKER_DURATION" envDefault:"1s"`
}

func GetConfig() (*Cfg, error) {
	c := &Cfg{}
	if err := env.Parse(c); err != nil {
		return nil, err
	}
	return c, nil
}
