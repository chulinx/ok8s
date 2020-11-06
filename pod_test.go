package ok8s

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

var (
	p       = NewPod(NewTestClientSet())
	podName = "nginx"
	podInstance = NewPodAll("/Users/lisong/.kube/ack-devops.conf")
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

func TestPodType_GetPodsByLabel(t *testing.T) {
	pod := &PodType{MultiClientSets: MultiClientSets{
		ClientSets: NewTestClientSet(),
	},
	}
	plist, err := pod.GetByLabel(ns, "web=true")
	fmt.Println(plist, err)
}

func TestPodType_GetPodMetrics(t *testing.T) {
	podInstance.GetMetrics("thanos", "federation-scraper-ack-devops-0")
}

func TestPodType_GetPodCpu(t *testing.T) {
	cpuUsage:=podInstance.GetSinglePodCpu("thanos", "federation-scraper-ack-devops-0")
	fmt.Println(cpuUsage)
}

func TestPodType_GetRequestResourceAndGetLimitResource(t *testing.T) {
	fmt.Println(podInstance.GetRequestResource("zx","nginx-6dfd4b854c-j46jr"))
	fmt.Println(podInstance.GetLimitResource("zx","nginx-6dfd4b854c-j46jr"))

}
