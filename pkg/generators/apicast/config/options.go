package config

import (
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
)

func NewEnvOptions(spec saasv1alpha1.ApicastEnvironmentSpec, env string) pod.Options {
	opts := pod.Options{}

	opts.AddEnvvar("APICAST_CONFIGURATION_LOADER").Unpack("lazy")
	opts.AddEnvvar("APICAST_CONFIGURATION_CACHE").Unpack(spec.Config.ConfigurationCache)
	opts.AddEnvvar("APICAST_EXTENDED_METRICS").Unpack("true")
	opts.AddEnvvar("THREESCALE_DEPLOYMENT_ENV").Unpack(env)
	opts.AddEnvvar("THREESCALE_PORTAL_ENDPOINT").Unpack(spec.Config.ThreescalePortalEndpoint)
	opts.AddEnvvar("APICAST_LOG_LEVEL").Unpack(spec.Config.LogLevel)
	opts.AddEnvvar("APICAST_OIDC_LOG_LEVEL").Unpack(spec.Config.OIDCLogLevel)
	opts.AddEnvvar("APICAST_RESPONSE_CODES").Unpack("true")

	return opts
}
