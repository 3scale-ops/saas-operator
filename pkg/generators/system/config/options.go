package config

import (
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pod"
)

// Options holds configuration for system app and sidekiq pods
type Options struct {
	AMPRelease                    pod.EnvVarValue `env:"AMP_RELEASE"`
	ForceSSL                      pod.EnvVarValue `env:"FORCE_SSL"`
	ProviderPlan                  pod.EnvVarValue `env:"PROVIDER_PLAN"`
	SSLCertDir                    pod.EnvVarValue `env:"SSL_CERT_DIR"`
	SandboxProxyOpensslVeridyMode pod.EnvVarValue `env:"THREESCALE_SANDBOX_PROXY_OPENSSL_VERIFY_MODE"`
	Superdomain                   pod.EnvVarValue `env:"THREESCALE_SUPERDOMAIN"`

	RailsEnvironment pod.EnvVarValue `env:"RAILS_ENV"`
	RailsLogLevel    pod.EnvVarValue `env:"RAILS_LOG_LEVEL"`
	RailsLogToStdout pod.EnvVarValue `env:"RAILS_LOG_TO_STDOUT"`

	SphinxBindAddress pod.EnvVarValue `env:"THINKING_SPHINX_ADDRESS"`
	SphinxPort        pod.EnvVarValue `env:"THINKING_SPHINX_PORT"`

	SeedMasterAccessToken pod.EnvVarValue `env:"MASTER_ACCESS_TOKEN" secret:"system-seed"`
	SeedMasterDomain      pod.EnvVarValue `env:"MASTER_DOMAIN"`
	SeedMasterUser        pod.EnvVarValue `env:"MASTER_USER" secret:"system-seed"`
	SeedMasterPassword    pod.EnvVarValue `env:"MASTER_PASSWORD" secret:"system-seed"`
	SeedAdminAccessToken  pod.EnvVarValue `env:"ADMIN_ACCESS_TOKEN" secret:"system-seed"`
	SeedAdminUser         pod.EnvVarValue `env:"USER_LOGIN" secret:"system-seed"`
	SeedAdminPassword     pod.EnvVarValue `env:"USER_PASSWORD" secret:"system-seed"`
	SeedAdminEmail        pod.EnvVarValue `env:"USER_EMAIL"`
	SeedTenantName        pod.EnvVarValue `env:"TENANT_NAME"`

	DatabaseURL pod.EnvVarValue `env:"DATABASE_URL" secret:"system-database"`

	MemcachedServers pod.EnvVarValue `env:"MEMCACHE_SERVERS"`

	RecaptchaPublicKey  pod.EnvVarValue `env:"RECAPTCHA_PUBLIC_KEY" secret:"system-recaptcha"`
	RecaptchaPrivateKey pod.EnvVarValue `env:"RECAPTCHA_PRIVATE_KEY" secret:"system-recaptcha"`

	EventsHookPassword pod.EnvVarValue `env:"EVENTS_SHARED_SECRET" secret:"system-events-hook"`

	RedisURL                     pod.EnvVarValue `env:"REDIS_URL"`
	RedisMessageBusURL           pod.EnvVarValue `env:"MESSAGE_BUS_REDIS_URL"`
	RedisActionCableURL          pod.EnvVarValue `env:"ACTION_CABLE_REDIS_URL"`
	RedisNamespace               pod.EnvVarValue `env:"REDIS_NAMESPACE"`
	RedisMessageBusNamespace     pod.EnvVarValue `env:"MESSAGE_BUS_REDIS_NAMESPACE"`
	RedisSentinelHosts           pod.EnvVarValue `env:"REDIS_SENTINEL_HOSTS"`
	RedisSentinelRole            pod.EnvVarValue `env:"REDIS_SENTINEL_ROLE"`
	RedisMessageBusSentinelHosts pod.EnvVarValue `env:"MESSAGE_BUS_REDIS_SENTINEL_HOSTS"`
	RedisMessageBusSentinelRole  pod.EnvVarValue `env:"MESSAGE_BUS_REDIS_SENTINEL_ROLE"`

	SMTPAddress           pod.EnvVarValue `env:"SMTP_ADDRESS"`
	SMPTUserName          pod.EnvVarValue `env:"SMTP_USER_NAME" secret:"system-smtp"`
	SMTPPassword          pod.EnvVarValue `env:"SMTP_PASSWORD" secret:"system-smtp"`
	SMTPPort              pod.EnvVarValue `env:"SMTP_PORT"`
	SMPTAuthentication    pod.EnvVarValue `env:"SMTP_AUTHENTICATION"`
	SMTPOpensslVerifyMode pod.EnvVarValue `env:"SMTP_OPENSSL_VERIFY_MODE"`
	SMTPSTARTTLSAuto      pod.EnvVarValue `env:"SMTP_STARTTLS_AUTO"`

	ApicastAccessToken pod.EnvVarValue `env:"APICAST_ACCESS_TOKEN" secret:"system-master-apicast"`

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

	AppSecretKeyBase                 pod.EnvVarValue `env:"SECRET_KEY_BASE" secret:"system-app"`
	AccessCode                       pod.EnvVarValue `env:"ACCESS_CODE" secret:"system-app"`
	SegmentDeletionToken             pod.EnvVarValue `env:"SEGMENT_DELETION_TOKEN" secret:"system-app"`
	SegmentDeletionWorkspace         pod.EnvVarValue `env:"SEGMENT_DELETION_WORKSPACE"`
	SegmentWriteKey                  pod.EnvVarValue `env:"SEGMENT_WRITE_KEY" secret:"system-app"`
	GithubClientID                   pod.EnvVarValue `env:"GITHUB_CLIENT_ID" secret:"system-app"`
	GithubClientSecret               pod.EnvVarValue `env:"GITHUB_CLIENT_SECRET" secret:"system-app"`
	PrometheusUser                   pod.EnvVarValue `env:"PROMETHEUS_USER" secret:"system-app"`
	PrometheusPassword               pod.EnvVarValue `env:"PROMETHEUS_PASSWORD" secret:"system-app"`
	RedHatCustomerPortalClientID     pod.EnvVarValue `env:"RH_CUSTOMER_PORTAL_CLIENT_ID" secret:"system-app"`
	RedHatCustomerPortalClientSecret pod.EnvVarValue `env:"RH_CUSTOMER_PORTAL_CLIENT_SECRET" secret:"system-app"`
	BugsnagAPIKey                    pod.EnvVarValue `env:"BUGSNAG_API_KEY" secret:"system-app"`
	DatabaseSecret                   pod.EnvVarValue `env:"DB_SECRET" secret:"system-app"`
}

