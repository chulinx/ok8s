package ok8s

import (
	"context"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"strconv"
	"strings"
	"time"
)

func NewClientSet(clientset *kubernetes.Clientset) *ClientSets {
	return &ClientSets{ClientSet: clientset}
}

func NewTestClientSet() ClientSets {
	return *NewClientSet(ClientSet("/Users/lisong/.kube/ack-devops.conf"))
}



func ClientSet(file string) *kubernetes.Clientset {
	// create the kubernetes clientSet
	clientset, err := kubernetes.NewForConfig(kubeConfig(file))
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

// MetricsClientSet use get pod cpu mem and other monitor metrics
func MetricsClientSet(file string) *versioned.Clientset  {
	clientSet,err := versioned.NewForConfig(kubeConfig(file))
	if err != nil {
		panic(err)
	}
	return clientSet
}

func DefaultTimeOut() context.Context {
	return TimeOut(time.Second * 3)
}

func TimeOut(time time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.TODO(), time)
	return ctx
}

func kubeConfig(file string) *rest.Config {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", file)
	if err != nil {
		panic(err.Error())
	}
	return config
}

func MetricsToInt(metric,sep string) int  {
	i,err :=strconv.Atoi(metric)
	if err != nil {
		i,err = strconv.Atoi(strings.Split(metric,sep)[0])
		if err != nil {
			return 0
		}
	}
	return i
}

func AllToMi(value string) string {
	switch  {
	case strings.Contains(value,"K") || strings.Contains(value,"Ki"):
		return fmt.Sprintf("%dMi", MetricsToInt(value, "Ki")/1000)
	case strings.Contains(value,"M") || strings.Contains(value,"Mi"):
		return value
	case strings.Contains(value,"G") || strings.Contains(value,"Gi"):
		return fmt.Sprintf("%dMi", MetricsToInt(value, "Ki")*1000)
	case strings.Contains(value,"T") || strings.Contains(value,"Ti"):
		return fmt.Sprintf("%dMi", MetricsToInt(value, "Ki")*1000*1000)
	default:
		fmt.Println("value format not support")
	}
	return ""
}