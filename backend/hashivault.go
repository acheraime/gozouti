package backend

import "fmt"

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
