package backend

import "fmt"

type HashiVaultBackend struct {
	Type  TLSBackendType
	Token string
}

func NewHashiVaultBackend() (Backend, error) {
	b := HashiVaultBackend{Type: BackendHashiVault}

	return b, nil
}

func (h HashiVaultBackend) verify() error {
	return nil
}

func (h HashiVaultBackend) Publish() error {
	fmt.Println("publishing certs to hashi backend")
	return nil
}
