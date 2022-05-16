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
	"fmt"

	"github.com/3scale/saas-operator/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	TwemproxyPodSyncLabelKey   string = fmt.Sprintf("%s/twemproxyconfig.sync", GroupVersion.Group)
	TwemproxySyncAnnotationKey string = fmt.Sprintf("%s/twemproxyconfig.configmap-hash", GroupVersion.Group)

	twemproxyDefaultGrafanaDashboard defaultGrafanaDashboardSpec = defaultGrafanaDashboardSpec{
		SelectorKey:   pointer.StringPtr("monitoring-key"),
		SelectorValue: pointer.StringPtr("middleware"),
	}
)

// TwemproxyConfigSpec defines the desired state of TwemproxyConfig
type TwemproxyConfigSpec struct {
	// SentinelURI is the redis URI of sentinel. This is required as TewmproxyConfig
	// will obtain the info about available redis servers from sentinel.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	SentinelURIs []string `json:"sentinelURIs,omitempty"`
	// ServerPools is the list of Twemproxy server pools
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ServerPools []TwemproxyServerPool `json:"serverPools"`
	// ReconcileServerPools is a flag that allows to deactivate
	// the reconcile of the contents of the managed ConfigMap. This is
	// useful in an emergency, to fix something manually. The re-sync
	// logic will still work whenever the contents of the ConfigMap
	// are changed, even if they are manually changed.
	// This switch defaults to "true".
	ReconcileServerPools *bool `json:"reconcileServerPools,omitempty"`
	// Configures the Grafana Dashboard for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	GrafanaDashboard *GrafanaDashboardSpec `json:"grafanaDashboard,omitempty"`
}

func (spec *TwemproxyConfigSpec) Default() {
	for idx := range spec.ServerPools {
		spec.ServerPools[idx].Default()
	}
	if spec.ReconcileServerPools == nil {
		spec.ReconcileServerPools = pointer.Bool(true)
	}
	spec.GrafanaDashboard = InitializeGrafanaDashboardSpec(spec.GrafanaDashboard, twemproxyDefaultGrafanaDashboard)
}

type TwemproxyServerPool struct {
	// The name of the server pool
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Name string `json:"name"`
	// The topology of the servers within the server pool. This
	// field describes the association of logical shards to physical
	// shards.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Topology []ShardedRedisTopology `json:"topology"`
	// The address to bind to. Format is ip:port
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	BindAddress string `json:"bindAddress"`
	// Timeout to stablish connection with the servers in the
	// server pool
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Timeout int `json:"timeout"`
	// Max number of pending connections in the queue
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	TCPBacklog int `json:"tcpBacklog"`
	// Connect to all servers in the pool during startup
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	PreConnect bool `json:"preConnect"`
	// Target defines which are the servers that will be configured
	// as backend redis servers for the Twemproxy configuration. The
	// default is to target masters, but read-write slaves can be configured
	// through this option. If read-write slaves are configured but there are none
	// available, the config will fall back to masters. The masters never fall back
	// to slaves though and will just wait for sentinel triggered failovers to solve
	// the unavailability.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:validation:Enum=masters;slaves-rw
	// +optional
	Target *TargetRedisServers `json:"target,omitempty"`
}

func (pool *TwemproxyServerPool) Default() {
	if pool.Target == nil {
		t := Masters
		pool.Target = &t
	}
}

type TargetRedisServers string

const (
	Masters  TargetRedisServers = "masters"
	SlavesRW TargetRedisServers = "slaves-rw"
)

type ShardedRedisTopology struct {
	// The name of the locigal shard
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ShardName string `json:"shardName"`
	// The physical shard where the logical one is stored.
	// This name should match the shard names monitored by
	// Sentinel.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	PhysicalShard string `json:"physicalShard"`
}

// TwemproxyConfigStatus defines the observed state of TwemproxyConfig
type TwemproxyConfigStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TwemproxyConfig is the Schema for the twemproxyconfigs API
type TwemproxyConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TwemproxyConfigSpec   `json:"spec,omitempty"`
	Status TwemproxyConfigStatus `json:"status,omitempty"`
}

func (tc *TwemproxyConfig) Default() {
	tc.Spec.Default()
}

func (tc *TwemproxyConfig) PodSyncSelector() client.MatchingLabels {
	return client.MatchingLabels{
		TwemproxyPodSyncLabelKey: util.ObjectKey(tc).Name,
	}
}

//+kubebuilder:object:root=true

// TwemproxyConfigList contains a list of TwemproxyConfig
type TwemproxyConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TwemproxyConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TwemproxyConfig{}, &TwemproxyConfigList{})
}
