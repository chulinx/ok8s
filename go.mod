module github.com/chulinx/ok8s

go 1.15

require (
	k8s.io/api v0.18.1
	k8s.io/apimachinery v0.18.1
	k8s.io/client-go v0.18.1
)

replace (
	k8s.io/api => k8s.io/api v0.18.1
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.1
	k8s.io/client-go => k8s.io/client-go v0.18.1
)
