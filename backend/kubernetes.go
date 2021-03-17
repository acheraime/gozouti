package backend

import "fmt"

type KubernetesBackend struct {
	Type      TLSBackendType
	NameSpace string
}

func NewK8sBackend() (Backend, error) {
	b := KubernetesBackend{
		Type: Backendkubernetes,
	}

	return b, nil
}

func (k KubernetesBackend) verify() error {
	return nil
}

func (k KubernetesBackend) Publish() error {
	fmt.Println("publishing certs to k8s backend")
	return nil
}
