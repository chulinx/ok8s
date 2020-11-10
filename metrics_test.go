package ok8s

import (
	"fmt"
	"testing"
)

var metricsInstance = NewMetrics("/Users/lisong/.kube/ack-devops.conf")

func TestMetricsType_GetSinglePodCpu(t *testing.T) {
	fmt.Println(metricsInstance.GetSinglePodCpu("thanos","thanos-store-sharding-4-1"))
}
