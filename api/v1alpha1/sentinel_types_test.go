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
	"testing"

	"github.com/3scale/saas-operator/pkg/redis/client"
	redis "github.com/3scale/saas-operator/pkg/redis/server"
	"github.com/3scale/saas-operator/pkg/redis/sharded"
	"github.com/3scale/saas-operator/pkg/util"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSentinelStatus_ShardedCluster(t *testing.T) {
	type fields struct {
		Sentinels       []string
		MonitoredShards MonitoredShards
	}
	type args struct {
		ctx  context.Context
		pool *redis.ServerPool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *sharded.Cluster
		wantErr bool
	}{
		{
			name: "Generates a sharded.Cluster resource from the sentinel status",
			fields: fields{
				Sentinels: []string{"127.0.0.1:26379"},
				MonitoredShards: []MonitoredShard{
					{Name: "shard01",
						Servers: map[string]RedisServerDetails{
							"srv1": {Role: client.Master, Address: "127.0.0.1:1000", Config: map[string]string{"save": ""}},
							"srv2": {Role: client.Slave, Address: "127.0.0.1:2000", Config: map[string]string{"slave-read-only": "yes"}},
						}},
					{Name: "shard02",
						Servers: map[string]RedisServerDetails{
							"srv3": {Role: client.Master, Address: "127.0.0.1:3000", Config: map[string]string{}},
							"srv4": {Role: client.Slave, Address: "127.0.0.1:4000", Config: map[string]string{}},
						}},
				},
			},
			args: args{
				ctx: context.TODO(),
				pool: redis.NewServerPool(
					redis.MustNewServer("redis://127.0.0.1:1000", util.Pointer("srv1")),
					redis.MustNewServer("redis://127.0.0.1:2000", util.Pointer("srv2")),
					redis.MustNewServer("redis://127.0.0.1:3000", util.Pointer("srv3")),
					redis.MustNewServer("redis://127.0.0.1:4000", util.Pointer("srv4")),
					redis.MustNewServer("redis://127.0.0.1:26379", util.Pointer("sentinel")),
				),
			},
			want: &sharded.Cluster{
				Shards: []*sharded.Shard{
					{Name: "shard01",
						Servers: []*sharded.RedisServer{
							sharded.NewRedisServerFromParams(redis.MustNewServer("redis://127.0.0.1:1000", util.Pointer("srv1")), client.Master, map[string]string{"save": ""}),
							sharded.NewRedisServerFromParams(redis.MustNewServer("redis://127.0.0.1:2000", util.Pointer("srv2")), client.Slave, map[string]string{"slave-read-only": "yes"}),
						}},
					{Name: "shard02",
						Servers: []*sharded.RedisServer{
							sharded.NewRedisServerFromParams(redis.MustNewServer("redis://127.0.0.1:3000", util.Pointer("srv3")), client.Master, map[string]string{}),
							sharded.NewRedisServerFromParams(redis.MustNewServer("redis://127.0.0.1:4000", util.Pointer("srv4")), client.Slave, map[string]string{}),
						}},
				},
				Sentinels: []*sharded.SentinelServer{
					sharded.NewSentinelServerFromParams(redis.MustNewServer("redis://127.0.0.1:26379", util.Pointer("sentinel"))),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &SentinelStatus{
				Sentinels:       tt.fields.Sentinels,
				MonitoredShards: tt.fields.MonitoredShards,
			}
			got, err := ss.ShardedCluster(tt.args.ctx, tt.args.pool)
			if (err != nil) != tt.wantErr {
				t.Errorf("SentinelStatus.ShardedCluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want, cmpopts.IgnoreUnexported(sharded.Cluster{}, sharded.Shard{}, redis.Server{})); len(diff) > 0 {
				t.Errorf("SentinelStatus.ShardedCluster() = diff %s", diff)
			}
		})
	}
}
