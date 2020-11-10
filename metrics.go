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

func (m *MetricsType) GetRequestResource(namespace, name string) map[string]int {
	var cpu, mem, storage int
	result := make(map[string]int)
	ok, kResource := m.Get(namespace, name)
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

func (m *MetricsType) GetLimitResource(namespace, name string) map[string]int {
	var cpu, mem, storage int
	result := make(map[string]int)
	ok, kResource := m.Get(namespace, name)
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



// GetMetrics return a map.
// This contain cpu(m)、mem(Mi) and storage
func (m *MetricsType) GetMetrics(namespace, name string) map[string]int {
	result := make(map[string]int)
	var cpu, mem, storage int
	podMetric, _ := m.MClientSets.MetricsV1beta1().PodMetricses(namespace).Get(DefaultTimeOut(), name, metav1.GetOptions{})
	for _, container := range podMetric.Containers {
		cpu = cpu + MetricsToInt(container.Usage.Cpu().String(), "m")
		mem = mem + MetricsToInt(container.Usage.Memory().String(), "Ki")/1000
		storage = storage + MetricsToInt(container.Usage.Storage().String(), "")
	}
	result["cpu"], result["mem"], result["storage"] = cpu, mem, storage
	return result
}

// GetContainerCpu return a pod all containers cpu(map[string]int)
func (m *MetricsType) GetContainerCpu(namespace, name string) (map[string]int, error) {
	containerCpu := make(map[string]int)
	podMetric, err := m.MClientSets.MetricsV1beta1().PodMetricses(namespace).Get(DefaultTimeOut(), name, metav1.GetOptions{})
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
func (m *MetricsType) GetSinglePodCpu(namespace, name string) (cpu int) {
	return m.GetMetrics(namespace, name)["cpu"]
}

// GetSinglePodMem return mem(int)
func (m *MetricsType) GetSinglePodMem(namespace, name string) int {
	return m.GetMetrics(namespace, name)["mem"]
}
