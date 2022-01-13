package backend

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"path/filepath"

	"github.com/acheraime/gozouti/utils"
	traefikScheme "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/generated/clientset/versioned/scheme"
	v1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
	"google.golang.org/api/container/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type K8sSecret v1.Secret
type KubernetesBackend struct {
	Type        TLSBackendType
	NameSpace   string
	K8sContext  string
	ProjectID   string
	K8sProvider K8sProvider
	ClusterName string
	client      *kubernetes.Clientset
	apiExClient *clientset.Clientset
	config      *api.Config
}

func NewK8sBackend(config BackendConfig) (Backend, error) {
	b := KubernetesBackend{
		Type:        Backendkubernetes,
		K8sProvider: *config.K8sProvider,
		ClusterName: *config.K8sClusterName,
	}

	if config.DestNameSpace != nil {
		b.NameSpace = *config.DestNameSpace
	}

	if config.ProjectID != nil {
		b.ProjectID = *config.ProjectID
	}

	if err := b.build(); err != nil {
		log.Println("unable to build the backend" + err.Error())
		return nil, err
	}

	return &b, nil
}

func (k *KubernetesBackend) build() error {
	// Build configures and set a client
	// to interact with
	ctx := context.Background()
	// set k8s configuration
	if err := k.setK8sConfig(ctx); err != nil {
		return err
	}

	return nil
}

func (k KubernetesBackend) Publish() error {
	fmt.Println("publishing certs to k8s backend")
	return nil
}

func (k *KubernetesBackend) setK8sConfig(ctx context.Context) error {
	// Barebone configuration structure
	// uppon which we will build info to
	// connect to k8s clusters.
	k.config = &api.Config{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters:   map[string]*api.Cluster{},
		AuthInfos:  map[string]*api.AuthInfo{},
		Contexts:   map[string]*api.Context{},
	}

	switch k.K8sProvider {
	case "gcp":
		svc, err := container.NewService(ctx)
		if err != nil {
			return err
		}
		if err := k.buildGCPConfig(ctx, svc); err != nil {
			return err
		}
	case "docker-desktop":
		if err := k.setK8sClientFromFile(); err != nil {
			return err
		}
	}
	return nil
}

func (k *KubernetesBackend) buildGCPConfig(ctx context.Context, svc *container.Service) error {
	// Retrieve the list of gke clusters for the project ID
	res, err := svc.Projects.Zones.Clusters.List(k.ProjectID, "-").Context(ctx).Do()
	if err != nil {
		return err
	}

	if res.Clusters == nil {
		return fmt.Errorf("There's no k8s cluster in project %s", k.ProjectID)
	}

	var cluster *container.Cluster
	for _, c := range res.Clusters {
		//clusterName := fmt.Sprintf("gke_%s_%s_%s", k.ProjectID, c.Zone, c.Name)
		if c.Name == k.ClusterName {
			cluster = c
			break
		}
	}

	if cluster == nil {
		return fmt.Errorf("Cluster %s was not found in project %s", k.ClusterName, k.ProjectID)
	}

	// Master certificate
	cert, err := base64.StdEncoding.DecodeString(cluster.MasterAuth.ClusterCaCertificate)
	if err != nil {
		return err
	}
	// Populate the config object with information from this cluster
	k.config.Clusters[k.ClusterName] = &api.Cluster{
		CertificateAuthorityData: cert,
		Server:                   "https://" + cluster.Endpoint,
	}
	// Contexts
	k.config.Contexts[k.ClusterName] = &api.Context{
		Cluster:  k.ClusterName,
		AuthInfo: k.ClusterName,
	}

	k.config.AuthInfos[k.ClusterName] = &api.AuthInfo{
		AuthProvider: &api.AuthProviderConfig{
			Name: "gcp",
			Config: map[string]string{
				"scopes": "https://www.googleapis.com/auth/cloud-platform",
			},
		},
	}

	cfg, err := k.configFromContext(k.ClusterName)
	if err != nil {
		return err
	}

	if err := k.setK8sClient(cfg); err != nil {
		return err
	}

	return nil
}

