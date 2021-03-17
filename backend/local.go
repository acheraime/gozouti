package backend

import (
	"fmt"
	"os"
)

type LocalBackend struct {
	Type           TLSBackendType
	DestinationDir string
}

func NewLocalBackend(dir string) (Backend, error) {
	b := LocalBackend{
		Type:           BackendLocal,
		DestinationDir: dir,
	}

	if err := b.verify(); err != nil {
		return b, err
	}

	return b, nil
}

func (l LocalBackend) verify() error {
	if err := checkDir(l.DestinationDir); err != nil {
		return err
	}

	return nil
}

func (l LocalBackend) Publish() error {
	fmt.Println("sending certs to destination")
	return nil
}

func checkDir(dir string) error {
	finfo, err := os.Stat(dir)
	if err != nil {
		return err
	}

	if !finfo.IsDir() {
		return fmt.Errorf("%s is not a directory", dir)
	}

	return nil
}
