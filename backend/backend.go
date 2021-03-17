package backend

type TLSBackendType string

const (
	Backendkubernetes TLSBackendType = "kubernetes"
	BackendRedis      TLSBackendType = "redis"
	BackendHashiVault TLSBackendType = "hashivault"
	BackendLocal      TLSBackendType = "local"
)

type Backend interface {
	verify() error
	Publish() error
}

func NewBackend(backendType TLSBackendType) (Backend, error) {
	var backend Backend

	switch backendType {
	case BackendLocal:
		return NewLocalBackend("test")
	case BackendHashiVault:
	default:
	}
	return backend, nil
}
