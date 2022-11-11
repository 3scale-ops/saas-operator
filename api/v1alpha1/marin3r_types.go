package v1alpha1

import (
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
)

// SidecarPort defines port for the Marin3r sidecar container
type SidecarPort struct {
	// Port name
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Name string `json:"name"`
	// Port value
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Port int32 `json:"port"`
}

// Marin3rSidecarSpec defines the marin3r sidecar for the component
type Marin3rSidecarSpec struct {
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
	// Compute Resources required by this container.
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
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	EnvoyResources []EnvoyDynamicConfig `json:"envoyResources,omitempty"`
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

// +kubebuilder:validation:MinProperties:=1
// +kubebuilder:validation:MaxProperties:=1
type EnvoyDynamicConfig struct {
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ListenerHttp *ListenerHttp `json:"listenerHttp,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	RouteConfiguration *RouteConfiguration `json:"routeConfiguration,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Cluster *Cluster `json:"cluster,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Runtime *Runtime `json:"runtime,omitempty"`
}

type EnvoyDynamicConfigMeta struct {
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Name string `json:"name"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:default:=v1
	// +optional
	GeneratorVersion *string `json:"generatorVersion,omitempty"`
}

func (meta *EnvoyDynamicConfigMeta) GetName() string {
	return meta.Name
}

func (meta *EnvoyDynamicConfigMeta) GetGeneratorVersion() string {
	return *meta.GeneratorVersion
}

type ListenerHttp struct {
	EnvoyDynamicConfigMeta `json:",inline"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Port uint32 `json:"port"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	RouteConfigName string `json:"routeConfigName"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	CertificateSecretName *string `json:"certificateSecretName,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	RateLimitOptions *RateLimitOptions `json:"rateLimitOptions,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	DefaultHostForHttp10 *string `json:"defaultHostForHttp10,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:default:=false
	// +optional
	EnableHttp2 *bool `json:"enableHttp2,omitempty"`
}

type RateLimitOptions struct {
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Domain string `json:"domain"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:default:=false
	// +optional
	FailureModeDeny *bool `json:"failureModeDeny,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Timeout time.Duration `json:"timeout"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	RateLimitService string `json:"rateLimitService"`
}

type Cluster struct {
	EnvoyDynamicConfigMeta `json:",inline"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Host string `json:"host"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Port uint32 `json:"port"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:default:=false
	// +optional
	IsHttp2 *bool `json:"isHttp2"`
}

type RouteConfiguration struct {
	EnvoyDynamicConfigMeta `json:",inline"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	VirtualHosts []runtime.RawExtension `json:"virtualHosts"`
}

type Runtime struct {
	EnvoyDynamicConfigMeta `json:",inline"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ListenerNames []string `json:"listenerNames"`
}
