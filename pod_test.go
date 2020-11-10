package ok8s

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

var (
	p       = NewPod(NewTestClientSet())
	podName = "nginx"
	podInstance = NewPodFunc("/Users/lisong/.kube/ack-devops.conf")
	pod     = corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: podName},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				corev1.Container{
					Name:  podName,
					Image: "nginx:latest",
					Ports: []corev1.ContainerPort{
						corev1.ContainerPort{
							Name:          "web",
							ContainerPort: 80,
							Protocol:      corev1.ProtocolTCP,
						},
					},
				},
			},
		},
	}
)

func TestPodType_Create(t *testing.T) {
	ok, err := p.Create(ns, pod)
	AssertError(ok, err, t)
}


