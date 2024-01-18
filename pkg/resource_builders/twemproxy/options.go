package twemproxy

import (
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
)

const (
	TwemproxyConfigFile = "/etc/twemproxy/nutcracker.yml"
)

// NewOptions generates the configuration for the Twemproxy sidecar
func NewOptions(spec saasv1alpha1.TwemproxySpec) *pod.Options {
	opts := pod.NewOptions()

	opts.Unpack(TwemproxyConfigFile).IntoEnvvar("TWEMPROXY_CONFIG_FILE")
	opts.Unpack(spec.Options.MetricsPort, ":%d").IntoEnvvar("TWEMPROXY_METRICS_ADDRESS")
	opts.Unpack(spec.Options.StatsInterval.Milliseconds()).IntoEnvvar("TWEMPROXY_STATS_INTERVAL")
	opts.Unpack(spec.Options.LogLevel).IntoEnvvar("TWEMPROXY_LOG_LEVEL")

	return opts
}
