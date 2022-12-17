package processor

import (
	"github.com/fl64/echo/internal/app-http/models"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Processor struct {
	prom    *prometheus.Registry
	metrics ProcessorMetrics
	msg     string
}

type ProcessorMetrics struct {
	operationDuration *prometheus.HistogramVec
	operationsCount   *prometheus.CounterVec
}

func NewProcessor(msg string, prom *prometheus.Registry) *Processor {
	p := &Processor{
		prom: prom,
		msg:  msg,
	}
	p.initMetrics()

	return p
}

func (p *Processor) initMetrics() {
	p.metrics.operationDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   "echo",
		Subsystem:   "operation",
		Name:        "duration",
		Help:        "A histogram of latencies for operation duration",
		ConstLabels: nil,
		Buckets:     prometheus.DefBuckets,
	}, []string{})
	p.metrics.operationsCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "echo",
		Subsystem:   "operation",
		Name:        "count_total",
		Help:        "A operations counter",
		ConstLabels: nil,
	}, []string{})
	if p.prom != nil {
		p.prom.MustRegister(p.metrics.operationsCount)
		p.prom.MustRegister(p.metrics.operationDuration)
	}

}

func (p *Processor) GetRequest(r *http.Request) (*models.Req, error) {
	req := &models.Req{
		Host:       r.Host,
		URL:        r.URL.String(),
		Method:     r.Method,
		Headers:    r.Header,
		RemoteAddr: r.RemoteAddr,
	}
	buf, err := io.ReadAll(r.Body)
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

func (p Processor) Do(r *http.Request) (interface{}, error) {
	startTime := time.Now()
	var result interface{}
	switch p.msg {
	case "":
		res := &models.Info{}
		res.HostData = make(map[string]string)
		res.HostData["hostname"], _ = os.Hostname()
		res.HostData["args"] = strings.Join(os.Args, ";")
		res.Envs = *p.GetEnvs()
		req, err := p.GetRequest(r)
		if err != nil {
			return nil, err
		}
		res.Req = req
		result = res
	default:
		res := &models.Msg{Msg: os.ExpandEnv(p.msg)}
		result = res
	}

	p.metrics.operationsCount.WithLabelValues().Inc()
	p.metrics.operationDuration.WithLabelValues().Observe(time.Since(startTime).Seconds())
	return result, nil
}

func (p Processor) DoMessage(r *http.Request) (interface{}, error) {
	startTime := time.Now()
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
	p.metrics.operationsCount.WithLabelValues().Inc()
	p.metrics.operationDuration.WithLabelValues().Observe(time.Since(startTime).Seconds())
	return result, nil
}
