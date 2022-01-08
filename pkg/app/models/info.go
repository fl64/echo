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
	Routes     *[]Route          `json:"routes"`
}

type Route struct {
	Dst      string `json:"dst,omitempty"`
	Gateway  string `json:"gateway,omitempty"`
	Dev      string `json:"dev,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Metric   int    `json:"metric,omitempty"`
}
