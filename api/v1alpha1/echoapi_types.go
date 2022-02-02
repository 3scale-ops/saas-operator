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
	echoapiDefaultImage defaultImageSpec = defaultImageSpec{
		Name:       pointer.StringPtr("quay.io/3scale/echoapi"),
		Tag:        pointer.StringPtr("v1.0.3"),
		PullPolicy: (*corev1.PullPolicy)(pointer.StringPtr(string(corev1.PullIfNotPresent))),
	}
	echoapiDefaultReplicas int32                              = 2
	echoapiDefaultHPA      defaultHorizontalPodAutoscalerSpec = defaultHorizontalPodAutoscalerSpec{
		MinReplicas:         pointer.Int32Ptr(2),
		MaxReplicas:         pointer.Int32Ptr(4),
		ResourceUtilization: pointer.Int32Ptr(90),
		ResourceName:        pointer.StringPtr("cpu"),
	}
	echoapiDefaultPDB defaultPodDisruptionBudgetSpec = defaultPodDisruptionBudgetSpec{
		MaxUnavailable: util.IntStrPtr(intstr.FromInt(1)),
	}
	echoapiDefaultResources defaultResourceRequirementsSpec = defaultResourceRequirementsSpec{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("75m"),
			corev1.ResourceMemory: resource.MustParse("64Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("150m"),
			corev1.ResourceMemory: resource.MustParse("128Mi"),
		},
	}
	echoapiDefaultLivenessProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: pointer.Int32Ptr(25),
		TimeoutSeconds:      pointer.Int32Ptr(2),
		PeriodSeconds:       pointer.Int32Ptr(20),
		SuccessThreshold:    pointer.Int32Ptr(1),
		FailureThreshold:    pointer.Int32Ptr(5),
	}
	echoapiDefaultReadinessProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: pointer.Int32Ptr(25),
		TimeoutSeconds:      pointer.Int32Ptr(2),
		PeriodSeconds:       pointer.Int32Ptr(20),
		SuccessThreshold:    pointer.Int32Ptr(1),
		FailureThreshold:    pointer.Int32Ptr(5),
	}
	echoapiDefaultMarin3rSpec     defaultMarin3rSidecarSpec  = defaultMarin3rSidecarSpec{}
	echoapiDefaultNLBLoadBalancer defaultNLBLoadBalancerSpec = defaultNLBLoadBalancerSpec{
		ProxyProtocol:                 pointer.BoolPtr(true),
		CrossZoneLoadBalancingEnabled: pointer.BoolPtr(true),
	}
)

// EchoAPISpec defines the desired state of echoapi
type EchoAPISpec struct {
	// Image specification for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Image *ImageSpec `json:"image,omitempty"`
	// Configures the Grafana Dashboard for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
	// Resource requirements for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	HPA *HorizontalPodAutoscalerSpec `json:"hpa,omitempty"`
	// Number of replicas (ignored if hpa is enabled) for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	PDB *PodDisruptionBudgetSpec `json:"pdb,omitempty"`
	// Horizontal Pod Autoscaler for the component
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
	// Marin3r configures the Marin3r sidecars for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Marin3r *Marin3rSidecarSpec `json:"marin3r,omitempty"`
	// Configures the AWS Network load balancer for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	LoadBalancer *NLBLoadBalancerSpec `json:"loadBalancer,omitempty"`
	// The external endpoint/s for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Endpoint Endpoint `json:"endpoint"`
	// Describes node affinity scheduling rules for the pod.
	// +optional
	NodeAffinity *corev1.NodeAffinity `json:"nodeAffinity,omitempty" protobuf:"bytes,1,opt,name=nodeAffinity"`
	// If specified, the pod's tolerations.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty" protobuf:"bytes,22,opt,name=tolerations"`
}

// Default implements defaulting for EchoAPI
func (spec *EchoAPISpec) Default() {

	spec.Image = InitializeImageSpec(spec.Image, echoapiDefaultImage)
	spec.HPA = InitializeHorizontalPodAutoscalerSpec(spec.HPA, echoapiDefaultHPA)
	spec.Replicas = intOrDefault(spec.Replicas, &echoapiDefaultReplicas)
	spec.PDB = InitializePodDisruptionBudgetSpec(spec.PDB, echoapiDefaultPDB)
	spec.Resources = InitializeResourceRequirementsSpec(spec.Resources, echoapiDefaultResources)
	spec.LivenessProbe = InitializeProbeSpec(spec.LivenessProbe, echoapiDefaultLivenessProbe)
	spec.ReadinessProbe = InitializeProbeSpec(spec.ReadinessProbe, echoapiDefaultReadinessProbe)
	spec.Marin3r = InitializeMarin3rSidecarSpec(spec.Marin3r, echoapiDefaultMarin3rSpec)
	spec.LoadBalancer = InitializeNLBLoadBalancerSpec(spec.LoadBalancer, echoapiDefaultNLBLoadBalancer)
}

// EchoAPIStatus defines the observed state of EchoAPI
type EchoAPIStatus struct{}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// EchoAPI is the Schema for the echoapis API
type EchoAPI struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EchoAPISpec   `json:"spec,omitempty"`
	Status EchoAPIStatus `json:"status,omitempty"`
}

// Default implements defaulting for the EchoAPI resource
func (e *EchoAPI) Default() {
	e.Spec.Default()
}

// +kubebuilder:object:root=true

// EchoAPIList contains a list of echoapi
type EchoAPIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EchoAPI `json:"items"`
}

// GetItem returns a client.Objectfrom a EchoAPIList
func (bl *EchoAPIList) GetItem(idx int) client.Object {
	return &bl.Items[idx]
}

// CountItems returns the item count in EchoAPIList.Items
func (bl *EchoAPIList) CountItems() int {
	return len(bl.Items)
}

func init() {
	SchemeBuilder.Register(&EchoAPI{}, &EchoAPIList{})
}
