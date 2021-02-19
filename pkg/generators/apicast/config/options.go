package config

import (
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pod"
)

// EnvOptions holds configuration for the sphinx pods
type EnvOptions struct {
	ApicastConfigurationLoader pod.EnvVarValue `env:"APICAST_CONFIGURATION_LOADER"`
	ApicastConfigurationCache  pod.EnvVarValue `env:"APICAST_CONFIGURATION_CACHE"`
	ApicastExtendedMetrics     pod.EnvVarValue `env:"APICAST_EXTENDED_METRICS"`
	ThreeScaleDeploymentEnv    pod.EnvVarValue `env:"THREESCALE_DEPLOYMENT_ENV"`
	ThreescalePortalEndpoint   pod.EnvVarValue `env:"THREESCALE_PORTAL_ENDPOINT"`
	ApicastLogLevel            pod.EnvVarValue `env:"APICAST_LOG_LEVEL"`
	ApicastOIDCLogLevel        pod.EnvVarValue `env:"APICAST_OIDC_LOG_LEVEL"`
	ApicastResponseCodes       pod.EnvVarValue `env:"APICAST_RESPONSE_CODES"`
}

// NewEnvOptions returns an Options struct for the given saasv1alpha1.ApicastEnvironmentSpec
func NewEnvOptions(spec saasv1alpha1.ApicastEnvironmentSpec, env string) EnvOptions {
	opts := EnvOptions{
		ApicastConfigurationLoader: &pod.ClearTextValue{Value: "lazy"},
		ApicastConfigurationCache:  &pod.ClearTextValue{Value: fmt.Sprintf("%d", spec.Config.ConfigurationCache)},
		ApicastExtendedMetrics:     &pod.ClearTextValue{Value: "true"},
		ThreeScaleDeploymentEnv:    &pod.ClearTextValue{Value: env},
		ThreescalePortalEndpoint:   &pod.ClearTextValue{Value: spec.Config.ThreescalePortalEndpoint},
		ApicastLogLevel:            &pod.ClearTextValue{Value: *spec.Config.LogLevel},
		ApicastOIDCLogLevel:        &pod.ClearTextValue{Value: *spec.Config.OIDCLogLevel},
		ApicastResponseCodes:       &pod.ClearTextValue{Value: "true"},
	}
	return opts
}
