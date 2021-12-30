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
	"time"

	"github.com/3scale/saas-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

const (
	// SentinelPort is the port where sentinel process listens
	SentinelPort uint32 = 26379
)

// bitnami/redis-sentinel:4.0.11-debian-9-r110
var (
	sentinelDefaultReplicas int32            = 3
	sentinelDefaultImage    defaultImageSpec = defaultImageSpec{
		Name:       pointer.StringPtr("bitnami/redis-sentinel"),
		Tag:        pointer.StringPtr("4.0.11-debian-9-r110"),
		PullPolicy: (*corev1.PullPolicy)(pointer.StringPtr(string(corev1.PullIfNotPresent))),
	}
	sentinelDefaultResources defaultResourceRequirementsSpec = defaultResourceRequirementsSpec{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("75m"),
			corev1.ResourceMemory: resource.MustParse("64Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("150m"),
			corev1.ResourceMemory: resource.MustParse("128Mi"),
		},
	}
	sentinelDefaultProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: pointer.Int32Ptr(25),
		TimeoutSeconds:      pointer.Int32Ptr(1),
		PeriodSeconds:       pointer.Int32Ptr(10),
		SuccessThreshold:    pointer.Int32Ptr(1),
		FailureThreshold:    pointer.Int32Ptr(3),
	}
	sentinelDefaultPDB defaultPodDisruptionBudgetSpec = defaultPodDisruptionBudgetSpec{
		MaxUnavailable: util.IntStrPtr(intstr.FromInt(1)),
	}

	sentinelDefaultGrafanaDashboard defaultGrafanaDashboardSpec = defaultGrafanaDashboardSpec{
		SelectorKey:   pointer.StringPtr("monitoring-key"),
		SelectorValue: pointer.StringPtr("middleware"),
	}
	sentinelDefaultStorageSize            string        = "10Mi"
	sentinelDefaultMetricsRefreshInterval time.Duration = 30 * time.Second
)

// SentinelConfig defines configuration options for the component
type SentinelConfig struct {
	// Monitored shards indicates the redis servers that form
	// part of each shard monitored by sentinel
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	MonitoredShards map[string][]string `json:"monitoredShards,"`
	// StorageClass is the storage class to be used for
	// the persistent sentinel config file where the shards
	// state is stored
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	StorageClass *string `json:"storageClass,omitempty"`
	// StorageSize is the storage size to  provision for
	// the persistent sentinel config file where the shards
	// state is stored
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	StorageSize *resource.Quantity `json:"storageSize,omitempty"`
	// MetricsRefreshInterval determines the refresh interval for gahtering
	// metrics from sentinel
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	MetricsRefreshInterval *time.Duration `json:"metricsRefreshInterval,omitempty"`
}

// Default sets default values for any value not specifically set in the AutoSSLConfig struct
func (cfg *SentinelConfig) Default() {
	if cfg.StorageSize == nil {
		size := resource.MustParse(sentinelDefaultStorageSize)
		cfg.StorageSize = &size
	}

	if cfg.MetricsRefreshInterval == nil {
		cfg.MetricsRefreshInterval = &sentinelDefaultMetricsRefreshInterval
	}
}

// SentinelSpec defines the desired state of Sentinel
type SentinelSpec struct {
	// Image specification for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Image *ImageSpec `json:"image,omitempty"`
	// Number of replicas (ignored if hpa is enabled) for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
	// Pod Disruption Budget for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	PDB *PodDisruptionBudgetSpec `json:"pdb,omitempty"`
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
	// Describes node affinity scheduling rules for the pod.
	// +optional
	NodeAffinity *corev1.NodeAffinity `json:"nodeAffinity,omitempty" protobuf:"bytes,1,opt,name=nodeAffinity"`
	// If specified, the pod's tolerations.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty" protobuf:"bytes,22,opt,name=tolerations"`
	// Config configures the sentinel process
	Config *SentinelConfig `json:"config"`
}

// Default implements defaulting for the Sentinel resource
func (s *Sentinel) Default() {

	s.Spec.Image = InitializeImageSpec(s.Spec.Image, sentinelDefaultImage)
	s.Spec.Replicas = intOrDefault(s.Spec.Replicas, &sentinelDefaultReplicas)
	s.Spec.PDB = InitializePodDisruptionBudgetSpec(s.Spec.PDB, sentinelDefaultPDB)
	s.Spec.Resources = InitializeResourceRequirementsSpec(s.Spec.Resources, sentinelDefaultResources)
	s.Spec.LivenessProbe = InitializeProbeSpec(s.Spec.LivenessProbe, sentinelDefaultProbe)
	s.Spec.ReadinessProbe = InitializeProbeSpec(s.Spec.ReadinessProbe, sentinelDefaultProbe)
	s.Spec.GrafanaDashboard = InitializeGrafanaDashboardSpec(s.Spec.GrafanaDashboard, sentinelDefaultGrafanaDashboard)
	s.Spec.Config.Default()
}

// SentinelStatus defines the observed state of Sentinel
type SentinelStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Sentinel is the Schema for the sentinels API
type Sentinel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SentinelSpec   `json:"spec,omitempty"`
	Status SentinelStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SentinelList contains a list of Sentinel
type SentinelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Sentinel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Sentinel{}, &SentinelList{})
}