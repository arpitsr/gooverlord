package utils

import (
	"context"
	"fmt"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8SDiscovery struct{}

var k8sDiscovery *K8SDiscovery
var onceDiscover sync.Once

func (kn *K8SDiscovery) GetIndexNodes() []string {
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		panic(err.Error())
	}

	// Create a new Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Get the pod IPs for the deployment "my-deployment" in the "default" namespace
	service := "meilisearch"
	namespace := "default"
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(),
		metav1.ListOptions{
			LabelSelector: fmt.Sprintf("app.kubernetes.io/finder=%s", service),
		})
	if err != nil {
		panic(err.Error())
	}

	var ips = make([]string, 0)
	for _, pod := range pods.Items {
		fmt.Println(pod.Status.PodIP)
		ips = append(ips, pod.Status.PodIP)
	}
	return ips
}

func NewK8SDiscovery() *K8SDiscovery {
	onceDiscover.Do(func() {
		k8sDiscovery = &K8SDiscovery{}
	})
	return k8sDiscovery
}
