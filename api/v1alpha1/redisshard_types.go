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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

var (
	redisShardDefaultImage defaultImageSpec = defaultImageSpec{
		Name:       pointer.StringPtr("redis"),
		Tag:        pointer.StringPtr("4.0.11-alpine"),
		PullPolicy: (*corev1.PullPolicy)(pointer.StringPtr(string(corev1.PullIfNotPresent))),
	}
	redisShardDefaultMasterIndex int32 = 0
	RedisShardDefaultReplicas    int32 = 3
)

// RedisShardSpec defines the desired state of RedisShard
type RedisShardSpec struct {
	// Image specification for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Image *ImageSpec `json:"image,omitempty"`
	// MasterIndex is the StatefulSet Pod index of the redis server
	// with the master role. The other Pods are slaves of the master one.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	MasterIndex *int32 `json:"masterIndex,omitempty"`
}

// Default implements defaulting for the Sentinel resource
func (rs *RedisShard) Default() {

	rs.Spec.Image = InitializeImageSpec(rs.Spec.Image, redisShardDefaultImage)
	rs.Spec.MasterIndex = intOrDefault(rs.Spec.MasterIndex, &redisShardDefaultMasterIndex)
}

type RedisShardNodes struct {
	// Master is the node that acts as master role in the redis shard
	// +operator-sdk:csv:customresourcedefinitions:type=status
	// +optional
	Master *string `json:"master,omitempty"`
	// Slaves are the nodes that act as master role in the redis shard
	// +operator-sdk:csv:customresourcedefinitions:type=status
	// +optional
	Slaves []string `json:"slaves,omitempty"`
}

// RedisShardStatus defines the observed state of RedisShard
type RedisShardStatus struct {
	// ShardNodes describes the nodes in the redis shard
	// +operator-sdk:csv:customresourcedefinitions:type=status
	// +optional
	ShardNodes *RedisShardNodes `json:"shardNodes,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RedisShard is the Schema for the redisshards API
// +kubebuilder:printcolumn:JSONPath=".status.shardNodes.master",name=Master,type=string
// +kubebuilder:printcolumn:JSONPath=".status.shardNodes.slaves",name=Slaves,type=string
type RedisShard struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RedisShardSpec   `json:"spec,omitempty"`
	Status RedisShardStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RedisShardList contains a list of RedisShard
type RedisShardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RedisShard `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RedisShard{}, &RedisShardList{})
}