func (k *KubernetesBackend) setK8sClientFromFile() error {
	var kubefile string
	if home := utils.HomeDir(); home != "" {
		kubefile = filepath.Join(home, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("https://kubernetes.docker.internal:6443", kubefile)

	if err != nil {
		return err
	}

	if err := k.setK8sClient(config); err != nil {
		return err
	}

	return nil

}

func (k KubernetesBackend) configFromContext(clusterName string) (*rest.Config, error) {
	if k.config.Clusters == nil {
		return nil, fmt.Errorf("No confguration found for cluster: %s", clusterName)
	}

	cfg, err := clientcmd.NewNonInteractiveClientConfig(*k.config, clusterName,
		&clientcmd.ConfigOverrides{CurrentContext: clusterName}, nil).ClientConfig()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (k *KubernetesBackend) setK8sClient(cfg *rest.Config) error {
	// setK8sClient helper method use the configuration
	// to expose a clientset object suitable to
	// interact with k8s api server

	k8s, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}

	k.client = k8s
	k.apiExClient, err = clientset.NewForConfig(cfg)
	if err != nil {
		return err
	}

	return nil
}

// CreateSecret is a helper func that
// abstracts creation of a secret
// in a k8s cluster
func (k KubernetesBackend) createSecret(ctx context.Context, namespace string, secretobj *v1.Secret) error {
	secret, err := k.client.CoreV1().Secrets(namespace).Create(ctx, secretobj, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Secret created: %s\n", secret.Name)

	return nil
}

func (k KubernetesBackend) updateSecret(ctx context.Context, namespace string, secret *v1.Secret) error {
	updated, err := k.client.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Secret updated: %s\n", updated.Name)

	return nil
}

func (k KubernetesBackend) createCustom(ctx context.Context, middleware *v1alpha1.Middleware) error {
	_ = traefikScheme.AddToScheme(scheme.Scheme)
	response := &v1alpha1.Middleware{}
	var middlewaresResource = schema.GroupVersionResource{Group: "traefik.containo.us", Version: "v1alpha1", Resource: "middlewares"}
	options := &metav1.GetOptions{}
	options.ResourceVersion = middlewaresResource.String()
	options.Kind = "Middleware"

	// err := k.apiExClient.RESTClient().Post().
	// 	Resource("middlewares").
	// 	Namespace(k.NameSpace).
	// 	VersionedParams(&metav1.CreateOptions{}, traefikScheme.ParameterCodec).
	// 	Body(middleware).
	// 	Do(ctx).
	// 	Into(response)
	// if err != nil {
	// 	return err
	// }
	err := k.apiExClient.RESTClient().Get().Resource("middlewares").VersionedParams(options, traefikScheme.ParameterCodec).Do(ctx).Into(response)
	if err != nil {
		return err
	}
	fmt.Print(string(response.Name))
	return nil
}

func (k KubernetesBackend) CreateTraefikMiddleWare(middleware *v1alpha1.Middleware) error {
	if err := k.createCustom(context.TODO(), middleware); err != nil {
		return err
	}

	return nil
}

func (k KubernetesBackend) Migrate(cert, key []byte, secretName string) error {
	secretData := map[string][]byte{
		"tls.crt": cert,
		"tls.key": key,
	}

	secret := &v1.Secret{
		Type: "kubernetes.io/tls",
		Data: secretData,
	}
	secret.SetName(secretName)

	fmt.Printf("Migrating %s\n", secretName)
	if err := k.createSecret(context.TODO(), k.NameSpace, secret); err != nil {
		if errors.IsAlreadyExists(err) {
			// If the secret exists already
			// let's update it
			if err = k.updateSecret(context.TODO(), k.NameSpace, secret); err != nil {
				return err
			}
		}
		return err
	}

	return nil
}

func (k KubernetesBackend) Test() bool {
	ctx := context.Background()
	// As a test we will iterate over all the clusters found
	// and display the pods we find.
	// If no error we'll return true
	if k.config.Clusters == nil {
		log.Println("no cluster found")
		return false
	}

	pods, err := k.client.CoreV1().Pods(k.NameSpace).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return false
	}

	fmt.Printf("============= List of pods found on cluster %s ==============\n", k.ClusterName)
	for _, p := range pods.Items {
		fmt.Printf("%s \n", p.Name)
	}
	fmt.Println("============================================================")
	return true
}

func (h KubernetesBackend) GetType() TLSBackendType {
	return h.Type
}
