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
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestBackupStatusList_Less(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		bsl  BackupStatusList
		args args
		want bool
	}{
		{
			name: "0 < 1",
			bsl: []BackupStatus{
				{
					Shard:        "shard01",
					ScheduledFor: metav1.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					Shard:        "shard02",
					ScheduledFor: metav1.Date(2023, time.August, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			args: args{i: 0, j: 1},
			want: true,
		},
		{
			name: "0 < 1 (same timestamp)",
			bsl: []BackupStatus{
				{
					Shard:        "shard01",
					ScheduledFor: metav1.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					Shard:        "shard02",
					ScheduledFor: metav1.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			args: args{i: 0, j: 1},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bsl.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("BackupStatusList.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShardedRedisBackupStatus_FindBackup(t *testing.T) {
	type fields struct {
		Backups BackupStatusList
	}
	type args struct {
		shard string
		state BackupState
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *BackupStatus
		wantPos int
	}{
		{
			name: "Returns latest pending backup",
			fields: fields{
				Backups: []BackupStatus{
					{
						Shard:        "shard02",
						ScheduledFor: metav1.Date(2023, time.August, 1, 0, 1, 0, 0, time.UTC),
						State:        BackupPendingState,
					},
					{
						Shard:        "shard01",
						ScheduledFor: metav1.Date(2023, time.August, 1, 0, 1, 0, 0, time.UTC),
						State:        BackupPendingState,
					},
					{
						Shard:        "shard02",
						ScheduledFor: metav1.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC),
						State:        BackupRunningState,
					},
					{
						Shard:        "shard01",
						ScheduledFor: metav1.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC),
						State:        BackupRunningState,
					},
				},
			},
			args: args{shard: "shard02", state: BackupPendingState},
			want: &BackupStatus{
				Shard:        "shard02",
				ScheduledFor: metav1.Date(2023, time.August, 1, 0, 1, 0, 0, time.UTC),
				State:        BackupPendingState,
			},
			wantPos: 0,
		},
		{
			name: "Returns latest running backup",
			fields: fields{
				Backups: []BackupStatus{
					{
						Shard:        "shard02",
						ScheduledFor: metav1.Date(2023, time.August, 1, 0, 1, 0, 0, time.UTC),
						State:        BackupPendingState,
					},
					{
						Shard:        "shard01",
						ScheduledFor: metav1.Date(2023, time.August, 1, 0, 1, 0, 0, time.UTC),
						State:        BackupPendingState,
					},
					{
						Shard:        "shard02",
						ScheduledFor: metav1.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC),
						State:        BackupRunningState,
					},
					{
						Shard:        "shard01",
						ScheduledFor: metav1.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC),
						State:        BackupRunningState,
					},
				},
			},
			args: args{shard: "shard02", state: BackupRunningState},
			want: &BackupStatus{
				Shard:        "shard02",
				ScheduledFor: metav1.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC),
				State:        BackupRunningState,
			},
			wantPos: 2,
		},
		{
			name: "Returns a nil",
			fields: fields{
				Backups: []BackupStatus{
					{
						Shard:        "shard02",
						ScheduledFor: metav1.Date(2023, time.August, 1, 0, 1, 0, 0, time.UTC),
						State:        BackupPendingState,
					},
				},
			},
			args:    args{shard: "shard02", state: BackupRunningState},
			want:    nil,
			wantPos: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := &ShardedRedisBackupStatus{
				Backups: tt.fields.Backups,
			}
			gotBackup, gotPos := status.FindLastBackup(tt.args.shard, tt.args.state)
			if !reflect.DeepEqual(gotBackup, tt.want) {
				t.Errorf("ShardedRedisBackupStatus.FindBackup() = %v, want %v", gotBackup, tt.want)
			}
			if gotPos != tt.wantPos {
				t.Errorf("ShardedRedisBackupStatus.FindBackup() = %v, want %v", gotPos, tt.wantPos)
			}
		})
	}
}

func TestShardedRedisBackupStatus_DeleteBackup(t *testing.T) {
	type fields struct {
		Backups BackupStatusList
	}
	type args struct {
		pos int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   BackupStatusList
	}{
		{
			name: "Deletes the backup at the given position from the list",
			fields: fields{
				Backups: BackupStatusList{
					{
						Shard:        "shard02",
						ScheduledFor: metav1.Date(2023, time.August, 1, 0, 1, 0, 0, time.UTC),
						State:        BackupPendingState,
					},
					{
						Shard:        "shard01",
						ScheduledFor: metav1.Date(2023, time.August, 1, 0, 1, 0, 0, time.UTC),
						State:        BackupPendingState,
					},
					{
						Shard:        "shard02",
						ScheduledFor: metav1.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC),
						State:        BackupRunningState,
					},
					{
						Shard:        "shard01",
						ScheduledFor: metav1.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC),
						State:        BackupRunningState,
					},
				},
			},
			args: args{
				pos: 2,
			},
			want: BackupStatusList{
				{
					Shard:        "shard02",
					ScheduledFor: metav1.Date(2023, time.August, 1, 0, 1, 0, 0, time.UTC),
					State:        BackupPendingState,
				},
				{
					Shard:        "shard01",
					ScheduledFor: metav1.Date(2023, time.August, 1, 0, 1, 0, 0, time.UTC),
					State:        BackupPendingState,
				},
				{
					Shard:        "shard01",
					ScheduledFor: metav1.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC),
					State:        BackupRunningState,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := &ShardedRedisBackupStatus{
				Backups: tt.fields.Backups,
			}
			status.DeleteBackup(tt.args.pos)
			if !reflect.DeepEqual(status.Backups, tt.want) {
				t.Errorf("ShardedRedisBackupStatus.FindBackup() = %v, want %v", status.Backups, tt.want)
			}
		})
	}
}
