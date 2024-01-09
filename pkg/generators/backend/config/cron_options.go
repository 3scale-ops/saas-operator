package config

import (
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
)

// CronOptions holds configuration for the cron pods
type CronOptions struct {
	RackEnv                   pod.EnvVarValue `env:"RACK_ENV"`
	ConfigRedisProxy          pod.EnvVarValue `env:"CONFIG_REDIS_PROXY"`
	ConfigRedisSentinelHosts  pod.EnvVarValue `env:"CONFIG_REDIS_SENTINEL_HOSTS"`
	ConfigRedisSentinelRole   pod.EnvVarValue `env:"CONFIG_REDIS_SENTINEL_ROLE"`
	ConfigQueuesMasterName    pod.EnvVarValue `env:"CONFIG_QUEUES_MASTER_NAME"`
	ConfigQueuesSentinelHosts pod.EnvVarValue `env:"CONFIG_QUEUES_SENTINEL_HOSTS"`
	ConfigQueuesSentinelRole  pod.EnvVarValue `env:"CONFIG_QUEUES_SENTINEL_ROLE"`
	ConfigHoptoadService      pod.EnvVarValue `env:"CONFIG_HOPTOAD_SERVICE" secret:"backend-error-monitoring"`
	ConfigHoptoadAPIKey       pod.EnvVarValue `env:"CONFIG_HOPTOAD_API_KEY" secret:"backend-error-monitoring"`
}

// NewCronOptions returns a CronOptions struct for the given saasv1alpha1.BackendSpec
func NewCronOptions(spec saasv1alpha1.BackendSpec) CronOptions {
	opts := CronOptions{
		RackEnv:                   &pod.ClearTextValue{Value: *spec.Config.RackEnv},
		ConfigRedisProxy:          &pod.ClearTextValue{Value: spec.Config.RedisStorageDSN},
		ConfigRedisSentinelHosts:  &pod.ClearTextValue{Value: ""},
		ConfigRedisSentinelRole:   &pod.ClearTextValue{Value: ""},
		ConfigQueuesMasterName:    &pod.ClearTextValue{Value: spec.Config.RedisQueuesDSN},
		ConfigQueuesSentinelHosts: &pod.ClearTextValue{Value: ""},
		ConfigQueuesSentinelRole:  &pod.ClearTextValue{Value: ""},
	}

	if spec.Config.ErrorMonitoringService != nil && spec.Config.ErrorMonitoringKey != nil {
		opts.ConfigHoptoadService = &pod.SecretValue{Value: *spec.Config.ErrorMonitoringService}
		opts.ConfigHoptoadAPIKey = &pod.SecretValue{Value: *spec.Config.ErrorMonitoringKey}
	}

	return opts
}
