/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"github.com/3scale/saas-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

var (
	apicastDefaultReplicas int32            = 2
	apicastDefaultImage    defaultImageSpec = defaultImageSpec{
		Name:       pointer.StringPtr("quay.io/3scale/apicast-cloud-hosted"),
		Tag:        pointer.StringPtr("latest"),
		PullPolicy: (*corev1.PullPolicy)(pointer.StringPtr(string(corev1.PullIfNotPresent))),
	}
	apicastDefaultLoadBalancer defaultLoadBalancerSpec = defaultLoadBalancerSpec{
		ProxyProtocol:                 pointer.BoolPtr(true),
		CrossZoneLoadBalancingEnabled: pointer.BoolPtr(true),
		ConnectionDrainingEnabled:     pointer.BoolPtr(true),
		ConnectionDrainingTimeout:     pointer.Int32Ptr(60),
		HealthcheckHealthyThreshold:   pointer.Int32Ptr(2),
		HealthcheckUnhealthyThreshold: pointer.Int32Ptr(2),
		HealthcheckInterval:           pointer.Int32Ptr(5),
		HealthcheckTimeout:            pointer.Int32Ptr(3),
	}
	apicastDefaultResources defaultResourceRequirementsSpec = defaultResourceRequirementsSpec{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("500m"),
			corev1.ResourceMemory: resource.MustParse("64Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("1"),
			corev1.ResourceMemory: resource.MustParse("128Mi"),
		},
	}
	apicastDefaultHPA defaultHorizontalPodAutoscalerSpec = defaultHorizontalPodAutoscalerSpec{
		MinReplicas:         pointer.Int32Ptr(2),
		MaxReplicas:         pointer.Int32Ptr(4),
		ResourceUtilization: pointer.Int32Ptr(90),
		ResourceName:        pointer.StringPtr("cpu"),
	}
	apicastDefaultLivenessProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: pointer.Int32Ptr(5),
		TimeoutSeconds:      pointer.Int32Ptr(5),
		PeriodSeconds:       pointer.Int32Ptr(10),
		SuccessThreshold:    pointer.Int32Ptr(1),
		FailureThreshold:    pointer.Int32Ptr(3),
	}
	apicastDefaultReadinessProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: pointer.Int32Ptr(5),
		TimeoutSeconds:      pointer.Int32Ptr(5),
		PeriodSeconds:       pointer.Int32Ptr(30),
		SuccessThreshold:    pointer.Int32Ptr(1),
		FailureThreshold:    pointer.Int32Ptr(3),
	}
	apicastDefaultPDB defaultPodDisruptionBudgetSpec = defaultPodDisruptionBudgetSpec{
		MaxUnavailable: util.IntStrPtr(intstr.FromInt(1)),
	}
	apicastDefaultGrafanaDashboard defaultGrafanaDashboardSpec = defaultGrafanaDashboardSpec{
		SelectorKey:   pointer.StringPtr("monitoring-key"),
		SelectorValue: pointer.StringPtr("middleware"),
	}
	apicastDefaultMarin3rSpec  defaultMarin3rSidecarSpec = defaultMarin3rSidecarSpec{}
	apicastDefaultLogLevel     string                    = "warn"
	apicastDefaultOIDCLogLevel string                    = "warn"
)

// ApicastSpec defines the desired state of Apicast
type ApicastSpec struct {
	// Configures the staging Apicast environment
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Staging ApicastEnvironmentSpec `json:"staging"`
	// Configures the production Apicast environment
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Production ApicastEnvironmentSpec `json:"production"`
	// Configures the Grafana Dashboard for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	GrafanaDashboard *GrafanaDashboardSpec `json:"grafanaDashboard,omitempty"`
}

// Default implements defaulting for ApicastSpec
func (spec *ApicastSpec) Default() {

	spec.Staging.Default()
	spec.Production.Default()
	spec.GrafanaDashboard = InitializeGrafanaDashboardSpec(spec.GrafanaDashboard, apicastDefaultGrafanaDashboard)

}

// ResolveCanarySpec modifies the BackendSpec given the provided canary configuration
func (spec *ApicastSpec) ResolveCanarySpec(canary *Canary) (*ApicastSpec, error) {
	canarySpec := &ApicastSpec{}
	if err := canary.PatchSpec(spec, canarySpec); err != nil {
		return nil, err
	}

	if canary.ImageName != nil {
		canarySpec.Staging.Image.Name = canary.ImageName
		canarySpec.Production.Image.Name = canary.ImageName
	}
	if canary.ImageTag != nil {
		canarySpec.Staging.Image.Tag = canary.ImageTag
		canarySpec.Production.Image.Tag = canary.ImageTag
	}

	canarySpec.Staging.Replicas = canary.Replicas
	canarySpec.Production.Replicas = canary.Replicas

	// Call Default() on the resolved canary spec to apply
	// defaulting to potentially added fields
	canarySpec.Default()

	return canarySpec, nil
}

