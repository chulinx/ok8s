package ok8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MetricsType struct {
	MetricsClientSets
	PodType
}

func NewMetrics(kubeConfig string) *MetricsType {
	return &MetricsType{
		MetricsClientSets: MetricsClientSets{MClientSets: MetricsClientSet(kubeConfig)},
		PodType: *NewPodFunc(kubeConfig),
	}
}

// default unit cpu is,m mem is Mi
func (m *MetricsType) GetRequestResource(namespace, name string) map[string]int64 {
	var cpu, mem, storage int64
	result := make(map[string]int64)
	ok, kResource := m.Get(namespace, name)
	if !ok {
		return result
	}
	for _, container := range kResource.Pod.Spec.Containers {
		cpu = cpu +container.Resources.Requests.Cpu().MilliValue()
		mem = mem + container.Resources.Requests.Memory().Value()/1024/1024
		storage = storage + container.Resources.Requests.Storage().Value()
	}
	result["cpu"], result["mem"], result["storage"] = cpu, mem, storage
	return result
}

func (m *MetricsType) GetLimitResource(namespace, name string) map[string]int64 {
	var cpu, mem, storage int64
	result := make(map[string]int64)
	ok, kResource := m.Get(namespace, name)
	if !ok {
		return result
	}
	for _, container := range kResource.Pod.Spec.Containers {
		cpu = cpu +container.Resources.Limits.Cpu().MilliValue()
		mem = mem + container.Resources.Limits.Memory().Value()/1024/1024
		storage = storage + container.Resources.Limits.Storage().Value()
	}
	result["cpu"], result["mem"], result["storage"] = cpu, mem, storage
	return result
}

// GetPodMetrics return a map.
// This contain cpu(m)、mem(Mi) and storage
func (m *MetricsType) GetPodMetrics(namespace, name string) map[string]int64 {
	result := make(map[string]int64)
	var cpu, mem, storage int64
	podMetric, _ := m.MClientSets.MetricsV1beta1().PodMetricses(namespace).Get(DefaultTimeOut(), name, metav1.GetOptions{})
	for _, container := range podMetric.Containers {
		cpu = cpu + container.Usage.Cpu().MilliValue()
		mem = mem + container.Usage.Memory().Value()/1024/1024
		storage = storage + container.Usage.Storage().Value()
	}
	result["cpu"], result["mem"], result["storage"] = cpu, mem, storage
	return result
}

// GetContainerCpu return a pod all containers cpu(map[string]int)
func (m *MetricsType) GetContainerCpu(namespace, name string) (map[string]int64, error) {
	containerCpu := make(map[string]int64)
	podMetric, err := m.MClientSets.MetricsV1beta1().PodMetricses(namespace).Get(DefaultTimeOut(), name, metav1.GetOptions{})
	if err != nil {
		return containerCpu, err
	}
	for _, container := range podMetric.Containers {
		containerName := container.Name
		if containerName != "" {
			containerCpu[containerName] = container.Usage.Cpu().MilliValue()
		}
	}
	return containerCpu, nil
}

// GetSinglePodCpuUsage return pod cpu(int)
func (m *MetricsType) GetSinglePodCpuUsage(namespace, name string) (cpu int64) {
	return m.GetPodMetrics(namespace, name)["cpu"]
}

// GetSinglePodMemUsage return mem(int)
func (m *MetricsType) GetSinglePodMemUsage(namespace, name string) int64 {
	return m.GetPodMetrics(namespace, name)["mem"]
}
