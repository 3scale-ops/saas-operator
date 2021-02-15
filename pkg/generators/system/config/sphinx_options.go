package config

import (
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pod"
)

// SphinxOptions holds configuration for the sphinx pods
type SphinxOptions struct {
	SphinxBindAddress       pod.EnvVarValue `env:"THINKING_SPHINX_ADDRESS"`
	SphinxPort              pod.EnvVarValue `env:"THINKING_SPHINX_PORT"`
	SphinxPIDFile           pod.EnvVarValue `env:"THINKING_SPHINX_PID_FILE"`
	SphinxConfigurationFile pod.EnvVarValue `env:"THINKING_SPHINX_CONFIGURATION_FILE"`
	DeltaIndexInterval      pod.EnvVarValue `env:"DELTA_INDEX_INTERVAL"`
	FUllReindexInterval     pod.EnvVarValue `env:"FULL_REINDEX_INTERVAL"`

	RailsEnvironment pod.EnvVarValue `env:"RAILS_ENV"`

	DatabaseURL pod.EnvVarValue `env:"DATABASE_URL" secret:"system-database"`

	RedisURL                     pod.EnvVarValue `env:"REDIS_URL"`
	RedisMessageBusURL           pod.EnvVarValue `env:"MESSAGE_BUS_REDIS_URL"`
	RedisNamespace               pod.EnvVarValue `env:"REDIS_NAMESPACE"`
	RedisMessageBusNamespace     pod.EnvVarValue `env:"MESSAGE_BUS_REDIS_NAMESPACE"`
	RedisSentinelHosts           pod.EnvVarValue `env:"REDIS_SENTINEL_HOSTS"`
	RedisSentinelRole            pod.EnvVarValue `env:"REDIS_SENTINEL_ROLE"`
	RedisMessageBusSentinelHosts pod.EnvVarValue `env:"MESSAGE_BUS_REDIS_SENTINEL_HOSTS"`
	RedisMessageBusSentinelRole  pod.EnvVarValue `env:"MESSAGE_BUS_REDIS_SENTINEL_ROLE"`
}

// NewSphinxOptions returns an Options struct for the given saasv1alpha1.SystemSpec
func NewSphinxOptions(spec saasv1alpha1.SystemSpec) SphinxOptions {
	opts := SphinxOptions{
		SphinxBindAddress:            &pod.ClearTextValue{Value: *spec.Sphinx.Config.Thinking.BindAddress},
		SphinxPort:                   &pod.ClearTextValue{Value: fmt.Sprintf("%d", *spec.Sphinx.Config.Thinking.Port)},
		SphinxPIDFile:                &pod.ClearTextValue{Value: *spec.Sphinx.Config.Thinking.PIDFile},
		SphinxConfigurationFile:      &pod.ClearTextValue{Value: *spec.Sphinx.Config.Thinking.ConfigFile},
		DeltaIndexInterval:           &pod.ClearTextValue{Value: fmt.Sprintf("%d", *spec.Sphinx.Config.DeltaIndexInterval)},
		FUllReindexInterval:          &pod.ClearTextValue{Value: fmt.Sprintf("%d", *spec.Sphinx.Config.FullReindexInterval)},
		RailsEnvironment:             &pod.ClearTextValue{Value: *spec.Config.Rails.Environment},
		DatabaseURL:                  &pod.SecretValue{Value: spec.Config.DatabaseDSN},
		RedisURL:                     &pod.ClearTextValue{Value: spec.Config.Redis.DSN},
		RedisMessageBusURL:           &pod.ClearTextValue{Value: spec.Config.Redis.MessageBusDSN},
		RedisNamespace:               &pod.ClearTextValue{Value: ""},
		RedisMessageBusNamespace:     &pod.ClearTextValue{Value: ""},
		RedisSentinelHosts:           &pod.ClearTextValue{Value: ""},
		RedisSentinelRole:            &pod.ClearTextValue{Value: ""},
		RedisMessageBusSentinelHosts: &pod.ClearTextValue{Value: ""},
		RedisMessageBusSentinelRole:  &pod.ClearTextValue{Value: ""},
	}
	return opts
}
