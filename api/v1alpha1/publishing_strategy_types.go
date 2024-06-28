package v1alpha1

import (
	"reflect"
	"sort"

	"github.com/3scale-ops/basereconciler/util"
	envoyconfig "github.com/3scale-ops/saas-operator/pkg/resource_builders/envoyconfig/descriptor"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

type PublishingStrategiesReconcileMode string

var (
	PublishingStrategiesReconcileModeMerge   PublishingStrategiesReconcileMode = "Merge"
	PublishingStrategiesReconcileModeReplace PublishingStrategiesReconcileMode = "Replace"
)

type PublishingStrategies struct {
	// PublishingStrategiesReconcileMode specifies if the list of strategies
	// should be merged with the defautls or replace them entirely. Allowed values
	// are "Merge" or "Replace". "Replace" strategy should be used to enable 2 strategies
	// at the same time for a single endpoint.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Mode *PublishingStrategiesReconcileMode `json:"mode,omitempty"`
	// Endpoints holds the list of publishing strategies for each workload endpoint.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Endpoints []PublishingStrategy `json:"endpoints,omitempty"`
}

func (ps *PublishingStrategies) Default() {
	if ps.Mode == nil {
		ps.Mode = util.Pointer(PublishingStrategiesReconcileModeMerge)
	}
	if ps.Endpoints == nil {
		ps.Endpoints = []PublishingStrategy{}
	}
}

// InitializePublishingStrategies initializes a PublishingStrategies struct
func InitializePublishingStrategies(spec *PublishingStrategies) *PublishingStrategies {
	if spec == nil {
		new := &PublishingStrategies{}
		new.Default()
		return new
	} else {
		spec.Default()
		return spec
	}
}

type Strategy string

const (
	SimpleStrategy         Strategy = "Simple"
	Marin3rSidecarStrategy Strategy = "Marin3rSidecar"
)

type PublishingStrategy struct {
	// Strategy defines the type of publishing strategy
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Strategy Strategy `json:"strategy"`
	// EndpointName defines the endpoint affected by this publishing strategy
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	EndpointName string `json:"name"`
	// Simple holds configuration for the Simple publishing strategy
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Simple *Simple `json:"simple,omitempty"`
	// Marin3rSidecar holds configuration for the Marin3rSidecar publishing strategy
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Marin3rSidecar *Marin3rSidecarSpec `json:"marin3rSidecar,omitempty"`
	// Create explicitely tells the controller that this is a new endpoint that
	// should be added. Default is false, causing the controller to error when seeing
	// an unknown endpoint.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Create *bool `json:"create,omitempty"`
}

type ServiceType string

const (
	ServiceTypeClusterIP ServiceType = "ClusterIP"
	ServiceTypeNLB       ServiceType = "NLB"
	ServiceTypeELB       ServiceType = "ELB"
)

type Simple struct {
	// ServiceType defines the type of k8s Service to use for exposing
	// the service to its consumers
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ServiceType *ServiceType `json:"serviceType,omitempty"`
	// ExternalDnsHostnames defines the hostnames that ExternalDNS
	// should configure records for external consumners to reach the service
	// Only works with Services of type NLB/ELB
	ExternalDnsHostnames []string `json:"externalDnsHostnames,omitempty"`
	// ServiceNameOverride allows the user to override the generated
	// Service name
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ServiceNameOverride *string `json:"serviceName,omitempty"`
	// ServicePortsOverride allows the user to override the ports
	// of a Service. It's a replace operation, so specify all the
	// required ports.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ServicePortsOverride []corev1.ServicePort `json:"servicePorts,omitempty"`
	// Classic LB configuration
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ElasticLoadBalancerConfig *ElasticLoadBalancerSpec `json:"elasticLoadBalancerConfig,omitempty"`
	// NLB configuration
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	NetworkLoadBalancerConfig *NetworkLoadBalancerSpec `json:"networkLoadBalancerConfig,omitempty"`
}

func (s *Simple) Default() {
	if s.ServiceType == nil {
		s.ServiceType = util.Pointer(ServiceTypeClusterIP)
	}
}

// Marin3rSidecarSpec defines the marin3r sidecar for the component
type Marin3rSidecarSpec struct {
	*Simple `json:",inline"`
	// The NodeID that identifies the Envoy sidecar to the DiscoveryService
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	NodeID *string `json:"nodeID,omitempty"`
	// The Envoy API version to use
	// +kubebuilder:validation:Enum=v3
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	EnvoyAPIVersion *string `json:"envoyAPIVersion,omitempty"`
	// The Envoy iamge to use
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	EnvoyImage *string `json:"envoyImage,omitempty"`
	// The ports that the sidecar exposes
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Ports []SidecarPort `json:"ports,omitempty"`
	// Compute Resources required by the sidecar container.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Resources *ResourceRequirementsSpec `json:"resources,omitempty"`
	// The port where Marin3r's shutdown manager listens
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ShutdownManagerPort *uint32 `json:"shtdnmgrPort,omitempty"`
	// Extra containers to sync with the shutdown manager upon pod termination
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ShutdownManagerExtraLifecycleHooks []string `json:"shtdnmgrExtraLifecycleHooks"`
	// Extra annotations to pass the Pod to further configure the sidecar container.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ExtraPodAnnotations map[string]string `json:"extraPodAnnotations,omitempty"`
	// Envoy dynamic configuration. Populating this field causes the operator
	// to create a Marin3r EnvoyConfig resource, so Marin3r must be installed
	// in the cluster.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	EnvoyDynamicConfig MapOfEnvoyDynamicConfig `json:"dynamicConfigs,omitempty"`
}

type defaultMarin3rSidecarSpec struct {
	Ports               []SidecarPort
	Resources           defaultResourceRequirementsSpec
	ExtraPodAnnotations map[string]string
}

// Default sets default values for any value not specifically set in the ResourceRequirementsSpec struct
func (spec *Marin3rSidecarSpec) Default(def defaultMarin3rSidecarSpec) {
	if spec.Ports == nil {
		spec.Ports = def.Ports
	}

	if spec.Resources == nil {
		if !reflect.DeepEqual(def.Resources, defaultResourceRequirementsSpec{}) {
			spec.Resources = &ResourceRequirementsSpec{}
			spec.Resources.Default(def.Resources)
		}
	} else {
		spec.Resources.Default(def.Resources)
	}

	// spec.Resources = InitializeResourceRequirementsSpec(spec.Resources, def.Resources)
	if spec.ExtraPodAnnotations == nil {
		spec.ExtraPodAnnotations = def.ExtraPodAnnotations
	}
}

// IsDeactivated true if the field is set with the deactivated value (empty struct)
func (spec *Marin3rSidecarSpec) IsDeactivated() bool {
	return reflect.DeepEqual(spec, &Marin3rSidecarSpec{})
}

// InitializeMarin3rSidecarSpec initializes a ResourceRequirementsSpec struct
func InitializeMarin3rSidecarSpec(spec *Marin3rSidecarSpec, def defaultMarin3rSidecarSpec) *Marin3rSidecarSpec {
	if spec == nil {
		new := &Marin3rSidecarSpec{}
		new.Default(def)
		return new
	}
	if !spec.IsDeactivated() {
		copy := spec.DeepCopy()
		copy.Default(def)
		return copy
	}
	return spec
}

// SidecarPort defines port for the Marin3r sidecar container
type SidecarPort struct {
	// Port name
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Name string `json:"name"`
	// Port value
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Port int32 `json:"port"`
}

type MapOfEnvoyDynamicConfig map[string]EnvoyDynamicConfig

// AsList transforms from the map in the external API to the list of elements
// that the internal API expects.
func (mapofconfs MapOfEnvoyDynamicConfig) AsList() []envoyconfig.EnvoyDynamicConfigDescriptor {

	list := make([]envoyconfig.EnvoyDynamicConfigDescriptor, 0, len(mapofconfs))

	for name, conf := range mapofconfs {
		list = append(list, conf.DeepCopy().AsEnvoyDynamicConfigDescriptor(name))
	}

	// ensure consistent order of configs
	sort.Slice(list, func(a, b int) bool {
		return list[a].GetName() < list[b].GetName()
	})

	return list
}

// +kubebuilder:validation:MinProperties:=2
// +kubebuilder:validation:MaxProperties:=2
type EnvoyDynamicConfig struct {
	// hidden field
	Name string `json:"-"`
	// GeneratorVersion specifies the version of a given template.
	// "v1" is the default.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:default:=v1
	// +optional
	GeneratorVersion *string `json:"generatorVersion,omitempty"`
	// ListenerHttp contains options for an HTTP/HTTPS listener
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ListenerHttp *ListenerHttp `json:"listenerHttp,omitempty"`
	// RouteConfiguration contains options for an Envoy route_configuration
	// protobuffer message
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	RouteConfiguration *RouteConfiguration `json:"routeConfiguration,omitempty"`
	// Cluster contains options for an Envoy cluster protobuffer message
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Cluster *Cluster `json:"cluster,omitempty"`
	// Runtime contains options for an Envoy runtime protobuffer message
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Runtime *Runtime `json:"runtime,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	RawConfig *RawConfig `json:"rawConfig,omitempty"`
}

// AsEnvoyDynamicConfigDescriptor converts the external API type into the internal EnvoyDynamicConfigDescriptor
// interface. The name field is populated with the parameter passed to the function.
func (config *EnvoyDynamicConfig) AsEnvoyDynamicConfigDescriptor(name string) envoyconfig.EnvoyDynamicConfigDescriptor {
	config.Name = name
	return config
}

func (config *EnvoyDynamicConfig) GetName() string {
	return config.Name
}

// GetGeneratorVersion returns the template's version
func (config *EnvoyDynamicConfig) GetGeneratorVersion() string {
	return *config.GeneratorVersion
}

func (config *EnvoyDynamicConfig) GetOptions() interface{} {
	if config.ListenerHttp != nil {
		return config.ListenerHttp
	} else if config.RouteConfiguration != nil {
		return config.RouteConfiguration
	} else if config.Cluster != nil {
		return config.Cluster
	} else if config.Runtime != nil {
		return config.Runtime
	} else if config.RawConfig != nil {
		return config.RawConfig
	}

	return nil
}

// ListenerHttp contains options for an HTTP/HTTPS listener
type ListenerHttp struct {
	// The port where the listener listens for new connections
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Port uint32 `json:"port"`
	// Whether proxy protocol should be enabled or not. Defaults to true.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:default:=true
	// +optional
	ProxyProtocol *bool `json:"proxyProtocol,omitempty"`
	// The name of the RouteConfiguration to use in the listener
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	RouteConfigName string `json:"routeConfigName"`
	// The name of the Secret containing a valid certificate. If unset
	// the listener will be http, if set https
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	CertificateSecretName *string `json:"certificateSecretName,omitempty"`
	// Rate limit options for the ratelimit filter of the HTTP connection
	// manager
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	RateLimitOptions *RateLimitOptions `json:"rateLimitOptions,omitempty"`
	// If this filed is set, http 1.0 will be enabled and this will be
	// the default hostname to use.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	DefaultHostForHttp10 *string `json:"defaultHostForHttp10,omitempty"`
	// Enable http2 in the listener.Disabled by default.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:default:=false
	// +optional
	EnableHttp2 *bool `json:"enableHttp2,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:default:=true
	// +optional
	// Allow headers with underscores
	AllowHeadersWithUnderscores *bool `json:"allowHeadersWithUnderscores,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	// Max connection duration. If unset no max connection duration will be applied.
	MaxConnectionDuration *metav1.Duration `json:"maxConnectionDuration,omitempty"`
}

// RateLimitOptions contains options for the ratelimit filter of the
// http connection manager
type RateLimitOptions struct {
	// The rate limit domain
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Domain string `json:"domain"`
	// Whether to allow requests or not if the rate limit service
	// is unavailable
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:default:=false
	// +optional
	FailureModeDeny *bool `json:"failureModeDeny,omitempty"`
	// Max time to wait for a response from the rate limit service
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:validation:Format:=duration
	Timeout metav1.Duration `json:"timeout"`
	// Location of the rate limit service. Must point to one of the
	// defined clusters.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	RateLimitCluster string `json:"rateLimitCluster"`
}

// Cluster contains options for an Envoy cluster protobuffer message
type Cluster struct {
	// The upstream host
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Host string `json:"host"`
	// The upstream port
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Port uint32 `json:"port"`
	// Specifies if the upstream cluster is http2 or not (default).
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:default:=false
	// +optional
	IsHttp2 *bool `json:"isHttp2"`
}

// RouteConfiguration contains options for an Envoy route_configuration
// protobuffer message
type RouteConfiguration struct {
	// The virtual_hosts definitions for this route configuration.
	// Virtual hosts must be specified using directly Envoy's API
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	VirtualHosts []runtime.RawExtension `json:"virtualHosts"`
}

// Runtime contains options for an Envoy runtime protobuffer message
type Runtime struct {
	// The list of listeners to apply overload protection limits to
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ListenerNames []string `json:"listenerNames"`
}

// RawConfig is a struct with methods to manage a
// configuration defined using directly the Envoy config API
type RawConfig struct {
	// Type is the type url for the protobuf message
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:validation:Enum=listener;routeConfiguration;cluster;runtime
	Type string `json:"type"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// Allows defining configuration using directly envoy's config API.
	// WARNING: no validation of this field's value is performed before
	// writting the custom resource to etcd.
	Value runtime.RawExtension `json:"value"`
}

// UPGRADE CODE
// TODO: delete after upgrade release

func PublishingStrategyGenerator(bldrs ...WorkloadPublishingStrategyBuilder) *PublishingStrategies {
	endpoints := []PublishingStrategy{}
	for _, bldr := range bldrs {
		if endpoint := bldr.Build(); endpoint != nil {
			endpoints = append(endpoints, *endpoint)
		}
	}
	return &PublishingStrategies{
		Mode:      util.Pointer(PublishingStrategiesReconcileModeMerge),
		Endpoints: endpoints,
	}
}

type WorkloadPublishingStrategyBuilder struct {
	Name                 string
	ServiceNameOverride  *string
	Endpoint             *Endpoint
	Marin3r              *Marin3rSidecarSpec
	ELB                  *ElasticLoadBalancerSpec
	NLB                  *NetworkLoadBalancerSpec
	ServicePortOverrides []corev1.ServicePort
}

func (gen WorkloadPublishingStrategyBuilder) Build() *PublishingStrategy {
	var out *PublishingStrategy

	if gen.Marin3r != nil {
		// generate a Marin3rSidecar strategy
		out = &PublishingStrategy{
			Strategy:       Marin3rSidecarStrategy,
			EndpointName:   gen.Name,
			Marin3rSidecar: gen.Marin3r,
		}
		gen.Marin3r.Simple = &Simple{}

		if gen.Endpoint != nil && len(gen.Endpoint.DNS) > 0 {
			out.Marin3rSidecar.Simple.ExternalDnsHostnames = gen.Endpoint.DNS
		}

		if gen.ServiceNameOverride != nil {
			out.Marin3rSidecar.Simple.ServiceNameOverride = gen.ServiceNameOverride
		}

		if gen.ELB != nil {
			out.Marin3rSidecar.ServiceType = util.Pointer(ServiceTypeELB)
			out.Marin3rSidecar.Simple.ElasticLoadBalancerConfig = gen.ELB
		}

		if gen.NLB != nil {
			out.Marin3rSidecar.Simple.ServiceType = util.Pointer(ServiceTypeNLB)
			out.Marin3rSidecar.Simple.NetworkLoadBalancerConfig = gen.NLB
		}

		if len(gen.ServicePortOverrides) > 0 {
			out.Marin3rSidecar.Simple.ServicePortsOverride = gen.ServicePortOverrides
		}

	} else {
		// generate a Simple strategy
		out = &PublishingStrategy{
			Strategy:     SimpleStrategy,
			EndpointName: gen.Name,
			Simple: &Simple{
				ServiceType: util.Pointer(ServiceTypeClusterIP),
			},
		}

		if gen.Endpoint != nil && len(gen.Endpoint.DNS) > 0 {
			out.Simple.ExternalDnsHostnames = gen.Endpoint.DNS
		}

		if gen.ServiceNameOverride != nil {
			out.Simple.ServiceNameOverride = gen.ServiceNameOverride
		}

		if gen.ELB != nil {
			out.Simple.ServiceType = util.Pointer(ServiceTypeELB)
			out.Simple.ElasticLoadBalancerConfig = gen.ELB
		}

		if gen.NLB != nil {
			out.Simple.ServiceType = util.Pointer(ServiceTypeNLB)
			out.Simple.NetworkLoadBalancerConfig = gen.NLB
		}

	}

	return out
}
