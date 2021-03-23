package migrator

import (
	"fmt"

	"github.com/acheraime/certutils/backend"
	"github.com/acheraime/certutils/utils"
)

type Migrator struct {
	sourceDir        string
	migrateFromDir   bool
	inCert           string
	inKey            string
	excluded         []string
	backendProvider  *backend.K8sProvider
	backendProjectID *string
	backendCluster   *string
	backendNamespace *string
	backendType      backend.TLSBackendType
	backend          backend.Backend
	backendDirectory *string
	errors           []error
}

func NewMigrator(bType backend.TLSBackendType) (Migrator, error) {
	m := Migrator{backendType: bType}

	return m, nil
}

func (m *Migrator) SetSourceDir(inDir *string) {
	if inDir != nil {
		m.sourceDir = *inDir
		m.migrateFromDir = true
	}

}

func (m *Migrator) SetCertFiles(key, cert string) {
	m.inCert = cert
	m.inKey = key
}

func (m *Migrator) SetProjectID(pID *string) {
	m.backendProjectID = pID
}

func (m *Migrator) SetK8sCluster(c *string) {
	m.backendCluster = c
}

func (m *Migrator) SetK8sNamespace(n *string) {
	m.backendNamespace = n
}

func (m *Migrator) SetBackendProvider(p string) {
	// Set default provider
	if p == "" {
		p = "gcp"
	}
	provider := backend.K8sProvider(p)

	m.backendProvider = &provider
}

func (m *Migrator) SetBackendDirectory(d *string) {
	m.backendDirectory = d
}

func (m Migrator) validate() error {
	switch m.backendType {
	case backend.Backendkubernetes:
		if !m.backendProvider.IsValid() {
			return fmt.Errorf("kubernetes provider %s is not supported", *m.backendProvider)
		}
		// Check cluster name
		switch *m.backendProvider {
		case "dockker-desktop":
			*m.backendCluster = "docker-desktop"
		default:
			if m.backendCluster == nil {
				return fmt.Errorf("a cluster name is required with kubernetes backend")
			}
			if *m.backendProvider == "gcp" && m.backendProjectID == nil {
				return fmt.Errorf("a projectID is required for GCP")
			}
		}

	}
	if m.migrateFromDir {
		if err := utils.CheckDir(m.sourceDir); err != nil {
			return fmt.Errorf("source directory '%s' is not valid: %s", m.sourceDir, err.Error())
		}
	}
	return nil
}

func (m *Migrator) Migrate() error {
	if err := m.validate(); err != nil {
		return err
	}
	backendConfig := backend.BackendConfig{
		K8sClusterName: m.backendCluster,
		K8sProvider:    m.backendProvider,
		ProjectID:      m.backendProjectID,
		DestNameSpace:  m.backendNamespace,
		LocalDir:       m.backendDirectory,
	}
	// set backend
	b, err := backend.NewBackend(m.backendType, backendConfig)
	if err != nil {
		return err
	}
	m.backend = b
	// Get a runner
	runner, err := NewRunner(*m)
	if err != nil {
		return err
	}

	// Run the migration
	if err := runner.Run(); err != nil {
		return err
	}
	return nil
}
