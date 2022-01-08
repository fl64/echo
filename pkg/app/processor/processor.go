package processor

import (
	"echo-http/pkg/app/models"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

type Processor struct {
}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p Processor) GetInfo(r *http.Request) (result *models.Info, err error) {
	result = &models.Info{}
	// getting request info
	result.Host = r.Host
	result.Method = r.Method
	result.Headers = r.Header
	result.URL = r.URL.String()
	result.Envs = make(map[string]string)

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	result.Body = string(buf)
	// getting ENV info
	for _, env := range os.Environ() {
		envData := strings.Split(env, "=")
		result.Envs[envData[0]] = envData[1]
	}

	result.HostData = make(map[string]string)
	result.HostData["hostname"], _ = os.Hostname()
	result.HostData["args"] = strings.Join(os.Args, ";")

	// getting interfaces info
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			result.Ips = append(result.Ips, net.IP.String(ip))
		}
	}

	result.RemoteAddr = r.RemoteAddr
	return
}
