package traefik

import (
	"fmt"

	"gopkg.in/yaml.v3"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	apiVersion1 = "traefik.containo.us/v1alpha1"
)

type ObjectMeta struct {
	ApiVersion string `yaml:"apiVersion,omitempty"`
	Kind       string `yaml:"kind,omitempty"`
}

// type Middleware v1alpha1.Middleware
type TypeMeta struct {
	Meta MetaData `yaml:"metadata,omitempty"`
}

type MetaData struct {
	Name      string `yaml:"name,omitempty"`
	Label     string `yaml:"label,omitempty"`
	Namespace string `yaml:"namespace,omitempty"`
}

type Middleware struct {
	ObjectMeta `yaml:",inline"`
	TypeMeta   `yaml:",inline"`

	Spec MiddlewareSpec `yaml:"spec"`
}

// type MiddlewareSpec v1alpha1.MiddlewareSpec

type MiddlewareSpec struct {
	RedirectRegex RedirectRegex `yaml:"redirectRegex,omitempty" json:"redirectRegex"`
	Chain         `yaml:"chain,omitempty" json:"chain,omitempty"`
}

type MiddlewareRef struct {
	Name      string `yaml:"name,omitempty"`
	Namespace string `yaml:"namespace,omitempty"`
}

//type Chain v1alpha1.Chain

type Chain struct {
	Middlewares []MiddlewareRef `yaml:"middlewares" json:"middlewares,omitempty"`
}

// type RedirectRegex dynamic.RedirectRegex

type RedirectRegex struct {
	Regex       string `yaml:"regex,omitempty" json:"regex,omitempty"`
	Replacement string `yaml:"replacement,omitempty" json:"replacement,omitempty"`
	Permanent   bool   `yaml:"permanent,omitempty" json:"permanent,omitempty"`
}

func (o *ObjectMeta) SetGroupVersion() {
	o.ApiVersion = "traefik.containo.us/v1alpha1"
	o.Kind = "Middleware"
}

var middlewaresKind = schema.GroupVersionKind{Group: "traefik.containo.us", Version: "v1alpha1", Kind: "middlewares"}
var middlewaresResource = schema.GroupVersionResource{Group: "traefik.containo.us", Version: "v1alpha1", Resource: "middlewares"}

func NewRegexRedirect(name, namespace, regex, replacement string, permanent bool) (Middleware, error) {
	m := Middleware{}
	m.SetGroupVersion()
	m.Meta.Name = name
	m.Meta.Namespace = namespace
	m.Spec = MiddlewareSpec{
		RedirectRegex: RedirectRegex{
			Regex:       regex,
			Replacement: replacement,
			Permanent:   permanent,
		},
	}

	return m, nil
}

func NewChain(name, namespace string, middlewares []Middleware) (*Middleware, error) {
	middlewareRef := []MiddlewareRef{}
	for _, m := range middlewares {
		middlewareRef = append(middlewareRef, MiddlewareRef{Name: m.Meta.Name})
	}
	c := new(Middleware)
	c.Meta.Name = name
	c.Meta.Namespace = namespace
	c.SetGroupVersion()
	c.Spec = MiddlewareSpec{
		Chain: Chain{
			Middlewares: middlewareRef,
		},
	}

	return c, nil
}

func (m Middleware) String() string {
	o, err := yaml.Marshal(&m)
	if err != nil {
		fmt.Println("fail to parse struct: " + err.Error())
	}

	return string(o)
}

func (m Middleware) Dump() {
	fmt.Println(string(m.String()))

}
