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
	"k8s.io/apimachinery/pkg/util/intstr"
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

// HTTPProbeSpec specifies configuration for an HTTP probe
type HTTPProbeSpec struct {
	// Number of seconds after the container has started before liveness probes are initiated
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	InitialDelaySeconds *int `json:"initialDelaySeconds,omitempty"`
	// Number of seconds after which the probe times out
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	TimeoutSeconds *int `json:"timeoutSeconds,omitempty"`
	// How often (in seconds) to perform the probe
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	PeriodSeconds *int `json:"periodSeconds,omitempty"`
	// Minimum consecutive successes for the probe to be considered successful after having failed
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	SuccessThreshold *int `json:"successThreshold,omitempty"`
	// Minimum consecutive failures for the probe to be considered failed after having succeeded
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	FailureThreshold *int `json:"failureThreshold,omitempty"`
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
	ConnectionDrainingTimeout *int `json:"connectionDrainingTimeout,omitempty"`
	// Sets the healthy threshold for the load balancer
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ConnectionHealthcheckHealthyThreshold *int `json:"connectionHealthcheckHealthyThreshold,omitempty"`
	// Sets the unhealthy threshold for the load balancer
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ConnectionHealthcheckUnhealthyThreshold *int `json:"connectionHealthcheckUnhealthyThreshold,omitempty"`
	// Sets the interval between health checks
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ConnectionHealthcheckInterval *int `json:"connectionHealthcheckInterval,omitempty"`
	// Sets the timeout for the health check
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ConnectionHealthcheckTimeout *int `json:"connectionHealthcheckTimeout,omitempty"`
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
	ResourceUtilization *int `json:"resourceUtilization,omitempty"`
}
