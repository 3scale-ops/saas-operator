package config

import (
	"github.com/3scale-ops/basereconciler/util"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
)

// QueOptions holds configuration for the que pods
type QueOptions struct {
	RailsEnvironment        pod.EnvVarValue `env:"RAILS_ENV"`
	RailsLogLevel           pod.EnvVarValue `env:"RAILS_LOG_LEVEL"`
	RailsLogToStdOut        pod.EnvVarValue `env:"RAILS_LOG_TO_STDOUT"`
	DatabaseURL             pod.EnvVarValue `env:"DATABASE_URL" secret:"zync"`
	SecretKeyBase           pod.EnvVarValue `env:"SECRET_KEY_BASE" secret:"zync"`
	ZyncAuthenticationToken pod.EnvVarValue `env:"ZYNC_AUTHENTICATION_TOKEN" secret:"zync"`
	BugsnagAPIKey           pod.EnvVarValue `env:"BUGSNAG_API_KEY" secret:"zync"`
	BugsnagReleaseStage     pod.EnvVarValue `env:"BUGSNAG_RELEASE_STAGE"`
}

// NewQueOptions returns an Options struct for the given saasv1alpha1.ZyncSpec
func NewQueOptions(spec saasv1alpha1.ZyncSpec) QueOptions {
	opts := QueOptions{
		RailsEnvironment: &pod.ClearTextValue{Value: *spec.Config.Rails.Environment},
		RailsLogLevel:    &pod.ClearTextValue{Value: *spec.Config.Rails.LogLevel},
		RailsLogToStdOut: &pod.ClearTextValue{Value: "true"},

		DatabaseURL:             &pod.SecretValue{Value: spec.Config.DatabaseDSN},
		SecretKeyBase:           &pod.SecretValue{Value: spec.Config.SecretKeyBase},
		ZyncAuthenticationToken: &pod.SecretValue{Value: spec.Config.ZyncAuthToken},
	}

	if spec.Config.Bugsnag.Enabled() {
		opts.BugsnagAPIKey = &pod.SecretValue{Value: spec.Config.Bugsnag.APIKey}

		if spec.Config.Bugsnag.ReleaseStage != nil {
			opts.BugsnagReleaseStage = &pod.ClearTextValue{Value: *spec.Config.Bugsnag.ReleaseStage}
		}

	} else {
		opts.BugsnagAPIKey = &pod.SecretValue{Value: saasv1alpha1.SecretReference{Override: util.Pointer("")}}
	}

	return opts
}
