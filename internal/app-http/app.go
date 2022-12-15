package app_http

import (
	"context"
	"fmt"
	"github.com/fl64/echo/internal/app-http/api"
	"github.com/fl64/echo/internal/app-http/api/middleware"
	"github.com/fl64/echo/internal/app-http/handlers"
	"github.com/fl64/echo/internal/app-http/processor"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

type httpSrv struct {
	addr string
	srv  *http.Server
}

type httpsSrv struct {
	httpSrv
	keyFile string
	crtFile string
}

type App struct {
	http       httpSrv
	https      httpsSrv
	prom       *prometheus.Registry
	respStatus *atomic.Int32
	PodNS      string
	PodName    string
	SleepDelay time.Duration
}

const (
	Namespace = "echo"
)

func NewApp(addr, addrTLS, crt, key string, prom *prometheus.Registry, podNS, podName string, sleepDelay time.Duration) *App {
	return &App{
		http: httpSrv{
			addr: addr,
		},
		https: httpsSrv{
			httpSrv: httpSrv{
				addr: addrTLS,
			},
			crtFile: crt,
			keyFile: key,
		},
		prom:       prom,
		respStatus: &atomic.Int32{},
		PodNS:      podNS,
		PodName:    podName,
		SleepDelay: sleepDelay,
	}
}

func (a *App) Run(ctx context.Context) error {

	// watching for annotations
	//a.disaster.Store(false)
	cfg, err := rest.InClusterConfig()
	if err == nil {
		clientset, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			log.Fatalf("Error building kubernetes clientset: %v", err)
		}
		go func() {
			log.Infof("Run annotation checker")
			for {
				pod, err := clientset.CoreV1().Pods(a.PodNS).Get(context.TODO(), a.PodName, metav1.GetOptions{})
				if err != nil {
					log.Errorf("Can't get pod: %v", err)
				}
				if metav1.HasAnnotation(pod.ObjectMeta, "disaster") {
					if statusStr, ok := pod.Annotations["fl64.io/status"]; ok {
						status, err := strconv.Atoi(statusStr)
						if err != nil {
							a.respStatus.Store(int32(status))
							continue
						}
					}
				}
				a.respStatus.Store(200)
				time.Sleep(a.SleepDelay)
			}

		}()
	} else {
		log.Warn("Not in cluster")
	}

	go func() {
		<-ctx.Done()
		ctxShutDown := context.Background()
		ctxShutDown, cancel := context.WithTimeout(ctxShutDown, time.Second*5)
		defer func() {
			cancel()
		}()

		if a.https.srv != nil {
			if err := a.https.srv.Shutdown(ctxShutDown); err != nil {
				log.Fatalf("https server Shutdown Failed:%s", err)
			} else {
				log.Info("Https server stopped")
			}
		}

		if err := a.http.srv.Shutdown(ctxShutDown); err != nil {
			log.Fatalf("http server Shutdown Failed:%s", err)
		} else {
			log.Info("Http server stopped")
		}

	}()

	log.Info("Starting app ...")
	p := processor.NewProcessor(a.prom)
	h := handlers.NewHandler(p, a.respStatus)
	r := api.CreateRoutes(h)
	m := middleware.NewMiddleware(a.prom)
	r.Use(m.Metrics)
	r.Use(m.Logging)

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Infof("Working dir: %s", pwd)

	_, errKey := os.Stat(a.https.keyFile)
	_, errCrt := os.Stat(a.https.crtFile)

	a.http.srv = &http.Server{
		Addr:    a.http.addr,
		Handler: r,
	}

	// if cert/key exist -> run https
	if errKey == nil && errCrt == nil {
		a.https.srv = &http.Server{
			Addr:    a.https.addr,
			Handler: r,
		}
		log.Infof("Starting https on %s", a.https.addr)
		go func() {
			err := a.https.srv.ListenAndServeTLS(a.https.crtFile, a.https.keyFile)
			if err != nil && err != http.ErrServerClosed {
				log.Errorf("Can't serve https: %+v", err)
			}
		}()
	}

	log.Infof("Starting http on %s", a.http.addr)
	err = a.http.srv.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return err
	}

	log.Info("App stopped")
	return nil

}
