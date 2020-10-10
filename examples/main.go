package main

import (
	"fmt"
	"github.com/chulinx/ok8s"
)

func main() {
	nsName := "test"
	client := ok8s.NewTestClientSet()

	// Create
	namespaceInstance := ok8s.NewNameSpace(client)
	namespaceInstance.Create(nsName)

	deploymentInstance := ok8s.NewDeployment(client)
	deploymentInstance.Create(nsName, deployStruct)

	svcInstance := ok8s.NewService(client)
	svcInstance.Create(nsName, serviceStruct)

	ingressInstance := ok8s.NewIngress(client)
	ingressInstance.Create(nsName, ingressStruct)

	// Get
	_, deployment := deploymentInstance.Get(nsName, resourceName)
	FormatPrint("Deployment", deployment.Deployment.Name)

	_, service := svcInstance.Get(nsName, resourceName)
	FormatPrint("Service", service.Service.Name)

	_, ingress := ingressInstance.Get(nsName, resourceName)
	FormatPrint("Ingress", ingress.Ingress.Name)

	// Delete
	deploymentInstance.Delete(nsName, resourceName)
	if deploymentInstance.IsExits(nsName, resourceName) {
		fmt.Println("Not delete")
	} else {
		fmt.Println("Delete Success")
	}

	svcInstance.Delete(nsName, resourceName)
	if svcInstance.IsExits(nsName, resourceName) {
		fmt.Println("Not delete")
	} else {
		fmt.Println("Delete Success")
	}

	ingressInstance.Delete(nsName, resourceName)
	if ingressInstance.IsExits(nsName, resourceName) {
		fmt.Println("Not delete")
	} else {
		fmt.Println("Delete Success")
	}

	namespaceInstance.Delete(nsName)
}

func FormatPrint(Type, name string) {
	fmt.Printf("Resource: %s     Name: %s\n", Type, name)
}
