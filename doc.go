// Package ok8s help you manger kubernetes resource
//
// A High encapsulation Kubernetes Golang-api base client-go
//
// As quick start
//		nsName := "test"
//		client := ok8s.NewTestClientSet()
//
//Create deployment
//		namespaceInstance := ok8s.NewNameSpace(client)
//		namespaceInstance.Create(nsName)
// Get deployment
//		_, deployment := deploymentInstance.Get(nsName, resourceName)
//		FormatPrint("Deployment", deployment.Deployment.Name)
// Delete deployment
//		deploymentInstance.Delete(nsName, resourceName)
//		if deploymentInstance.IsExits(nsName, resourceName) {
//			fmt.Println("Not delete")
//		} else {
//			fmt.Println("Delete Success")
//		}
//
package ok8s


