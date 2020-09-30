package kapi


//
//func TestDeployment_WatchAdd(t *testing.T) {
//	deploy := &Deployment{ClientSets{ClientSet: default_clientset.ClientSet() }}
//	eventFuncs := new(cache.ResourceEventHandlerFuncs)
//	eventFuncs.DeleteFunc= func(obj interface{}) {
//		deployment := obj.(*appsv1.Deployment)
//		log.Printf("delete deployment %s",deployment.Name)
//		fmt.Println(deployment.Labels,deployment.Spec.Template.Labels)
//	}
//	eventFuncs.AddFunc= func(obj interface{}) {
//		deployment := obj.(*appsv1.Deployment)
//		log.Printf("add deployment %s",deployment.Name)
//		fmt.Println(deployment.Labels,deployment.Spec.Template.Labels)
//	}
//	deploy.Watch("echo-server",*eventFuncs)
//}
