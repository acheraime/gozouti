package migrator

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
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
	var data = map[string]TLSData{}
	var buildData = func(cert, key string) error {
		// Check if cert or key files are part of
		// the exclusion list
		if m.migrator.excluded != nil {
			for _, exclusion := range m.migrator.excluded {
				exclusionPath := m.buildSourceFilePath(exclusion)
				if exclusionPath == cert || exclusionPath == key {
					fmt.Println("Will not migrate: " + exclusion)
					return nil
				}
			}
		}

		// certificate bytes
		certBytes, err := ioutil.ReadFile(cert)
		if err != nil {
			return err
		}
		// private key bytes
		keyBytes, err := ioutil.ReadFile(key)
		if err != nil {
			return err
		}

		// Validate cert and get dns names
		c, err := tls.ParseCertificate(certBytes)
		if err != nil {
			return err
		}
		if !c.IsValid() {
			fmt.Println("Certificate is not valid")
		}

		certName := secretNameFromDNS(c.Certificate.DNSNames)
		data[certName] = TLSData{
			cert: certBytes,
			key:  keyBytes,
		}

		return nil
	}

	if !m.migrator.migrateFromDir {
		// Migrate cert and key files provided
		// as arguments to the migrate command
		if err := buildData(m.migrator.inCert, m.migrator.inKey); err != nil {
			return err
		}

	} else {
		// walk directory and create cert, key pair
		dirFiles, err := ioutil.ReadDir(m.migrator.sourceDir)
		if err != nil {
			return err
		}
		dmap := map[string]string{}
		for _, file := range dirFiles {
			if file.IsDir() {
				continue
			}
			if isKey(file.Name()) {
				continue
			}
			// this is a cert file
			for _, nf := range dirFiles {
				if splitExt(file.Name()) == splitExt(nf.Name()) && isKey(nf.Name()) {
					dmap[file.Name()] = nf.Name()
				}
			}

		}

		for c, k := range dmap {
			certPath := m.buildSourceFilePath(c)
			keyPath := m.buildSourceFilePath(k)
			fmt.Println("processing -> " + c)
			if err := buildData(certPath, keyPath); err != nil {
				return err
			}
		}
	}

	m.data = data

	return nil
}

func (m Runner) buildSourceFilePath(f string) string {
	return filepath.Join(m.migrator.sourceDir, f)
}

func (m Runner) Run() error {
	if m.data == nil {
		return fmt.Errorf("no tls data found")
	}

	fmt.Println("Starting migration")
	for name, cert := range m.data {
		if err := m.migrator.backend.Migrate(cert.cert, cert.key, name); err != nil {
			return err
		}
	}
	fmt.Println("migration complete")
	return nil
}

func secretNameFromDNS(names []string) string {
	var secretName string
	for _, name := range names {
		if strings.HasPrefix(name, "*") {
			secretName = strings.Replace(name, "*", secretPrefix, 1)
			break
		}

		secretName = name
	}

	secretName = strings.ReplaceAll(secretName, ".", "-") + "-" + secretSuffix

	return secretName
}

func isPem(in string) bool {
	switch filepath.Ext(in) {
	case ".pem", ".crt", ".cer", ".key":
		return true
	default:
		return false
	}
}

func isKey(in string) bool {
	return isPem(in) && filepath.Ext(in) == ".key"
}

func splitExt(f string) string {
	return strings.Split(f, ".")[0]
}
