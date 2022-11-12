package models

import (
	"net/http"
)

type Info struct {
	Req      *Req              `json:"request"`
	HostData map[string]string `json:"hostdata"`
	Envs
}

type Envs struct {
	Env map[string]string `json:"env"`
}

type Req struct {
	Host       string      `json:"host"`
	URL        string      `json:"url"`
	Method     string      `json:"method"`
	Headers    http.Header `json:"headers"`
	Body       string      `json:"body"`
	RemoteAddr string      `json:"remoteaddr"`
}
