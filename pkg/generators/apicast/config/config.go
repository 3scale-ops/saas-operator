package config

const (
	ApicastConfigurationLoader string = "APICAST_CONFIGURATION_LOADER"
	ApicastConfigurationCache  string = "APICAST_CONFIGURATION_CACHE"
	ApicastExtendedMetrics     string = "APICAST_EXTENDED_METRICS"
	ThreeScaleDeploymentEnv    string = "THREESCALE_DEPLOYMENT_ENV"
	ThreescalePortalEndpoint   string = "THREESCALE_PORTAL_ENDPOINT"
	ApicastLogLevel            string = "APICAST_LOG_LEVEL"
	ApicastOIDCLogLevel        string = "APICAST_OIDC_LOG_LEVEL"
	ApicastResponseCodes       string = "APICAST_RESPONSE_CODES"
)

var Default map[string]string = map[string]string{
	ApicastConfigurationLoader: "lazy",
	ApicastExtendedMetrics:     "true",
	ApicastResponseCodes:       "true",
}
