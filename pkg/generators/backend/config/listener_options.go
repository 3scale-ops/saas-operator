package config

import (
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/generators/seed"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
)

// NewListenerOptions returns listener options for the given saasv1alpha1.BackendSpec
func NewListenerOptions(spec saasv1alpha1.BackendSpec) pod.Options {
	opts := pod.Options{}

	opts.AddEnvvar("RACK_ENV").Unpack(spec.Config.RackEnv)
	opts.AddEnvvar("CONFIG_REDIS_PROXY").Unpack(spec.Config.RedisStorageDSN)
	opts.AddEnvvar("CONFIG_REDIS_SENTINEL_HOSTS").Unpack("")
	opts.AddEnvvar("CONFIG_REDIS_SENTINEL_ROLE").Unpack("")
	opts.AddEnvvar("CONFIG_QUEUES_MASTER_NAME").Unpack(spec.Config.RedisQueuesDSN)
	opts.AddEnvvar("CONFIG_QUEUES_SENTINEL_HOSTS").Unpack("")
	opts.AddEnvvar("CONFIG_QUEUES_SENTINEL_ROLE").Unpack("")
	opts.AddEnvvar("CONFIG_MASTER_SERVICE_ID").Unpack(spec.Config.MasterServiceID)
	opts.AddEnvvar("CONFIG_REQUEST_LOGGERS").Unpack(spec.Listener.Config.LogFormat)
	opts.AddEnvvar("CONFIG_REDIS_ASYNC").Unpack(spec.Listener.Config.RedisAsync)
	opts.AddEnvvar("LISTENER_WORKERS").Unpack(spec.Listener.Config.ListenerWorkers)
	opts.AddEnvvar("CONFIG_LEGACY_REFERRER_FILTERS").Unpack(spec.Listener.Config.LegacyReferrerFilters)
	opts.AddEnvvar("CONFIG_LISTENER_PROMETHEUS_METRICS_ENABLED").Unpack("true")
	opts.AddEnvvar("CONFIG_INTERNAL_API_USER").AsSecretRef(BackendInternalApiSecret).WithSeedKey(seed.BackendInternalApiUser).
		Unpack(spec.Config.InternalAPIUser)
	opts.AddEnvvar("CONFIG_INTERNAL_API_PASSWORD").AsSecretRef(BackendInternalApiSecret).WithSeedKey(seed.BackendInternalApiPassword).
		Unpack(spec.Config.InternalAPIPassword)
	opts.AddEnvvar("CONFIG_HOPTOAD_SERVICE").AsSecretRef(BackendErrorMonitoringSecret).WithSeedKey(seed.BackendErrorMonitoringService).
		Unpack(spec.Config.ErrorMonitoringService)
	opts.AddEnvvar("CONFIG_HOPTOAD_API_KEY").AsSecretRef(BackendErrorMonitoringSecret).WithSeedKey(seed.BackendErrorMonitoringApiKey).
		Unpack(spec.Config.ErrorMonitoringKey)

	return opts
}
