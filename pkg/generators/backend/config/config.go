package config

import "github.com/3scale/saas-operator/pkg/generators/common_blocks/secrets"

const (
	RackEnv                                string = "RACK_ENV"
	ConfigRedisProxy                       string = "CONFIG_REDIS_PROXY"
	ConfigRedisSentinelHosts               string = "CONFIG_REDIS_SENTINEL_HOSTS"
	ConfigRedisSentinelRole                string = "CONFIG_REDIS_SENTINEL_ROLE"
	ConfigQueuesMasterName                 string = "CONFIG_QUEUES_MASTER_NAME"
	ConfigQueuesSentinelHosts              string = "CONFIG_QUEUES_SENTINEL_HOSTS"
	ConfigQueuesSentinelRole               string = "CONFIG_QUEUES_SENTINEL_ROLE"
	ConfigMasterServiceID                  string = "CONFIG_MASTER_SERVICE_ID"
	ConfigRequestLoggers                   string = "CONFIG_REQUEST_LOGGERS"
	ConfigRedisAsync                       string = "CONFIG_REDIS_ASYNC"
	ListenerWorkers                        string = "LISTENER_WORKERS"
	ConfigLegacyReferrerFilters            string = "CONFIG_LEGACY_REFERRER_FILTERS"
	ConfigListenerPrometheusMetricsEnabled string = "CONFIG_LISTENER_PROMETHEUS_METRICS_ENABLED"
	ConfigWorkersLoggerFormatter           string = "CONFIG_WORKERS_LOGGER_FORMATTER"
	ConfigWorkerPrometheusMetricsPort      string = "CONFIG_WORKER_PROMETHEUS_METRICS_ENABLED"
	ConfigWorkerPrometheusMetricsEnabled   string = "CONFIG_WORKER_PROMETHEUS_METRICS_PORT"
	ConfigInternalAPIUser                  string = "CONFIG_INTERNAL_API_USER"
	ConfigInternalAPIPassword              string = "CONFIG_INTERNAL_API_PASSWORD"
	ConfigEventsHook                       string = "CONFIG_EVENTS_HOOK"
	ConfigEventsHookSharedSecret           string = "CONFIG_EVENTS_HOOK_SHARED_SECRET"
	ConfigHoptoadService                   string = "CONFIG_HOPTOAD_SERVICE"
	ConfigHoptoadAPIKey                    string = "CONFIG_HOPTOAD_API_KEY"
	SystemEventsHookSecretName             string = "backend-system-events-hook"
	InternalAPISecretName                  string = "backend-internal-api"
	ErrorMonitoringSecretName              string = "backend-error-monitoring"
)

var ListenerDefault map[string]string = map[string]string{
	ConfigRedisSentinelHosts:               "",
	ConfigRedisSentinelRole:                "",
	ConfigQueuesSentinelHosts:              "",
	ConfigQueuesSentinelRole:               "",
	ConfigListenerPrometheusMetricsEnabled: "true",
}

var WorkerDefault map[string]string = map[string]string{
	ConfigRedisSentinelHosts:             "",
	ConfigRedisSentinelRole:              "",
	ConfigQueuesSentinelHosts:            "",
	ConfigQueuesSentinelRole:             "",
	ConfigWorkerPrometheusMetricsEnabled: "true",
	ConfigWorkerPrometheusMetricsPort:    "9421",
}

var CronDefault map[string]string = map[string]string{
	ConfigRedisSentinelHosts:  "",
	ConfigRedisSentinelRole:   "",
	ConfigQueuesSentinelHosts: "",
	ConfigQueuesSentinelRole:  "",
}

var SecretDefinitions secrets.SecretConfigurations = secrets.SecretConfigurations{
	{
		SecretName: SystemEventsHookSecretName,
		ConfigOptions: map[string]string{
			ConfigEventsHook:             "/spec/config/systemEventsHookURL",
			ConfigEventsHookSharedSecret: "/spec/config/systemEventsHookPassword",
		},
	},
	{
		SecretName: InternalAPISecretName,
		ConfigOptions: map[string]string{
			ConfigInternalAPIUser:     "/spec/config/internalAPIUser",
			ConfigInternalAPIPassword: "/spec/config/internalAPIPassword",
		},
	},
	{
		SecretName: ErrorMonitoringSecretName,
		ConfigOptions: map[string]string{
			ConfigHoptoadService: "/spec/config/errorMonitoringService",
			ConfigHoptoadAPIKey:  "/spec/config/errorMonitoringKey",
		},
	},
}
