package config

import (
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
)

// NewOptions returns cors-proxy options the given saasv1alpha1.CORSProxySpec
func NewOptions(spec saasv1alpha1.CORSProxySpec) pod.Options {
	opts := pod.Options{}
	opts.Unpack(spec.Config.SystemDatabaseDSN).IntoEnvvar("DATABASE_URL").AsSecretRef("cors-proxy-system-database")
	return opts
}
