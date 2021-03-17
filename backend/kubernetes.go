package backend

import (
	"context"
	"encoding/base64"
	"fmt"

	"google.golang.org/api/container/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd/api"
)

type KubernetesBackend struct {
	Type          TLSBackendType
	NameSpace     string
	K8sContext    string
	ProjectID     string
	CloudProvider string
	client        *kubernetes.Clientset
	config        *api.Config
}

func NewK8sBackend() (Backend, error) {
	b := KubernetesBackend{
		Type: Backendkubernetes,
	}

	return b, nil
}

func (k KubernetesBackend) build() error {
	ctx := context.Background()

	return nil
}

func (k KubernetesBackend) Publish() error {
	fmt.Println("publishing certs to k8s backend")
	return nil
}

func (k KubernetesBackend) getk8sconfig(ctx context.Context) (*api.Config, error) {
	svc, err := container.NewService(ctx)
	if err != nil {
		return nil, err
	}

	// Bare bone configuration structure
	config := api.Config{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters:   map[string]*api.Cluster{},
		AuthInfos:  map[string]*api.AuthInfo{},
		Contexts:   map[string]*api.Context{},
	}

	switch k.CloudProvider {
	case "gcp":

	}
	return &config, nil
}

func getK8sclient(conf api.Config)

func (k *KubernetesBackend) buildGCPConfig(ctx context.Context, svc *container.Service) error {
	// Retrieve the list of gke clusters for the project ID
	res, err := svc.Projects.Zones.Clusters.List(k.ProjectID, "-").Context(ctx).Do()
	if err != nil {
		return err
	}

	// Iterate over the clusters and start populating the api
	for _, c := range res.Clusters {
		// Cluster name
		clusterName := fmt.Sprintf("gke_%s_%s_%s", k.ProjectID, c.Zone, c.Name)
		cert, err := base64.StdEncoding.DecodeString(c.MasterAuth.ClientCertificate)
		if err != nil {
			return err
		}
		// Populate the config object with information from this cluster
		k.config.Clusters[clusterName] = &api.Cluster{
			CertificateAuthorityData: cert,
			Server:                   "https://" + c.Endpoint,
		}
		// Contexts
		k.config.Contexts[clusterName] = &api.Context{
			Cluster:  clusterName,
			AuthInfo: clusterName,
		}

		k.config.AuthInfos[clusterName] = &api.AuthInfo{
			AuthProvider: &api.AuthProviderConfig{
				Name: "gcp",
				Config: map[string]string{
					"scopes": "https://www.googleapis.com/auth/cloud-platform",
				},
			},
		}
	}

	return nil
}
