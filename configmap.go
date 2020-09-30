package kapi

import (
	"context"
	"fmt"
	apicorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"time"
)

type ConfigMap struct {
	ClientSets
	KResource
}


func (c *ClientSets) Prefix(namespace string) interface{} {
	return c.ClientSet.CoreV1().ConfigMaps(namespace)
}

func (c *ConfigMap) Create(namespace string, resource interface{}) (bool, error) {
	c.KResource.configmap = resource.(apicorev1.ConfigMap)
	_, err := c.Prefix(namespace).(corev1.ConfigMapInterface).Create(context.TODO(), &c.KResource.configmap, metav1.CreateOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *ConfigMap)Update(namespace string,resource interface{}) bool {
	c.KResource.configmap = resource.(apicorev1.ConfigMap)
	_, err :=c.Prefix(namespace).(corev1.ConfigMapInterface).Update(context.TODO(),&c.KResource.configmap,metav1.UpdateOptions{})
	if err != nil {
		return false
	}
	return true
}


func (c *ConfigMap) Get(namespace,name string) (bool,KResource) {
	cm,err := c.Prefix(namespace).(corev1.ConfigMapInterface).Get(context.TODO(),name,metav1.GetOptions{})
	c.configmap = *cm
	if err != nil {
		return false, c.KResource
	}
	return true,c.KResource
}

func (c *ConfigMap)List(namespace string) (KResource, error)  {
	cms,err :=  c.Prefix(namespace).(corev1.ConfigMapInterface).List(context.TODO(),metav1.ListOptions{})
	c.configmapList = *cms
	if err != nil {
		return c.KResource, err
	}
	return c.KResource, nil
}

func (c *ConfigMap)Delete(namespace, name string) bool  {
	err := c.Prefix(namespace).(corev1.ConfigMapInterface).Delete(context.TODO(),name,metav1.DeleteOptions{})
	if err != nil {
		return false
	}
	return true
}

func (d *ConfigMap) Watch(namespace string, eventFun cache.ResourceEventHandlerFuncs) {
	watchList := cache.NewListWatchFromClient(d.ClientSet.CoreV1().RESTClient(),
		"configmaps", namespace, fields.Everything())
	fmt.Println(watchList.List(metav1.ListOptions{}))
	_, controller := cache.NewInformer(watchList,
		&apicorev1.ConfigMap{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    eventFun.AddFunc,
			DeleteFunc: eventFun.DeleteFunc,
			UpdateFunc: eventFun.UpdateFunc,
		})
	stop := make(chan struct{})
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}