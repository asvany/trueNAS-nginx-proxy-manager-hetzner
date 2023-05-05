package main

import (
	"context"
	"log"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// corev1 "k8s.io/api/core/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetPodID() string {
	var config *rest.Config
	var err error

	namespace := os.Getenv("POD_NAMESPACE")
	if namespace == "" {
		namespace = "ix-nginx-proxy-manager"
	}
	log.Println("namespace is: ", namespace)

	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = "/etc/rancher/k3s/k3s.yaml"
	}
	log.Println("kubeconfig path is: ", kubeconfigPath)

	if _, err := os.Stat(kubeconfigPath); !os.IsNotExist(err) {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			panic(err.Error())
		}
	} else {
		config = nil
		log.Fatalln("kubeconfig not found")

	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("error building kubernetes clientset: ", err.Error())

	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	if len(pods.Items) == 0 {
		log.Fatalln("no pods found")
	} else {
		if len(pods.Items) > 1 {
			log.Println("more than one pod found")
		}
	}

	pod := pods.Items[0]

	log.Println("Pod name:", pod.Name)

	pod_data, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), pod.Name, metav1.GetOptions{})
	if err != nil {
		log.Println("can't get pod ip")
		panic(err.Error())
	}

	return pod_data.Status.PodIP
}
