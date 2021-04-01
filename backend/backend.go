package backend

import (
	"errors"

	v1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
)

type TLSBackendType string

const (
	Backendkubernetes TLSBackendType = "kubernetes"
	BackendRedis      TLSBackendType = "redis"
	BackendHashiVault TLSBackendType = "hashivault"
	BackendLocal      TLSBackendType = "local"
)

type K8sProvider string

const (
	ProviderGCP           K8sProvider = "gcp"
	ProviderDockerDesktop K8sProvider = "docker-desktop"
	ProviderAWS           K8sProvider = "aws"
	ProviderAzure         K8sProvider = "azure"
)

func (p K8sProvider) IsValid() bool {
	switch p {
	case ProviderGCP, ProviderDockerDesktop, ProviderAWS, ProviderAzure:
		return true
	default:
		return false
	}
}

type Backend interface {
	build() error
	Publish() error
	Test() bool
	Migrate([]byte, []byte, string) error
	GetType() TLSBackendType
	CreateTraefikMiddleWare(*v1alpha1.Middleware) error
}

type BackendConfig struct {
	LocalDir       *string
	K8sClusterName *string
	K8sProvider    *K8sProvider
	ProjectID      *string
	DestNameSpace  *string
}

func NewBackend(backendType TLSBackendType, cfg BackendConfig) (Backend, error) {
	switch backendType {
	case BackendLocal:
		return NewLocalBackend(cfg)
	case BackendHashiVault:
		return NewHashiVaultBackend(cfg)
	case Backendkubernetes:
		return NewK8sBackend(cfg)
	default:
		return nil, errors.New("backend not supported")
	}
}
