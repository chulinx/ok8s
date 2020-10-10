package ok8s

import (
	"fmt"
	apicorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"time"
)

type ConfigMapType struct {
	ClientSets
	KResource
}

func NewConfigMap(cs ClientSets) K8sApi {
	c := &ConfigMapType{
		ClientSets: cs,
	}
	return NewK8(c)
}

func (c *ClientSets) Prefix(namespace string) interface{} {
	return c.ClientSet.CoreV1().ConfigMaps(namespace)
}

// The resource struct is apicorev1.ConfigMapType
func (c *ConfigMapType) Create(namespace string, resource interface{}) (bool, error) {
	c.KResource.Configmap = resource.(apicorev1.ConfigMap)
	_, err := c.Prefix(namespace).(corev1.ConfigMapInterface).Create(DefaultTimeOut(), &c.KResource.Configmap, metav1.CreateOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

// The resource struct is apicorev1.ConfigMapType
func (c *ConfigMapType) Update(namespace string, resource interface{}) bool {
	c.KResource.Configmap = resource.(apicorev1.ConfigMap)
	_, err := c.Prefix(namespace).(corev1.ConfigMapInterface).Update(DefaultTimeOut(), &c.KResource.Configmap, metav1.UpdateOptions{})
	if err != nil {
		return false
	}
	return true
}

func (c *ConfigMapType) IsExits(namespace, name string) bool {
	configmap, err := c.Prefix(namespace).(corev1.ConfigMapInterface).Get(DefaultTimeOut(), name, metav1.GetOptions{})
	if configmap.Name == "" && err != nil {
		return false
	}
	return true
}

// Get return apicorev1.ConfigMapType
func (c *ConfigMapType) Get(namespace, name string) (bool, KResource) {
	cm, err := c.Prefix(namespace).(corev1.ConfigMapInterface).Get(DefaultTimeOut(), name, metav1.GetOptions{})
	c.Configmap = *cm
	if err != nil {
		return false, c.KResource
	}
	return true, c.KResource
}

// List return multiple apicorev1.ConfigMapType
func (c *ConfigMapType) List(namespace string) (KResource, error) {
	cms, err := c.Prefix(namespace).(corev1.ConfigMapInterface).List(DefaultTimeOut(), metav1.ListOptions{})
	c.ConfigmapList = *cms
	if err != nil {
		return c.KResource, err
	}
	return c.KResource, nil
}

func (c *ConfigMapType) Delete(namespace, name string) (bool, error) {
	err := c.Prefix(namespace).(corev1.ConfigMapInterface).Delete(DefaultTimeOut(), name, metav1.DeleteOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (d *ConfigMapType) Watch(namespace string, eventFuncs cache.ResourceEventHandlerFuncs) {
	watchList := cache.NewListWatchFromClient(d.ClientSet.CoreV1().RESTClient(),
		"configmaps", namespace, fields.Everything())
	fmt.Println(watchList.List(metav1.ListOptions{}))
	_, controller := cache.NewInformer(watchList,
		&d.KResource.Configmap,
		time.Second*0,
		eventFuncs,
	)
	stop := make(chan struct{})
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}