// NewOptions returns an Options struct for the given saasv1alpha1.SystemSpec
func NewOptions(spec saasv1alpha1.SystemSpec) Options {
	opts := Options{
		AMPRelease:                    &pod.ClearTextValue{Value: *spec.Config.AMPRelease},
		ForceSSL:                      &pod.ClearTextValue{Value: fmt.Sprintf("%t", *spec.Config.ForceSSL)},
		ProviderPlan:                  &pod.ClearTextValue{Value: *spec.Config.ThreescaleProviderPlan},
		SSLCertDir:                    &pod.ClearTextValue{Value: *spec.Config.SSLCertsDir},
		SandboxProxyOpensslVeridyMode: &pod.ClearTextValue{Value: *spec.Config.SandboxProxyOpensslVerifyMode},
		Superdomain:                   &pod.ClearTextValue{Value: *spec.Config.ThreescaleSuperdomain},

		RailsEnvironment: &pod.ClearTextValue{Value: *spec.Config.Rails.Environment},
		RailsLogLevel:    &pod.ClearTextValue{Value: *spec.Config.Rails.LogLevel},
		RailsLogToStdout: &pod.ClearTextValue{Value: "true"},

		SphinxBindAddress: &pod.ClearTextValue{Value: *spec.Sphinx.Config.Thinking.BindAddress},
		SphinxPort:        &pod.ClearTextValue{Value: fmt.Sprintf("%d", *spec.Sphinx.Config.Thinking.Port)},

		SeedMasterAccessToken: &pod.SecretValue{Value: spec.Config.Seed.MasterAccessToken},
		SeedMasterDomain:      &pod.ClearTextValue{Value: spec.Config.Seed.MasterDomain},
		SeedMasterUser:        &pod.SecretValue{Value: spec.Config.Seed.MasterUser},
		SeedMasterPassword:    &pod.SecretValue{Value: spec.Config.Seed.MasterPassword},
		SeedAdminAccessToken:  &pod.SecretValue{Value: spec.Config.Seed.AdminAccessToken},
		SeedAdminUser:         &pod.SecretValue{Value: spec.Config.Seed.AdminUser},
		SeedAdminPassword:     &pod.SecretValue{Value: spec.Config.Seed.AdminPassword},
		SeedAdminEmail:        &pod.ClearTextValue{Value: spec.Config.Seed.AdminEmail},
		SeedTenantName:        &pod.ClearTextValue{Value: spec.Config.Seed.TenantName},

		DatabaseURL: &pod.SecretValue{Value: spec.Config.DatabaseDSN},

		MemcachedServers: &pod.ClearTextValue{Value: spec.Config.MemcachedServers},

		RecaptchaPublicKey:  &pod.SecretValue{Value: spec.Config.Recaptcha.PublicKey},
		RecaptchaPrivateKey: &pod.SecretValue{Value: spec.Config.Recaptcha.PrivateKey},

		EventsHookPassword: &pod.SecretValue{Value: spec.Config.EventsSharedSecret},

		RedisURL:                     &pod.ClearTextValue{Value: spec.Config.Redis.QueuesDSN},
		RedisMessageBusURL:           &pod.ClearTextValue{Value: spec.Config.Redis.MessageBusDSN},
		RedisActionCableURL:          &pod.ClearTextValue{Value: ""},
		RedisNamespace:               &pod.ClearTextValue{Value: ""},
		RedisMessageBusNamespace:     &pod.ClearTextValue{Value: ""},
		RedisSentinelHosts:           &pod.ClearTextValue{Value: ""},
		RedisSentinelRole:            &pod.ClearTextValue{Value: ""},
		RedisMessageBusSentinelHosts: &pod.ClearTextValue{Value: ""},
		RedisMessageBusSentinelRole:  &pod.ClearTextValue{Value: ""},

		SMTPAddress:           &pod.ClearTextValue{Value: spec.Config.SMTP.Address},
		SMPTUserName:          &pod.SecretValue{Value: spec.Config.SMTP.User},
		SMTPPassword:          &pod.SecretValue{Value: spec.Config.SMTP.Password},
		SMTPPort:              &pod.ClearTextValue{Value: fmt.Sprintf("%d", spec.Config.SMTP.Port)},
		SMPTAuthentication:    &pod.ClearTextValue{Value: spec.Config.SMTP.AuthProtocol},
		SMTPOpensslVerifyMode: &pod.ClearTextValue{Value: spec.Config.SMTP.OpenSSLVerifyMode},
		SMTPSTARTTLSAuto:      &pod.ClearTextValue{Value: fmt.Sprintf("%t", spec.Config.SMTP.STARTTLSAuto)},

		ApicastAccessToken: &pod.SecretValue{Value: spec.Config.ApicastAccessToken},

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
		PrometheusUser:                   &pod.SecretValue{Value: spec.Config.Metrics.User},
		PrometheusPassword:               &pod.SecretValue{Value: spec.Config.Metrics.Password},
		RedHatCustomerPortalClientID:     &pod.SecretValue{Value: spec.Config.RedHatCustomerPortal.ClientID},
		RedHatCustomerPortalClientSecret: &pod.SecretValue{Value: spec.Config.RedHatCustomerPortal.ClientSecret},
		DatabaseSecret:                   &pod.SecretValue{Value: spec.Config.DatabaseSecret},
	}

	if spec.Config.Bugsnag.Enabled() {
		opts.BugsnagAPIKey = &pod.SecretValue{Value: spec.Config.Bugsnag.APIKey}
	} else {
		opts.BugsnagAPIKey = &pod.ClearTextValue{Value: ""}
	}

	return opts
}
