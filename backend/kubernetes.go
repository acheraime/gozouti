package backend

import (
	"context"
	"fmt"

	"google.golang.org/api/container/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd/api"
)

type KubernetesBackend struct {
	Type       TLSBackendType
	NameSpace  string
	K8sContext string
	ProjectID  string
	Provider   string
	client     *kubernetes.Clientset
}

func NewK8sBackend() (Backend, error) {
	b := KubernetesBackend{
		Type: Backendkubernetes,
	}

	return b, nil
}

func (k KubernetesBackend) build() error {
	return nil
}

func (k KubernetesBackend) Publish() error {
	fmt.Println("publishing certs to k8s backend")
	return nil
}

func getk8sconfig(ctx context.Context) (*api.Config, error) {
	svc, err := container.NewService(ctx)
	if err != nil {
		return err
	}

	return nil, nil
}

func getK8sclient(conf api.Config)
