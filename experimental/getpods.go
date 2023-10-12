package main

// If we want to enable clustering, this could be useful in maintining cluster state.
// Go version of get_pods.rb, with help from Derek and ChatGPTs

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	serviceName := "proxysql"

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, _ := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", serviceName),
	})
	for _, pod := range pods.Items {
		fmt.Printf("%s - %s\n", pod.Name, pod.Status.PodIP)
	}
}
