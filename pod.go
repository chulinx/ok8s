package ok8s

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"time"
)

// The MultiClientSets contain k8s clients and metrics clients
type MultiClientSets struct {
	ClientSets
	MetricsClientSets
}

type PodType struct {
	MultiClientSets
	KResource
}

// NewPodAll can use all methods
func NewPodAll(kubeConfig string) *PodType {
	return &PodType{
		MultiClientSets: MultiClientSets{
			MetricsClientSets: MetricsClientSets{MClientSets: MetricsClientSet("/Users/lisong/.kube/ack-devops.conf")},
			ClientSets:        ClientSets{ClientSet: ClientSet(kubeConfig)},
		},
	}
}

// NewPod is Pod's K8sApi interface
func NewPod(cs ClientSets) K8sApi {
	pod := &PodType{
		MultiClientSets: MultiClientSets{ClientSets: cs},
	}
	return NewK8(pod)
}

func (p *PodType) Prefix(namespace string) interface{} {
	return p.ClientSet.CoreV1().Pods(namespace)
}

func (p *PodType) Create(namespace string, resource interface{}) (bool, error) {
	p.Pod = resource.(corev1.Pod)
	_, err := p.Prefix(namespace).(v1.PodInterface).Create(DefaultTimeOut(), &p.Pod, metav1.CreateOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (p *PodType) Get(namespace, name string) (bool, KResource) {
	pod, err := p.Prefix(namespace).(v1.PodInterface).Get(DefaultTimeOut(), name, metav1.GetOptions{})
	if pod != nil {
		p.KResource.Pod = *pod
	}
	if err != nil {
		return false, p.KResource
	}
	return true, p.KResource
}

func (p *PodType) IsExits(namespace, name string) bool {
	pod, err := p.Prefix(namespace).(v1.PodInterface).Get(DefaultTimeOut(), name, metav1.GetOptions{})
	if pod.Name == "" && err != nil {
		return false
	}
	return true
}

func (p *PodType) List(namespace string) (KResource, error) {
	pods, err := p.Prefix(namespace).(v1.PodInterface).List(DefaultTimeOut(), metav1.ListOptions{})
	if pods != nil && len(pods.Items) > 0 {
		p.KResource.PodList = *pods
	}
	if err != nil {
		return p.KResource, err
	}
	return p.KResource, nil
}

func (p *PodType) Delete(namespace, name string) (bool, error) {
	err := p.Prefix(namespace).(v1.PodInterface).Delete(DefaultTimeOut(), name, metav1.DeleteOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (p *PodType) Update(namespace string, resource interface{}) bool {
	p.Pod = resource.(corev1.Pod)
	_, err := p.Prefix(namespace).(v1.PodInterface).Update(DefaultTimeOut(), &p.Pod, metav1.UpdateOptions{})
	if err != nil {
		return false
	}
	return true
}

func (p *PodType) Watch(namespace string, eventFuncs cache.ResourceEventHandlerFuncs) {
	watchList := cache.NewListWatchFromClient(p.ClientSet.AppsV1().RESTClient(),
		"pods", namespace, fields.Everything())
	fmt.Println(watchList.List(metav1.ListOptions{}))
	_, controller := cache.NewInformer(watchList,
		&p.KResource.Pod,
		time.Second*0,
		eventFuncs,
	)
	stop := make(chan struct{})
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}

func (p *PodType) GetRequestResource(namespace, name string) map[string]int {
	var cpu, mem, storage int
	result := make(map[string]int)
	ok, kResource := p.Get(namespace, name)
	if !ok {
		return result
	}
	for _, container := range kResource.Pod.Spec.Containers {
		cpu = cpu + MetricsToInt(container.Resources.Requests.Cpu().String(), "m")
		mem = mem + MetricsToInt(AllToMi(container.Resources.Requests.Memory().String()), "Mi")
		storage = storage + MetricsToInt(container.Resources.Requests.Storage().String(), "")
	}
	result["cpu"], result["mem"], result["storage"] = cpu, mem, storage
	return result
}

func (p *PodType) GetLimitResource(namespace, name string) map[string]int {
	var cpu, mem, storage int
	result := make(map[string]int)
	ok, kResource := p.Get(namespace, name)
	if !ok {
		return result
	}
	for _, container := range kResource.Pod.Spec.Containers {
		cpu = cpu + MetricsToInt(container.Resources.Limits.Cpu().String(), "m")
		mem = mem + MetricsToInt(AllToMi(container.Resources.Limits.Memory().String()), "Mi")
		storage = storage + MetricsToInt(container.Resources.Limits.Storage().String(), "")
	}
	result["cpu"], result["mem"], result["storage"] = cpu, mem, storage
	return result
}

func (p *PodType) GetByLabel(namespace string, label string) (*corev1.PodList, error) {
	podList, err := p.Prefix(namespace).(v1.PodInterface).List(DefaultTimeOut(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return podList, err
	}
	return podList, nil
}

// GetMetrics return a map.
// This contain cpu(m)、mem(Mi) and storage
func (p *PodType) GetMetrics(namespace, name string) map[string]int {
	result := make(map[string]int)
	var cpu, mem, storage int
	podMetric, _ := p.MClientSets.MetricsV1beta1().PodMetricses(namespace).Get(DefaultTimeOut(), name, metav1.GetOptions{})
	for _, container := range podMetric.Containers {
		cpu = cpu + MetricsToInt(container.Usage.Cpu().String(), "m")
		mem = mem + MetricsToInt(container.Usage.Memory().String(), "Ki")/1000
		storage = storage + MetricsToInt(container.Usage.Storage().String(), "")
	}
	result["cpu"], result["mem"], result["storage"] = cpu, mem, storage
	return result
}

// GetContainerCpu return a pod all containers cpu(map[string]int)
func (p *PodType) GetContainerCpu(namespace, name string) (map[string]int, error) {
	containerCpu := make(map[string]int)
	podMetric, err := p.MClientSets.MetricsV1beta1().PodMetricses(namespace).Get(DefaultTimeOut(), name, metav1.GetOptions{})
	if err != nil {
		return containerCpu, err
	}
	for _, container := range podMetric.Containers {
		containerName := container.Name
		if containerName != "" {
			containerCpu[containerName] = MetricsToInt(container.Usage.Cpu().String(), "m")
		}
	}
	return containerCpu, nil
}

// GetSinglePodCpu return pod cpu(int)
func (p *PodType) GetSinglePodCpu(namespace, name string) (cpu int) {
	return p.GetMetrics(namespace, name)["cpu"]
}

// GetSinglePodMem return mem(int)
func (p *PodType) GetSinglePodMem(namespace, name string) int {
	return p.GetMetrics(namespace, name)["mem"]
}
