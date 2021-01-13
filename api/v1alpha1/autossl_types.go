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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

var (
	defaultReplicas int32            = 2
	defaultImage    defaultImageSpec = defaultImageSpec{
		Name: pointer.StringPtr("quay.io/3scale/autossl"),
		Tag:  pointer.StringPtr("latest"),
	}
	defaultLoadBalancer defaultLoadBalancerSpec = defaultLoadBalancerSpec{
		ProxyProtocol:                 pointer.BoolPtr(true),
		CrossZoneLoadBalancingEnabled: pointer.BoolPtr(true),
		ConnectionDrainingEnabled:     pointer.BoolPtr(true),
		ConnectionDrainingTimeout:     pointer.Int32Ptr(60),
		HealthcheckHealthyThreshold:   pointer.Int32Ptr(2),
		HealthcheckUnhealthyThreshold: pointer.Int32Ptr(2),
		HealthcheckInterval:           pointer.Int32Ptr(5),
		HealthcheckTimeout:            pointer.Int32Ptr(3),
	}
	defaultResources defaultResourceRequirementsSpec = defaultResourceRequirementsSpec{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("75m"),
			corev1.ResourceMemory: resource.MustParse("64Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("150m"),
			corev1.ResourceMemory: resource.MustParse("128Mi"),
		},
	}
	defaultHPA defaultHorizontalPodAutoscalerSpec = defaultHorizontalPodAutoscalerSpec{
		MinReplicas:         pointer.Int32Ptr(2),
		MaxReplicas:         pointer.Int32Ptr(4),
		ResourceUtilization: pointer.Int32Ptr(90),
		ResourceName:        pointer.StringPtr("cpu"),
	}
	defaultProbe defaultHTTPProbeSpec = defaultHTTPProbeSpec{
		InitialDelaySeconds: pointer.Int32Ptr(25),
		TimeoutSeconds:      pointer.Int32Ptr(1),
		PeriodSeconds:       pointer.Int32Ptr(10),
		SuccessThreshold:    pointer.Int32Ptr(1),
		FailureThreshold:    pointer.Int32Ptr(3),
	}
	defaultPDB defaultPodDisruptionBudgetSpec = defaultPodDisruptionBudgetSpec{
		MaxUnavailable: intstrPtr(intstr.IntOrString{Type: intstr.Int, IntVal: 1}),
	}

	defaultGrafanaDashboard defaultGrafanaDashboardSpec = defaultGrafanaDashboardSpec{
		SelectorKey:   pointer.StringPtr("monitoring-key"),
		SelectorValue: pointer.StringPtr("middleware"),
	}
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
	LivenessProbe *HTTPProbeSpec `json:"livenessProbe,omitempty"`
	// Readiness probe for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ReadinessProbe *HTTPProbeSpec `json:"readinessProbe,omitempty"`
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
}

// Default implements defaulting for the AutoSSL resource
func (r *AutoSSL) Default() {

	r.Spec.Image = InitializeImageSpec(r.Spec.Image, defaultImage)
	r.Spec.HPA = InitializeHorizontalPodAutoscalerSpec(r.Spec.HPA, defaultHPA)

	if r.Spec.HPA.IsDeactivated() {
		r.Spec.Replicas = intOrDefault(r.Spec.Replicas, &defaultReplicas)
	} else {
		r.Spec.Replicas = nil
	}

	r.Spec.PDB = InitializePodDisruptionBudgetSpec(r.Spec.PDB, defaultPDB)
	r.Spec.Resources = InitializeResourceRequirementsSpec(r.Spec.Resources, defaultResources)
	r.Spec.LivenessProbe = InitializeHTTPProbeSpec(r.Spec.LivenessProbe, defaultProbe)
	r.Spec.ReadinessProbe = InitializeHTTPProbeSpec(r.Spec.ReadinessProbe, defaultProbe)
	r.Spec.LoadBalancer = InitializeLoadBalancerSpec(r.Spec.LoadBalancer, defaultLoadBalancer)
	r.Spec.GrafanaDashboard = InitializeGrafanaDashboardSpec(r.Spec.GrafanaDashboard, defaultGrafanaDashboard)
	r.Spec.Config.Default()
}

// AutoSSLConfig defines configuration options for the component
type AutoSSLConfig struct {
	// ets the nginx log level
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
func (spec *AutoSSLConfig) Default() {
	spec.ACMEStaging = pointer.BoolPtr(false)
	spec.RedisPort = pointer.Int32Ptr(6379)
	spec.LogLevel = pointer.StringPtr("warn")
	if spec.DomainWhitelist == nil {
		spec.DomainWhitelist = []string{}
	}
	if spec.DomainBlacklist == nil {
		spec.DomainBlacklist = []string{}
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
