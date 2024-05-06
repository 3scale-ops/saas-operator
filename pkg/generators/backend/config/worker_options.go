package config

import (
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/generators/seed"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
)

// NewWorkerOptions returns worker options for the given saasv1alpha1.BackedSpec
func NewWorkerOptions(spec saasv1alpha1.BackendSpec) pod.Options {
	opts := pod.Options{}

	opts.AddEnvvar("RACK_ENV").Unpack(spec.Config.RackEnv)
	opts.AddEnvvar("CONFIG_REDIS_PROXY").Unpack(spec.Config.RedisStorageDSN)
	opts.AddEnvvar("CONFIG_REDIS_SENTINEL_HOSTS").Unpack("")
	opts.AddEnvvar("CONFIG_REDIS_SENTINEL_ROLE").Unpack("")
	opts.AddEnvvar("CONFIG_QUEUES_MASTER_NAME").Unpack(spec.Config.RedisQueuesDSN)
	opts.AddEnvvar("CONFIG_QUEUES_SENTINEL_HOSTS").Unpack("")
	opts.AddEnvvar("CONFIG_QUEUES_SENTINEL_ROLE").Unpack("")
	opts.AddEnvvar("CONFIG_MASTER_SERVICE_ID").Unpack(spec.Config.MasterServiceID)
	opts.AddEnvvar("CONFIG_REDIS_ASYNC").Unpack(spec.Worker.Config.RedisAsync)
	opts.AddEnvvar("CONFIG_WORKERS_LOGGER_FORMATTER").Unpack(spec.Worker.Config.LogFormat)
	opts.AddEnvvar("CONFIG_WORKER_PROMETHEUS_METRICS_ENABLED").Unpack("true")
	opts.AddEnvvar("CONFIG_WORKER_PROMETHEUS_METRICS_PORT").Unpack("9421")
	opts.AddEnvvar("CONFIG_EVENTS_HOOK").AsSecretRef(BackendSystemEventsSecret).WithSeedKey(seed.SystemEventsHookURL).
		Unpack(spec.Config.SystemEventsHookURL)
	opts.AddEnvvar("CONFIG_EVENTS_HOOK_SHARED_SECRET").AsSecretRef(BackendSystemEventsSecret).WithSeedKey(seed.SystemEventsHookSharedSecret).
		Unpack(spec.Config.SystemEventsHookPassword)
	opts.AddEnvvar("CONFIG_HOPTOAD_SERVICE").AsSecretRef(BackendErrorMonitoringSecret).WithSeedKey(seed.BackendErrorMonitoringService).
		Unpack(spec.Config.ErrorMonitoringService)
	opts.AddEnvvar("CONFIG_HOPTOAD_API_KEY").AsSecretRef(BackendErrorMonitoringSecret).WithSeedKey(seed.BackendErrorMonitoringApiKey).
		Unpack(spec.Config.ErrorMonitoringKey)

	return opts
}
