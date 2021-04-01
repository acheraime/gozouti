package backend

import (
	"fmt"

	v1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
)

type HashiVaultBackend struct {
	Type  TLSBackendType
	Token string
}

func NewHashiVaultBackend(config BackendConfig) (Backend, error) {
	b := HashiVaultBackend{Type: BackendHashiVault}

	return &b, nil
}

func (h HashiVaultBackend) build() error {
	return nil
}

func (h HashiVaultBackend) Publish() error {
	fmt.Println("publishing certs to hashi backend")
	return nil
}

func (h HashiVaultBackend) Test() bool {
	return true
}

func (h HashiVaultBackend) Migrate(cert, key []byte, certName string) error {
	return fmt.Errorf("Not implemented")
}

func (h HashiVaultBackend) GetType() TLSBackendType {
	return h.Type
}

func (h HashiVaultBackend) CreateTraefikMiddleWare(middleware *v1alpha1.Middleware) error {
	return fmt.Errorf("not implemented")
}
