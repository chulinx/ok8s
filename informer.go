package kapi

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"time"
)

func NewInformer(sets *kubernetes.Clientset)  {
	factory :=informers.NewSharedInformerFactoryWithOptions(sets,time.Second,informers.WithNamespace("zx"))
	informer :=  factory.Core().V1().Pods().Informer()
	stop := make(chan struct{})
	defer close(stop)
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			p := obj.(*v1.Pod)
			fmt.Printf("add pod %s\n",p.Name)
		},
		DeleteFunc: func(obj interface{}) {
			p := obj.(*v1.Pod)
			fmt.Printf("add pod %s\n",p.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			p := oldObj.(*v1.Pod)
			fmt.Printf("add pod %s\n",p.Name)
		},
	})
	informer.Run(stop)
}