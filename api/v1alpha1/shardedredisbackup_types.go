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
	"reflect"
	"sort"

	"github.com/3scale/saas-operator/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	backupHistoryLimit int32 = 10
)

// ShardedRedisBackupSpec defines the desired state of ShardedRedisBackup
type ShardedRedisBackupSpec struct {
	SentinelRef string `json:"sentinelRef"`
	Schedule    string `json:"schedule"`
	Timeout     string `json:"timeout"`
	//+optional
	HistoryLimit *int32 `json:"historyLimit,omitempty"`
}

// Default implements defaulting for ShardedRedisBackuppec
func (spec *ShardedRedisBackupSpec) Default() {

	spec.HistoryLimit = intOrDefault(spec.HistoryLimit, &backupHistoryLimit)
}

// ShardedRedisBackupStatus defines the observed state of ShardedRedisBackup
type ShardedRedisBackupStatus struct {
	//+optional
	Backups BackupStatusList `json:"backups,omitempty"`
}

func (status *ShardedRedisBackupStatus) AddBackup(b BackupStatus) {
	status.Backups = append(status.Backups, b)
	sort.Sort(sort.Reverse(status.Backups))
}

func (status *ShardedRedisBackupStatus) FindLastBackup(shardName string, state BackupState) *BackupStatus {
	// backups expected to be ordered from newer to oldest
	for i, b := range status.Backups {
		if b.Shard == shardName && b.State == state {
			return &status.Backups[i]
		}
	}
	return nil
}

func (status *ShardedRedisBackupStatus) GetRunningBackups() []*BackupStatus {
	list := []*BackupStatus{}
	for i, b := range status.Backups {
		if b.State == BackupRunningState {
			list = append(list, &status.Backups[i])
		}
	}
	return list
}

func (status *ShardedRedisBackupStatus) ApplyHistoryLimit(limit int32, shards []string) bool {
	truncated := make([][]BackupStatus, len(shards))
	for idx, shard := range shards {
		var count int32 = 0
		truncated[idx] = make([]BackupStatus, 0, limit)
		for _, bs := range status.Backups {
			if count == limit {
				break
			}
			if bs.Shard == shard {
				truncated[idx] = append(truncated[idx], bs)
				count++
			}
		}
	}

	var joint BackupStatusList = util.ConcatSlices(truncated)
	sort.Sort(sort.Reverse(joint))

	if !reflect.DeepEqual(joint, status.Backups) {
		status.Backups = joint
		return true
	}

	return false
}

type BackupStatusList []BackupStatus

func (bsl BackupStatusList) Len() int { return len(bsl) }
func (bsl BackupStatusList) Less(i, j int) bool {
	a := fmt.Sprintf("%d-%s", bsl[i].ScheduledFor.UTC().UnixMilli(), bsl[i].Shard)
	b := fmt.Sprintf("%d-%s", bsl[j].ScheduledFor.UTC().UnixMilli(), bsl[j].Shard)
	return a < b
}
func (bsl BackupStatusList) Swap(i, j int) { bsl[i], bsl[j] = bsl[j], bsl[i] }

type BackupState string

type BackupStatus struct {
	Shard string `json:"shard"`
	//+optional
	ServerAlias *string `json:"serverAlias,omitempty"`
	//+optional
	ServerID     *string     `json:"serverID,omitempty"`
	ScheduledFor metav1.Time `json:"scheduledFor"`
	//+optional
	StartedAt *metav1.Time `json:"startedAt,omitempty"`
	Message   string       `json:"message"`
	State     BackupState  `json:"state"`
}

const (
	BackupPendingState   BackupState = "Pending"
	BackupRunningState   BackupState = "Running"
	BackupCompletedState BackupState = "Completed"
	BackupFailedState    BackupState = "Failed"
	BackupUnknownState   BackupState = "Unknown"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ShardedRedisBackup is the Schema for the shardedredisbackups API
type ShardedRedisBackup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ShardedRedisBackupSpec   `json:"spec,omitempty"`
	Status ShardedRedisBackupStatus `json:"status,omitempty"`
}

// Default implements defaulting for the Sentinel resource
func (srb *ShardedRedisBackup) Default() {
	srb.Spec.Default()
}

//+kubebuilder:object:root=true

// ShardedRedisBackupList contains a list of ShardedRedisBackup
type ShardedRedisBackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ShardedRedisBackup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ShardedRedisBackup{}, &ShardedRedisBackupList{})
}
