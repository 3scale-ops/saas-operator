package config

import (
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/generators/seed"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
)

// NewQueOptions returns Zync Que options for the given saasv1alpha1.ZyncSpec
func NewQueOptions(spec saasv1alpha1.ZyncSpec) pod.Options {
	opts := pod.Options{}

	opts.AddEnvvar("RAILS_ENV").Unpack(spec.Config.Rails.Environment)
	opts.AddEnvvar("RAILS_LOG_LEVEL").Unpack(spec.Config.Rails.LogLevel)
	opts.AddEnvvar("RAILS_LOG_TO_STDOUT").Unpack("true")
	opts.AddEnvvar("DATABASE_URL").AsSecretRef(ZyncSecret).WithSeedKey(seed.ZyncDatabaseUrl).
		Unpack(spec.Config.DatabaseDSN)
	opts.AddEnvvar("SECRET_KEY_BASE").AsSecretRef(ZyncSecret).WithSeedKey(seed.ZyncSecretKeyBase).
		Unpack(spec.Config.SecretKeyBase)
	opts.AddEnvvar("ZYNC_AUTHENTICATION_TOKEN").AsSecretRef(ZyncSecret).WithSeedKey(seed.ZyncAuthToken).
		Unpack(spec.Config.ZyncAuthToken)
	opts.AddEnvvar("BUGSNAG_API_KEY").AsSecretRef(ZyncSecret).EmptyIf(!spec.Config.Bugsnag.Enabled()).WithSeedKey(seed.ZyncBugsnagApiKey).
		Unpack(spec.Config.Bugsnag.APIKey)
	opts.AddEnvvar("BUGSNAG_RELEASE_STAGE").Unpack(spec.Config.Bugsnag.ReleaseStage)

	return opts
}
