package backend

import (
	"fmt"
	"log"

	"github.com/acheraime/certutils/utils"
)

type LocalBackend struct {
	Type           TLSBackendType
	DestinationDir string
}

func NewLocalBackend(config BackendConfig) (Backend, error) {
	b := LocalBackend{
		Type:           BackendLocal,
		DestinationDir: config.LocalDir,
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
