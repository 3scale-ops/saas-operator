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
	"strconv"
	"strings"

	"github.com/3scale/saas-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

var (
	redisShardDefaultImage defaultImageSpec = defaultImageSpec{
		Name:       pointer.String("redis"),
		Tag:        pointer.String("4.0.11-alpine"),
		PullPolicy: (*corev1.PullPolicy)(pointer.String(string(corev1.PullIfNotPresent))),
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
	// SlaveCount is the number of redis slaves
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	SlaveCount *int32 `json:"slaveCount,omitempty"`
}

// Default implements defaulting for RedisShardSpec
func (spec *RedisShardSpec) Default() {

	spec.Image = InitializeImageSpec(spec.Image, redisShardDefaultImage)
	spec.MasterIndex = intOrDefault(spec.MasterIndex, &redisShardDefaultMasterIndex)
	spec.SlaveCount = intOrDefault(spec.SlaveCount, pointer.Int32(RedisShardDefaultReplicas-1))
}

type RedisShardNodes struct {
	// Master is the node that acts as master role in the redis shard
	// +operator-sdk:csv:customresourcedefinitions:type=status
	// +optional
	Master map[string]string `json:"master,omitempty"`
	// Slaves are the nodes that act as master role in the redis shard
	// +operator-sdk:csv:customresourcedefinitions:type=status
	// +optional
	Slaves map[string]string `json:"slaves,omitempty"`
}

func (rsn *RedisShardNodes) MasterHostPort() string {
	for _, hostport := range rsn.Master {
		return hostport
	}
	return ""
}

func (rsn *RedisShardNodes) GetNodeByPodIndex(podIndex int) (string, string) {
	nodes := util.MergeMaps(map[string]string{}, rsn.Master, rsn.Slaves)

	for alias, hostport := range nodes {
		i := alias[strings.LastIndex(alias, "-")+1:]
		index, _ := strconv.Atoi(i)
		if index == podIndex {
			return alias, hostport
		}
	}

	return "", ""
}

func (rsn *RedisShardNodes) GetHostPortByPodIndex(podIndex int) string {
	_, hostport := rsn.GetNodeByPodIndex(podIndex)
	return hostport
}

func (rsn *RedisShardNodes) GetAliasByPodIndex(podIndex int) string {
	alias, _ := rsn.GetNodeByPodIndex(podIndex)
	return alias
}

func (rsn *RedisShardNodes) GetIndexByHostPort(hostport string) int {
	nodes := util.MergeMaps(map[string]string{}, rsn.Master, rsn.Slaves)
	for alias, hp := range nodes {
		if hostport == hp {
			i := alias[strings.LastIndex(alias, "-")+1:]
			index, _ := strconv.Atoi(i)
			return index
		}
	}
	return -1
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

// Default implements defaulting for the RedisShard resource
func (rs *RedisShard) Default() {
	rs.Spec.Default()
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
