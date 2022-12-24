package k8s

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"strconv"
	"sync/atomic"
)

const statusAnnotation = "status"

type k8sClient struct {
	podNS              string
	podName            string
	clientset          *kubernetes.Clientset
	httpResponseStatus *atomic.Int32
}

func NewK8sClient(podNS, podName string, httpResponseStatus *atomic.Int32) (*k8sClient, error) {
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
		clientset:          clientset,
		httpResponseStatus: httpResponseStatus,
	}, err
}

func (k *k8sClient) setStatus(status int32) {
	if k.httpResponseStatus.Load() == status {
		return // nothing to do
	}
	log.Infof("Change http response status code from %d to %d", k.httpResponseStatus.Load(), status)
	k.httpResponseStatus.Store(int32(status))
}

func (k *k8sClient) checkStatusFromAnnotation(annotations map[string]string) {
	if statusStr, ok := annotations[statusAnnotation]; ok {
		status, err := strconv.Atoi(statusStr)
		if err == nil && status >= 100 && status <= 999 {
			k.setStatus(int32(status))
			return
		}
	}
	k.setStatus(200)
}

func (k *k8sClient) Run(ctx context.Context) {
	log.Infof("Run annotation checker")

	selector := fields.ParseSelectorOrDie("metadata.name==" + k.podName + ",status.phase==" + string(v1.PodRunning))
	watchlist := cache.NewListWatchFromClient(
		k.clientset.CoreV1().RESTClient(),
		string(v1.ResourcePods),
		k.podNS,
		selector, //fields.Everything(),
	)

	_, controller := cache.NewInformer(
		watchlist,
		&v1.Pod{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				log.Infof("Pod added %s %s", pod.Namespace, pod.Name)
				k.checkStatusFromAnnotation(pod.Annotations)
			},
			DeleteFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				log.Infof("Pod deleted %s %s", pod.Namespace, pod.Name)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				pod := newObj.(*v1.Pod)
				log.Infof("Pod updated %s %s", pod.Namespace, pod.Name)
				k.checkStatusFromAnnotation(pod.Annotations)

			},
		},
	)
	controller.Run(ctx.Done())
}
