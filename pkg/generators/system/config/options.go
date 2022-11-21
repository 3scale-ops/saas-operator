package config

import (
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/resource_builders/pod"
	"k8s.io/utils/pointer"
)

// Options holds configuration for system app and sidekiq pods
type Options struct {
	ForceSSL                      pod.EnvVarValue `env:"FORCE_SSL"`
	ProviderPlan                  pod.EnvVarValue `env:"PROVIDER_PLAN"`
	SSLCertDir                    pod.EnvVarValue `env:"SSL_CERT_DIR"`
	SandboxProxyOpensslVerifyMode pod.EnvVarValue `env:"THREESCALE_SANDBOX_PROXY_OPENSSL_VERIFY_MODE"`
	Superdomain                   pod.EnvVarValue `env:"THREESCALE_SUPERDOMAIN"`

	RailsEnvironment pod.EnvVarValue `env:"RAILS_ENV"`
	RailsLogLevel    pod.EnvVarValue `env:"RAILS_LOG_LEVEL"`
	RailsLogToStdout pod.EnvVarValue `env:"RAILS_LOG_TO_STDOUT"`

	SphinxAddress pod.EnvVarValue `env:"THINKING_SPHINX_ADDRESS"`
	SphinxPort    pod.EnvVarValue `env:"THINKING_SPHINX_PORT"`

	DatabaseURL pod.EnvVarValue `env:"DATABASE_URL" secret:"system-database"`

	MemcachedServers pod.EnvVarValue `env:"MEMCACHE_SERVERS"`

	RecaptchaPublicKey  pod.EnvVarValue `env:"RECAPTCHA_PUBLIC_KEY" secret:"system-recaptcha"`
	RecaptchaPrivateKey pod.EnvVarValue `env:"RECAPTCHA_PRIVATE_KEY" secret:"system-recaptcha"`

	EventsHookPassword pod.EnvVarValue `env:"EVENTS_SHARED_SECRET" secret:"system-events-hook"`

	RedisURL           pod.EnvVarValue `env:"REDIS_URL"`
	RedisNamespace     pod.EnvVarValue `env:"REDIS_NAMESPACE"`
	RedisSentinelHosts pod.EnvVarValue `env:"REDIS_SENTINEL_HOSTS"`
	RedisSentinelRole  pod.EnvVarValue `env:"REDIS_SENTINEL_ROLE"`

	SMTPAddress           pod.EnvVarValue `env:"SMTP_ADDRESS"`
	SMTPUserName          pod.EnvVarValue `env:"SMTP_USER_NAME" secret:"system-smtp"`
	SMTPPassword          pod.EnvVarValue `env:"SMTP_PASSWORD" secret:"system-smtp"`
	SMTPPort              pod.EnvVarValue `env:"SMTP_PORT"`
	SMTPAuthentication    pod.EnvVarValue `env:"SMTP_AUTHENTICATION"`
	SMTPOpensslVerifyMode pod.EnvVarValue `env:"SMTP_OPENSSL_VERIFY_MODE"`
	SMTPSTARTTLSAuto      pod.EnvVarValue `env:"SMTP_STARTTLS_AUTO"`

	MappingServiceAccessToken pod.EnvVarValue `env:"APICAST_ACCESS_TOKEN" secret:"system-master-apicast"`

	ZyncAuthenticationToken pod.EnvVarValue `env:"ZYNC_AUTHENTICATION_TOKEN" secret:"system-zync"`

	BackendRedisURL            pod.EnvVarValue `env:"BACKEND_REDIS_URL"`
	BackendRedisSentinelHosts  pod.EnvVarValue `env:"BACKEND_REDIS_SENTINEL_HOSTS"`
	BackendRedisSentinelRole   pod.EnvVarValue `env:"BACKEND_REDIS_SENTINEL_ROLE"`
	ApicastBackendRootEndpoint pod.EnvVarValue `env:"APICAST_BACKEND_ROOT_ENDPOINT"`
	BackendRoute               pod.EnvVarValue `env:"BACKEND_ROUTE"`
	BackendPublicURL           pod.EnvVarValue `env:"BACKEND_PUBLIC_URL"`
	BackendInternalAPIUser     pod.EnvVarValue `env:"CONFIG_INTERNAL_API_USER" secret:"system-backend"`
	BackendInternalAPIPassword pod.EnvVarValue `env:"CONFIG_INTERNAL_API_PASSWORD" secret:"system-backend"`

	AssetsAWSAccessKeyID     pod.EnvVarValue `env:"AWS_ACCESS_KEY_ID" secret:"system-multitenant-assets-s3"`
	AssetsAWSSecretAccessKey pod.EnvVarValue `env:"AWS_SECRET_ACCESS_KEY" secret:"system-multitenant-assets-s3"`
	AssetsAWSBucket          pod.EnvVarValue `env:"AWS_BUCKET"`
	AssetsAWSRegion          pod.EnvVarValue `env:"AWS_REGION"`
	AssetsHost               pod.EnvVarValue `env:"RAILS_ASSET_HOST"`

	AppSecretKeyBase                 pod.EnvVarValue `env:"SECRET_KEY_BASE" secret:"system-app"`
	AccessCode                       pod.EnvVarValue `env:"ACCESS_CODE" secret:"system-app"`
	SegmentDeletionToken             pod.EnvVarValue `env:"SEGMENT_DELETION_TOKEN" secret:"system-app"`
	SegmentDeletionWorkspace         pod.EnvVarValue `env:"SEGMENT_DELETION_WORKSPACE"`
	SegmentWriteKey                  pod.EnvVarValue `env:"SEGMENT_WRITE_KEY" secret:"system-app"`
	GithubClientID                   pod.EnvVarValue `env:"GITHUB_CLIENT_ID" secret:"system-app"`
	GithubClientSecret               pod.EnvVarValue `env:"GITHUB_CLIENT_SECRET" secret:"system-app"`
	RedHatCustomerPortalClientID     pod.EnvVarValue `env:"RH_CUSTOMER_PORTAL_CLIENT_ID" secret:"system-app"`
	RedHatCustomerPortalClientSecret pod.EnvVarValue `env:"RH_CUSTOMER_PORTAL_CLIENT_SECRET" secret:"system-app"`
	BugsnagAPIKey                    pod.EnvVarValue `env:"BUGSNAG_API_KEY" secret:"system-app"`
	DatabaseSecret                   pod.EnvVarValue `env:"DB_SECRET" secret:"system-app"`
}

