package config

import (
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/generators/seed"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
)

type Secret string

func (s Secret) String() string { return string(s) }

const (
	MappingServiceSystemMasterAccessTokenSecret Secret = "mapping-service-system-master-access-token"
)

// NewOptions returns mapping-service options for the given saasv1alpha1.CORSProxySpec
func NewOptions(spec saasv1alpha1.MappingServiceSpec) pod.Options {
	opts := pod.Options{}

	opts.AddEnvvar("MASTER_ACCESS_TOKEN").AsSecretRef(MappingServiceSystemMasterAccessTokenSecret).WithSeedKey(seed.SystemApicastAccessToken).
		Unpack(spec.Config.SystemAdminToken)
	opts.AddEnvvar("API_HOST").Unpack(spec.Config.APIHost)
	opts.AddEnvvar("APICAST_CONFIGURATION_LOADER").Unpack("lazy")
	opts.AddEnvvar("APICAST_LOG_LEVEL").Unpack(spec.Config.LogLevel)
	opts.AddEnvvar("PREVIEW_BASE_DOMAIN").Unpack(spec.Config.PreviewBaseDomain)

	return opts
}
