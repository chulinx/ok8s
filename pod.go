package ok8s

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"time"
)

type PodType struct {
	ClientSets
	KResource
}

func NewPod(cs ClientSets) K8sApi  {
	pod := &PodType{
		ClientSets:cs,
	}
	return NewK8(pod)
}

func (p *PodType) Prefix(namespace string) interface{} {
	return p.ClientSet.CoreV1().Pods(namespace)
}

func (p *PodType)Create(namespace string, resource interface{}) (bool, error)  {
	p.Pod = resource.(corev1.Pod)
	_,err :=p.Prefix(namespace).(v1.PodInterface).Create(DefaultTimeOut(),&p.Pod,metav1.CreateOptions{})
	if err != nil {
		return false, err
	}
	return true,nil
}

func (p *PodType)Get(namespace, name string) (bool, KResource)  {
	pod, err := p.Prefix(namespace).(v1.PodInterface).Get(DefaultTimeOut(), name, metav1.GetOptions{})
	if pod != nil {
		p.KResource.Pod = *pod
	}
	if err != nil {
		return false, p.KResource
	}
	return true, p.KResource
}

func (p *PodType)IsExits(namespace, name string) bool  {
	pod, err := p.Prefix(namespace).(v1.PodInterface).Get(DefaultTimeOut(), name, metav1.GetOptions{})
	if pod.Name == "" && err != nil {
		return false
	}
	return true
}

func (p *PodType)List(namespace string) (KResource, error)  {
	pods, err := p.Prefix(namespace).(v1.PodInterface).List(DefaultTimeOut(), metav1.ListOptions{})
	if pods != nil && len(pods.Items) > 0 {
		p.KResource.PodList = *pods
	}
	if err != nil {
		return p.KResource, err
	}
	return p.KResource, nil
}

func (p *PodType)Delete(namespace, name string) (bool, error)  {
	err := p.Prefix(namespace).(v1.PodInterface).Delete(DefaultTimeOut(), name, metav1.DeleteOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (p *PodType)Update(namespace string, resource interface{}) bool  {
	p.Pod = resource.(corev1.Pod)
	_, err := p.Prefix(namespace).(v1.PodInterface).Update(DefaultTimeOut(), &p.Pod, metav1.UpdateOptions{})
	if err != nil {
		return false
	}
	return true
}

func (p *PodType)Watch(namespace string, eventFuncs cache.ResourceEventHandlerFuncs)  {
	watchList := cache.NewListWatchFromClient(p.ClientSet.AppsV1().RESTClient(),
		"pods", namespace, fields.Everything())
	fmt.Println(watchList.List(metav1.ListOptions{}))
	_, controller := cache.NewInformer(watchList,
		&p.KResource.Pod,
		time.Second*0,
		eventFuncs,
	)
	stop := make(chan struct{})
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}

func (p *PodType)GetPodsByLabel(namespace string,label string) (*corev1.PodList,error) {
	podList,err := p.Prefix(namespace).(v1.PodInterface).List(DefaultTimeOut(),metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return podList,err
	}
	return podList,nil
}