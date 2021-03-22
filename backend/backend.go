package backend

import "errors"

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
}

type BackendConfig struct {
	LocalDir       string
	K8sClusterName string
	K8sProvider    K8sProvider
	ProjectID      string
	DestNameSpace  string
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
