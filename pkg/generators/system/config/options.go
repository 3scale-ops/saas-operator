package config

import (
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
)

func NewOptions(spec saasv1alpha1.SystemSpec) pod.Options {
	opts := pod.Options{}

	opts.Unpack(spec.Config.ForceSSL).IntoEnvvar("FORCE_SSL")
	opts.Unpack(spec.Config.ThreescaleProviderPlan).IntoEnvvar("PROVIDER_PLAN")
	opts.Unpack(spec.Config.SSLCertsDir).IntoEnvvar("SSL_CERT_DIR")
	opts.Unpack(spec.Config.SandboxProxyOpensslVerifyMode).IntoEnvvar("THREESCALE_SANDBOX_PROXY_OPENSSL_VERIFY_MODE")
	opts.Unpack(spec.Config.ThreescaleSuperdomain).IntoEnvvar("THREESCALE_SUPERDOMAIN")

	opts.Unpack(spec.Config.Rails.Environment).IntoEnvvar("RAILS_ENV")
	opts.Unpack(spec.Config.Rails.LogLevel).IntoEnvvar("RAILS_LOG_LEVEL")
	opts.Unpack("true").IntoEnvvar("RAILS_LOG_TO_STDOUT")

	opts.Unpack(spec.Config.SearchServer.Host).IntoEnvvar("THINKING_SPHINX_ADDRESS")
	opts.Unpack(spec.Config.SearchServer.Port).IntoEnvvar("THINKING_SPHINX_PORT")
	opts.Unpack(spec.Config.SearchServer.BatchSize).IntoEnvvar("THINKING_SPHINX_BATCH_SIZE")

	opts.Unpack(spec.Config.DatabaseDSN).IntoEnvvar("DATABASE_URL").AsSecretRef("system-database")

	opts.Unpack(spec.Config.MemcachedServers).IntoEnvvar("MEMCACHE_SERVERS")

	opts.Unpack(spec.Config.Recaptcha.PublicKey).IntoEnvvar("RECAPTCHA_PUBLIC_KEY").AsSecretRef("system-recaptcha")
	opts.Unpack(spec.Config.Recaptcha.PrivateKey).IntoEnvvar("RECAPTCHA_PRIVATE_KEY").AsSecretRef("system-recaptcha")

	opts.Unpack(spec.Config.EventsSharedSecret).IntoEnvvar("EVENTS_SHARED_SECRET").AsSecretRef("system-events-hook")

	opts.Unpack(spec.Config.Redis.QueuesDSN).IntoEnvvar("REDIS_URL")
	opts.Unpack("").IntoEnvvar("REDIS_NAMESPACE")
	opts.Unpack("").IntoEnvvar("REDIS_SENTINEL_HOSTS")
	opts.Unpack("").IntoEnvvar("REDIS_SENTINEL_ROLE")

	opts.Unpack(spec.Config.SMTP.Address).IntoEnvvar("SMTP_ADDRESS")
	opts.Unpack(spec.Config.SMTP.User).IntoEnvvar("SMTP_USER_NAME").AsSecretRef("system-smtp")
	opts.Unpack(spec.Config.SMTP.Password).IntoEnvvar("SMTP_PASSWORD").AsSecretRef("system-smtp")
	opts.Unpack(spec.Config.SMTP.Port).IntoEnvvar("SMTP_PORT")
	opts.Unpack(spec.Config.SMTP.AuthProtocol).IntoEnvvar("SMTP_AUTHENTICATION")
	opts.Unpack(spec.Config.SMTP.OpenSSLVerifyMode).IntoEnvvar("SMTP_OPENSSL_VERIFY_MODE")
	opts.Unpack(spec.Config.SMTP.STARTTLS).IntoEnvvar("SMTP_STARTTLS")
	opts.Unpack(spec.Config.SMTP.STARTTLSAuto).IntoEnvvar("SMTP_STARTTLS_AUTO")

	opts.Unpack(spec.Config.MappingServiceAccessToken).IntoEnvvar("APICAST_ACCESS_TOKEN").AsSecretRef("system-master-apicast")

	opts.Unpack(spec.Config.Zync.Endpoint).IntoEnvvar("ZYNC_ENDPOINT")
	opts.Unpack(spec.Config.Zync.AuthToken).IntoEnvvar("ZYNC_AUTHENTICATION_TOKEN").AsSecretRef("system-zync")

	opts.Unpack(spec.Config.Backend.RedisDSN).IntoEnvvar("BACKEND_REDIS_URL")
	opts.Unpack("").IntoEnvvar("BACKEND_REDIS_SENTINEL_HOSTS")
	opts.Unpack("").IntoEnvvar("BACKEND_REDIS_SENTINEL_ROLE")
	opts.Unpack(spec.Config.Backend.InternalEndpoint).IntoEnvvar("BACKEND_URL")
	opts.Unpack(spec.Config.Backend.ExternalEndpoint).IntoEnvvar("BACKEND_PUBLIC_URL")
	opts.Unpack(spec.Config.Backend.InternalAPIUser).IntoEnvvar("CONFIG_INTERNAL_API_USER").AsSecretRef("system-backend")
	opts.Unpack(spec.Config.Backend.InternalAPIPassword).IntoEnvvar("CONFIG_INTERNAL_API_PASSWORD").AsSecretRef("system-backend")

	opts.Unpack(spec.Config.Assets.AccessKey).IntoEnvvar("AWS_ACCESS_KEY_ID").AsSecretRef("system-multitenant-assets-s3")
	opts.Unpack(spec.Config.Assets.SecretKey).IntoEnvvar("AWS_SECRET_ACCESS_KEY").AsSecretRef("system-multitenant-assets-s3")
	opts.Unpack(spec.Config.Assets.Bucket).IntoEnvvar("AWS_BUCKET")
	opts.Unpack(spec.Config.Assets.Region).IntoEnvvar("AWS_REGION")
	opts.Unpack(spec.Config.Assets.Host).IntoEnvvar("RAILS_ASSET_HOST")

	opts.Unpack(spec.Config.SecretKeyBase).IntoEnvvar("SECRET_KEY_BASE").AsSecretRef("system-app")
	opts.Unpack(spec.Config.AccessCode).IntoEnvvar("ACCESS_CODE").AsSecretRef("system-app")
	opts.Unpack(spec.Config.Segment.DeletionToken).IntoEnvvar("SEGMENT_DELETION_TOKEN").AsSecretRef("system-app")
	opts.Unpack(spec.Config.Segment.DeletionWorkspace).IntoEnvvar("SEGMENT_DELETION_WORKSPACE")
	opts.Unpack(spec.Config.Segment.WriteKey).IntoEnvvar("SEGMENT_WRITE_KEY").AsSecretRef("system-app")
	opts.Unpack(spec.Config.Github.ClientID).IntoEnvvar("GITHUB_CLIENT_ID").AsSecretRef("system-app")
	opts.Unpack(spec.Config.Github.ClientSecret).IntoEnvvar("GITHUB_CLIENT_SECRET").AsSecretRef("system-app")
	opts.Unpack(spec.Config.RedHatCustomerPortal.ClientID).IntoEnvvar("RH_CUSTOMER_PORTAL_CLIENT_ID").AsSecretRef("system-app")
	opts.Unpack(spec.Config.RedHatCustomerPortal.ClientSecret).IntoEnvvar("RH_CUSTOMER_PORTAL_CLIENT_SECRET").AsSecretRef("system-app")
	opts.Unpack(spec.Config.RedHatCustomerPortal.Realm).IntoEnvvar("RH_CUSTOMER_PORTAL_REALM")
	opts.Unpack(spec.Config.Bugsnag.APIKey).IntoEnvvar("BUGSNAG_API_KEY").AsSecretRef("system-app").EmptyIf(!spec.Config.Bugsnag.Enabled())
	opts.Unpack(spec.Config.Bugsnag.ReleaseStage).IntoEnvvar("BUGSNAG_RELEASE_STAGE")
	opts.Unpack(spec.Config.DatabaseSecret).IntoEnvvar("DB_SECRET").AsSecretRef("system-app")

	return opts
}
