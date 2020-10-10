package kapi

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"testing"
)

var (
	ns = "test"
	d = NewDeployment(NewTestClientSet())
	replicas = int32(1)
	deploy = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "web"},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "test",
							Image: "nginx:latest",
							Ports: []v1.ContainerPort{
								{
									Name:          "web",
									ContainerPort: 80,
									Protocol:      v1.ProtocolTCP,
								},
							},
						},
					},
				},
			},
		},
	}
)

func TestDeployment_Create(t *testing.T) {
	ok, err := d.Create(ns, deploy)
	AssertError(ok, err, t)
}

func TestDeployment_Update(t *testing.T) {
	replicas = int32(2)
	ok := d.Update(ns,deploy)
	Assert(ok,t)
}

func TestDeployment_IsExits(t *testing.T) {
	ok := d.IsExits(ns,"web")
	Assert(ok,t)
}



func TestDeployment_Get(t *testing.T) {
	ok,krs := d.Get(ns,"web")
	if krs.deployment.Name != "" {
		Assert(ok,t)
	}
}

func TestDeployment_Delete(t *testing.T) {
	ok, err := d.Delete(ns, "web")
	AssertError(ok, err, t)
}

func TestDeployment_Watch(t *testing.T) {
	watchFuncs := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Printf("Create deployment %s\n", obj.(*appsv1.Deployment).Name)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Printf("Delete deployment %s\n", obj.(*appsv1.Deployment).Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Printf("Update deployment %s\n", oldObj.(*appsv1.Deployment).Name)
		},
	}
	d.Watch("test", watchFuncs)
}

// 断言assert
func Assert(ok bool,t *testing.T)  {
	if !ok {
		t.Fatalf("Failed")
	}else {
		t.Logf("Success")
	}
}


func AssertError(ok bool, err error, t *testing.T) {
	if !ok {
		t.Fatalf(err.Error())
	} else {
		t.Logf("Success")
	}
}
