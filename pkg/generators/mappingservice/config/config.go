package config

import "github.com/3scale/saas-operator/pkg/generators/common_blocks/secrets"

const (
	MasterAccessToken          string = "MASTER_ACCESS_TOKEN"
	APIHost                    string = "API_HOST"
	ApicastConfigurationLoader string = "APICAST_CONFIGURATION_LOADER"
	ApicastLogLevel            string = "APICAST_LOG_LEVEL"
	PreviewBaseDomain          string = "PREVIEW_BASE_DOMAIN"
)

var Default map[string]string = map[string]string{
	ApicastConfigurationLoader: "lazy",
}

var SecretDefinitions secrets.SecretConfigurations = secrets.SecretConfigurations{
	{
		SecretName: "mapping-service-system-master-access-token",
		ConfigOptions: map[string]string{
			MasterAccessToken: "/spec/config/systemAdminToken",
		},
	},
}
