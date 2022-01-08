package models

import "net/http"

type Info struct {
	Req      *Req              `json:"request"`
	Envs     *Envs             `json:"envs"`
	HostData map[string]string `json:"hostdata"`
	Routes   *[]Route          `json:"routes"`
	Ifaces   *[]IFace          `json:"ifaces"`
}

type Envs struct {
	Env map[string]string `json:"env"`
}

type Route struct {
	Dst      string   `json:"dst,omitempty"`
	Gateway  string   `json:"gateway,omitempty"`
	Dev      string   `json:"dev,omitempty"`
	Protocol string   `json:"protocol,omitempty"`
	Metric   int      `json:"metric,omitempty"`
	Flags    []string `json:"flags,omitempty"`
}

type IFace struct {
	Ifindex   int      `json:"ifindex,omitempty"`
	Ifname    string   `json:"ifname,omitempty"`
	Flags     []string `json:"flags,omitempty"`
	MTU       int      `json:"mtu,omitempty"`
	Operstate string   `json:"operstate,omitempty"`
	Group     string   `json:"group,omitempty"`
	LinkType  string   `json:"link_type,omitempty"`
	Address   string   `json:"address,omitempty"`
	Broadcast string   `json:"broadcast,omitempty"`
	AddrInfo  []Addr   `json:"addr_info,omitempty"`
}

type Addr struct {
	Family            string `json:"family,omitempty"`
	Local             string `json:"local,omitempty"`
	Prefixlen         int    `json:"prefixlen,omitempty"`
	Broadcast         string `json:"broadcast,omitempty"`
	Scope             string `json:"scope,omitempty"`
	Label             string `json:"label,omitempty"`
	ValidLifeTime     int64  `json:"valid_life_time,omitempty"`
	PreferredLifeTime int64  `json:"preferred_life_time,omitempty"`
}

type Req struct {
	Host       string      `json:"host"`
	URL        string      `json:"url"`
	Method     string      `json:"method"`
	Headers    http.Header `json:"headers"`
	Body       string      `json:"body"`
	RemoteAddr string      `json:"remoteaddr"`
}
