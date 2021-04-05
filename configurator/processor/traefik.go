package processor

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/acheraime/certutils/backend"
	"github.com/acheraime/certutils/configurator/traefik"
	"github.com/acheraime/certutils/utils"
)

type TraefikRedirect struct {
	resources RedirectResources
	buffer    bytes.Buffer
	parseURL  bool
	baseHost  string
	backend   backend.Backend
	namespace string
	outDir    string
	alias     string
}

type TRedirectConfig struct {
	Alias     string
	Namespace string
	OutputDir string
	BaseHost  string
}

func NewTraefikRedirect(cfg TRedirectConfig, input [][]string, b backend.Backend) (Processor, error) {
	if b.GetType() != backend.Backendkubernetes {
		return nil, fmt.Errorf("traefik configuration requires a kubernetes backend type")
	}

	if cfg.Namespace == "" {
		cfg.Namespace = "default"
	}

	tr := TraefikRedirect{
		backend:   b,
		namespace: cfg.Namespace,
		outDir:    cfg.OutputDir,
		baseHost:  cfg.BaseHost,
		alias:     cfg.Alias,
	}

	resources, err := NewRedirectResources(input, true, cfg.Alias, cfg.BaseHost)
	if err != nil {
		return nil, err
	}

	tr.resources = resources
	return tr, nil
}

func (t TraefikRedirect) build(o io.Writer) error {
	if t.resources.Resources == nil {
		return fmt.Errorf("nothing to generate: %s", t.resources.Resources)
	}

	// reset the buffer
	t.buffer.Reset()
	fileHeader += fmt.Sprintf("#### Timestamp: %s\n", time.Now().Format(time.RFC1123))
	t.buffer.Write([]byte(fileHeader + "\n\n"))

	middlewares := []traefik.Middleware{}

	regexTemplate, err := template.New("regex").Parse(`^https://(.*).({{ .URLHost }}){{ .Regex }}?$$`)
	if err != nil {
		return err
	}
	replTemplate, err := template.New("replacement").Parse(`https://${1}.${2}{{ .Replacement }}`)
	if err != nil {
		return err
	}

	for _, r := range t.resources.Resources {
		var regexStr bytes.Buffer
		if err := regexTemplate.Execute(&regexStr, r); err != nil {
			return err
		}
		var replStr bytes.Buffer
		if err := replTemplate.Execute(&replStr, r); err != nil {
			return err
		}

		m, err := traefik.NewRegexRedirect(r.Name, t.namespace, regexStr.String(), replStr.String(), true)
		if err != nil {
			return err
		}

		t.buffer.Write([]byte("---\n"))
		t.buffer.Write([]byte(m.String() + "\n\n"))

		middlewares = append(middlewares, m)
	}

	// create a chain of middlewares
	redirectName := t.alias + "-redirects"
	chain, err := traefik.NewChain(redirectName, t.namespace, middlewares)
	if err != nil {
		return err
	}

	t.buffer.Write([]byte("---\n"))
	t.buffer.Write([]byte(chain.String() + "\n\n"))
	// WriteTo implicitly call
	// flush
	t.buffer.WriteTo(o)

	return nil
}

func (t TraefikRedirect) DryRun(dest string) error {

	return t.build(os.Stdout)
}

func (t TraefikRedirect) Generate() error {
	if err := utils.CheckDir(t.outDir); err != nil {
		return err
	}

	// Default file name
	fileName := fmt.Sprintf("%s-generated-%s.yaml", t.alias, time.Now().Format("20060102150405"))
	outPath := filepath.Join(t.outDir, fileName)
	f, err := os.OpenFile(outPath, os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.build(f)
}
