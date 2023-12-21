package config

import (
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
)

// NewWorkerOptions returns worker options for the given saasv1alpha1.BackedSpec
func NewWorkerOptions(spec saasv1alpha1.BackendSpec) pod.Options {
	opts := pod.Options{}

	opts.Unpack(spec.Config.RackEnv).IntoEnvvar("RACK_ENV")
	opts.Unpack(spec.Config.RedisStorageDSN).IntoEnvvar("CONFIG_REDIS_PROXY")
	opts.Unpack("").IntoEnvvar("CONFIG_REDIS_SENTINEL_HOSTS")
	opts.Unpack("").IntoEnvvar("CONFIG_REDIS_SENTINEL_ROLE")
	opts.Unpack(spec.Config.RedisQueuesDSN).IntoEnvvar("CONFIG_QUEUES_MASTER_NAME")
	opts.Unpack("").IntoEnvvar("CONFIG_QUEUES_SENTINEL_HOSTS")
	opts.Unpack("").IntoEnvvar("CONFIG_QUEUES_SENTINEL_ROLE")
	opts.Unpack(spec.Config.MasterServiceID).IntoEnvvar("CONFIG_MASTER_SERVICE_ID")
	opts.Unpack(spec.Worker.Config.RedisAsync).IntoEnvvar("CONFIG_REDIS_ASYNC")
	opts.Unpack(spec.Worker.Config.LogFormat).IntoEnvvar("CONFIG_WORKERS_LOGGER_FORMATTER")
	opts.Unpack("true").IntoEnvvar("CONFIG_WORKER_PROMETHEUS_METRICS_ENABLED")
	opts.Unpack("9421").IntoEnvvar("CONFIG_WORKER_PROMETHEUS_METRICS_PORT")
	opts.Unpack(spec.Config.SystemEventsHookURL).IntoEnvvar("CONFIG_EVENTS_HOOK").AsSecretRef("backend-system-events-hook")
	opts.Unpack(spec.Config.SystemEventsHookPassword).IntoEnvvar("CONFIG_EVENTS_HOOK_SHARED_SECRET").AsSecretRef("backend-system-events-hook")
	opts.Unpack(spec.Config.ErrorMonitoringService).IntoEnvvar("CONFIG_HOPTOAD_SERVICE").AsSecretRef("backend-error-monitoring")
	opts.Unpack(spec.Config.ErrorMonitoringKey).IntoEnvvar("CONFIG_HOPTOAD_API_KEY").AsSecretRef("backend-error-monitoring")

	return opts
}
