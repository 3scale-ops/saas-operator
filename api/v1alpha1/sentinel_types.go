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
	"context"
	"sort"
	"time"

	"github.com/3scale-ops/basereconciler/util"
	"github.com/3scale-ops/saas-operator/pkg/redis/client"
	redis "github.com/3scale-ops/saas-operator/pkg/redis/server"
	"github.com/3scale-ops/saas-operator/pkg/redis/sharded"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	// SentinelPort is the port where sentinel process listens
	SentinelPort uint32 = 26379
)

// bitnami/redis-sentinel:4.0.11-debian-9-r110
var (
	SentinelDefaultReplicas int32            = 3
	SentinelDefaultQuorum   int              = 2
	sentinelDefaultImage    defaultImageSpec = defaultImageSpec{
		Name:       util.Pointer("bitnami/redis-sentinel"),
		Tag:        util.Pointer("4.0.11-debian-9-r110"),
		PullPolicy: (*corev1.PullPolicy)(util.Pointer(string(corev1.PullIfNotPresent))),
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
		InitialDelaySeconds: util.Pointer[int32](25),
		TimeoutSeconds:      util.Pointer[int32](1),
		PeriodSeconds:       util.Pointer[int32](10),
		SuccessThreshold:    util.Pointer[int32](1),
		FailureThreshold:    util.Pointer[int32](3),
	}
	sentinelDefaultPDB defaultPodDisruptionBudgetSpec = defaultPodDisruptionBudgetSpec{
		MaxUnavailable: util.Pointer(intstr.FromInt(1)),
	}

	sentinelDefaultGrafanaDashboard defaultGrafanaDashboardSpec = defaultGrafanaDashboardSpec{
		SelectorKey:   util.Pointer("monitoring-key"),
		SelectorValue: util.Pointer("middleware"),
	}
	sentinelDefaultStorageSize            string        = "10Mi"
	sentinelDefaultMetricsRefreshInterval time.Duration = 30 * time.Second
)

// SentinelConfig defines configuration options for the component
type SentinelConfig struct {
	// Monitored shards indicates the redis servers that form
	// part of each shard monitored by sentinel
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	MonitoredShards map[string][]string `json:"monitoredShards,"`
	// ClusterTopology indicates the redis servers that form
	// part of each shard monitored by sentinel
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ClusterTopology map[string]map[string]string `json:"clusterTopology,"`
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

// Default implements defaulting for SentinelSpec
func (spec *SentinelSpec) Default() {

	spec.Image = InitializeImageSpec(spec.Image, sentinelDefaultImage)
	spec.Replicas = intOrDefault(spec.Replicas, &SentinelDefaultReplicas)
	spec.PDB = InitializePodDisruptionBudgetSpec(spec.PDB, sentinelDefaultPDB)
	spec.Resources = InitializeResourceRequirementsSpec(spec.Resources, sentinelDefaultResources)
	spec.LivenessProbe = InitializeProbeSpec(spec.LivenessProbe, sentinelDefaultProbe)
	spec.ReadinessProbe = InitializeProbeSpec(spec.ReadinessProbe, sentinelDefaultProbe)
	spec.GrafanaDashboard = InitializeGrafanaDashboardSpec(spec.GrafanaDashboard, sentinelDefaultGrafanaDashboard)
	spec.Config.Default()
}

// SentinelStatus defines the observed state of Sentinel
type SentinelStatus struct {
	// Addresses of the sentinel instances currently running
	// +operator-sdk:csv:customresourcedefinitions:type=status
	// +optional
	Sentinels []string `json:"sentinels,omitempty"`
	// MonitoredShards is the list of shards that the Sentinel
	// resource is currently monitoring
	// +operator-sdk:csv:customresourcedefinitions:type=status
	// +optional
	MonitoredShards MonitoredShards `json:"monitoredShards,omitempty"`
}

// ShardedCluster returns a *sharded.Cluster struct from the information reported by the sentinel status instead
// of directly contacting sentinel/redis to gather the state of the cluster. This avoids calls to sentinel/redis
// but is less robust as it depends entirely on the Sentinel controller working properly and without delays.
// As of now, this is used in the SharededRedisBackup controller but not in the TwemproxyConfig controller.
func (ss *SentinelStatus) ShardedCluster(ctx context.Context, pool *redis.ServerPool) (*sharded.Cluster, error) {

	// have a list of sentinels but must provide a map
	// TODO: at some point change the SentinelStatus.Sentinels to also have a map and avoid this
	msentinel := make(map[string]string, len(ss.Sentinels))
	for _, s := range ss.Sentinels {
		msentinel[s] = "redis://" + s
	}

	shards := make([]*sharded.Shard, 0, len(ss.MonitoredShards))
	// generate slice of shards from status
	for _, s := range ss.MonitoredShards {
		servers := make([]*sharded.RedisServer, 0, len(s.Servers))
		for _, rsd := range s.Servers {
			srv, err := pool.GetServer("redis://"+rsd.Address, nil)
			if err != nil {
				return nil, err
			}
			servers = append(servers, sharded.NewRedisServerFromParams(srv, rsd.Role, rsd.Config))
		}
		sort.Slice(servers, func(i, j int) bool {
			return servers[i].ID() < servers[j].ID()
		})
		shards = append(shards, sharded.NewShardFromServers(s.Name, pool, servers...))
	}

	cluster, err := sharded.NewShardedCluster(ctx, pool, msentinel, shards...)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

type MonitoredShards []MonitoredShard

// MonitoredShards implements sort.Interface based on the Name field.
func (ms MonitoredShards) Len() int           { return len(ms) }
func (ms MonitoredShards) Less(i, j int) bool { return ms[i].Name < ms[j].Name }
func (ms MonitoredShards) Swap(i, j int)      { ms[i], ms[j] = ms[j], ms[i] }

// MonitoredShard contains information of one of the shards
// monitored by the Sentinel resource
type MonitoredShard struct {
	// Name is the name of the redis shard
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Name string `json:"name"`
	// Server is a map intended to store configuration information
	// of each of the RedisServer instances that belong to the MonitoredShard
	// +operator-sdk:csv:customresourcedefinitions:type=status
	// +optional
	Servers map[string]RedisServerDetails `json:"servers,omitempty"`
}

type RedisServerDetails struct {
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Role client.Role `json:"role"`
	// +operator-sdk:csv:customresourcedefinitions:type=status
	// +optional
	Address string `json:"address,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=status
	// +optional
	Config map[string]string `json:"config,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=status
	// +optional
	Info map[string]string `json:"info,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:JSONPath=".status.sentinels",name=Sentinels,type=string
// +kubebuilder:printcolumn:JSONPath=".status.monitoredShards",name=Shards,type=string
// Sentinel is the Schema for the sentinels API
type Sentinel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SentinelSpec   `json:"spec,omitempty"`
	Status SentinelStatus `json:"status,omitempty"`
}

// Default implements defaulting for the Sentinel resource
func (s *Sentinel) Default() {
	s.Spec.Default()
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
