package config

import (
	"fmt"
	"strings"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/resource_builders/pod"
)

const (
	leACMEStagingEndpoint = "https://acme-staging-v02.api.letsencrypt.org/directory"
)

// Options holds configuration for the sphinx pods
type Options struct {
	ACMEStaging          pod.EnvVarValue `env:"ACME_STAGING"`
	ContactEmail         pod.EnvVarValue `env:"CONTACT_EMAIL"`
	ProxyEndpoint        pod.EnvVarValue `env:"PROXY_ENDPOINT"`
	StorageAdapter       pod.EnvVarValue `env:"STORAGE_ADAPTER"`
	RedisHost            pod.EnvVarValue `env:"REDIS_HOST"`
	RedisPort            pod.EnvVarValue `env:"REDIS_PORT"`
	VerificationEndpoint pod.EnvVarValue `env:"VERIFICATION_ENDPOINT"`
	LogLevel             pod.EnvVarValue `env:"LOG_LEVEL"`
	DomainWhitelist      pod.EnvVarValue `env:"DOMAIN_WHITELIST"`
	DomainBlacklist      pod.EnvVarValue `env:"DOMAIN_BLACKLIST"`
}

// NewOptions returns an Options struct for the given saasv1alpha1.AutoSSLSpec
func NewOptions(spec saasv1alpha1.AutoSSLSpec) Options {
	opts := Options{
		ACMEStaging: &pod.ClearTextValue{Value: func() string {
			if *spec.Config.ACMEStaging {
				return leACMEStagingEndpoint
			}
			return ""
		}()},
		ContactEmail:         &pod.ClearTextValue{Value: spec.Config.ContactEmail},
		ProxyEndpoint:        &pod.ClearTextValue{Value: spec.Config.ProxyEndpoint},
		StorageAdapter:       &pod.ClearTextValue{Value: "redis"},
		RedisHost:            &pod.ClearTextValue{Value: spec.Config.RedisHost},
		RedisPort:            &pod.ClearTextValue{Value: fmt.Sprintf("%v", *spec.Config.RedisPort)},
		VerificationEndpoint: &pod.ClearTextValue{Value: spec.Config.VerificationEndpoint},
		LogLevel:             &pod.ClearTextValue{Value: *spec.Config.LogLevel},
		DomainWhitelist:      &pod.ClearTextValue{Value: strings.Join(spec.Config.DomainWhitelist, ",")},
		DomainBlacklist:      &pod.ClearTextValue{Value: strings.Join(spec.Config.DomainBlacklist, ",")},
	}

	return opts
}
