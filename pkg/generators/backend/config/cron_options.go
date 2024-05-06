package config

import (
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/generators/seed"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
)

// NewCronOptions returns cron options for the given saasv1alpha1.BackendSpec
func NewCronOptions(spec saasv1alpha1.BackendSpec) pod.Options {
	opts := pod.Options{}

	opts.AddEnvvar("RACK_ENV").Unpack(spec.Config.RackEnv)
	opts.AddEnvvar("CONFIG_REDIS_PROXY").Unpack(spec.Config.RedisStorageDSN)
	opts.AddEnvvar("CONFIG_REDIS_SENTINEL_HOSTS").Unpack("")
	opts.AddEnvvar("CONFIG_REDIS_SENTINEL_ROLE").Unpack("")
	opts.AddEnvvar("CONFIG_QUEUES_MASTER_NAME").Unpack(spec.Config.RedisQueuesDSN)
	opts.AddEnvvar("CONFIG_QUEUES_SENTINEL_HOSTS").Unpack("")
	opts.AddEnvvar("CONFIG_QUEUES_SENTINEL_ROLE").Unpack("")
	opts.AddEnvvar("CONFIG_HOPTOAD_SERVICE").AsSecretRef(BackendErrorMonitoringSecret).WithSeedKey(seed.BackendErrorMonitoringService).
		Unpack(spec.Config.ErrorMonitoringService)
	opts.AddEnvvar("CONFIG_HOPTOAD_API_KEY").AsSecretRef(BackendErrorMonitoringSecret).WithSeedKey(seed.BackendErrorMonitoringApiKey).
		Unpack(spec.Config.ErrorMonitoringKey)

	return opts
}
