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
	mappingserviceDefaultReplicas int32            = 2
	mappingserviceDefaultImage    defaultImageSpec = defaultImageSpec{
		Name:       util.Pointer("quay.io/3scale/apicast-cloud-hosted"),
		Tag:        util.Pointer("latest"),
		PullPolicy: (*corev1.PullPolicy)(util.Pointer(string(corev1.PullIfNotPresent))),
	}
	mappingserviceDefaultResources defaultResourceRequirementsSpec = defaultResourceRequirementsSpec{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("500m"),
			corev1.ResourceMemory: resource.MustParse("64Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("1"),
			corev1.ResourceMemory: resource.MustParse("128Mi"),
		},
	}
	mappingserviceDefaultHPA defaultHorizontalPodAutoscalerSpec = defaultHorizontalPodAutoscalerSpec{
		MinReplicas:         util.Pointer[int32](2),
		MaxReplicas:         util.Pointer[int32](4),
		ResourceUtilization: util.Pointer[int32](90),
		ResourceName:        util.Pointer("cpu"),
	}
	mappingserviceLivenessDefaultProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: util.Pointer[int32](5),
		TimeoutSeconds:      util.Pointer[int32](5),
		PeriodSeconds:       util.Pointer[int32](10),
		SuccessThreshold:    util.Pointer[int32](1),
		FailureThreshold:    util.Pointer[int32](3),
	}
	mappingserviceReadinessDefaultProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: util.Pointer[int32](5),
		TimeoutSeconds:      util.Pointer[int32](5),
		PeriodSeconds:       util.Pointer[int32](30),
		SuccessThreshold:    util.Pointer[int32](1),
		FailureThreshold:    util.Pointer[int32](3),
	}
	mappingserviceDefaultPDB defaultPodDisruptionBudgetSpec = defaultPodDisruptionBudgetSpec{
		MaxUnavailable: util.Pointer(intstr.FromInt(1)),
	}

	mappingserviceDefaultGrafanaDashboard defaultGrafanaDashboardSpec = defaultGrafanaDashboardSpec{
		SelectorKey:   util.Pointer("monitoring-key"),
		SelectorValue: util.Pointer("middleware"),
	}
	mappingserviceDefaultLogLevel string = "warn"
)

// MappingServiceSpec defines the desired state of MappingService
type MappingServiceSpec struct {
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
	Config MappingServiceConfig `json:"config"`
	// Describes node affinity scheduling rules for the pod.
	// +optional
	NodeAffinity *corev1.NodeAffinity `json:"nodeAffinity,omitempty" protobuf:"bytes,1,opt,name=nodeAffinity"`
	// If specified, the pod's tolerations.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty" protobuf:"bytes,22,opt,name=tolerations"`
	// Describes how the services provided by this workload are exposed to its consumers
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	PublishingStrategies *PublishingStrategies `json:"publishingStrategies,omitempty"`
}

// Default implements defaulting for MappingServiceSpec
func (spec *MappingServiceSpec) Default() {

	spec.Image = InitializeImageSpec(spec.Image, mappingserviceDefaultImage)
	spec.HPA = InitializeHorizontalPodAutoscalerSpec(spec.HPA, mappingserviceDefaultHPA)
	spec.Replicas = intOrDefault(spec.Replicas, &mappingserviceDefaultReplicas)
	spec.PDB = InitializePodDisruptionBudgetSpec(spec.PDB, mappingserviceDefaultPDB)
	spec.Resources = InitializeResourceRequirementsSpec(spec.Resources, mappingserviceDefaultResources)
	spec.LivenessProbe = InitializeProbeSpec(spec.LivenessProbe, mappingserviceLivenessDefaultProbe)
	spec.ReadinessProbe = InitializeProbeSpec(spec.ReadinessProbe, mappingserviceReadinessDefaultProbe)
	spec.GrafanaDashboard = InitializeGrafanaDashboardSpec(spec.GrafanaDashboard, mappingserviceDefaultGrafanaDashboard)
	spec.Config.Default()
}

// MappingServiceConfig configures app behavior for MappingService
type MappingServiceConfig struct {
	// System endpoint to fetch proxy configs from
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	APIHost string `json:"apiHost"`
	// Base domain to replace the proxy configs base domain
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	PreviewBaseDomain *string `json:"previewBaseDomain,omitempty"`
	// Openresty log level
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuiler:validation:Enum=debug;info;notice;warn;error;crit;alert;emerg
	// +optional
	LogLevel *string `json:"logLevel,omitempty"`
	// External Secret common configuration
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ExternalSecret ExternalSecret `json:"externalSecret,omitempty"`
	// A reference to the secret holding the system admin token
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	SystemAdminToken SecretReference `json:"systemAdminToken"`
}

// Default sets default values for any value not specifically set in the MappingServiceConfig struct
func (cfg *MappingServiceConfig) Default() {
	cfg.LogLevel = stringOrDefault(cfg.LogLevel, util.Pointer(mappingserviceDefaultLogLevel))
	cfg.ExternalSecret.SecretStoreRef = InitializeExternalSecretSecretStoreReferenceSpec(cfg.ExternalSecret.SecretStoreRef, defaultExternalSecretSecretStoreReference)
	cfg.ExternalSecret.RefreshInterval = durationOrDefault(cfg.ExternalSecret.RefreshInterval, &defaultExternalSecretRefreshInterval)
}

// MappingServiceStatus defines the observed state of MappingService
type MappingServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MappingService is the Schema for the mappingservices API
type MappingService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MappingServiceSpec   `json:"spec,omitempty"`
	Status MappingServiceStatus `json:"status,omitempty"`
}

// Default implements defaulting for the MappingService resource
func (ms *MappingService) Default() {
	ms.Spec.Default()
}

// +kubebuilder:object:root=true

// MappingServiceList contains a list of MappingService
type MappingServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MappingService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MappingService{}, &MappingServiceList{})
}