// ApicastEnvironmentSpec is the configuration for an Apicast environment
type ApicastEnvironmentSpec struct {
	// Image specification for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Image *ImageSpec `json:"image,omitempty"`
	// Pod Disruption Budget for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	PDB *PodDisruptionBudgetSpec `json:"pdb,omitempty"`
	// Horizontal Pod Autoscaler for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	HPA *HorizontalPodAutoscalerSpec `json:"hpa,omitempty"`
	// Number of replicas (ignored if hpa is enabled) for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
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
	// Application specific configuration options for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Config ApicastConfig `json:"config"`
	// The external endpoint/s for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Endpoint Endpoint `json:"endpoint"`
	// Marin3r configures the Marin3r sidecars for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Marin3r *Marin3rSidecarSpec `json:"marin3r,omitempty"`
	// Configures the AWS load balancer for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	LoadBalancer *LoadBalancerSpec `json:"loadBalancer,omitempty"`
	// Describes node affinity scheduling rules for the pod.
	// +optional
	NodeAffinity *corev1.NodeAffinity `json:"nodeAffinity,omitempty" protobuf:"bytes,1,opt,name=nodeAffinity"`
	// If specified, the pod's tolerations.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty" protobuf:"bytes,22,opt,name=tolerations"`
	// Canary defines spec changes for the canary Deployment. If
	// left unset the canary Deployment wil not be created.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Canary *Canary `json:"canary,omitempty"`
}

// Default implements defaulting for the each apicast environment
func (spec *ApicastEnvironmentSpec) Default() {

	spec.Image = InitializeImageSpec(spec.Image, apicastDefaultImage)
	spec.HPA = InitializeHorizontalPodAutoscalerSpec(spec.HPA, apicastDefaultHPA)
	spec.Replicas = intOrDefault(spec.Replicas, &apicastDefaultReplicas)
	spec.PDB = InitializePodDisruptionBudgetSpec(spec.PDB, apicastDefaultPDB)
	spec.Resources = InitializeResourceRequirementsSpec(spec.Resources, apicastDefaultResources)
	spec.LivenessProbe = InitializeProbeSpec(spec.LivenessProbe, apicastDefaultLivenessProbe)
	spec.ReadinessProbe = InitializeProbeSpec(spec.ReadinessProbe, apicastDefaultReadinessProbe)
	spec.LoadBalancer = InitializeLoadBalancerSpec(spec.LoadBalancer, apicastDefaultLoadBalancer)
	spec.Marin3r = InitializeMarin3rSidecarSpec(spec.Marin3r, apicastDefaultMarin3rSpec)
	spec.Config.Default()
}

// ApicastConfig configures app behavior for Apicast
type ApicastConfig struct {
	// Apicast configurations cache TTL
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ConfigurationCache int32 `json:"configurationCache"`
	// Endpoint to request proxy configurations to
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ThreescalePortalEndpoint string `json:"threescalePortalEndpoint"`
	// Openresty log level
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:validation:Enum=debug;info;notice;warn;error;crit;alert;emerg
	// +optional
	LogLevel *string `json:"logLevel,omitempty"`
	// OpenID Connect integration log level
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:validation:Enum=debug;info;notice;warn;error;crit;alert;emerg
	// +optional
	OIDCLogLevel *string `json:"oidcLogLevel,omitempty"`
}

// Default sets default values for any value not specifically set in the ApicastConfig struct
func (cfg *ApicastConfig) Default() {
	cfg.LogLevel = stringOrDefault(cfg.LogLevel, pointer.StringPtr(apicastDefaultLogLevel))
	cfg.OIDCLogLevel = stringOrDefault(cfg.OIDCLogLevel, pointer.StringPtr(apicastDefaultOIDCLogLevel))

}

// ApicastStatus defines the observed state of Apicast
type ApicastStatus struct{}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Apicast is the Schema for the apicasts API
type Apicast struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApicastSpec   `json:"spec,omitempty"`
	Status ApicastStatus `json:"status,omitempty"`
}

// Default implements defaulting for the Apicast resource
func (a *Apicast) Default() {
	a.Spec.Default()
}

// +kubebuilder:object:root=true

// ApicastList contains a list of Apicast
type ApicastList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Apicast `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Apicast{}, &ApicastList{})
}
