package ok8s

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
	/* For DeploymentType k8s.io/api/apps/v1 v1.deployment
	 * For Configmap k8s.io/api/core/v1  v1.configmap
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
	//Labels(DeploymentType, namespace string) map[string]map[string]string
	// Watch a resource
	/* Example:
	*  eventFuncs := new(cache.ResourceEventHandlerFuncs)
	*  eventFuncs.DeleteFunc= func(obj interface{}) {
	*  		DeploymentType := obj.(*appsv1.DeploymentType)
	*  		log.Printf("delete DeploymentType %s",DeploymentType.Name)
	*  		fmt.Println(Deployment.Labels,DeploymentType.Spec.Template.Labels)
	*  }
	*  eventFuncs.AddFunc= func(obj interface{}) {
	*  		DeploymentType := obj.(*appsv1.DeploymentType)
	*  		log.Printf("add DeploymentType %s",DeploymentType.Name)
	*  		fmt.Println(DeploymentType.Labels,DeploymentType.Spec.Template.Labels)
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
//
// Unified all resource
type KResource struct {
	Deployment     appsv1.Deployment
	DeploymentList appsv1.DeploymentList
	Configmap      apicorev1.ConfigMap
	ConfigmapList  apicorev1.ConfigMapList
	Service        apicorev1.Service
	ServiceList    apicorev1.ServiceList
	Secret         apicorev1.Secret
	SecretList     apicorev1.SecretList
	Ingress        v1beta1.Ingress
	IngressList    v1beta1.IngressList
}

func NewK8(inter K8sApi) K8sApi {
	return inter
}
