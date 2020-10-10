package kapi

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"time"
)

type Deployment struct {
	ClientSets
	KResource
}

func NewDeployment(cs ClientSets) K8sApi {
	d := &Deployment{
		ClientSets: cs,
	}
	return NewK8(d)
}

func (d *Deployment) Prefix(namespace string) interface{} {
	return d.ClientSet.AppsV1().Deployments(namespace)
}

func (d *Deployment) Create(namespace string, resource interface{}) (bool, error) {
	d.deployment = resource.(appsv1.Deployment)
	_, err := d.Prefix(namespace).(v1.DeploymentInterface).Create(DefaultTimeOut(), &d.deployment, metav1.CreateOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (d *Deployment)IsExits(namespace, name string) bool  {
	deploy,err :=d.Prefix(namespace).(v1.DeploymentInterface).Get(DefaultTimeOut(),name,metav1.GetOptions{})
	if deploy == nil && err != nil {
		return false
	}
	return true
}

func (d *Deployment) Delete(namespace, name string) (bool, error) {
	err := d.Prefix(namespace).(v1.DeploymentInterface).Delete(DefaultTimeOut(), name, metav1.DeleteOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (d *Deployment) List(namespace string) (KResource, error) {
	deploys, err := d.Prefix(namespace).(v1.DeploymentInterface).List(DefaultTimeOut(), metav1.ListOptions{})
	if deploys != nil && len(deploys.Items) > 0 {
		d.KResource.deploymentList = *deploys
	}
	if err != nil {
		return d.KResource, err
	}
	return d.KResource, nil
}

func (d *Deployment) Get(namespace, name string) (bool, KResource) {
	deploy, err := d.Prefix(namespace).(v1.DeploymentInterface).Get(DefaultTimeOut(), name, metav1.GetOptions{})
	if deploy != nil {
		d.KResource.deployment = *deploy
	}
	if err != nil {
		return false, d.KResource
	}
	return true, d.KResource
}

func (d *Deployment) Labels(deployment, namespace string) map[string]map[string]string {
	labels := make(map[string]map[string]string)
	deploy, err := d.Prefix(namespace).(v1.DeploymentInterface).Get(DefaultTimeOut(), deployment, metav1.GetOptions{})
	if err != nil {
		return labels
	}
	labels["deployment"] = deploy.Labels
	labels["pod"] = deploy.Spec.Template.Labels
	return labels
}

func (d *Deployment) Update(namespace string, resource interface{}) bool {
	d.deployment = resource.(appsv1.Deployment)
	_, err := d.Prefix(namespace).(v1.DeploymentInterface).Update(DefaultTimeOut(), &d.deployment, metav1.UpdateOptions{})
	if err != nil {
		return false
	}
	return true
}

func (d *Deployment) Watch(namespace string, eventFuncs cache.ResourceEventHandlerFuncs) {
	watchList := cache.NewListWatchFromClient(d.ClientSet.AppsV1().RESTClient(),
		"deployments", namespace, fields.Everything())
	fmt.Println(watchList.List(metav1.ListOptions{}))
	_, controller := cache.NewInformer(watchList,
		&d.KResource.deployment,
		time.Second*0,
		eventFuncs,
	)
	stop := make(chan struct{})
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}
