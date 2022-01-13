package configurator

import (
	"bytes"
	"fmt"

	"github.com/acheraime/gozouti/backend"
	"github.com/acheraime/gozouti/configurator/input"
	"github.com/acheraime/gozouti/configurator/processor"
)

type ConfigType string

const (
	RedirectConfig ConfigType = "redirect"
)

type Platform string

const (
	NginxPlatform   Platform = "nginx"
	TraefikPlatform Platform = "traefik"
)

type InputType string

const (
	CSVInput InputType = "csv"
)

type Options struct {
	In       string
	Out      string
	InType   InputType
	Type     ConfigType
	Platform Platform
	// Traefik redirect specific configuration
	RedirectAlias       string
	RedirectNamespace   string
	RedirectBaseHostURL string
	RedirectRewriteHost bool
	DryRun              bool
	buffer              bytes.Buffer
	dupes               []map[string]string
}

func NewConfigurator(options Options) (processor.Processor, error) {
	var proc processor.Processor

	if options.InType != CSVInput {
		return proc, fmt.Errorf("Input type not supported")
	}

	in, err := input.CSVInputReader(options.In)
	if err != nil {
		return nil, err
	}

	switch options.Platform {
	case NginxPlatform:
		switch options.Type {
		case RedirectConfig:
			proc, err = processor.NewNginxRedirect(in, options.Out, true)
			if err != nil {
				return proc, err
			}
		}
	case TraefikPlatform:
		// configure a kubernetes backend
		var cluster = "docker-desktop"
		var projectID = "fusion-dev-163815"
		var provider = backend.ProviderDockerDesktop
		var namespace = "default"
		backendConfig := backend.BackendConfig{
			K8sClusterName: &cluster,
			K8sProvider:    &provider,
			DestNameSpace:  &namespace,
			ProjectID:      &projectID,
		}
		// set backend
		b, err := backend.NewBackend(backend.Backendkubernetes, backendConfig)
		if err != nil {
			return nil, err
		}
		proc, err = processor.NewTraefikRedirect(processor.TRedirectConfig{
			Alias:       options.RedirectAlias,
			Namespace:   options.RedirectNamespace,
			BaseHost:    options.RedirectBaseHostURL,
			OutputDir:   options.Out,
			RewriteHost: options.RedirectRewriteHost,
		}, in, b)
	default:
		return nil, fmt.Errorf("patform %s is not currently supported", options.Platform)
	}

	return proc, nil
}
