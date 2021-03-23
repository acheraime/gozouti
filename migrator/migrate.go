package migrator

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/acheraime/certutils/tls"
)

const (
	secretPrefix = "star"
	secretSuffix = "ssl"
)

type TLSData struct {
	cert, key []byte
}
type Runner struct {
	migrator Migrator
	data     map[string]TLSData
}

func NewRunner(migrator Migrator) (*Runner, error) {
	r := Runner{
		migrator: migrator,
	}

	if err := r.setData(); err != nil {
		return nil, err
	}

	return &r, nil
}

func (m *Runner) setData() error {
	if !m.migrator.migrateFromDir {
		// certificate bytes
		certBytes, err := ioutil.ReadFile(m.migrator.inCert)
		if err != nil {
			return err
		}
		// private key bytes
		keyBytes, err := ioutil.ReadFile(m.migrator.inKey)
		if err != nil {
			return err
		}

		// Validate cert and get dns names
		c, err := tls.ParseCertificate(certBytes)
		if err != nil {
			return err
		}
		if c.IsValid() {
			log.Print("Certificate is valid")
		}

		certName := secretNameFromDNS(c.Certificate.DNSNames)
		m.data = map[string]TLSData{
			certName: TLSData{
				cert: certBytes,
				key:  keyBytes,
			},
		}
	}

	return nil
}

func (m Runner) Run() error {
	log.Println("Starting migration")
	for name, cert := range m.data {
		if err := m.migrator.backend.Migrate(cert.cert, cert.key, name); err != nil {
			return err
		}
	}
	log.Println("migration comple")
	return nil
}

func secretNameFromDNS(names []string) string {
	var secretName string
	for _, name := range names {
		if strings.HasPrefix(name, "*") {
			secretName = strings.Replace(name, "*", secretPrefix, 1)
			break
		}
	}
	secretName = strings.ReplaceAll(secretName, ".", "-") + "-" + secretSuffix

	return secretName
}
