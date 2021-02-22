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
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	zyncDefaultImage defaultImageSpec = defaultImageSpec{
		Name:       pointer.StringPtr("quay.io/3scale/zync"),
		Tag:        pointer.StringPtr("nightly"),
		PullPolicy: (*corev1.PullPolicy)(pointer.StringPtr(string(corev1.PullIfNotPresent))),
	}
	zyncDefaultGrafanaDashboard defaultGrafanaDashboardSpec = defaultGrafanaDashboardSpec{
		SelectorKey:   pointer.StringPtr("monitoring-key"),
		SelectorValue: pointer.StringPtr("middleware"),
	}
	zyncDefaultConfigRailsEnvironment string                             = "development"
	zyncDefaultConfigRailsLogLevel    string                             = "info"
	zyncDefaultConfigRailsMaxThreads  int32                              = 10
	zyncDefaultConfigBugsnagSpec      BugsnagSpec                        = BugsnagSpec{}
	zyncDefaultAPIHPA                 defaultHorizontalPodAutoscalerSpec = defaultHorizontalPodAutoscalerSpec{
		MinReplicas:         pointer.Int32Ptr(2),
		MaxReplicas:         pointer.Int32Ptr(4),
		ResourceUtilization: pointer.Int32Ptr(90),
		ResourceName:        pointer.StringPtr("cpu"),
	}
	zyncDefaultAPIPDB defaultPodDisruptionBudgetSpec = defaultPodDisruptionBudgetSpec{
		MaxUnavailable: util.IntStrPtr(intstr.FromInt(1)),
	}
	zyncDefaultAPIReplicas  int32                           = 2
	zyncDefaultAPIResources defaultResourceRequirementsSpec = defaultResourceRequirementsSpec{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("250m"),
			corev1.ResourceMemory: resource.MustParse("250Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("750m"),
			corev1.ResourceMemory: resource.MustParse("512Mi"),
		},
	}
	zyncDefaultAPILivenessProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: pointer.Int32Ptr(10),
		TimeoutSeconds:      pointer.Int32Ptr(30),
		PeriodSeconds:       pointer.Int32Ptr(10),
		SuccessThreshold:    pointer.Int32Ptr(1),
		FailureThreshold:    pointer.Int32Ptr(3),
	}
	zyncDefaultAPIReadinessProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: pointer.Int32Ptr(30),
		TimeoutSeconds:      pointer.Int32Ptr(10),
		PeriodSeconds:       pointer.Int32Ptr(10),
		SuccessThreshold:    pointer.Int32Ptr(1),
		FailureThreshold:    pointer.Int32Ptr(3),
	}
	zyncDefaultQueHPA defaultHorizontalPodAutoscalerSpec = defaultHorizontalPodAutoscalerSpec{
		MinReplicas:         pointer.Int32Ptr(2),
		MaxReplicas:         pointer.Int32Ptr(4),
		ResourceUtilization: pointer.Int32Ptr(90),
		ResourceName:        pointer.StringPtr("cpu"),
	}
	zyncDefaultQuePDB defaultPodDisruptionBudgetSpec = defaultPodDisruptionBudgetSpec{
		MaxUnavailable: util.IntStrPtr(intstr.FromInt(1)),
	}
	zyncDefaultQueReplicas  int32                           = 2
	zyncDefaultQueResources defaultResourceRequirementsSpec = defaultResourceRequirementsSpec{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("250m"),
			corev1.ResourceMemory: resource.MustParse("250Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("750m"),
			corev1.ResourceMemory: resource.MustParse("512Mi"),
		},
	}
	zyncDefaultQueLivenessProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: pointer.Int32Ptr(10),
		TimeoutSeconds:      pointer.Int32Ptr(30),
		PeriodSeconds:       pointer.Int32Ptr(10),
		SuccessThreshold:    pointer.Int32Ptr(1),
		FailureThreshold:    pointer.Int32Ptr(3),
	}
	zyncDefaultQueReadinessProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: pointer.Int32Ptr(30),
		TimeoutSeconds:      pointer.Int32Ptr(10),
		PeriodSeconds:       pointer.Int32Ptr(10),
		SuccessThreshold:    pointer.Int32Ptr(1),
		FailureThreshold:    pointer.Int32Ptr(3),
	}
)

// ZyncSpec defines the desired state of Zync
type ZyncSpec struct {
	// Image specification for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Image *ImageSpec `json:"image,omitempty"`
	// Application specific configuration options for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Config ZyncConfig `json:"config"`
	// Configures the Grafana Dashboard for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	GrafanaDashboard *GrafanaDashboardSpec `json:"grafanaDashboard,omitempty"`
	// Configures the main zync api component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	API *APISpec `json:"api,omitempty"`
	// Configures the zync que component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Que *QueSpec `json:"que,omitempty"`
}

