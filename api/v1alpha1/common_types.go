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
	"reflect"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	// Finalizer is the finalizer string for resoures in the saas group
	Finalizer string = "saas.3scale.net"
)

// ImageSpec defines the image for the component
type ImageSpec struct {
	// Docker repository of the image
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Name *string `json:"name,omitempty"`
	// Image tag
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Tag *string `json:"tag,omitempty"`
	// Name of the Secret that holds quay.io credentials to access
	// the image repository
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	PullSecretName *string `json:"pullSecretName,omitempty"`
}

type defaultImageSpec struct {
	Name, Tag, PullSecretName *string
}

// Default sets default values for any value not specifically set in the ImageSpec struct
func (spec *ImageSpec) Default(def defaultImageSpec) {
	spec.Name = stringOrDefault(spec.Name, def.Name)
	spec.Tag = stringOrDefault(spec.Tag, def.Tag)
	spec.PullSecretName = stringOrDefault(spec.PullSecretName, def.PullSecretName)
}

// IsDeactivated true if the field is set with the deactivated value (empty struct)
func (spec *ImageSpec) IsDeactivated() bool { return false }

// InitializeImageSpec initializes a HTTPProbeSpec struct
func InitializeImageSpec(spec *ImageSpec, def defaultImageSpec) *ImageSpec {
	if spec == nil {
		new := &ImageSpec{}
		new.Default(def)
		return new
	}
	copy := spec.DeepCopy()
	copy.Default(def)
	return copy
}

// HTTPProbeSpec specifies configuration for an HTTP probe
type HTTPProbeSpec struct {
	// Number of seconds after the container has started before liveness probes are initiated
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	InitialDelaySeconds *int32 `json:"initialDelaySeconds,omitempty"`
	// Number of seconds after which the probe times out
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	TimeoutSeconds *int32 `json:"timeoutSeconds,omitempty"`
	// How often (in seconds) to perform the probe
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	PeriodSeconds *int32 `json:"periodSeconds,omitempty"`
	// Minimum consecutive successes for the probe to be considered successful after having failed
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	SuccessThreshold *int32 `json:"successThreshold,omitempty"`
	// Minimum consecutive failures for the probe to be considered failed after having succeeded
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	FailureThreshold *int32 `json:"failureThreshold,omitempty"`
}

type defaultHTTPProbeSpec struct {
	InitialDelaySeconds, TimeoutSeconds, PeriodSeconds,
	SuccessThreshold, FailureThreshold *int32
}

// Default sets default values for any value not specifically set in the HTTPProbeSpec struct
func (spec *HTTPProbeSpec) Default(def defaultHTTPProbeSpec) {
	spec.InitialDelaySeconds = intOrDefault(spec.InitialDelaySeconds, def.InitialDelaySeconds)
	spec.TimeoutSeconds = intOrDefault(spec.TimeoutSeconds, def.TimeoutSeconds)
	spec.PeriodSeconds = intOrDefault(spec.PeriodSeconds, def.PeriodSeconds)
	spec.SuccessThreshold = intOrDefault(spec.SuccessThreshold, def.SuccessThreshold)
	spec.FailureThreshold = intOrDefault(spec.FailureThreshold, def.FailureThreshold)
}

// IsDeactivated true if the field is set with the deactivated value (empty struct)
func (spec *HTTPProbeSpec) IsDeactivated() bool {
	if reflect.DeepEqual(spec, &HTTPProbeSpec{}) {
		return true
	}
	return false
}

