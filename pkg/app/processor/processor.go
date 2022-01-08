package processor

import (
	"bufio"
	"echo-http/pkg/app/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Processor struct {
}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) GetRoutes() (*[]models.Route, error) {
	routes := make([]models.Route, 1)
	out, err := exec.Command("ip", "-j", "r").Output()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(out, &routes)
	if err != nil {
		return nil, err
	}
	return &routes, nil
}

func (p *Processor) GetIFaces() (*[]models.IFace, error) {
	ifaces := make([]models.IFace, 1)
	out, err := exec.Command("ip", "-j", "a").Output()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(out, &ifaces)
	if err != nil {
		return nil, err
	}
	return &ifaces, nil
}

func (p *Processor) GetMounts() (*models.Mounts, error) {
	mounts := make(models.Mounts, 0)
	out, err := exec.Command("mount").Output()
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(strings.NewReader(string(out)))
	for sc.Scan() {
		mounts = append(mounts, sc.Text())
	}
	return &mounts, nil
}

func (p *Processor) GetResolvConf() (*models.ResolvConf, error) {
	resolvconf := make(models.ResolvConf, 0)
	out, err := exec.Command("cat", "/etc/resolv.conf").Output()
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(strings.NewReader(string(out)))
	for sc.Scan() {
		if strings.HasPrefix(sc.Text(), "#") {
			continue
		}
		if len(sc.Text()) == 0 {
			continue
		}
		resolvconf = append(resolvconf, sc.Text())
	}
	return &resolvconf, nil
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
	return &models.Envs{Env: envs}
}

func (p Processor) GetAll(r *http.Request) (*models.Info, error) {
	result := &models.Info{}
	result.HostData = make(map[string]string)
	result.HostData["hostname"], _ = os.Hostname()
	result.HostData["args"] = strings.Join(os.Args, ";")
	result.Envs = p.GetEnvs()

	resolvconf, err := p.GetResolvConf()
	if err != nil {
		return nil, err
	}
	result.ResolvConf = resolvconf

	req, err := p.GetRequest(r)
	if err != nil {
		return nil, err
	}
	result.Req = req
	routes, err := p.GetRoutes()
	if err != nil {
		return nil, err
	}
	result.Routes = routes
	ifs, err := p.GetIFaces()
	if err != nil {
		return nil, err
	}
	result.Ifaces = ifs
	mounts, err := p.GetMounts()
	if err != nil {
		return nil, err
	}
	result.Mounts = mounts
	return result, nil
}