// Default implements defaulting for the Zync resource
func (z *Zync) Default() {

	z.Spec.Image = InitializeImageSpec(z.Spec.Image, zyncDefaultImage)
	z.Spec.Config.Default()
	if z.Spec.API == nil {
		z.Spec.API = &APISpec{}
	}
	z.Spec.API.Default()
	if z.Spec.Que == nil {
		z.Spec.Que = &QueSpec{}
	}
	z.Spec.Que.Default()
	z.Spec.GrafanaDashboard = InitializeGrafanaDashboardSpec(z.Spec.GrafanaDashboard, zyncDefaultGrafanaDashboard)
}

// APISpec is the configuration for main Zync api component
type APISpec struct {
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
}

// Default implements defaulting for the each main zync api component
func (spec *APISpec) Default() {

	spec.HPA = InitializeHorizontalPodAutoscalerSpec(spec.HPA, zyncDefaultAPIHPA)

	if spec.HPA.IsDeactivated() {
		spec.Replicas = intOrDefault(spec.Replicas, &zyncDefaultAPIReplicas)
	} else {
		spec.Replicas = nil
	}

	spec.PDB = InitializePodDisruptionBudgetSpec(spec.PDB, zyncDefaultAPIPDB)
	spec.Resources = InitializeResourceRequirementsSpec(spec.Resources, zyncDefaultAPIResources)
	spec.LivenessProbe = InitializeProbeSpec(spec.LivenessProbe, zyncDefaultAPILivenessProbe)
	spec.ReadinessProbe = InitializeProbeSpec(spec.ReadinessProbe, zyncDefaultAPIReadinessProbe)
}

// QueSpec is the configuration for Zync que
type QueSpec struct {
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
}

// Default implements defaulting for the each zync que
func (spec *QueSpec) Default() {

	spec.HPA = InitializeHorizontalPodAutoscalerSpec(spec.HPA, zyncDefaultQueHPA)

	if spec.HPA.IsDeactivated() {
		spec.Replicas = intOrDefault(spec.Replicas, &zyncDefaultQueReplicas)
	} else {
		spec.Replicas = nil
	}

	spec.PDB = InitializePodDisruptionBudgetSpec(spec.PDB, zyncDefaultQuePDB)
	spec.Resources = InitializeResourceRequirementsSpec(spec.Resources, zyncDefaultQueResources)
	spec.LivenessProbe = InitializeProbeSpec(spec.LivenessProbe, zyncDefaultQueLivenessProbe)
	spec.ReadinessProbe = InitializeProbeSpec(spec.ReadinessProbe, zyncDefaultQueReadinessProbe)
}

// ZyncConfig configures app behavior for Zync
type ZyncConfig struct {
	// Rails configuration options for zync components
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Rails *ZyncRailsSpec `json:"rails,omitempty"`
	// A reference to the secret holding the database DSN
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	DatabaseDSN SecretReference `json:"databaseDSN"`
	// A reference to the secret holding the secret-key-base
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	SecretKeyBase SecretReference `json:"secretKeyBase"`
	// A reference to the secret holding the zync authentication token
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ZyncAuthToken SecretReference `json:"zyncAuthToken"`
	// Options for configuring Bugsnag integration
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Bugsnag *BugsnagSpec `json:"bugsnag,omitempty"`
}

// Default sets default values for any value not specifically set in the ZyncConfig struct
func (cfg *ZyncConfig) Default() {
	if cfg.Rails == nil {
		cfg.Rails = &ZyncRailsSpec{}
	}
	if cfg.Bugsnag == nil {
		cfg.Bugsnag = &zyncDefaultConfigBugsnagSpec
	}
	cfg.Rails.Default()
}

// ZyncRailsSpec configures rails for system components
type ZyncRailsSpec struct {
	// Rails environment
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Environment *string `json:"environment,omitempty"`
	// Rails log level (debug, info, warn, error, fatal or unknown)
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:validation:Enum=debug;info;warn;error;fatal;unknown
	// +optional
	LogLevel *string `json:"logLevel,omitempty"`
	// Rails max threads (only applies to api)
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	MaxThreads *int32 `json:"maxThreads,omitempty"`
}

// Default applies defaults for ZyncRailsSpec
func (zrs *ZyncRailsSpec) Default() {
	zrs.Environment = pointer.StringPtr(zyncDefaultConfigRailsEnvironment)
	zrs.LogLevel = pointer.StringPtr(zyncDefaultConfigRailsLogLevel)
	zrs.MaxThreads = pointer.Int32Ptr(zyncDefaultConfigRailsMaxThreads)
}

// ZyncStatus defines the observed state of Zync
type ZyncStatus struct{}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Zync is the Schema for the zyncs API
type Zync struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ZyncSpec   `json:"spec,omitempty"`
	Status ZyncStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ZyncList contains a list of Zync
type ZyncList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Zync `json:"items"`
}

// GetItem returns a client.Objectfrom a ZyncList
func (bl *ZyncList) GetItem(idx int) client.Object {
	return &bl.Items[idx]
}

// CountItems returns the item count in ZyncList.Items
func (bl *ZyncList) CountItems() int {
	return len(bl.Items)
}

func init() {
	SchemeBuilder.Register(&Zync{}, &ZyncList{})
}
