package config

import (
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/resource_builders/pod"
)

// SphinxOptions holds configuration for the sphinx pods
type SphinxOptions struct {
	SphinxBindAddress       pod.EnvVarValue `env:"THINKING_SPHINX_ADDRESS"`
	SphinxPort              pod.EnvVarValue `env:"THINKING_SPHINX_PORT"`
	SphinxPIDFile           pod.EnvVarValue `env:"THINKING_SPHINX_PID_FILE"`
	SphinxConfigurationFile pod.EnvVarValue `env:"THINKING_SPHINX_CONFIGURATION_FILE"`
	SphinxBatchSize         pod.EnvVarValue `env:"THINKING_SPHINX_BATCH_SIZE"`

	RailsEnvironment pod.EnvVarValue `env:"RAILS_ENV"`
	SecretKeyBase    pod.EnvVarValue `env:"SECRET_KEY_BASE"`

	DatabaseURL pod.EnvVarValue `env:"DATABASE_URL" secret:"system-database"`

	RedisURL           pod.EnvVarValue `env:"REDIS_URL"`
	RedisNamespace     pod.EnvVarValue `env:"REDIS_NAMESPACE"`
	RedisSentinelHosts pod.EnvVarValue `env:"REDIS_SENTINEL_HOSTS"`
	RedisSentinelRole  pod.EnvVarValue `env:"REDIS_SENTINEL_ROLE"`
}

// NewSphinxOptions returns an Options struct for the given saasv1alpha1.SystemSpec
func NewSphinxOptions(spec saasv1alpha1.SystemSpec) SphinxOptions {
	opts := SphinxOptions{
		SphinxBindAddress:       &pod.ClearTextValue{Value: *spec.Sphinx.Config.Thinking.BindAddress},
		SphinxPort:              &pod.ClearTextValue{Value: fmt.Sprintf("%d", *spec.Sphinx.Config.Thinking.Port)},
		SphinxPIDFile:           &pod.ClearTextValue{Value: *spec.Sphinx.Config.Thinking.PIDFile},
		SphinxConfigurationFile: &pod.ClearTextValue{Value: *spec.Sphinx.Config.Thinking.ConfigFile},
		SphinxBatchSize:         &pod.ClearTextValue{Value: fmt.Sprintf("%d", *spec.Sphinx.Config.Thinking.BatchSize)},
		RailsEnvironment:        &pod.ClearTextValue{Value: *spec.Config.Rails.Environment},
		SecretKeyBase:           &pod.ClearTextValue{Value: "dummy"}, // https://github.com/rails/rails/issues/32947
		DatabaseURL:             &pod.SecretValue{Value: spec.Config.DatabaseDSN},
		RedisURL:                &pod.ClearTextValue{Value: spec.Config.Redis.QueuesDSN},
		RedisNamespace:          &pod.ClearTextValue{Value: ""},
		RedisSentinelHosts:      &pod.ClearTextValue{Value: ""},
		RedisSentinelRole:       &pod.ClearTextValue{Value: ""},
	}
	return opts
}
