package autossl

import (
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators"
)

const (
	component string = "autossl"
)

// Generator configures the generators for AutoSSL
type Generator struct {
	generators.BaseOptions
	Spec saasv1alpha1.AutoSSLSpec
}

// NewAutoSSLGenerator returns a new Options struct
func NewAutoSSLGenerator(instance, namespace string, spec saasv1alpha1.AutoSSLSpec) Generator {
	return Generator{
		BaseOptions: generators.BaseOptions{
			Component:    component,
			InstanceName: instance,
			Namespace:    namespace,
		},
		Spec: spec,
	}
}
