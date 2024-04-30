package config

import (
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/generators/seed"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
)

type Secret string

func (s Secret) String() string { return string(s) }

const (
	SystemDatabaseSecret            Secret = "system-database"
	SystemRecaptchaSecret           Secret = "system-recaptcha"
	SystemEventsHookSecret          Secret = "system-events-hook"
	SystemSmptSecret                Secret = "system-smtp"
	SystemMasterApicastSecret       Secret = "system-master-apicast"
	SystemZyncSecret                Secret = "system-zync"
	SystemBackendSecret             Secret = "system-backend"
	SystemMultitenantAssetsS3Secret Secret = "system-multitenant-assets-s3"
	SystemAppSecret                 Secret = "system-app"
)

func NewOptions(spec saasv1alpha1.SystemSpec) pod.Options {
	opts := pod.Options{}

	opts.AddEnvvar("FORCE_SSL").Unpack(spec.Config.ForceSSL)
	opts.AddEnvvar("PROVIDER_PLAN").Unpack(spec.Config.ThreescaleProviderPlan)
	opts.AddEnvvar("SSL_CERT_DIR").Unpack(spec.Config.SSLCertsDir)
	opts.AddEnvvar("THREESCALE_SANDBOX_PROXY_OPENSSL_VERIFY_MODE").Unpack(spec.Config.SandboxProxyOpensslVerifyMode)
	opts.AddEnvvar("THREESCALE_SUPERDOMAIN").Unpack(spec.Config.ThreescaleSuperdomain)

	opts.AddEnvvar("RAILS_ENV").Unpack(spec.Config.Rails.Environment)
	opts.AddEnvvar("RAILS_LOG_LEVEL").Unpack(spec.Config.Rails.LogLevel)
	opts.AddEnvvar("RAILS_LOG_TO_STDOUT").Unpack("true")

	opts.AddEnvvar("THINKING_SPHINX_ADDRESS").Unpack(spec.Config.SearchServer.Host)
	opts.AddEnvvar("THINKING_SPHINX_PORT").Unpack(spec.Config.SearchServer.Port)
	opts.AddEnvvar("THINKING_SPHINX_BATCH_SIZE").Unpack(spec.Config.SearchServer.BatchSize)

	opts.AddEnvvar("DATABASE_URL").AsSecretRef(SystemDatabaseSecret).WithSeedKey(seed.SystemDatabaseDsn).
		Unpack(spec.Config.DatabaseDSN)

	opts.AddEnvvar("MEMCACHE_SERVERS").Unpack(spec.Config.MemcachedServers)

	opts.AddEnvvar("RECAPTCHA_PUBLIC_KEY").AsSecretRef(SystemRecaptchaSecret).WithSeedKey(seed.SystemRecaptchaPublicKey).
		Unpack(spec.Config.Recaptcha.PublicKey)
	opts.AddEnvvar("RECAPTCHA_PRIVATE_KEY").AsSecretRef(SystemRecaptchaSecret).WithSeedKey(seed.SystemRecaptchaPrivateKey).
		Unpack(spec.Config.Recaptcha.PrivateKey)

	opts.AddEnvvar("EVENTS_SHARED_SECRET").AsSecretRef(SystemEventsHookSecret).WithSeedKey(seed.SystemEventsHookSharedSecret).
		Unpack(spec.Config.EventsSharedSecret)

	opts.AddEnvvar("REDIS_URL").Unpack(spec.Config.Redis.QueuesDSN)
	opts.AddEnvvar("REDIS_NAMESPACE").Unpack("")
	opts.AddEnvvar("REDIS_SENTINEL_HOSTS").Unpack("")
	opts.AddEnvvar("REDIS_SENTINEL_ROLE").Unpack("")

	opts.AddEnvvar("SMTP_ADDRESS").Unpack(spec.Config.SMTP.Address)
	opts.AddEnvvar("SMTP_USER_NAME").AsSecretRef(SystemSmptSecret).WithSeedKey(seed.SystemSmtpUser).
		Unpack(spec.Config.SMTP.User)
	opts.AddEnvvar("SMTP_PASSWORD").AsSecretRef(SystemSmptSecret).WithSeedKey(seed.SystemSmtpPassword).
		Unpack(spec.Config.SMTP.Password)
	opts.AddEnvvar("SMTP_PORT").Unpack(spec.Config.SMTP.Port)
	opts.AddEnvvar("SMTP_AUTHENTICATION").Unpack(spec.Config.SMTP.AuthProtocol)
	opts.AddEnvvar("SMTP_OPENSSL_VERIFY_MODE").Unpack(spec.Config.SMTP.OpenSSLVerifyMode)
	opts.AddEnvvar("SMTP_STARTTLS").Unpack(spec.Config.SMTP.STARTTLS)
	opts.AddEnvvar("SMTP_STARTTLS_AUTO").Unpack(spec.Config.SMTP.STARTTLSAuto)

	opts.AddEnvvar("APICAST_ACCESS_TOKEN").AsSecretRef(SystemMasterApicastSecret).WithSeedKey(seed.SystemApicastAccessToken).
		Unpack(spec.Config.MappingServiceAccessToken)

	opts.AddEnvvar("ZYNC_ENDPOINT").Unpack(spec.Config.Zync.Endpoint)
	opts.AddEnvvar("ZYNC_AUTHENTICATION_TOKEN").AsSecretRef(SystemZyncSecret).WithSeedKey(seed.ZyncAuthToken).
		Unpack(spec.Config.Zync.AuthToken)

	opts.AddEnvvar("BACKEND_REDIS_URL").Unpack(spec.Config.Backend.RedisDSN)
	opts.AddEnvvar("BACKEND_REDIS_SENTINEL_HOSTS").Unpack("")
	opts.AddEnvvar("BACKEND_REDIS_SENTINEL_ROLE").Unpack("")
	opts.AddEnvvar("BACKEND_URL").Unpack(spec.Config.Backend.InternalEndpoint)
	opts.AddEnvvar("BACKEND_PUBLIC_URL").Unpack(spec.Config.Backend.ExternalEndpoint)
	opts.AddEnvvar("CONFIG_INTERNAL_API_USER").AsSecretRef(SystemBackendSecret).WithSeedKey(seed.BackendInternalApiUser).
		Unpack(spec.Config.Backend.InternalAPIUser)
	opts.AddEnvvar("CONFIG_INTERNAL_API_PASSWORD").AsSecretRef(SystemBackendSecret).WithSeedKey(seed.BackendInternalApiPassword).
		Unpack(spec.Config.Backend.InternalAPIPassword)

	opts.AddEnvvar("AWS_ACCESS_KEY_ID").AsSecretRef(SystemMultitenantAssetsS3Secret).WithSeedKey(seed.SystemAssetsS3AwsAccessKey).
		Unpack(spec.Config.Assets.AccessKey)
	opts.AddEnvvar("AWS_SECRET_ACCESS_KEY").AsSecretRef(SystemMultitenantAssetsS3Secret).WithSeedKey(seed.SystemAssetsS3AwsSecretKey).
		Unpack(spec.Config.Assets.SecretKey)
	opts.AddEnvvar("AWS_BUCKET").Unpack(spec.Config.Assets.Bucket)
	opts.AddEnvvar("AWS_REGION").Unpack(spec.Config.Assets.Region)
	opts.AddEnvvar("AWS_S3_HOSTNAME").Unpack(spec.Config.Assets.S3Endpoint)
	opts.AddEnvvar("RAILS_ASSET_HOST").Unpack(spec.Config.Assets.Host)

	opts.AddEnvvar("SECRET_KEY_BASE").AsSecretRef(SystemAppSecret).WithSeedKey(seed.SystemSecretKeyBase).
		Unpack(spec.Config.SecretKeyBase)
	opts.AddEnvvar("ACCESS_CODE").AsSecretRef(SystemAppSecret).WithSeedKey(seed.SystemAccessCode).
		Unpack(spec.Config.AccessCode)
	opts.AddEnvvar("SEGMENT_DELETION_TOKEN").AsSecretRef(SystemAppSecret).WithSeedKey(seed.SystemSegmentDeletionToken).
		Unpack(spec.Config.Segment.DeletionToken)
	opts.AddEnvvar("SEGMENT_DELETION_WORKSPACE").Unpack(spec.Config.Segment.DeletionWorkspace)
	opts.AddEnvvar("SEGMENT_WRITE_KEY").AsSecretRef(SystemAppSecret).WithSeedKey(seed.SystemSegmentWriteKey).
		Unpack(spec.Config.Segment.WriteKey)
	opts.AddEnvvar("GITHUB_CLIENT_ID").AsSecretRef(SystemAppSecret).WithSeedKey(seed.SystemGithubClientId).
		Unpack(spec.Config.Github.ClientID)
	opts.AddEnvvar("GITHUB_CLIENT_SECRET").AsSecretRef(SystemAppSecret).WithSeedKey(seed.SystemGithubClientSecret).
		Unpack(spec.Config.Github.ClientSecret)
	opts.AddEnvvar("RH_CUSTOMER_PORTAL_CLIENT_ID").AsSecretRef(SystemAppSecret).WithSeedKey(seed.SystemRHCustomerPortalClientId).
		Unpack(spec.Config.RedHatCustomerPortal.ClientID)
	opts.AddEnvvar("RH_CUSTOMER_PORTAL_CLIENT_SECRET").AsSecretRef(SystemAppSecret).WithSeedKey(seed.SystemRHCustomerPortalClientSecret).
		Unpack(spec.Config.RedHatCustomerPortal.ClientSecret)
	opts.AddEnvvar("RH_CUSTOMER_PORTAL_REALM").Unpack(spec.Config.RedHatCustomerPortal.Realm)
	opts.AddEnvvar("BUGSNAG_API_KEY").AsSecretRef(SystemAppSecret).WithSeedKey(seed.SystemBugsnagApiKey).EmptyIf(!spec.Config.Bugsnag.Enabled()).
		Unpack(spec.Config.Bugsnag.APIKey)
	opts.AddEnvvar("BUGSNAG_RELEASE_STAGE").Unpack(spec.Config.Bugsnag.ReleaseStage)
	opts.AddEnvvar("DB_SECRET").AsSecretRef(SystemAppSecret).WithSeedKey(seed.SystemDatabaseSecret).
		Unpack(spec.Config.DatabaseSecret)

	if spec.Config.Apicast != nil {
		opts.AddEnvvar("APICAST_STAGING_DOMAIN").Unpack(spec.Config.Apicast.StagingDomain)
		opts.AddEnvvar("APICAST_PRODUCTION_DOMAIN").Unpack(spec.Config.Apicast.ProductionDomain)
		opts.AddEnvvar("APICAST_CLOUD_HOSTED_REGISTRY_URL").Unpack(spec.Config.Apicast.CloudHostedRegistryURL)
		opts.AddEnvvar("APICAST_SELF_MANAGED_REGISTRY_URL").Unpack(spec.Config.Apicast.SelfManagedRegistryURL)
	}

	return opts
}
