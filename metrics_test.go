package ok8s

import (
	"fmt"
	"testing"
)

var metricsInstance = NewMetrics("/Users/lisong/.kube/ack-test-pro.conf")
var nstest = "japan-pos-demo"
var pName = "coupon-integration-85c7fc968b-mvgl7"

func TestMetricsType_GetSinglePodCpu(t *testing.T) {
	fmt.Println(metricsInstance.GetSinglePodCpu(nstest, pName))
}

func TestMetricsType_GetRequestResource(t *testing.T) {
	fmt.Println(metricsInstance.GetRequestResource(nstest, pName),metricsInstance.GetLimitResource(nstest, pName))
}