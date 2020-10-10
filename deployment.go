package ok8s

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"time"
)

type DeploymentType struct {
	ClientSets
	KResource
}

func NewDeployment(cs ClientSets) K8sApi {
	d := &DeploymentType{
		ClientSets: cs,
	}
	return NewK8(d)
}

func (d *DeploymentType) Prefix(namespace string) interface{} {
	return d.ClientSet.AppsV1().Deployments(namespace)
}

func (d *DeploymentType) Create(namespace string, resource interface{}) (bool, error) {
	d.Deployment = resource.(appsv1.Deployment)
	_, err := d.Prefix(namespace).(v1.DeploymentInterface).Create(DefaultTimeOut(), &d.Deployment, metav1.CreateOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (d *DeploymentType) IsExits(namespace, name string) bool {
	deploy, err := d.Prefix(namespace).(v1.DeploymentInterface).Get(DefaultTimeOut(), name, metav1.GetOptions{})
	if deploy.Name == "" && err != nil {
		return false
	}
	return true
}

func (d *DeploymentType) Delete(namespace, name string) (bool, error) {
	err := d.Prefix(namespace).(v1.DeploymentInterface).Delete(DefaultTimeOut(), name, metav1.DeleteOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (d *DeploymentType) List(namespace string) (KResource, error) {
	deploys, err := d.Prefix(namespace).(v1.DeploymentInterface).List(DefaultTimeOut(), metav1.ListOptions{})
	if deploys != nil && len(deploys.Items) > 0 {
		d.KResource.DeploymentList = *deploys
	}
	if err != nil {
		return d.KResource, err
	}
	return d.KResource, nil
}

func (d *DeploymentType) Get(namespace, name string) (bool, KResource) {
	deploy, err := d.Prefix(namespace).(v1.DeploymentInterface).Get(DefaultTimeOut(), name, metav1.GetOptions{})
	if deploy != nil {
		d.KResource.Deployment = *deploy
	}
	if err != nil {
		return false, d.KResource
	}
	return true, d.KResource
}

func (d *DeploymentType) Labels(deployment, namespace string) map[string]map[string]string {
	labels := make(map[string]map[string]string)
	deploy, err := d.Prefix(namespace).(v1.DeploymentInterface).Get(DefaultTimeOut(), deployment, metav1.GetOptions{})
	if err != nil {
		return labels
	}
	labels["DeploymentType"] = deploy.Labels
	labels["pod"] = deploy.Spec.Template.Labels
	return labels
}

func (d *DeploymentType) Update(namespace string, resource interface{}) bool {
	d.Deployment = resource.(appsv1.Deployment)
	_, err := d.Prefix(namespace).(v1.DeploymentInterface).Update(DefaultTimeOut(), &d.Deployment, metav1.UpdateOptions{})
	if err != nil {
		return false
	}
	return true
}

func (d *DeploymentType) Watch(namespace string, eventFuncs cache.ResourceEventHandlerFuncs) {
	watchList := cache.NewListWatchFromClient(d.ClientSet.AppsV1().RESTClient(),
		"deployments", namespace, fields.Everything())
	fmt.Println(watchList.List(metav1.ListOptions{}))
	_, controller := cache.NewInformer(watchList,
		&d.KResource.Deployment,
		time.Second*0,
		eventFuncs,
	)
	stop := make(chan struct{})
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}
