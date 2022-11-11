package cfg

import (
	"github.com/caarlos0/env/v6"
)

type Cfg struct {
	ServerAddr  string `env:"SERVER_ADDR" envDefault:":8000"`
	MetricsAddr string `env:"METRICS_ADDR" envDefault:":8001"`
	TLSKeyFile  string `env:"TLS_KEY_FILE" envDefault:"tls.key"`
	TLSCrtFile  string `env:"TLS_CRT_FILE" envDefault:"tls.crt"`
}

func GetConfig() (*Cfg, error) {
	c := &Cfg{}
	if err := env.Parse(c); err != nil {
		return nil, err
	}
	return c, nil
}
