package config

import (
	"fmt"
	"strconv"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pod"
)

// WorkerOptions holds configuration for the worker pods
type WorkerOptions struct {
	RackEnv                              pod.EnvVarValue `env:"RACK_ENV"`
	ConfigRedisProxy                     pod.EnvVarValue `env:"CONFIG_REDIS_PROXY"`
	ConfigRedisSentinelHosts             pod.EnvVarValue `env:"CONFIG_REDIS_SENTINEL_HOSTS"`
	ConfigRedisSentinelRole              pod.EnvVarValue `env:"CONFIG_REDIS_SENTINEL_ROLE"`
	ConfigQueuesMasterName               pod.EnvVarValue `env:"CONFIG_QUEUES_MASTER_NAME"`
	ConfigQueuesSentinelHosts            pod.EnvVarValue `env:"CONFIG_QUEUES_SENTINEL_HOSTS"`
	ConfigQueuesSentinelRole             pod.EnvVarValue `env:"CONFIG_QUEUES_SENTINEL_ROLE"`
	ConfigMasterServiceID                pod.EnvVarValue `env:"CONFIG_MASTER_SERVICE_ID"`
	ConfigRedisAsync                     pod.EnvVarValue `env:"CONFIG_REDIS_ASYNC"`
	ConfigWorkersLoggerFormatter         pod.EnvVarValue `env:"CONFIG_WORKERS_LOGGER_FORMATTER"`
	ConfigWorkerPrometheusMetricsEnabled pod.EnvVarValue `env:"CONFIG_WORKER_PROMETHEUS_METRICS_PORT"`
	ConfigWorkerPrometheusMetricsPort    pod.EnvVarValue `env:"CONFIG_WORKER_PROMETHEUS_METRICS_ENABLED"`
	ConfigInternalAPIUser                pod.EnvVarValue `env:"CONFIG_INTERNAL_API_USER" secret:"backend-internal-api"`
	ConfigInternalAPIPassword            pod.EnvVarValue `env:"CONFIG_INTERNAL_API_PASSWORD" secret:"backend-internal-api"`
	ConfigEventsHook                     pod.EnvVarValue `env:"CONFIG_EVENTS_HOOK" secret:"backend-system-events-hook"`
	ConfigEventsHookSharedSecret         pod.EnvVarValue `env:"CONFIG_EVENTS_HOOK_SHARED_SECRET" secret:"backend-system-events-hook"`
	ConfigHoptoadService                 pod.EnvVarValue `env:"CONFIG_HOPTOAD_SERVICE" secret:"backend-error-monitoring"`
	ConfigHoptoadAPIKey                  pod.EnvVarValue `env:"CONFIG_HOPTOAD_API_KEY" secret:"backend-error-monitoring"`
}

// NewWorkerOptions returns an Options struct for the given saasv1alpha1.BackedSpec
func NewWorkerOptions(spec saasv1alpha1.BackendSpec) WorkerOptions {
	opts := WorkerOptions{
		RackEnv:                              &pod.ClearTextValue{Value: *spec.Config.RackEnv},
		ConfigRedisProxy:                     &pod.ClearTextValue{Value: spec.Config.RedisStorageDSN},
		ConfigRedisSentinelHosts:             &pod.ClearTextValue{Value: ""},
		ConfigRedisSentinelRole:              &pod.ClearTextValue{Value: ""},
		ConfigQueuesMasterName:               &pod.ClearTextValue{Value: spec.Config.RedisQueuesDSN},
		ConfigQueuesSentinelHosts:            &pod.ClearTextValue{Value: ""},
		ConfigQueuesSentinelRole:             &pod.ClearTextValue{Value: ""},
		ConfigMasterServiceID:                &pod.ClearTextValue{Value: fmt.Sprintf("%d", *spec.Config.MasterServiceID)},
		ConfigRedisAsync:                     &pod.ClearTextValue{Value: strconv.FormatBool(*spec.Worker.Config.RedisAsync)},
		ConfigWorkersLoggerFormatter:         &pod.ClearTextValue{Value: *spec.Worker.Config.LogFormat},
		ConfigWorkerPrometheusMetricsEnabled: &pod.ClearTextValue{Value: "true"},
		ConfigWorkerPrometheusMetricsPort:    &pod.ClearTextValue{Value: "9421"},

		ConfigInternalAPIUser:        &pod.SecretValue{Value: spec.Config.InternalAPIUser},
		ConfigInternalAPIPassword:    &pod.SecretValue{Value: spec.Config.InternalAPIPassword},
		ConfigEventsHook:             &pod.SecretValue{Value: spec.Config.SystemEventsHookURL},
		ConfigEventsHookSharedSecret: &pod.SecretValue{Value: spec.Config.SystemEventsHookPassword},
	}

	if spec.Config.ErrorMonitoringService != nil && spec.Config.ErrorMonitoringKey != nil {
		opts.ConfigHoptoadService = &pod.SecretValue{Value: *spec.Config.ErrorMonitoringService}
		opts.ConfigHoptoadAPIKey = &pod.SecretValue{Value: *spec.Config.ErrorMonitoringKey}
	}

	return opts
}
