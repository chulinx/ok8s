package ok8s

import (
	"context"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func BuildConfig(kubeConfigPath ...string) (*rest.Config, error) {
	// 尝试使用InClusterConfig
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	// 如果在集群外部运行，
	if len(kubeConfigPath) > 0 && len(kubeConfigPath[0]) > 0 {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath[0])
		if err != nil {
			return nil, fmt.Errorf("error building kubeconfig: %s", err)
		}

		return config, nil
	}
	// 如果没有提供kubeconfig路径，则尝试使用$HOME/.kube/config文件
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("error building kubeconfig: %s", err)
	}
	return config, nil
}

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
func MetricsClientSet(file string) *versioned.Clientset {
	clientSet, err := versioned.NewForConfig(kubeConfig(file))
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

func MetricsToInt(metric, sep string) int {
	i, err := strconv.Atoi(metric)
	if err != nil {
		i, err = strconv.Atoi(strings.Split(metric, sep)[0])
		if err != nil {
			return 0
		}
	}
	return i
}

func AllToMi(value string) string {
	switch {
	case strings.Contains(value, "K") || strings.Contains(value, "Ki"):
		return fmt.Sprintf("%dMi", MetricsToInt(value, "Ki")/1000)
	case strings.Contains(value, "M") || strings.Contains(value, "Mi"):
		return value
	case strings.Contains(value, "G") || strings.Contains(value, "Gi"):
		return fmt.Sprintf("%dMi", MetricsToInt(value, "Ki")*1000)
	case strings.Contains(value, "T") || strings.Contains(value, "Ti"):
		return fmt.Sprintf("%dMi", MetricsToInt(value, "Ki")*1000*1000)
	default:
		fmt.Println("value format not support")
	}
	return ""
}
