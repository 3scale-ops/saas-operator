package twemproxy

import (
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/resource_builders/pod"
)

const (
	TwemproxyConfigFile = "/etc/twemproxy/nutcracker.yml"
)

// TwemproxyOptions holds configuration for the Twemproxy sidecar
type TwemproxyOptions struct {
	ConfigFile     pod.EnvVarValue `env:"TWEMPROXY_CONFIG_FILE"`
	MetricsAddress pod.EnvVarValue `env:"TWEMPROXY_METRICS_ADDRESS"`
	StatsInterval  pod.EnvVarValue `env:"TWEMPROXY_STATS_INTERVAL"`
	LogLevel       pod.EnvVarValue `env:"TWEMPROXY_LOG_LEVEL"`
}

// NewTwemproxyOptions returns a NewTwemproxyOptions struct for the given saasv1alpha1.BackendSpec
func NewTwemproxyOptions(spec saasv1alpha1.TwemproxySpec) TwemproxyOptions {
	opts := TwemproxyOptions{
		ConfigFile:     &pod.ClearTextValue{Value: TwemproxyConfigFile},
		MetricsAddress: &pod.ClearTextValue{Value: fmt.Sprintf(":%d", *spec.Options.MetricsPort)},
		StatsInterval:  &pod.ClearTextValue{Value: fmt.Sprintf("%d", spec.Options.StatsInterval.Milliseconds())},
		LogLevel:       &pod.ClearTextValue{Value: fmt.Sprintf("%d", *spec.Options.LogLevel)},
	}
	return opts
}
