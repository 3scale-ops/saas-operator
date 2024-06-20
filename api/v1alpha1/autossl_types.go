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
	"github.com/3scale-ops/basereconciler/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	autosslDefaultReplicas int32            = 2
	autosslDefaultImage    defaultImageSpec = defaultImageSpec{
		Name:       util.Pointer("quay.io/3scale/autossl"),
		Tag:        util.Pointer("latest"),
		PullPolicy: (*corev1.PullPolicy)(util.Pointer(string(corev1.PullIfNotPresent))),
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
		MinReplicas:         util.Pointer[int32](2),
		MaxReplicas:         util.Pointer[int32](4),
		ResourceUtilization: util.Pointer[int32](90),
		ResourceName:        util.Pointer("cpu"),
	}
	autosslDefaultProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: util.Pointer[int32](25),
		TimeoutSeconds:      util.Pointer[int32](1),
		PeriodSeconds:       util.Pointer[int32](10),
		SuccessThreshold:    util.Pointer[int32](1),
		FailureThreshold:    util.Pointer[int32](3),
	}
	autosslDefaultPDB defaultPodDisruptionBudgetSpec = defaultPodDisruptionBudgetSpec{
		MaxUnavailable: util.Pointer(intstr.FromInt(1)),
	}

	autosslDefaultGrafanaDashboard defaultGrafanaDashboardSpec = defaultGrafanaDashboardSpec{
		SelectorKey:   util.Pointer("monitoring-key"),
		SelectorValue: util.Pointer("middleware"),
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
	// Configures the Grafana Dashboard for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	GrafanaDashboard *GrafanaDashboardSpec `json:"grafanaDashboard,omitempty"`
	// Application specific configuration options for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Config AutoSSLConfig `json:"config"`
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
	// Describes how the services provided by this workload are exposed to its consumers
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	PublishingStrategies *PublishingStrategies `json:"publishingStrategies,omitempty"`
	// The external endpoint/s for the component
	// DEPRECATED
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Endpoint Endpoint `json:"endpoint"`
	// Configures the AWS load balancer for the component
	// DEPRECATED
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	LoadBalancer *ElasticLoadBalancerSpec `json:"loadBalancer,omitempty"`
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
	spec.LoadBalancer = InitializeElasticLoadBalancerSpec(spec.LoadBalancer, DefaultElasticLoadBalancerSpec)
	spec.GrafanaDashboard = InitializeGrafanaDashboardSpec(spec.GrafanaDashboard, autosslDefaultGrafanaDashboard)
	spec.PublishingStrategies = InitializePublishingStrategies(spec.PublishingStrategies)
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
	cfg.ACMEStaging = boolOrDefault(cfg.ACMEStaging, util.Pointer(autosslDefaultACMEStaging))
	cfg.RedisPort = intOrDefault(cfg.RedisPort, util.Pointer[int32](autosslDefaultRedisPort))
	cfg.LogLevel = stringOrDefault(cfg.LogLevel, util.Pointer(autosslDefaultLogLevel))
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
