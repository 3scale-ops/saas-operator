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

	opts.AddEnvvar("TWEMPROXY_CONFIG_FILE").Unpack(TwemproxyConfigFile)
	opts.AddEnvvar("TWEMPROXY_METRICS_ADDRESS").Unpack(spec.Options.MetricsPort, ":%d")
	opts.AddEnvvar("TWEMPROXY_STATS_INTERVAL").Unpack(spec.Options.StatsInterval.Milliseconds())
	opts.AddEnvvar("TWEMPROXY_LOG_LEVEL").Unpack(spec.Options.LogLevel)

	return opts
}
