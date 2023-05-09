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
	autosslDefaultReplicas int32            = 2
	autosslDefaultImage    defaultImageSpec = defaultImageSpec{
		Name:       pointer.String("quay.io/3scale/autossl"),
		Tag:        pointer.String("latest"),
		PullPolicy: (*corev1.PullPolicy)(pointer.String(string(corev1.PullIfNotPresent))),
	}
	autosslDefaultLoadBalancer defaultLoadBalancerSpec = defaultLoadBalancerSpec{
		ProxyProtocol:                 pointer.Bool(true),
		CrossZoneLoadBalancingEnabled: pointer.Bool(true),
		ConnectionDrainingEnabled:     pointer.Bool(true),
		ConnectionDrainingTimeout:     pointer.Int32(60),
		HealthcheckHealthyThreshold:   pointer.Int32(2),
		HealthcheckUnhealthyThreshold: pointer.Int32(2),
		HealthcheckInterval:           pointer.Int32(5),
		HealthcheckTimeout:            pointer.Int32(3),
	}
	autosslDefaultResources defaultResourceRequirementsSpec = defaultResourceRequirementsSpec{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("75m"),
			corev1.ResourceMemory: resource.MustParse("64Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("150m"),
			corev1.ResourceMemory: resource.MustParse("128Mi"),
		},
	}
	autosslDefaultHPA defaultHorizontalPodAutoscalerSpec = defaultHorizontalPodAutoscalerSpec{
		MinReplicas:         pointer.Int32(2),
		MaxReplicas:         pointer.Int32(4),
		ResourceUtilization: pointer.Int32(90),
		ResourceName:        pointer.String("cpu"),
	}
	autosslDefaultProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: pointer.Int32(25),
		TimeoutSeconds:      pointer.Int32(1),
		PeriodSeconds:       pointer.Int32(10),
		SuccessThreshold:    pointer.Int32(1),
		FailureThreshold:    pointer.Int32(3),
	}
	autosslDefaultPDB defaultPodDisruptionBudgetSpec = defaultPodDisruptionBudgetSpec{
		MaxUnavailable: util.IntStrPtr(intstr.FromInt(1)),
	}

	autosslDefaultGrafanaDashboard defaultGrafanaDashboardSpec = defaultGrafanaDashboardSpec{
		SelectorKey:   pointer.String("monitoring-key"),
		SelectorValue: pointer.String("middleware"),
	}
	autosslDefaultACMEStaging bool   = false
	autosslDefaultRedisPort   int32  = 6379
	autosslDefaultLogLevel    string = "warn"
)

// AutoSSLSpec defines the desired state of AutoSSL
type AutoSSLSpec struct {
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
	// Configures the AWS load balancer for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	LoadBalancer *LoadBalancerSpec `json:"loadBalancer,omitempty"`
	// Configures the Grafana Dashboard for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	GrafanaDashboard *GrafanaDashboardSpec `json:"grafanaDashboard,omitempty"`
	// Application specific configuration options for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Config AutoSSLConfig `json:"config"`
	// The external endpoint/s for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Endpoint Endpoint `json:"endpoint"`
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

// Default implements defaulting for AutoSSLSpec
func (spec *AutoSSLSpec) Default() {

	spec.Image = InitializeImageSpec(spec.Image, autosslDefaultImage)
	spec.HPA = InitializeHorizontalPodAutoscalerSpec(spec.HPA, autosslDefaultHPA)
	spec.Replicas = intOrDefault(spec.Replicas, &autosslDefaultReplicas)
	spec.PDB = InitializePodDisruptionBudgetSpec(spec.PDB, autosslDefaultPDB)
	spec.Resources = InitializeResourceRequirementsSpec(spec.Resources, autosslDefaultResources)
	spec.LivenessProbe = InitializeProbeSpec(spec.LivenessProbe, autosslDefaultProbe)
	spec.ReadinessProbe = InitializeProbeSpec(spec.ReadinessProbe, autosslDefaultProbe)
	spec.LoadBalancer = InitializeLoadBalancerSpec(spec.LoadBalancer, autosslDefaultLoadBalancer)
	spec.GrafanaDashboard = InitializeGrafanaDashboardSpec(spec.GrafanaDashboard, autosslDefaultGrafanaDashboard)
	spec.Config.Default()
}

// ResolveCanarySpec modifies the AutoSSLSpec given the provided canary configuration
func (spec *AutoSSLSpec) ResolveCanarySpec(canary *Canary) (*AutoSSLSpec, error) {
	canarySpec := &AutoSSLSpec{}
	canary.PatchSpec(spec, canarySpec)
	if canary.ImageName != nil {
		canarySpec.Image.Name = canary.ImageName
	}
	if canary.ImageTag != nil {
		canarySpec.Image.Tag = canary.ImageTag
	}
	canarySpec.Replicas = canary.Replicas

	// Call Default() on the resolved canary spec to apply
	// defaulting to potentially added fields
	canarySpec.Default()

	return canarySpec, nil
}

// AutoSSLConfig defines configuration options for the component
type AutoSSLConfig struct {
	// Sets the nginx log level
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	LogLevel *string `json:"logLevel,omitempty"`
	// Enables/disables the Let's Encrypt staging ACME endpoint
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ACMEStaging *bool `json:"acmeStaging,omitempty"`
	// Defines an email address for Let's Encrypt notifications
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ContactEmail string `json:"contactEmail"`
	// The endpoint to proxy_pass requests to
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ProxyEndpoint string `json:"proxyEndpoint"`
	// The endpoint used to validate if certificate generation is allowed
	// for the domain
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	VerificationEndpoint string `json:"verificationEndpoint"`
	// List of domains that will bypass domain verification
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	DomainWhitelist []string `json:"domainWhitelist,omitempty"`
	// List of domains that will never get autogenerated certificates
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	DomainBlacklist []string `json:"domainBlacklist,omitempty"`
	// Host for the redis database to store certificates
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	RedisHost string `json:"redisHost"`
	// Port for the redis database to store certificates
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	RedisPort *int32 `json:"redisPort,omitempty"`
}

// Default sets default values for any value not specifically set in the AutoSSLConfig struct
func (cfg *AutoSSLConfig) Default() {
	cfg.ACMEStaging = boolOrDefault(cfg.ACMEStaging, pointer.Bool(autosslDefaultACMEStaging))
	cfg.RedisPort = intOrDefault(cfg.RedisPort, pointer.Int32(autosslDefaultRedisPort))
	cfg.LogLevel = stringOrDefault(cfg.LogLevel, pointer.String(autosslDefaultLogLevel))
	if cfg.DomainWhitelist == nil {
		cfg.DomainWhitelist = []string{}
	}
	if cfg.DomainBlacklist == nil {
		cfg.DomainBlacklist = []string{}
	}
}

// AutoSSLStatus defines the observed state of AutoSSL
type AutoSSLStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// AutoSSL is the Schema for the autossls API
type AutoSSL struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AutoSSLSpec   `json:"spec,omitempty"`
	Status AutoSSLStatus `json:"status,omitempty"`
}

// Default implements defaulting for the AutoSSL resource
func (a *AutoSSL) Default() {
	a.Spec.Default()
}

// +kubebuilder:object:root=true

// AutoSSLList contains a list of AutoSSL
type AutoSSLList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AutoSSL `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AutoSSL{}, &AutoSSLList{})
}
