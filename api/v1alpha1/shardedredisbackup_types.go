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
	"time"

	"github.com/3scale/saas-operator/pkg/util"
	"github.com/dustin/go-humanize"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	AWSAccessKeyID_SecretKey     string = "AWS_ACCESS_KEY_ID"
	AWSSecretAccessKey_SecretKey string = "AWS_SECRET_ACCESS_KEY"
	BackupFile                   string = "redis_backup.rdb"

	// defaults
	backupHistoryLimit        int32  = 10
	backupDefaultTimeout      string = "10m"
	backupDefaultPollInterval string = "60s"
	backupDefaultSSHPort      uint32 = 22
	backupDefaultMinSize      string = "1 GB"
)

// ShardedRedisBackupSpec defines the desired state of ShardedRedisBackup
type ShardedRedisBackupSpec struct {
	SentinelRef string     `json:"sentinelRef"`
	Schedule    string     `json:"schedule"`
	DBFile      string     `json:"dbFile"`
	SSHOptions  SSHOptions `json:"sshOptions"`
	S3Options   S3Options  `json:"s3Options"`
	//+optional
	Timeout *metav1.Duration `json:"timeout"`
	//+optional
	HistoryLimit *int32 `json:"historyLimit,omitempty"`
	//+optional
	PollInterval *metav1.Duration `json:"pollInterval,omitempty"`
	// +optional
	MinSize *string `json:"minSize,omitempty"`
}

// Default implements defaulting for ShardedRedisBackuppec
func (spec *ShardedRedisBackupSpec) Default() {

	if spec.Timeout == nil {
		d, _ := time.ParseDuration(backupDefaultTimeout)
		spec.Timeout = &metav1.Duration{Duration: d}
	}
	if spec.PollInterval == nil {
		d, _ := time.ParseDuration(backupDefaultPollInterval)
		spec.PollInterval = &metav1.Duration{Duration: d}
	}
	spec.HistoryLimit = intOrDefault(spec.HistoryLimit, util.Pointer(backupHistoryLimit))
	spec.MinSize = stringOrDefault(spec.MinSize, util.Pointer(backupDefaultMinSize))
	spec.SSHOptions.Default()
}

func (spec *ShardedRedisBackupSpec) GetMinSize() (uint64, error) {
	if spec.MinSize == nil {
		return humanize.ParseBytes(backupDefaultMinSize)
	} else {
		return humanize.ParseBytes(*spec.MinSize)
	}
}

type SSHOptions struct {
	User                string                      `json:"user"`
	PrivateKeySecretRef corev1.LocalObjectReference `json:"privateKeySecretRef"`
	// +optional
	Port *uint32 `json:"port,omitempty"`
}

func (opts *SSHOptions) Default() {
	if opts.Port == nil {
		opts.Port = util.Pointer(backupDefaultSSHPort)
	}
}

type S3Options struct {
	Bucket               string                      `json:"bucket"`
	Path                 string                      `json:"path"`
	Region               string                      `json:"region"`
	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef"`
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