// NewOptions returns an Options struct for the given saasv1alpha1.SystemSpec
func NewOptions(spec saasv1alpha1.SystemSpec) Options {
	opts := Options{
		ForceSSL:                      &pod.ClearTextValue{Value: fmt.Sprintf("%t", *spec.Config.ForceSSL)},
		ProviderPlan:                  &pod.ClearTextValue{Value: *spec.Config.ThreescaleProviderPlan},
		SSLCertDir:                    &pod.ClearTextValue{Value: *spec.Config.SSLCertsDir},
		SandboxProxyOpensslVerifyMode: &pod.ClearTextValue{Value: *spec.Config.SandboxProxyOpensslVerifyMode},
		Superdomain:                   &pod.ClearTextValue{Value: *spec.Config.ThreescaleSuperdomain},

		RailsEnvironment: &pod.ClearTextValue{Value: *spec.Config.Rails.Environment},
		RailsLogLevel:    &pod.ClearTextValue{Value: *spec.Config.Rails.LogLevel},
		RailsLogToStdout: &pod.ClearTextValue{Value: "true"},

		SphinxAddress: &pod.ClearTextValue{Value: SystemSphinxServiceName},
		SphinxPort:    &pod.ClearTextValue{Value: fmt.Sprintf("%d", *spec.Sphinx.Config.Thinking.Port)},

		DatabaseURL: &pod.SecretValue{Value: spec.Config.DatabaseDSN},

		MemcachedServers: &pod.ClearTextValue{Value: spec.Config.MemcachedServers},

		RecaptchaPublicKey:  &pod.SecretValue{Value: spec.Config.Recaptcha.PublicKey},
		RecaptchaPrivateKey: &pod.SecretValue{Value: spec.Config.Recaptcha.PrivateKey},

		EventsHookPassword: &pod.SecretValue{Value: spec.Config.EventsSharedSecret},

		RedisURL:           &pod.ClearTextValue{Value: spec.Config.Redis.QueuesDSN},
		RedisNamespace:     &pod.ClearTextValue{Value: ""},
		RedisSentinelHosts: &pod.ClearTextValue{Value: ""},
		RedisSentinelRole:  &pod.ClearTextValue{Value: ""},

		SMTPAddress:           &pod.ClearTextValue{Value: spec.Config.SMTP.Address},
		SMTPUserName:          &pod.SecretValue{Value: spec.Config.SMTP.User},
		SMTPPassword:          &pod.SecretValue{Value: spec.Config.SMTP.Password},
		SMTPPort:              &pod.ClearTextValue{Value: fmt.Sprintf("%d", spec.Config.SMTP.Port)},
		SMTPAuthentication:    &pod.ClearTextValue{Value: spec.Config.SMTP.AuthProtocol},
		SMTPOpensslVerifyMode: &pod.ClearTextValue{Value: spec.Config.SMTP.OpenSSLVerifyMode},
		SMTPSTARTTLSAuto:      &pod.ClearTextValue{Value: fmt.Sprintf("%t", spec.Config.SMTP.STARTTLSAuto)},

		MappingServiceAccessToken: &pod.SecretValue{Value: spec.Config.MappingServiceAccessToken},

		ZyncAuthenticationToken: &pod.SecretValue{Value: spec.Config.ZyncAuthToken},

		BackendRedisURL:            &pod.ClearTextValue{Value: spec.Config.Backend.RedisDSN},
		BackendRedisSentinelHosts:  &pod.ClearTextValue{Value: ""},
		BackendRedisSentinelRole:   &pod.ClearTextValue{Value: ""},
		ApicastBackendRootEndpoint: &pod.ClearTextValue{Value: spec.Config.Backend.InternalEndpoint},
		BackendRoute:               &pod.ClearTextValue{Value: spec.Config.Backend.InternalEndpoint},
		BackendPublicURL:           &pod.ClearTextValue{Value: spec.Config.Backend.ExternalEndpoint},
		BackendInternalAPIUser:     &pod.SecretValue{Value: spec.Config.Backend.InternalAPIUser},
		BackendInternalAPIPassword: &pod.SecretValue{Value: spec.Config.Backend.InternalAPIPassword},

		AssetsAWSAccessKeyID:     &pod.SecretValue{Value: spec.Config.Assets.AccessKey},
		AssetsAWSSecretAccessKey: &pod.SecretValue{Value: spec.Config.Assets.SecretKey},
		AssetsAWSBucket:          &pod.ClearTextValue{Value: spec.Config.Assets.Bucket},
		AssetsAWSRegion:          &pod.ClearTextValue{Value: spec.Config.Assets.Region},

		AppSecretKeyBase:                 &pod.SecretValue{Value: spec.Config.SecretKeyBase},
		AccessCode:                       &pod.SecretValue{Value: spec.Config.AccessCode},
		SegmentDeletionToken:             &pod.SecretValue{Value: spec.Config.Segment.DeletionToken},
		SegmentDeletionWorkspace:         &pod.ClearTextValue{Value: spec.Config.Segment.DeletionWorkspace},
		SegmentWriteKey:                  &pod.SecretValue{Value: spec.Config.Segment.WriteKey},
		GithubClientID:                   &pod.SecretValue{Value: spec.Config.Github.ClientID},
		GithubClientSecret:               &pod.SecretValue{Value: spec.Config.Github.ClientSecret},
		RedHatCustomerPortalClientID:     &pod.SecretValue{Value: spec.Config.RedHatCustomerPortal.ClientID},
		RedHatCustomerPortalClientSecret: &pod.SecretValue{Value: spec.Config.RedHatCustomerPortal.ClientSecret},
		DatabaseSecret:                   &pod.SecretValue{Value: spec.Config.DatabaseSecret},
	}

	if spec.Config.Bugsnag.Enabled() {
		opts.BugsnagAPIKey = &pod.SecretValue{Value: spec.Config.Bugsnag.APIKey}
	} else {
		opts.BugsnagAPIKey = &pod.SecretValue{Value: saasv1alpha1.SecretReference{Override: pointer.StringPtr("")}}
	}

	if spec.Config.Assets.Host == nil {
		opts.AssetsHost = &pod.ClearTextValue{Value: ""}
	} else {
		opts.AssetsHost = &pod.ClearTextValue{Value: *spec.Config.Assets.Host}
	}

	return opts
}
