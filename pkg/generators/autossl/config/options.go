package config

import (
	"strings"

	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
)

const (
	leACMEStagingEndpoint = "https://acme-staging-v02.api.letsencrypt.org/directory"
)

func NewOptions(spec saasv1alpha1.AutoSSLSpec) pod.Options {
	opts := pod.Options{}

	opts.Unpack(func() string {
		if *spec.Config.ACMEStaging {
			return leACMEStagingEndpoint
		}
		return ""
	}()).IntoEnvvar("ACME_STAGING")
	opts.Unpack(spec.Config.ContactEmail).IntoEnvvar("CONTACT_EMAIL")
	opts.Unpack(spec.Config.ProxyEndpoint).IntoEnvvar("PROXY_ENDPOINT")
	opts.Unpack("redis").IntoEnvvar("STORAGE_ADAPTER")
	opts.Unpack(spec.Config.RedisHost).IntoEnvvar("REDIS_HOST")
	opts.Unpack(spec.Config.RedisPort).IntoEnvvar("REDIS_PORT")
	opts.Unpack(spec.Config.VerificationEndpoint).IntoEnvvar("VERIFICATION_ENDPOINT")
	opts.Unpack(spec.Config.LogLevel).IntoEnvvar("LOG_LEVEL")
	opts.Unpack(strings.Join(spec.Config.DomainWhitelist, ",")).IntoEnvvar("DOMAIN_WHITELIST")
	opts.Unpack(strings.Join(spec.Config.DomainBlacklist, ",")).IntoEnvvar("DOMAIN_BLACKLIST")

	return opts
}
