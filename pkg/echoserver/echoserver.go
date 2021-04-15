package echoserver

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type EchoServer struct {
	ListenAddr string
}

type reqInfo struct {
	URL        string            `json:"url"`
	Method     string            `json:"method"`
	Headers    http.Header       `json:"headers"`
	Body       string            `json:"body"`
	Envs       map[string]string `json:"env"`
	HostData   map[string]string `json:"hostdata"`
	Ips        []string          `json:"ipaddr"`
	RemoteAddr string            `json:"remoteaddr`
}

func (e *EchoServer) getInfo(r *http.Request) (result *reqInfo, err error) {
	result = &reqInfo{}
	// getting request info
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

func (e *EchoServer) echoHandlerJson(w http.ResponseWriter, r *http.Request) {
	Req, err := e.getInfo(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JsonReq, err := json.Marshal(Req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(JsonReq)
}

func (e *EchoServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		//log.Println(r.RequestURI)
		t1 := time.Now()
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.WithFields(log.Fields{
			"url":         r.URL.String(),
			"user-agent":  r.Header["User-Agent"],
			"method":      r.Method,
			"remote-addr": r.RemoteAddr,
			"duration":    t2.Sub(t1).String(),
		}).Info()
	})
}

func (e *EchoServer) Run() {
	log.SetFormatter(&log.JSONFormatter{})
	log.Printf("Starting server on %s", e.ListenAddr)
	router := mux.NewRouter()
	//router.HandleFunc("/", e.echoHandlerJson)
	//router.NotFoundHandler = http.HandlerFunc(e.echoHandlerJson)
	router.PathPrefix("/").Handler(http.HandlerFunc(e.echoHandlerJson))
	router.Use(e.loggingMiddleware)
	srv := &http.Server{
		Handler:      router,
		Addr:         e.ListenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
