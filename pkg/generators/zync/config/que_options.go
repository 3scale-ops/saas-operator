package config

import (
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
)

// NewQueOptions returns Zync Que options for the given saasv1alpha1.ZyncSpec
func NewQueOptions(spec saasv1alpha1.ZyncSpec) pod.Options {
	opts := pod.Options{}

	opts.Unpack(spec.Config.Rails.Environment).IntoEnvvar("RAILS_ENV")
	opts.Unpack(spec.Config.Rails.LogLevel).IntoEnvvar("RAILS_LOG_LEVEL")
	opts.Unpack("true").IntoEnvvar("RAILS_LOG_TO_STDOUT")
	opts.Unpack(spec.Config.DatabaseDSN).IntoEnvvar("DATABASE_URL").AsSecretRef("zync")
	opts.Unpack(spec.Config.SecretKeyBase).IntoEnvvar("SECRET_KEY_BASE").AsSecretRef("zync")
	opts.Unpack(spec.Config.ZyncAuthToken).IntoEnvvar("ZYNC_AUTHENTICATION_TOKEN").AsSecretRef("zync")
	opts.Unpack(spec.Config.Bugsnag.APIKey).IntoEnvvar("BUGSNAG_API_KEY").AsSecretRef("zync").EmptyIf(!spec.Config.Bugsnag.Enabled())
	opts.Unpack(spec.Config.Bugsnag.ReleaseStage).IntoEnvvar("BUGSNAG_RELEASE_STAGE")

	return opts
}
