package backend

import (
	"fmt"
	"log"

	"github.com/acheraime/certutils/utils"
	v1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
)

type LocalBackend struct {
	Type           TLSBackendType
	DestinationDir string
}

func NewLocalBackend(config BackendConfig) (Backend, error) {
	if config.LocalDir == nil {
		return nil, fmt.Errorf("please specify the local directory")
	}

	b := LocalBackend{
		Type:           BackendLocal,
		DestinationDir: *config.LocalDir,
	}

	if err := b.build(); err != nil {
		return b, err
	}

	return &b, nil
}

func (l LocalBackend) build() error {
	if err := utils.CheckDir(l.DestinationDir); err != nil {
		return err
	}

	return nil
}

func (l LocalBackend) Publish() error {
	fmt.Println("sending certs to destination")
	return nil
}

func (l LocalBackend) Test() bool {
	if err := utils.CheckDir(l.DestinationDir); err != nil {
		log.Println(err)
		return false
	}

	return true
}

func (l LocalBackend) Migrate(cert, key []byte, certName string) error {
	return fmt.Errorf("Not implemented")
}

func (l LocalBackend) GetType() TLSBackendType {
	return l.Type
}

func (l LocalBackend) CreateTraefikMiddleWare(middleware *v1alpha1.Middleware) error {
	return fmt.Errorf("not implemented")
}
