package ok8s

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type NameSpace struct {
	ClientSets
}

func NewNameSpace(cs ClientSets) *NameSpace {
	n := &NameSpace{
		ClientSets: cs,
	}
	return n
}

func (ns *NameSpace) Prefix() v12.NamespaceInterface {
	return ns.ClientSet.CoreV1().Namespaces()
}

func (ns *NameSpace) Create(namespace string) (bool, error) {
	_, err := ns.Prefix().Create(DefaultTimeOut(),
		&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}},
		metav1.CreateOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (ns *NameSpace) List() (v1.NamespaceList, error) {
	nslist, err := ns.Prefix().List(DefaultTimeOut(), metav1.ListOptions{})
	if err != nil {
		return v1.NamespaceList{}, err
	}
	return *nslist, nil
}

func (ns *NameSpace) IsExits(name string) bool {
	_, err := ns.Prefix().Get(DefaultTimeOut(), name, metav1.GetOptions{})
	if err != nil {
		return false
	}
	return true
}

func (ns *NameSpace) Delete(name string) (bool, error) {
	err := ns.Prefix().Delete(DefaultTimeOut(), name, metav1.DeleteOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (ns *NameSpace) Label(name string, labels map[string]string) (bool, error) {
	_, err := ns.Prefix().Update(DefaultTimeOut(),
		&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Labels: labels, Name: name}},
		metav1.UpdateOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (ns *NameSpace) ShowLabels(name string) map[string]string {
	n, err := ns.Prefix().Get(DefaultTimeOut(), name, metav1.GetOptions{})
	if err != nil {
		return map[string]string{}
	}
	return n.Labels
}
