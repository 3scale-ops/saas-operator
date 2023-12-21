package config

import (
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
)

// NewOptions returns mapping-service options for the given saasv1alpha1.CORSProxySpec
func NewOptions(spec saasv1alpha1.MappingServiceSpec) pod.Options {
	opts := pod.Options{}

	opts.Unpack(spec.Config.SystemAdminToken).IntoEnvvar("MASTER_ACCESS_TOKEN").AsSecretRef("mapping-service-system-master-access-token")
	opts.Unpack(spec.Config.APIHost).IntoEnvvar("API_HOST")
	opts.Unpack("lazy").IntoEnvvar("APICAST_CONFIGURATION_LOADER")
	opts.Unpack(spec.Config.LogLevel).IntoEnvvar("APICAST_LOG_LEVEL")
	opts.Unpack(spec.Config.PreviewBaseDomain).IntoEnvvar("PREVIEW_BASE_DOMAIN")

	return opts
}
