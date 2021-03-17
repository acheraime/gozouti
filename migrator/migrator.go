package migrator

import (
	"github.com/acheraime/certutils/backend"
)

type Migrator struct {
	Backend   backend.Backend
	SourceDir string
	Excluded  []string
}

func NewMigrator(bType backend.TLSBackendType) error {
	b, err := backend.NewBackend(bType)
	if err != nil {
		return err
	}
	b.Publish()

	return nil
}
