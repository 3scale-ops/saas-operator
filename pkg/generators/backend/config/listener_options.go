package config

import (
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/generators/seed"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
)

// NewListenerOptions returns listener options for the given saasv1alpha1.BackendSpec
func NewListenerOptions(spec saasv1alpha1.BackendSpec) pod.Options {
	opts := pod.Options{}

	opts.Unpack(spec.Config.RackEnv).IntoEnvvar("RACK_ENV")
	opts.Unpack(spec.Config.RedisStorageDSN).IntoEnvvar("CONFIG_REDIS_PROXY")
	opts.Unpack("").IntoEnvvar("CONFIG_REDIS_SENTINEL_HOSTS")
	opts.Unpack("").IntoEnvvar("CONFIG_REDIS_SENTINEL_ROLE")
	opts.Unpack(spec.Config.RedisQueuesDSN).IntoEnvvar("CONFIG_QUEUES_MASTER_NAME")
	opts.Unpack("").IntoEnvvar("CONFIG_QUEUES_SENTINEL_HOSTS")
	opts.Unpack("").IntoEnvvar("CONFIG_QUEUES_SENTINEL_ROLE")
	opts.Unpack(spec.Config.MasterServiceID).IntoEnvvar("CONFIG_MASTER_SERVICE_ID")
	opts.Unpack(spec.Listener.Config.LogFormat).IntoEnvvar("CONFIG_REQUEST_LOGGERS")
	opts.Unpack(spec.Listener.Config.RedisAsync).IntoEnvvar("CONFIG_REDIS_ASYNC")
	opts.Unpack(spec.Listener.Config.ListenerWorkers).IntoEnvvar("LISTENER_WORKERS")
	opts.Unpack(spec.Listener.Config.LegacyReferrerFilters).IntoEnvvar("CONFIG_LEGACY_REFERRER_FILTERS")
	opts.Unpack("true").IntoEnvvar("CONFIG_LISTENER_PROMETHEUS_METRICS_ENABLED")
	opts.Unpack(spec.Config.InternalAPIUser).IntoEnvvar("CONFIG_INTERNAL_API_USER").
		AsSecretRef(BackendInternalApiSecret).
		WithSeedKey(seed.BackendInternalApiUser)
	opts.Unpack(spec.Config.InternalAPIPassword).IntoEnvvar("CONFIG_INTERNAL_API_PASSWORD").
		AsSecretRef(BackendInternalApiSecret).
		WithSeedKey(seed.BackendInternalApiPassword)
	opts.Unpack(spec.Config.ErrorMonitoringService).IntoEnvvar("CONFIG_HOPTOAD_SERVICE").
		AsSecretRef(BackendErrorMonitoringSecret).
		WithSeedKey(seed.BackendErrorMonitoringService)
	opts.Unpack(spec.Config.ErrorMonitoringKey).IntoEnvvar("CONFIG_HOPTOAD_API_KEY").
		AsSecretRef(BackendErrorMonitoringSecret).
		WithSeedKey(seed.BackendErrorMonitoringApiKey)

	return opts
}
