package k8s

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strconv"
	"sync/atomic"
	"time"
)

const statusAnnotation = "status"

type k8sClient struct {
	podNS              string
	podName            string
	tickerDuration     time.Duration
	clientset          *kubernetes.Clientset
	httpResponseStatus *atomic.Int32
}

func NewK8sClient(podNS, podName string, tickerDuration time.Duration, httpResponseStatus *atomic.Int32) (*k8sClient, error) {
	var clientset *kubernetes.Clientset
	var err error
	cfg, err := rest.InClusterConfig()
	if err == nil {
		clientset, err = kubernetes.NewForConfig(cfg)
		if err != nil {
			return nil, fmt.Errorf("Error building kubernetes clientset: %v", err)
		}
	} else {
		return nil, fmt.Errorf("Not in cluster")
	}
	return &k8sClient{
		podNS:              podNS,
		podName:            podName,
		tickerDuration:     tickerDuration,
		clientset:          clientset,
		httpResponseStatus: httpResponseStatus,
	}, err
}

func (k *k8sClient) Run(ctx context.Context) {
	log.Infof("Run annotation checker")
	tk := time.NewTicker(k.tickerDuration)
	for {
		select {
		case <-ctx.Done():
			log.Info("K8s client Stopped")
			break
		case <-tk.C:
			pod, err := k.clientset.CoreV1().Pods(k.podNS).Get(context.TODO(), k.podName, metav1.GetOptions{})
			if err != nil {
				log.Errorf("Can't get pod: %v", err)
			}
			if metav1.HasAnnotation(pod.ObjectMeta, statusAnnotation) {
				if statusStr, ok := pod.Annotations[statusAnnotation]; ok {
					var err error
					var status int
					status, err = strconv.Atoi(statusStr)
					if err == nil {
						k.httpResponseStatus.Store(int32(status))
						continue
					} else {
						log.Errorf("can't convert status to int %v", err)
					}
				}
			}
			k.httpResponseStatus.Store(200)
		}
	}
}
