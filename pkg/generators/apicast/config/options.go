package config

import (
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
)

func NewEnvOptions(spec saasv1alpha1.ApicastEnvironmentSpec, env string) pod.Options {
	opts := pod.Options{}

	opts.Unpack("lazy").IntoEnvvar("APICAST_CONFIGURATION_LOADER")
	opts.Unpack(spec.Config.ConfigurationCache).IntoEnvvar("APICAST_CONFIGURATION_CACHE")
	opts.Unpack("true").IntoEnvvar("APICAST_EXTENDED_METRICS")
	opts.Unpack(env).IntoEnvvar("THREESCALE_DEPLOYMENT_ENV")
	opts.Unpack(spec.Config.ThreescalePortalEndpoint).IntoEnvvar("THREESCALE_PORTAL_ENDPOINT")
	opts.Unpack(spec.Config.LogLevel).IntoEnvvar("APICAST_LOG_LEVEL")
	opts.Unpack(spec.Config.OIDCLogLevel).IntoEnvvar("APICAST_OIDC_LOG_LEVEL")
	opts.Unpack("true").IntoEnvvar("APICAST_RESPONSE_CODES")

	return opts
}
