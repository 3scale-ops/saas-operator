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
	corsproxyDefaultReplicas int32            = 2
	corsproxyDefaultImage    defaultImageSpec = defaultImageSpec{
		Name:       pointer.StringPtr("quay.io/3scale/cors-proxy"),
		Tag:        pointer.StringPtr("latest"),
		PullPolicy: (*corev1.PullPolicy)(pointer.StringPtr(string(corev1.PullIfNotPresent))),
	}
	corsproxyDefaultResources defaultResourceRequirementsSpec = defaultResourceRequirementsSpec{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("75m"),
			corev1.ResourceMemory: resource.MustParse("64Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("150m"),
			corev1.ResourceMemory: resource.MustParse("128Mi"),
		},
	}
	corsproxyDefaultHPA defaultHorizontalPodAutoscalerSpec = defaultHorizontalPodAutoscalerSpec{
		MinReplicas:         pointer.Int32Ptr(2),
		MaxReplicas:         pointer.Int32Ptr(4),
		ResourceUtilization: pointer.Int32Ptr(90),
		ResourceName:        pointer.StringPtr("cpu"),
	}
	corsproxyDefaultProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: pointer.Int32Ptr(3),
		TimeoutSeconds:      pointer.Int32Ptr(1),
		PeriodSeconds:       pointer.Int32Ptr(10),
		SuccessThreshold:    pointer.Int32Ptr(1),
		FailureThreshold:    pointer.Int32Ptr(3),
	}
	corsproxyDefaultPDB defaultPodDisruptionBudgetSpec = defaultPodDisruptionBudgetSpec{
		MaxUnavailable: util.IntStrPtr(intstr.FromInt(1)),
	}

	corsproxyDefaultGrafanaDashboard defaultGrafanaDashboardSpec = defaultGrafanaDashboardSpec{
		SelectorKey:   pointer.StringPtr("monitoring-key"),
		SelectorValue: pointer.StringPtr("middleware"),
	}
)

// CORSProxySpec defines the desired state of CORSProxy
type CORSProxySpec struct {
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
	Config CORSProxyConfig `json:"config"`
}

// Default implements defaulting for the CORSProxy resource
func (a *CORSProxy) Default() {

	a.Spec.Image = InitializeImageSpec(a.Spec.Image, corsproxyDefaultImage)
	a.Spec.HPA = InitializeHorizontalPodAutoscalerSpec(a.Spec.HPA, corsproxyDefaultHPA)

	if a.Spec.HPA.IsDeactivated() {
		a.Spec.Replicas = intOrDefault(a.Spec.Replicas, &corsproxyDefaultReplicas)
	} else {
		a.Spec.Replicas = nil
	}

	a.Spec.PDB = InitializePodDisruptionBudgetSpec(a.Spec.PDB, corsproxyDefaultPDB)
	a.Spec.Resources = InitializeResourceRequirementsSpec(a.Spec.Resources, corsproxyDefaultResources)
	a.Spec.LivenessProbe = InitializeProbeSpec(a.Spec.LivenessProbe, corsproxyDefaultProbe)
	a.Spec.ReadinessProbe = InitializeProbeSpec(a.Spec.ReadinessProbe, corsproxyDefaultProbe)
	a.Spec.GrafanaDashboard = InitializeGrafanaDashboardSpec(a.Spec.GrafanaDashboard, corsproxyDefaultGrafanaDashboard)
	a.Spec.Config.Default()
}

// CORSProxyConfig defines configuration options for the component
type CORSProxyConfig struct {
	// System database connection string
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	SystemDatabaseDSN SecretReference `json:"systemDatabaseDSN"`
}

// Default sets default values for any value not specifically set in the CORSProxyConfig struct
func (cfg *CORSProxyConfig) Default() {}

// CORSProxyStatus defines the observed state of CORSProxy
type CORSProxyStatus struct{}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// CORSProxy is the Schema for the corsproxies API
type CORSProxy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CORSProxySpec   `json:"spec,omitempty"`
	Status CORSProxyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CORSProxyList contains a list of CORSProxy
type CORSProxyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CORSProxy `json:"items"`
}

// GetItem returns a client.Objectfrom a CORSProxyList
func (cpl *CORSProxyList) GetItem(idx int) client.Object {
	return &cpl.Items[idx]
}

// CountItems returns the item count in CORSProxyList.Items
func (cpl *CORSProxyList) CountItems() int {
	return len(cpl.Items)
}

func init() {
	SchemeBuilder.Register(&CORSProxy{}, &CORSProxyList{})
}
