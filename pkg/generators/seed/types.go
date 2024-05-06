package seed

type SeedKey string

func (s SeedKey) String() string { return string(s) }

const (
	// Backend
	BackendInternalApiUser        SeedKey = "backend-internal-api-user"
	BackendInternalApiPassword    SeedKey = "backend-internal-api-password"
	BackendErrorMonitoringService SeedKey = "backend-error-monitoring-service"
	BackendErrorMonitoringApiKey  SeedKey = "backend-error-monitoring-api-key"

	// System
	SystemDatabaseDsn                  SeedKey = "system-database-dsn"
	SystemRecaptchaPublicKey           SeedKey = "system-recaptcha-public-key"
	SystemRecaptchaPrivateKey          SeedKey = "system-recaptcha-private-key"
	SystemEventsHookURL                SeedKey = "system-events-url"
	SystemEventsHookSharedSecret       SeedKey = "system-events-shared-secret"
	SystemSmtpUser                     SeedKey = "system-smpt-user"
	SystemSmtpPassword                 SeedKey = "system-smpt-password"
	SystemMasterAccessToken            SeedKey = "system-master-access-token"
	SystemApicastAccessToken           SeedKey = "system-apicast-access-token"
	SystemAssetsS3AwsAccessKey         SeedKey = "system-assets-s3-aws-access-key"
	SystemAssetsS3AwsSecretKey         SeedKey = "system-assets-s3-aws-secret-key"
	SystemSecretKeyBase                SeedKey = "system-secret-key-base"
	SystemAccessCode                   SeedKey = "system-access-code"
	SystemSegmentDeletionToken         SeedKey = "system-segment-deletion-token"
	SystemSegmentWriteKey              SeedKey = "system-segment-write-key"
	SystemGithubClientId               SeedKey = "system-github-client-id"
	SystemGithubClientSecret           SeedKey = "system-github-client-secret"
	SystemRHCustomerPortalClientId     SeedKey = "system-rh-customer-portal-client-id"
	SystemRHCustomerPortalClientSecret SeedKey = "system-rh-customer-portal-client-secret"
	SystemBugsnagApiKey                SeedKey = "system-bugsnag-api-key"
	SystemDatabaseSecret               SeedKey = "system-database-secret"

	// Zync
	ZyncDatabaseUrl   SeedKey = "zync-database-url"
	ZyncSecretKeyBase SeedKey = "zync-secret-key-base"
	ZyncAuthToken     SeedKey = "zync-auth-token"
	ZyncBugsnagApiKey SeedKey = "zync-bugsnag-api-key"
)
