package config

import "github.com/3scale/saas-operator/pkg/generators/common_blocks/secrets"

const (
	DatabaseURL string = "DATABASE_URL"
)

var Default map[string]string = map[string]string{}

var SecretDefinitions secrets.SecretConfigurations = secrets.SecretConfigurations{
	{
		SecretName: "cors-proxy-system-database",
		ConfigOptions: map[string]string{
			DatabaseURL: "/spec/config/systemDatabaseDSN",
		},
	},
}
