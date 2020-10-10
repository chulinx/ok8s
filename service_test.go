package ok8s

import (
	"fmt"
	apicorev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/cache"
	"testing"
)

var (
	s       = NewService(NewTestClientSet())
	service = apicorev1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name: "web",
		},
		Spec: apicorev1.ServiceSpec{
			Selector: map[string]string{
				"app": "nginx",
			},
			Type: apicorev1.ServiceTypeClusterIP,
			Ports: []apicorev1.ServicePort{
				{
					Name:       "web",
					Protocol:   apicorev1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.IntOrString{IntVal: 80},
				},
			},
		},
	}
)

func TestService_Create(t *testing.T) {
	ok, err := s.Create(ns, service)
	AssertError(ok, err, t)
}

func TestService_IsExits(t *testing.T) {
	ok := s.IsExits(ns, "web")
	Assert(ok, t)
}

func TestService_Watch(t *testing.T) {
	watchFuncs := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Printf("Create ServiceType %s\n", obj.(*apicorev1.Service).Name)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Printf("Delete ServiceType %s\n", obj.(*apicorev1.Service).Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Printf("Update ServiceType %s\n", oldObj.(*apicorev1.Service).Name)
		},
	}
	s.Watch(ns, watchFuncs)
}
