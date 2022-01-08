package cfg

import (
	"github.com/caarlos0/env/v6"
)

type Cfg struct {
	Addr string `env:"LISTEN_ADDR" envDefault:":8000"`
}

func GetConfig() (*Cfg, error) {
	c := &Cfg{}
	if err := env.Parse(c); err != nil {
		return nil, err
	}
	return c, nil
}
