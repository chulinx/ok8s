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

type ServiceType struct {
	ClientSets
	KResource
}

func NewService(cs ClientSets) K8sApi {
	s := &ServiceType{
		ClientSets: cs,
	}
	return NewK8(s)
}

func (s *ServiceType) Prefix(namespace string) interface{} {
	return s.ClientSet.CoreV1().Services(namespace)
}

func (s *ServiceType) Create(namespace string, resource interface{}) (bool, error) {
	s.Service = resource.(apicorev1.Service)
	_, err := s.Prefix(namespace).(corev1.ServiceInterface).Create(DefaultTimeOut(), &s.Service, metav1.CreateOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *ServiceType) Get(namespace, name string) (bool, KResource) {
	service, err := s.Prefix(namespace).(corev1.ServiceInterface).Get(DefaultTimeOut(), name, metav1.GetOptions{})
	s.KResource.Service = *service
	if err != nil {
		return false, s.KResource
	}
	return true, s.KResource
}

func (s *ServiceType) IsExits(namespace, name string) bool {
	service, err := s.Prefix(namespace).(corev1.ServiceInterface).Get(DefaultTimeOut(), name, metav1.GetOptions{})
	if service.Name == "" && err != nil {
		return false
	}
	return true
}

func (s *ServiceType) Delete(namespace, name string) (bool, error) {
	err := s.Prefix(namespace).(corev1.ServiceInterface).Delete(DefaultTimeOut(), name, metav1.DeleteOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *ServiceType) List(namespace string) (KResource, error) {
	services, err := s.Prefix(namespace).(corev1.ServiceInterface).List(DefaultTimeOut(), metav1.ListOptions{})
	s.KResource.ServiceList = *services
	if err != nil {
		return s.KResource, err
	}
	return s.KResource, nil
}

func (s *ServiceType) Update(namespace string, resource interface{}) bool {
	s.Service = resource.(apicorev1.Service)
	_, err := s.Prefix(namespace).(corev1.ServiceInterface).Update(DefaultTimeOut(), &s.Service, metav1.UpdateOptions{})
	if err != nil {
		return false
	}
	return true
}

func (s *ServiceType) Watch(namespace string, eventFuncs cache.ResourceEventHandlerFuncs) {
	watchList := cache.NewListWatchFromClient(s.ClientSet.CoreV1().RESTClient(),
		"services", namespace, fields.Everything())
	fmt.Println(watchList.List(metav1.ListOptions{}))
	_, controller := cache.NewInformer(watchList,
		&s.KResource.Service,
		time.Second*0,
		eventFuncs,
	)
	stop := make(chan struct{})
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}
