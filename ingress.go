package ok8s

import (
	"fmt"
	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	networkv1beta1 "k8s.io/client-go/kubernetes/typed/networking/v1beta1"
	"k8s.io/client-go/tools/cache"
	"time"
)

type IngressType struct {
	ClientSets
	KResource
}

func NewIngress(cs ClientSets) K8sApi {
	return NewK8(&IngressType{
		ClientSets: cs,
	})
}

func (i *IngressType) Prefix(namespace string) interface{} {
	return i.ClientSet.NetworkingV1beta1().Ingresses(namespace)
}

func (i *IngressType) Create(namespace string, resource interface{}) (bool, error) {
	i.Ingress = resource.(v1beta1.Ingress)
	_, err := i.Prefix(namespace).(networkv1beta1.IngressInterface).Create(DefaultTimeOut(), &i.Ingress, metav1.CreateOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (i *IngressType) Get(namespace, name string) (bool, KResource) {
	ingress, err := i.Prefix(namespace).(networkv1beta1.IngressInterface).Get(DefaultTimeOut(), name, metav1.GetOptions{})
	if ingress != nil {
		i.KResource.Ingress = *ingress
	}
	if err != nil {
		return false, i.KResource
	}
	return true, i.KResource
}

func (i *IngressType) IsExits(namespace, name string) bool {
	ingress, err := i.Prefix(namespace).(networkv1beta1.IngressInterface).Get(DefaultTimeOut(), name, metav1.GetOptions{})
	if ingress.Name == "" && err != nil {
		return false
	}
	return true
}

func (i *IngressType) Delete(namespace, name string) (bool, error) {
	err := i.Prefix(namespace).(networkv1beta1.IngressInterface).Delete(DefaultTimeOut(), name, metav1.DeleteOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (i *IngressType) List(namespace string) (KResource, error) {
	ingress, err := i.Prefix(namespace).(corev1.SecretInterface).List(DefaultTimeOut(), metav1.ListOptions{})
	if ingress != nil && len(ingress.Items) > 0 {
		i.KResource.SecretList = *ingress
	}
	if err != nil {
		return i.KResource, err
	}
	return i.KResource, nil
}

func (i *IngressType) Update(namespace string, resource interface{}) bool {
	i.Secret = resource.(apicorev1.Secret)
	_, err := i.Prefix(namespace).(corev1.SecretInterface).Update(DefaultTimeOut(), &i.Secret, metav1.UpdateOptions{})
	if err != nil {
		return false
	}
	return true
}

func (i *IngressType) Watch(namespace string, eventFuncs cache.ResourceEventHandlerFuncs) {
	watchList := cache.NewListWatchFromClient(i.ClientSet.NetworkingV1beta1().RESTClient(),
		"ingresses", namespace, fields.Everything())
	fmt.Println(watchList.List(metav1.ListOptions{}))
	_, controller := cache.NewInformer(watchList,
		&i.KResource.Ingress,
		time.Second*0,
		eventFuncs,
	)
	stop := make(chan struct{})
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}
