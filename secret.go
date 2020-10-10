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

type SecretType struct {
	ClientSets
	KResource
}

func NewSecret(cs ClientSets) K8sApi {
	return NewK8(&SecretType{
		ClientSets: cs,
	})
}

func (s *SecretType) Prefix(namespace string) interface{} {
	return s.ClientSet.CoreV1().Secrets(namespace)
}

func (s *SecretType) Create(namespace string, resource interface{}) (bool, error) {
	s.Secret = resource.(apicorev1.Secret)
	_, err := s.Prefix(namespace).(corev1.SecretInterface).Create(DefaultTimeOut(), &s.Secret, metav1.CreateOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *SecretType) Get(namespace, name string) (bool, KResource) {
	secret, err := s.Prefix(namespace).(corev1.SecretInterface).Get(DefaultTimeOut(), name, metav1.GetOptions{})
	if secret != nil {
		s.KResource.Secret = *secret
	}
	if err != nil {
		return false, s.KResource
	}
	return true, s.KResource
}

func (s *SecretType) IsExits(namespace, name string) bool {
	secret, err := s.Prefix(namespace).(corev1.SecretInterface).Get(DefaultTimeOut(), name, metav1.GetOptions{})
	if secret.Name == "" && err != nil {
		return false
	}
	return true
}

func (s *SecretType) Delete(namespace, name string) (bool, error) {
	err := s.Prefix(namespace).(corev1.SecretInterface).Delete(DefaultTimeOut(), name, metav1.DeleteOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *SecretType) List(namespace string) (KResource, error) {
	secrets, err := s.Prefix(namespace).(corev1.SecretInterface).List(DefaultTimeOut(), metav1.ListOptions{})
	if secrets != nil && len(secrets.Items) > 0 {
		s.KResource.SecretList = *secrets
	}
	if err != nil {
		return s.KResource, err
	}
	return s.KResource, nil
}

func (s *SecretType) Update(namespace string, resource interface{}) bool {
	s.Secret = resource.(apicorev1.Secret)
	_, err := s.Prefix(namespace).(corev1.SecretInterface).Update(DefaultTimeOut(), &s.Secret, metav1.UpdateOptions{})
	if err != nil {
		return false
	}
	return true
}

func (s *SecretType) Watch(namespace string, eventFuncs cache.ResourceEventHandlerFuncs) {
	watchList := cache.NewListWatchFromClient(s.ClientSet.CoreV1().RESTClient(),
		"secrets", namespace, fields.Everything())
	fmt.Println(watchList.List(metav1.ListOptions{}))
	_, controller := cache.NewInformer(watchList,
		&s.KResource.Secret,
		time.Second*0,
		eventFuncs,
	)
	stop := make(chan struct{})
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}
