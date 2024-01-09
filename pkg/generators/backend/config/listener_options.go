package config

import (
	"fmt"
	"strconv"

	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
)

// ListenerOptions holds configuration for the listener pods
type ListenerOptions struct {
	RackEnv                                pod.EnvVarValue `env:"RACK_ENV"`
	ConfigRedisProxy                       pod.EnvVarValue `env:"CONFIG_REDIS_PROXY"`
	ConfigRedisSentinelHosts               pod.EnvVarValue `env:"CONFIG_REDIS_SENTINEL_HOSTS"`
	ConfigRedisSentinelRole                pod.EnvVarValue `env:"CONFIG_REDIS_SENTINEL_ROLE"`
	ConfigQueuesMasterName                 pod.EnvVarValue `env:"CONFIG_QUEUES_MASTER_NAME"`
	ConfigQueuesSentinelHosts              pod.EnvVarValue `env:"CONFIG_QUEUES_SENTINEL_HOSTS"`
	ConfigQueuesSentinelRole               pod.EnvVarValue `env:"CONFIG_QUEUES_SENTINEL_ROLE"`
	ConfigMasterServiceID                  pod.EnvVarValue `env:"CONFIG_MASTER_SERVICE_ID"`
	ConfigRequestLoggers                   pod.EnvVarValue `env:"CONFIG_REQUEST_LOGGERS"`
	ConfigRedisAsync                       pod.EnvVarValue `env:"CONFIG_REDIS_ASYNC"`
	ListenerWorkers                        pod.EnvVarValue `env:"LISTENER_WORKERS"`
	ConfigLegacyReferrerFilters            pod.EnvVarValue `env:"CONFIG_LEGACY_REFERRER_FILTERS"`
	ConfigListenerPrometheusMetricsEnabled pod.EnvVarValue `env:"CONFIG_LISTENER_PROMETHEUS_METRICS_ENABLED"`
	ConfigInternalAPIUser                  pod.EnvVarValue `env:"CONFIG_INTERNAL_API_USER" secret:"backend-internal-api"`
	ConfigInternalAPIPassword              pod.EnvVarValue `env:"CONFIG_INTERNAL_API_PASSWORD" secret:"backend-internal-api"`
	ConfigHoptoadService                   pod.EnvVarValue `env:"CONFIG_HOPTOAD_SERVICE" secret:"backend-error-monitoring"`
	ConfigHoptoadAPIKey                    pod.EnvVarValue `env:"CONFIG_HOPTOAD_API_KEY" secret:"backend-error-monitoring"`
}

// NewListenerOptions returns an Options struct for the given saasv1alpha1.BackendSpec
func NewListenerOptions(spec saasv1alpha1.BackendSpec) ListenerOptions {
	opts := ListenerOptions{
		RackEnv:                                &pod.ClearTextValue{Value: *spec.Config.RackEnv},
		ConfigRedisProxy:                       &pod.ClearTextValue{Value: spec.Config.RedisStorageDSN},
		ConfigRedisSentinelHosts:               &pod.ClearTextValue{Value: ""},
		ConfigRedisSentinelRole:                &pod.ClearTextValue{Value: ""},
		ConfigQueuesMasterName:                 &pod.ClearTextValue{Value: spec.Config.RedisQueuesDSN},
		ConfigQueuesSentinelHosts:              &pod.ClearTextValue{Value: ""},
		ConfigQueuesSentinelRole:               &pod.ClearTextValue{Value: ""},
		ConfigMasterServiceID:                  &pod.ClearTextValue{Value: fmt.Sprintf("%d", *spec.Config.MasterServiceID)},
		ConfigRequestLoggers:                   &pod.ClearTextValue{Value: *spec.Listener.Config.LogFormat},
		ConfigRedisAsync:                       &pod.ClearTextValue{Value: strconv.FormatBool(*spec.Listener.Config.RedisAsync)},
		ListenerWorkers:                        &pod.ClearTextValue{Value: fmt.Sprintf("%d", *spec.Listener.Config.ListenerWorkers)},
		ConfigLegacyReferrerFilters:            &pod.ClearTextValue{Value: strconv.FormatBool(*spec.Listener.Config.LegacyReferrerFilters)},
		ConfigListenerPrometheusMetricsEnabled: &pod.ClearTextValue{Value: "true"},

		ConfigInternalAPIUser:     &pod.SecretValue{Value: spec.Config.InternalAPIUser},
		ConfigInternalAPIPassword: &pod.SecretValue{Value: spec.Config.InternalAPIPassword},
	}

	if spec.Config.ErrorMonitoringService != nil && spec.Config.ErrorMonitoringKey != nil {
		opts.ConfigHoptoadService = &pod.SecretValue{Value: *spec.Config.ErrorMonitoringService}
		opts.ConfigHoptoadAPIKey = &pod.SecretValue{Value: *spec.Config.ErrorMonitoringKey}
	}

	return opts
}
