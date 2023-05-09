package v1alpha1

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

var (
	// Twemproxy defaults
	defaultTwemproxyLivenessProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: pointer.Int32(0),
		TimeoutSeconds:      pointer.Int32(1),
		PeriodSeconds:       pointer.Int32(5),
		SuccessThreshold:    pointer.Int32(1),
		FailureThreshold:    pointer.Int32(3),
	}
	defaultTwemproxyReadinessProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: pointer.Int32(0),
		TimeoutSeconds:      pointer.Int32(1),
		PeriodSeconds:       pointer.Int32(5),
		SuccessThreshold:    pointer.Int32(1),
		FailureThreshold:    pointer.Int32(3),
	}
	// TODO: add requirements
	defaultTwemproxyResources defaultResourceRequirementsSpec = defaultResourceRequirementsSpec{}
	defaultTwemproxyImage     defaultImageSpec                = defaultImageSpec{
		Name:       pointer.String("quay.io/3scale/twemproxy"),
		Tag:        pointer.String("v0.5.0"),
		PullPolicy: (*corev1.PullPolicy)(pointer.String(string(corev1.PullIfNotPresent))),
	}
	twemproxyDefaultLogLevel      int32           = 6
	twemproxyDefaultMetricsPort   int32           = 9151
	twemproxyDefaultStatsInterval metav1.Duration = metav1.Duration{Duration: 10 * time.Second}
)

// TwemproxySpec configures twemproxy sidecars
// to access a sharded redis
type TwemproxySpec struct {
	// Image specification for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Image *ImageSpec `json:"image,omitempty"`
	// Resource requirements for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Resources *ResourceRequirementsSpec `json:"resources,omitempty"`
	// Liveness probe for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	LivenessProbe *ProbeSpec `json:"livenessProbe,omitempty"`
	// Readiness probe for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ReadinessProbe *ProbeSpec `json:"readinessProbe,omitempty"`
	// TwemproxyConfigRef is a reference to a TwemproxyConfig
	// resource in the same Namespace
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	TwemproxyConfigRef string `json:"twemproxyConfigRef"`
	// Options
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Options *TwemproxyOptions `json:"options,omitempty"`
}

func (spec *TwemproxySpec) ConfigMapName() string {
	return spec.TwemproxyConfigRef
}

// Default implements defaulting for the each backend cron
func (spec *TwemproxySpec) Default() {

	spec.Image = InitializeImageSpec(spec.Image, defaultTwemproxyImage)
	spec.Resources = InitializeResourceRequirementsSpec(spec.Resources, defaultTwemproxyResources)
	spec.LivenessProbe = InitializeProbeSpec(spec.LivenessProbe, defaultTwemproxyLivenessProbe)
	spec.ReadinessProbe = InitializeProbeSpec(spec.ReadinessProbe, defaultTwemproxyReadinessProbe)
	if spec.Options == nil {
		spec.Options = &TwemproxyOptions{}
	}
	spec.Options.Default()
}

type TwemproxyOptions struct {
	// Set logging level to N. (default: 5, min: 0, max: 11)
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	LogLevel *int32 `json:"logLevel,omitempty"`
	// Set stats monitoring port to port.  (default: 22222)
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	MetricsPort *int32 `json:"metricsAddress,omitempty"`
	// Set stats aggregation interval in msec to interval.  (default: 30s)
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	StatsInterval *metav1.Duration `json:"statsInterval,omitempty"`
}

func (opts *TwemproxyOptions) Default() {
	opts.LogLevel = intOrDefault(opts.LogLevel, &twemproxyDefaultLogLevel)
	opts.MetricsPort = intOrDefault(opts.MetricsPort, &twemproxyDefaultMetricsPort)
	if opts.StatsInterval == nil {
		opts.StatsInterval = &twemproxyDefaultStatsInterval
	}
}
