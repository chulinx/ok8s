package kapi

import (
	appsv1 "k8s.io/api/apps/v1"
	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// K8sApi define resource operate method
type K8sApi interface {
	// The resource api interface
	Prefix(namespace string) interface{}
	// Create resource from namespace and resource struct
	Create(namespace string, resource interface{}) (bool, error)
	// Get a resource from namespace and resource name
	Get(namespace, name string) (bool, KResource)
	// Delete a resource from namespace and resource name
	Delete(namespace, name string) bool
	// List multiple resource from one namespace
	List(namespace string) (KResource, error)
	// Update a resource
	Update(namespace string, resource interface{}) bool
	//Labels(deployment, namespace string) map[string]map[string]string
	// Watch a resource
	/* Example:
	*  eventFuncs := new(cache.ResourceEventHandlerFuncs)
	*  eventFuncs.DeleteFunc= func(obj interface{}) {
	*  		deployment := obj.(*appsv1.Deployment)
	*  		log.Printf("delete deployment %s",deployment.Name)
	*  		fmt.Println(deployment.Labels,deployment.Spec.Template.Labels)
	*  }
	*  eventFuncs.AddFunc= func(obj interface{}) {
	*  		deployment := obj.(*appsv1.Deployment)
	*  		log.Printf("add deployment %s",deployment.Name)
	*  		fmt.Println(deployment.Labels,deployment.Spec.Template.Labels)
	*  }
	*  Watch("default",eventFuncs)
	 */
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
	configmap      apicorev1.ConfigMap
	configmapList  apicorev1.ConfigMapList
}
