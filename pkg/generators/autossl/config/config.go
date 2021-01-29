package config

const (
	ACMEStaging          string = "ACME_STAGING"
	ContactEmail         string = "CONTACT_EMAIL"
	ProxyEndpoint        string = "PROXY_ENDPOINT"
	StorageAdapter       string = "STORAGE_ADAPTER"
	RedisHost            string = "REDIS_HOST"
	RedisPort            string = "REDIS_PORT"
	VerificationEndpoint string = "VERIFICATION_ENDPOINT"
	LogLevel             string = "LOG_LEVEL"
	DomainWhitelist      string = "DOMAIN_WHITELIST"
	DomainBlacklist      string = "DOMAIN_BLACKLIST"
)

var Default map[string]string = map[string]string{
	StorageAdapter: "redis",
	ACMEStaging:    "",
}
