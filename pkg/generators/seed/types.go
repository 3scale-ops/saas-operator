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
	SystemEventsHookURL                SeedKey = "system-events-url" // this shouldn't be a secret
	SystemEventsHookSharedSecret       SeedKey = "system-events-shared-secret"
	SystemSmtpUser                     SeedKey = "system-smpt-user"
	SystemSmtpPassword                 SeedKey = "system-smpt-password"
	SystemMasterAccessToken            SeedKey = "system-master-access-token"
	SystemAssetsS3AwsAccessKey         SeedKey = "system-assets-s3-aws-access-key"
	SystemAssetsS3AwsSecretKey         SeedKey = "system-assets-s3-aws-secret-key"
	SystemSecretKeyBase                SeedKey = "system-secret-key-base"
	SystemAccessCode                   SeedKey = "system-access-code"
	SystemSegmentDeletionToken         SeedKey = "system-segment-deletion-token"
	SystemSegmentDeletionWorkspace     SeedKey = "system-segment-deletion-workspace"
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

// TODO: use this to generate a Secret from some input params
// var AutoGen map[SeedKey]string = map[SeedKey]string{
// 	BackendInternalApiUser:       "user",
// 	BackendInternalApiPassword:   "<generate#backend-internal-api-password>",
// 	SystemDatabaseDsn:            "mysql2://app:<generate#system-db-password>@<param#system-db-host>:3306/system_enterprise",
// 	SystemEventsHookURL:          "https://<param#system-app-host>/master/events/import",
// 	SystemEventsHookSharedSecret: "<generate#system-events-shared-secret>",
// 	SystemMasterAccessToken:      "<generate#system-master-token>",
// 	SystemAssetsS3AwsAccessKey:   "<param#s3-aws-access-key>",
// 	SystemAssetsS3AwsSecretKey:   "<param#s3-aws-secret-key>",
// 	SystemSecretKeyBase:          "<generate#system-key-base>",
// 	SystemAccessCode:             "<generate#system-access-code>",
// 	SystemDatabaseSecret:         "<generate#system-db-secret>",
// 	ZyncDatabaseUrl:              "postgresql://app:<generate#zync-db-password>@<param#zync-db-host>:5432/zync",
// 	ZyncSecretKeyBase:            "<generate#zync-secret-key-base>",
// 	ZyncAuthToken:                "<generate#zync-auth-token>",
// }
