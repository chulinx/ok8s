package kapi

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v12 "k8s.io/client-go/kubernetes/typed/core/v1"
)



type NameSpace struct {
	ClientSets
}

func NewNameSpace(clientset *kubernetes.Clientset) *NameSpace {
	return &NameSpace{
		ClientSets{ClientSet: clientset },
	}
}

func (ns *NameSpace)Prefix() v12.NamespaceInterface {
	return ns.ClientSet.CoreV1().Namespaces()
}

func (ns *NameSpace)List() (v1.NamespaceList,error)  {
	nslist,err:=ns.Prefix().List(context.TODO(),metav1.ListOptions{})
	if err != nil {
		return v1.NamespaceList{},err
	}
	return *nslist,nil
}

func (ns *NameSpace)IsExits(name string) bool  {
	_,err := ns.Prefix().Get(context.TODO(),name,metav1.GetOptions{})
	if err != nil {
		return false
	}
	return true
}

func (ns *NameSpace)Create(name string) ( bool,error) {
	_,err := ns.Prefix().Create(context.TODO(),
		&v1.Namespace{ObjectMeta:metav1.ObjectMeta{Name: name}},
		metav1.CreateOptions{})
	if err != nil {
		return false,err
	}
	return true,nil
}

func (ns *NameSpace)Delete(name string)(bool,error)  {
	err := ns.Prefix().Delete(context.TODO(),name,metav1.DeleteOptions{})
	if err != nil {
		return false, err
	}
	return true,nil
}

func (ns *NameSpace)Label(name string,labels map[string]string) (bool,error) {
	_,err :=ns.Prefix().Update(context.TODO(),
		&v1.Namespace{ObjectMeta:metav1.ObjectMeta{Labels: labels,Name:name }},
		metav1.UpdateOptions{})
	if err != nil {
		return false,err
	}
	return true,nil
}

func (ns *NameSpace)ShowLabels(name string) map[string]string {
	n,err := ns.Prefix().Get(context.TODO(),name,metav1.GetOptions{})
	if err != nil {
		return map[string]string{}
	}
	return n.Labels
}