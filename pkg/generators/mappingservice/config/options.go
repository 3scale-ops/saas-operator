package config

import (
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/resource_builders/pod"
)

// Options holds configuration for the mapping-service pod
type Options struct {
	MasterAccessToken          pod.EnvVarValue `env:"MASTER_ACCESS_TOKEN" secret:"mapping-service-system-master-access-token"`
	APIHost                    pod.EnvVarValue `env:"API_HOST"`
	ApicastConfigurationLoader pod.EnvVarValue `env:"APICAST_CONFIGURATION_LOADER"`
	ApicastLogLevel            pod.EnvVarValue `env:"APICAST_LOG_LEVEL"`
	PreviewBaseDomain          pod.EnvVarValue `env:"PREVIEW_BASE_DOMAIN"`
}

// NewOptions returns an Options struct for the given saasv1alpha1.CORSProxySpec
func NewOptions(spec saasv1alpha1.MappingServiceSpec) Options {
	opts := Options{
		MasterAccessToken:          &pod.SecretValue{Value: spec.Config.SystemAdminToken},
		APIHost:                    &pod.ClearTextValue{Value: spec.Config.APIHost},
		ApicastConfigurationLoader: &pod.ClearTextValue{Value: "lazy"},
		ApicastLogLevel:            &pod.ClearTextValue{Value: *spec.Config.LogLevel},
	}
	if spec.Config.PreviewBaseDomain != nil {
		opts.PreviewBaseDomain = &pod.ClearTextValue{Value: *spec.Config.PreviewBaseDomain}
	}

	return opts
}
