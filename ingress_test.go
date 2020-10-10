package ok8s

import (
	"k8s.io/api/networking/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
)

var (
	ingress = v1beta1.Ingress{
		ObjectMeta: v1.ObjectMeta{
			Name: "web",
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				{
					Host: "www.chulinx.com",
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Path: "/",
									Backend: v1beta1.IngressBackend{
										ServiceName: "web",
										ServicePort: intstr.IntOrString{IntVal: 80},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	i = NewIngress(NewTestClientSet())
)

func TestIngress_Create(t *testing.T) {
	ok, err := i.Create(ns, ingress)
	AssertError(ok, err, t)
}
