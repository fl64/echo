package processor

import (
	"echo-http/pkg/app/models"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Processor struct {
}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) GetRequest(r *http.Request) (*models.Req, error) {
	req := &models.Req{
		Host:       r.Host,
		URL:        r.URL.String(),
		Method:     r.Method,
		Headers:    r.Header,
		RemoteAddr: r.RemoteAddr,
	}
	buf, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}
	req.Body = string(buf)
	return req, nil
}

func (p Processor) GetEnvs() *models.Envs {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		envData := strings.Split(env, "=")
		envs[envData[0]] = envData[1]
	}
	return &models.Envs{
		Env: envs,
	}
}

func (p Processor) GetAll(r *http.Request) (*models.Info, error) {
	result := &models.Info{}
	result.HostData = make(map[string]string)
	result.HostData["hostname"], _ = os.Hostname()
	result.HostData["args"] = strings.Join(os.Args, ";")
	result.Envs = *p.GetEnvs()

	req, err := p.GetRequest(r)
	if err != nil {
		return nil, err
	}
	result.Req = req

	return result, nil
}
