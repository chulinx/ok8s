package kapi

import (
	appsv1 "k8s.io/api/apps/v1"
	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/api/networking/v1beta1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// K8sApi define resource operate method
type K8sApi interface {
	// The resource api interface
	Prefix(namespace string) interface{}
	// Create resource from namespace and resource struct
	/* For deployment k8s.io/api/apps/v1 v1.Deployment
	 * For configmap k8s.io/api/core/v1  v1.ConfigMap
	 */
	Create(namespace string, resource interface{}) (bool, error)
	// Get a resource from namespace and resource name
	Get(namespace, name string) (bool, KResource)
	// IsExits judge resource exits
	IsExits(namespace, name string) bool
	// Delete a resource from namespace and resource name
	Delete(namespace, name string) (bool, error)
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

// ClientSet kubernetes.Clientset
type ClientSets struct {
	ClientSet *kubernetes.Clientset
}

// Kubernetes Resource struct
type KResource struct {
	deployment     appsv1.Deployment
	deploymentList appsv1.DeploymentList
	configmap      apicorev1.ConfigMap
	configmapList  apicorev1.ConfigMapList
	service        apicorev1.Service
	serviceList    apicorev1.ServiceList
	secret 		   apicorev1.Secret
	secretList	   apicorev1.SecretList
	ingress 	   v1beta1.Ingress
	ingressList	   v1beta1.IngressList
}

func NewK8(inter K8sApi) K8sApi {
	return inter
}
