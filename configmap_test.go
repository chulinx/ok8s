package kapi

import (
	"fmt"
	apicorev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"testing"
)

var (
	cm = apicorev1.ConfigMap {
		ObjectMeta: v1.ObjectMeta{
			Name: "web-cm",
		},
		Data: map[string]string{
			"user":"root",
			"password":"123456",
		},
	}
	c = NewConfigMap(NewTestClientSet())
)

func TestConfigMap_Create(t *testing.T) {
	ok,err := c.Create(ns,cm)
	AssertError(ok,err,t)
}

func TestConfigMap_Watch(t *testing.T) {
	watchFuncs := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Printf("Create cm %s\n", obj.(*apicorev1.ConfigMap).Name)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Printf("Delete cm %s\n", obj.(*apicorev1.ConfigMap).Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Printf("Update cm %s\n", oldObj.(*apicorev1.ConfigMap).Name)
		},
	}
	c.Watch(ns,watchFuncs)
}
