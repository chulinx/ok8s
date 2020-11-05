package ok8s

import (
	"context"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

func NewClientSet(clientset *kubernetes.Clientset) *ClientSets {
	return &ClientSets{ClientSet: clientset}
}

func NewTestClientSet() ClientSets {
	return *NewClientSet(ClientSet("/Users/lisong/.kube/ack-devops.conf"))
}

func ClientSet(file string) *kubernetes.Clientset {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", file)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func DefaultTimeOut() context.Context {
	return TimeOut(time.Second * 3)
}

func TimeOut(time time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.TODO(), time)
	return ctx
}
