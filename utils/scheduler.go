package utils

import (
	"context"
	"fmt"
	"time"

	"com.ak.gooverlord/partitioner"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func getNodesFromK8sDeployment() []string {
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

func init() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for _ = range ticker.C {
			ips := getNodesFromK8sDeployment()
			go func(ips []string) {
				newCHR := partitioner.NewConsistenHashRing()
				for _, ip := range ips {
					ip = fmt.Sprintf("http://%s:7700", ip)
					newCHR.AddNode(ip)
				}
				partitioner.CHR.RWLock.Lock()
				defer partitioner.CHR.RWLock.Unlock()
				partitioner.CHR = newCHR
				fmt.Println(partitioner.CHR.RealNodesSet)
			}(ips)
		}
	}()
}
