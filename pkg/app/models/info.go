package models

import "net/http"

type Info struct {
	Host       string            `json:"host"`
	URL        string            `json:"url"`
	Method     string            `json:"method"`
	Headers    http.Header       `json:"headers"`
	Body       string            `json:"body"`
	Envs       map[string]string `json:"env"`
	HostData   map[string]string `json:"hostdata"`
	Ips        []string          `json:"ipaddr"`
	RemoteAddr string            `json:"remoteaddr"`
}