// InitializeHTTPProbeSpec initializes a HTTPProbeSpec struct
func InitializeHTTPProbeSpec(spec *HTTPProbeSpec, def defaultHTTPProbeSpec) *HTTPProbeSpec {
	if spec == nil {
		new := &HTTPProbeSpec{}
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

// LoadBalancerSpec configures the AWS load balancer for the component
type LoadBalancerSpec struct {
	// Enables/disbles use of proxy protocol in the load balancer
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ProxyProtocol *bool `json:"proxyProtocol,omitempty"`
	// Enables/disables cross zone load balancing
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	CrossZoneLoadBalancingEnabled *bool `json:"crossZoneLoadBalancingEnabled,omitempty"`
	// Enables/disables connection draining
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ConnectionDrainingEnabled *bool `json:"connectionDrainingEnabled,omitempty"`
	// Sets the timeout for connection draining
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ConnectionDrainingTimeout *int32 `json:"connectionDrainingTimeout,omitempty"`
	// Sets the healthy threshold for the load balancer
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	HealthcheckHealthyThreshold *int32 `json:"healthcheckHealthyThreshold,omitempty"`
	// Sets the unhealthy threshold for the load balancer
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	HealthcheckUnhealthyThreshold *int32 `json:"healthcheckUnhealthyThreshold,omitempty"`
	// Sets the interval between health checks
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	HealthcheckInterval *int32 `json:"healthcheckInterval,omitempty"`
	// Sets the timeout for the health check
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	HealthcheckTimeout *int32 `json:"healthcheckTimeout,omitempty"`
}

type defaultLoadBalancerSpec struct {
	ProxyProtocol, CrossZoneLoadBalancingEnabled, ConnectionDrainingEnabled               *bool
	ConnectionDrainingTimeout, HealthcheckHealthyThreshold, HealthcheckUnhealthyThreshold *int32
	HealthcheckInterval, HealthcheckTimeout                                               *int32
}

// Default sets default values for any value not specifically set in the HTTPProbeSpec struct
func (spec *LoadBalancerSpec) Default(def defaultLoadBalancerSpec) {
	spec.ProxyProtocol = boolOrDefault(spec.ProxyProtocol, def.ProxyProtocol)
	spec.CrossZoneLoadBalancingEnabled = boolOrDefault(spec.CrossZoneLoadBalancingEnabled, def.CrossZoneLoadBalancingEnabled)
	spec.ConnectionDrainingEnabled = boolOrDefault(spec.ConnectionDrainingEnabled, def.ConnectionDrainingEnabled)
	spec.ConnectionDrainingTimeout = intOrDefault(spec.ConnectionDrainingTimeout, def.ConnectionDrainingTimeout)
	spec.HealthcheckHealthyThreshold = intOrDefault(spec.HealthcheckHealthyThreshold, def.HealthcheckHealthyThreshold)
	spec.HealthcheckUnhealthyThreshold = intOrDefault(spec.HealthcheckUnhealthyThreshold, def.HealthcheckUnhealthyThreshold)
	spec.HealthcheckInterval = intOrDefault(spec.HealthcheckInterval, def.HealthcheckInterval)
	spec.HealthcheckTimeout = intOrDefault(spec.HealthcheckTimeout, def.HealthcheckTimeout)
}

// IsDeactivated true if the field is set with the deactivated value (empty struct)
func (spec *LoadBalancerSpec) IsDeactivated() bool { return false }

// InitializeLoadBalancerSpec initializes a LoadBalancerSpec struct
func InitializeLoadBalancerSpec(spec *LoadBalancerSpec, def defaultLoadBalancerSpec) *LoadBalancerSpec {
	if spec == nil {
		new := &LoadBalancerSpec{}
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

// GrafanaDashboardSpec configures the Grafana Dashboard for the component
type GrafanaDashboardSpec struct {
	// Label key used by grafana-operator for dashboard discovery
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	SelectorKey *string `json:"selectorKey,omitempty"`
	// Label value used by grafana-operator for dashboard discovery
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	SelectorValue *string `json:"selectorValue,omitempty"`
}

type defaultGrafanaDashboardSpec struct {
	SelectorKey, SelectorValue *string
}

// Default sets default values for any value not specifically set in the GrafanaDashboardSpec struct
func (spec *GrafanaDashboardSpec) Default(def defaultGrafanaDashboardSpec) {
	spec.SelectorKey = stringOrDefault(spec.SelectorKey, def.SelectorKey)
	spec.SelectorValue = stringOrDefault(spec.SelectorValue, def.SelectorValue)
}

// IsDeactivated true if the field is set with the deactivated value (empty struct)
func (spec *GrafanaDashboardSpec) IsDeactivated() bool {
	if reflect.DeepEqual(spec, &GrafanaDashboardSpec{}) {
		return true
	}
	return false
}

// InitializeGrafanaDashboardSpec initializes a GrafanaDashboardSpec struct
func InitializeGrafanaDashboardSpec(spec *GrafanaDashboardSpec, def defaultGrafanaDashboardSpec) *GrafanaDashboardSpec {
	if spec == nil {
		new := &GrafanaDashboardSpec{}
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

// Endpoint sets the external endpoint for the component
type Endpoint struct {
	// The list of dns records that will point to the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	DNS []string `json:"dns"`
}

// PodDisruptionBudgetSpec defines the PDB for the component
type PodDisruptionBudgetSpec struct {
	// An eviction is allowed if at least "minAvailable" pods selected by
	// "selector" will still be available after the eviction, i.e. even in the
	// absence of the evicted pod.  So for example you can prevent all voluntary
	// evictions by specifying "100%".
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	MinAvailable *intstr.IntOrString `json:"minAvailable,omitempty"`
	// An eviction is allowed if at most "maxUnavailable" pods selected by
	// "selector" are unavailable after the eviction, i.e. even in absence of
	// the evicted pod. For example, one can prevent all voluntary evictions
	// by specifying 0. This is a mutually exclusive setting with "minAvailable".
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty"`
}

type defaultPodDisruptionBudgetSpec struct {
	MinAvailable, MaxUnavailable *intstr.IntOrString
}

// Default sets default values for any value not specifically set in the PodDisruptionBudgetSpec struct
func (spec *PodDisruptionBudgetSpec) Default(def defaultPodDisruptionBudgetSpec) {
	if spec.MinAvailable == nil && spec.MaxUnavailable == nil {
		if def.MinAvailable != nil {
			spec.MinAvailable = def.MinAvailable
			spec.MaxUnavailable = nil
		} else if def.MaxUnavailable != nil {
			spec.MinAvailable = nil
			spec.MaxUnavailable = def.MaxUnavailable
		}
	}
}

// IsDeactivated true if the field is set with the deactivated value (empty struct)
func (spec *PodDisruptionBudgetSpec) IsDeactivated() bool {
	if reflect.DeepEqual(spec, &PodDisruptionBudgetSpec{}) {
		return true
	}
	return false
}

// InitializePodDisruptionBudgetSpec initializes a PodDisruptionBudgetSpec struct
func InitializePodDisruptionBudgetSpec(spec *PodDisruptionBudgetSpec, def defaultPodDisruptionBudgetSpec) *PodDisruptionBudgetSpec {
	if spec == nil {
		new := &PodDisruptionBudgetSpec{}
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

// HorizontalPodAutoscalerSpec defines the HPA for the component
type HorizontalPodAutoscalerSpec struct {
	// Lower limit for the number of replicas to which the autoscaler
	// can scale down.  It defaults to 1 pod.  minReplicas is allowed to be 0 if the
	// alpha feature gate HPAScaleToZero is enabled and at least one Object or External
	// metric is configured.  Scaling is active as long as at least one metric value is
	// available.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	MinReplicas *int32 `json:"minReplicas,omitempty"`
	// Upper limit for the number of replicas to which the autoscaler can scale up.
	// It cannot be less that minReplicas.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	MaxReplicas *int32 `json:"maxReplicas,omitempty"`
	// Target resource used to autoscale (cpu/memory)
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:validation:Enum=cpu;memory
	// +optional
	ResourceName *string `json:"resourceName,omitempty"`
	// A percentage indicating the target resource consumption used to autoscale
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ResourceUtilization *int32 `json:"resourceUtilization,omitempty"`
}

type defaultHorizontalPodAutoscalerSpec struct {
	MinReplicas, MaxReplicas, ResourceUtilization *int32
	ResourceName                                  *string
}

// Default sets default values for any value not specifically set in the PodDisruptionBudgetSpec struct
func (spec *HorizontalPodAutoscalerSpec) Default(def defaultHorizontalPodAutoscalerSpec) {
	spec.MinReplicas = intOrDefault(spec.MinReplicas, def.MinReplicas)
	spec.MaxReplicas = intOrDefault(spec.MaxReplicas, def.MaxReplicas)
	spec.ResourceName = stringOrDefault(spec.ResourceName, def.ResourceName)
	spec.ResourceUtilization = intOrDefault(spec.ResourceUtilization, def.ResourceUtilization)
}

// IsDeactivated true if the field is set with the deactivated value (empty struct)
func (spec *HorizontalPodAutoscalerSpec) IsDeactivated() bool {
	if reflect.DeepEqual(spec, &HorizontalPodAutoscalerSpec{}) {
		return true
	}
	return false
}

// InitializeHorizontalPodAutoscalerSpec initializes a HorizontalPodAutoscalerSpec struct
func InitializeHorizontalPodAutoscalerSpec(spec *HorizontalPodAutoscalerSpec, def defaultHorizontalPodAutoscalerSpec) *HorizontalPodAutoscalerSpec {
	if spec == nil {
		new := &HorizontalPodAutoscalerSpec{}
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

// ResourceRequirementsSpec defines the resource requirements for the component
type ResourceRequirementsSpec struct {
	// Limits describes the maximum amount of compute resources allowed.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Limits corev1.ResourceList `json:"limits,omitempty"`
	// Requests describes the minimum amount of compute resources required.
	// If Requests is omitted for a container, it defaults to Limits if that is explicitly specified,
	// otherwise to an implementation-defined value.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Requests corev1.ResourceList `json:"requests,omitempty"`
}

type defaultResourceRequirementsSpec struct {
	Limits, Requests corev1.ResourceList
}

// Default sets default values for any value not specifically set in the ResourceRequirementsSpec struct
func (spec *ResourceRequirementsSpec) Default(def defaultResourceRequirementsSpec) {
	if spec.Requests == nil {
		spec.Requests = def.Requests
	}
	if spec.Limits == nil {
		spec.Limits = def.Limits
	}
}

// IsDeactivated true if the field is set with the deactivated value (empty struct)
func (spec *ResourceRequirementsSpec) IsDeactivated() bool {
	if reflect.DeepEqual(spec, &ResourceRequirementsSpec{}) {
		return true
	}
	return false
}

// InitializeResourceRequirementsSpec initializes a ResourceRequirementsSpec struct
func InitializeResourceRequirementsSpec(spec *ResourceRequirementsSpec, def defaultResourceRequirementsSpec) *ResourceRequirementsSpec {
	if spec == nil {
		new := &ResourceRequirementsSpec{}
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

func stringOrDefault(value *string, defValue *string) *string {
	if value == nil {
		return defValue
	}
	return value
}

func intOrDefault(value *int32, defValue *int32) *int32 {
	if value == nil {
		return defValue
	}
	return value
}

func boolOrDefault(value *bool, defValue *bool) *bool {
	if value == nil {
		return defValue
	}
	return value
}

func intstrPtr(value intstr.IntOrString) *intstr.IntOrString {
	return &value
}
