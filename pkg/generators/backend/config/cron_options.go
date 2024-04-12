package config

import (
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/generators/seed"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
)

// NewCronOptions returns cron options for the given saasv1alpha1.BackendSpec
func NewCronOptions(spec saasv1alpha1.BackendSpec) pod.Options {
	opts := pod.Options{}

	opts.Unpack(spec.Config.RackEnv).IntoEnvvar("RACK_ENV")
	opts.Unpack(spec.Config.RedisStorageDSN).IntoEnvvar("CONFIG_REDIS_PROXY")
	opts.Unpack("").IntoEnvvar("CONFIG_REDIS_SENTINEL_HOSTS")
	opts.Unpack("").IntoEnvvar("CONFIG_REDIS_SENTINEL_ROLE")
	opts.Unpack(spec.Config.RedisQueuesDSN).IntoEnvvar("CONFIG_QUEUES_MASTER_NAME")
	opts.Unpack("").IntoEnvvar("CONFIG_QUEUES_SENTINEL_HOSTS")
	opts.Unpack("").IntoEnvvar("CONFIG_QUEUES_SENTINEL_ROLE")
	opts.Unpack(spec.Config.ErrorMonitoringService).IntoEnvvar("CONFIG_HOPTOAD_SERVICE").
		AsSecretRef(BackendErrorMonitoringSecret).
		WithSeedKey(seed.BackendErrorMonitoringService)
	opts.Unpack(spec.Config.ErrorMonitoringKey).IntoEnvvar("CONFIG_HOPTOAD_API_KEY").
		AsSecretRef(BackendErrorMonitoringSecret).
		WithSeedKey(seed.BackendErrorMonitoringApiKey)

	return opts
}
