package config

import (
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/resource_builders/pod"
)

// Options holds configuration for the cors-proxy pod
type Options struct {
	DatabaseURL pod.EnvVarValue `env:"DATABASE_URL" secret:"cors-proxy-system-database"`
}

// NewOptions returns an Options struct for the given saasv1alpha1.CORSProxySpec
func NewOptions(spec saasv1alpha1.CORSProxySpec) Options {
	return Options{
		DatabaseURL: &pod.SecretValue{Value: spec.Config.SystemDatabaseDSN},
	}
}
