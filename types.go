package kapi

import (
	appsv1 "k8s.io/api/apps/v1"
	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type K8sApier interface {
	Prefix(namespace string) interface{}
	Create(namespace string, resource interface{}) (bool, error)
	Get(namespace,name string) (bool,KResource)
	Delete(namespace, name string) bool
	List(namespace string) (KResource, error)
	Update(namespace string,resource interface{}) bool
	//Labels(deployment, namespace string) map[string]map[string]string
	Watch(namespace string, eventFun cache.ResourceEventHandlerFuncs)
}

func NewClientSet(clientset *kubernetes.Clientset) *ClientSets {
	return &ClientSets{ClientSet: clientset}
}

// ClientSet kubernetes.Clientset
type ClientSets struct {
	ClientSet *kubernetes.Clientset
}

// kubernetes Resource struct
type KResource struct {
	deployment     appsv1.Deployment
	deploymentList appsv1.DeploymentList
	configmap 	   apicorev1.ConfigMap
	configmapList  apicorev1.ConfigMapList
}
