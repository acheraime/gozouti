package backend

import "errors"

type TLSBackendType string

const (
	Backendkubernetes TLSBackendType = "kubernetes"
	BackendRedis      TLSBackendType = "redis"
	BackendHashiVault TLSBackendType = "hashivault"
	BackendLocal      TLSBackendType = "local"
)

type Backend interface {
	build() error
	Publish() error
}

func NewBackend(backendType TLSBackendType) (Backend, error) {
	switch backendType {
	case BackendLocal:
		return NewLocalBackend("test")
	case BackendHashiVault:
		return NewHashiVaultBackend()
	case Backendkubernetes:
		return NewK8sBackend()
	default:
		return nil, errors.New("backend not supported")
	}
}
