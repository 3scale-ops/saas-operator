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

	opts.AddEnvvar("ACME_STAGING").Unpack(func() string {
		if *spec.Config.ACMEStaging {
			return leACMEStagingEndpoint
		}
		return ""
	}())
	opts.AddEnvvar("CONTACT_EMAIL").Unpack(spec.Config.ContactEmail)
	opts.AddEnvvar("PROXY_ENDPOINT").Unpack(spec.Config.ProxyEndpoint)
	opts.AddEnvvar("STORAGE_ADAPTER").Unpack("redis")
	opts.AddEnvvar("REDIS_HOST").Unpack(spec.Config.RedisHost)
	opts.AddEnvvar("REDIS_PORT").Unpack(spec.Config.RedisPort)
	opts.AddEnvvar("VERIFICATION_ENDPOINT").Unpack(spec.Config.VerificationEndpoint)
	opts.AddEnvvar("LOG_LEVEL").Unpack(spec.Config.LogLevel)
	opts.AddEnvvar("DOMAIN_WHITELIST").Unpack(strings.Join(spec.Config.DomainWhitelist, ","))
	opts.AddEnvvar("DOMAIN_BLACKLIST").Unpack(strings.Join(spec.Config.DomainBlacklist, ","))

	return opts
}
