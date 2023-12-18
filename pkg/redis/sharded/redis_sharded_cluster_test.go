package sharded

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/3scale-ops/basereconciler/util"
	"github.com/3scale-ops/saas-operator/pkg/redis/client"
	redis "github.com/3scale-ops/saas-operator/pkg/redis/server"
	"github.com/go-test/deep"
)

func TestNewShardedCluster(t *testing.T) {
	type args struct {
		ctx        context.Context
		serverList map[string]map[string]string
		pool       *redis.ServerPool
	}
	tests := []struct {
		name    string
		args    args
		want    *Cluster
		wantErr bool
	}{
		{
			name: "Returns a new ShardedCluster object",
			args: args{
				ctx: context.TODO(),
				serverList: map[string]map[string]string{
					"shard00":  {"srv00-0": "redis://127.0.0.1:1000", "srv00-1": "redis://127.0.0.1:2000"},
					"shard01":  {"srv01-0": "redis://127.0.0.1:3000", "srv01-1": "redis://127.0.0.1:4000"},
					"sentinel": {"sentinel-0": "redis://127.0.0.1:5000", "sentinel-1": "redis://127.0.0.1:6000"},
				},
				pool: redis.NewServerPool(),
			},
			want: &Cluster{
				Shards: []*Shard{
					{
						Name: "shard00",
						Servers: []*RedisServer{
							NewRedisServerFromParams(
								redis.MustNewServer("redis://127.0.0.1:1000", util.Pointer("srv00-0")),
								client.Unknown,
								map[string]string{},
							),
							NewRedisServerFromParams(
								redis.MustNewServer("redis://127.0.0.1:2000", util.Pointer("srv00-1")),
								client.Unknown,
								map[string]string{},
							),
						},
						pool: redis.NewServerPool(
							redis.MustNewServer("redis://127.0.0.1:1000", util.Pointer("srv00-0")),
							redis.MustNewServer("redis://127.0.0.1:2000", util.Pointer("srv00-1")),
							redis.MustNewServer("redis://127.0.0.1:3000", util.Pointer("srv01-0")),
							redis.MustNewServer("redis://127.0.0.1:4000", util.Pointer("srv01-1")),
							redis.MustNewServer("redis://127.0.0.1:5000", util.Pointer("sentinel-0")),
							redis.MustNewServer("redis://127.0.0.1:6000", util.Pointer("sentinel-1")),
						),
					},
					{
						Name: "shard01",
						Servers: []*RedisServer{
							NewRedisServerFromParams(
								redis.MustNewServer("redis://127.0.0.1:3000", util.Pointer("srv01-0")),
								client.Unknown,
								map[string]string{},
							),
							NewRedisServerFromParams(
								redis.MustNewServer("redis://127.0.0.1:4000", util.Pointer("srv01-1")),
								client.Unknown,
								map[string]string{},
							),
						},
						pool: redis.NewServerPool(
							redis.MustNewServer("redis://127.0.0.1:1000", util.Pointer("srv00-0")),
							redis.MustNewServer("redis://127.0.0.1:2000", util.Pointer("srv00-1")),
							redis.MustNewServer("redis://127.0.0.1:3000", util.Pointer("srv01-0")),
							redis.MustNewServer("redis://127.0.0.1:4000", util.Pointer("srv01-1")),
							redis.MustNewServer("redis://127.0.0.1:5000", util.Pointer("sentinel-0")),
							redis.MustNewServer("redis://127.0.0.1:6000", util.Pointer("sentinel-1")),
						),
					},
				},
				Sentinels: []*SentinelServer{
					{Server: redis.MustNewServer("redis://127.0.0.1:5000", util.Pointer("sentinel-0"))},
					{Server: redis.MustNewServer("redis://127.0.0.1:6000", util.Pointer("sentinel-1"))},
				},
				pool: redis.NewServerPool(
					redis.MustNewServer("redis://127.0.0.1:1000", util.Pointer("srv00-0")),
					redis.MustNewServer("redis://127.0.0.1:2000", util.Pointer("srv00-1")),
					redis.MustNewServer("redis://127.0.0.1:3000", util.Pointer("srv01-0")),
					redis.MustNewServer("redis://127.0.0.1:4000", util.Pointer("srv01-1")),
					redis.MustNewServer("redis://127.0.0.1:5000", util.Pointer("sentinel-0")),
					redis.MustNewServer("redis://127.0.0.1:6000", util.Pointer("sentinel-1")),
				),
			},
			wantErr: false,
		},
		{
			name: "Returns error",
			args: args{
				ctx: context.TODO(),
				serverList: map[string]map[string]string{
					"shard00": {"srv00-0": "redis://127.0.0.1:1000", "srv00-1": "redis://127.0.0.1:2000"},
					"shard01": {"srv01-0": "127.0.0.1:3000", "srv01-1": "redis://127.0.0.1:4000"},
				},
				pool: redis.NewServerPool(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewShardedClusterFromTopology(tt.args.ctx, tt.args.serverList, tt.args.pool)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewShardedCluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("NewShardedCluster() got diff: %v", diff)
			}
		})
	}
}

func TestShardedCluster_GetShardNames(t *testing.T) {
	tests := []struct {
		name string
		sc   Cluster
		want []string
	}{
		{
			name: "Returns the shrard names as a slice of strings",
			sc: Cluster{
				Shards: []*Shard{
					{
						Name:    "shard00",
						Servers: []*RedisServer{},
					},
					{
						Name:    "shard01",
						Servers: []*RedisServer{},
					},
					{
						Name:    "shard02",
						Servers: []*RedisServer{},
					},
				},
			},
			want: []string{"shard00", "shard01", "shard02"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sc.GetShardNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShardedCluster.GetShardNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShardedCluster_LookupShardByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		sc   Cluster
		args args
		want *Shard
	}{
		{
			name: "Returns the shard of the given name",
			sc: Cluster{
				Shards: []*Shard{
					{
						Name: "shard00",
						Servers: []*RedisServer{
							NewRedisServerFromParams(
								redis.MustNewServer("redis://127.0.0.1:1000", util.Pointer("srv00-0")),
								client.Unknown,
								map[string]string{},
							),
						},
					},
					{
						Name: "shard01",
						Servers: []*RedisServer{
							NewRedisServerFromParams(
								redis.MustNewServer("redis://127.0.0.1:2000", util.Pointer("srv01-0")),
								client.Unknown,
								map[string]string{},
							),
						},
					},
				},
			},
			args: args{
				name: "shard01",
			},
			want: &Shard{
				Name: "shard01",
				Servers: []*RedisServer{
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.1:2000", util.Pointer("srv01-0")),
						client.Unknown,
						map[string]string{},
					),
				},
			},
		},
		{
			name: "Returns nil if not found",
			sc:   Cluster{},
			args: args{
				name: "shard01",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := deep.Equal(tt.sc.LookupShardByName(tt.args.name), tt.want); len(diff) > 0 {
				t.Errorf("ShardedCluster.LookupShardByName() got diff: %v", diff)
			}
		})
	}
}

func TestCluster_Discover(t *testing.T) {
	type fields struct {
		Shards    []*Shard
		Sentinels []*SentinelServer
		pool      *redis.ServerPool
	}
	type args struct {
		ctx     context.Context
		options []DiscoveryOption
	}
	tests := []struct {
		name    string
		fields  fields
		pool    *redis.ServerPool
		args    args
		wantErr bool
	}{
		{
			name: "Discovers roles for all cluster servers",
			fields: fields{
				Shards: []*Shard{
					{
						Name: "shard0",
						Servers: []*RedisServer{
							{
								Server: redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
									client.NewPredefinedRedisFakeResponse("role-master", nil),
								),
								Role:   client.Unknown,
								Config: map[string]string{},
							},
							{
								Server: redis.NewFakeServerWithFakeClient("127.0.0.1", "2000",
									client.NewPredefinedRedisFakeResponse("role-slave", nil),
								),
								Role:   client.Unknown,
								Config: map[string]string{},
							},
							{
								Server: redis.NewFakeServerWithFakeClient("127.0.0.1", "3000",
									client.NewPredefinedRedisFakeResponse("role-slave", nil),
								),
								Role:   client.Unknown,
								Config: map[string]string{},
							},
						},
						pool: redis.NewServerPool(),
					},
					{
						Name: "shard1",
						Servers: []*RedisServer{
							{
								Server: redis.NewFakeServerWithFakeClient("127.0.0.1", "4000",
									client.NewPredefinedRedisFakeResponse("role-slave", nil),
								),
								Role:   client.Unknown,
								Config: map[string]string{},
							},
							{
								Server: redis.NewFakeServerWithFakeClient("127.0.0.1", "5000",
									client.NewPredefinedRedisFakeResponse("role-master", nil),
								),
								Role:   client.Unknown,
								Config: map[string]string{},
							},
							{
								Server: redis.NewFakeServerWithFakeClient("127.0.0.1", "6000",
									client.NewPredefinedRedisFakeResponse("role-slave", nil),
								),
								Role:   client.Unknown,
								Config: map[string]string{},
							},
						},
						pool: redis.NewServerPool(),
					},
				},
				Sentinels: []*SentinelServer{},
				pool:      &redis.ServerPool{},
			},
			args: args{
				ctx:     context.TODO(),
				options: []DiscoveryOption{},
			},
			wantErr: false,
		},
		{
			name: "Returns error",
			fields: fields{
				Shards: []*Shard{
					{
						Name: "shard0",
						Servers: []*RedisServer{
							{
								Server: redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
									client.NewPredefinedRedisFakeResponse("role-master", errors.New("error")),
								),
								Role:   client.Unknown,
								Config: map[string]string{},
							},
						},
						pool: redis.NewServerPool(),
					},
				},
			},
			pool: &redis.ServerPool{},
			args: args{
				ctx:     context.TODO(),
				options: []DiscoveryOption{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cluster := Cluster{
				Shards:    tt.fields.Shards,
				Sentinels: tt.fields.Sentinels,
				pool:      tt.fields.pool,
			}
			if err := cluster.Discover(tt.args.ctx, tt.args.options...); (err != nil) != tt.wantErr {
				t.Errorf("Cluster.Discover() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCluster_SentinelDiscover(t *testing.T) {
	type fields struct {
		Shards    []*Shard
		Sentinels []*SentinelServer
		pool      *redis.ServerPool
	}
	type args struct {
		ctx  context.Context
		opts []DiscoveryOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Discovers roles for all cluster servers using sentinel",
			fields: fields{
				Shards: []*Shard{
					{
						Name: "shard0",
						Servers: []*RedisServer{
							{
								Server: redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
									client.NewPredefinedRedisFakeResponse("role-master", nil),
								),
								Role:   client.Unknown,
								Config: map[string]string{},
							},
							{
								Server: redis.NewFakeServerWithFakeClient("127.0.0.1", "2000",
									client.NewPredefinedRedisFakeResponse("role-slave", nil),
								),
								Role:   client.Unknown,
								Config: map[string]string{},
							},
							{
								Server: redis.NewFakeServerWithFakeClient("127.0.0.1", "3000",
									client.NewPredefinedRedisFakeResponse("role-slave", nil),
								),
								Role:   client.Unknown,
								Config: map[string]string{},
							},
						},
						pool: redis.NewServerPool(),
					},
					{
						Name: "shard1",
						Servers: []*RedisServer{
							{
								Server: redis.NewFakeServerWithFakeClient("127.0.0.1", "4000",
									client.NewPredefinedRedisFakeResponse("role-slave", nil),
								),
								Role:   client.Unknown,
								Config: map[string]string{},
							},
							{
								Server: redis.NewFakeServerWithFakeClient("127.0.0.1", "5000",
									client.NewPredefinedRedisFakeResponse("role-master", nil),
								),
								Role:   client.Unknown,
								Config: map[string]string{},
							},
							{
								Server: redis.NewFakeServerWithFakeClient("127.0.0.1", "6000",
									client.NewPredefinedRedisFakeResponse("role-slave", nil),
								),
								Role:   client.Unknown,
								Config: map[string]string{},
							},
						},
						pool: redis.NewServerPool(),
					},
				},
				Sentinels: []*SentinelServer{
					NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("sentinel-0", "1000",
						client.FakeResponse{
							// cmd: Ping
							InjectResponse: func() interface{} { return nil },
							InjectError:    func() error { return nil },
						},
						client.FakeResponse{
							// cmd: SentinelMasters()
							InjectResponse: func() interface{} {
								return []interface{}{
									[]interface{}{"name", "shard0", "ip", "127.0.0.1", "port", "1000"},
									[]interface{}{"name", "shard1", "ip", "127.0.0.1", "port", "5000"},
								}
							},
							InjectError: func() error { return nil },
						},
						client.FakeResponse{
							// cmd: SentinelMaster (shard0)
							InjectResponse: func() interface{} {
								return &client.SentinelMasterCmdResult{Name: "shard0", IP: "127.0.0.1", Port: 1000, Flags: "master"}
							},
							InjectError: func() error { return nil },
						},
						client.FakeResponse{
							// cmd: SentinelGetMasterAddrByName (shard0)
							InjectResponse: func() interface{} {
								return []string{"127.0.0.1", "1000"}
							},
							InjectError: func() error { return nil },
						},
						client.FakeResponse{
							// cmd: SentinelSlaves (shard0)
							InjectResponse: func() interface{} {
								return []interface{}{
									[]interface{}{
										"name", "127.0.0.1:2000",
										"ip", "127.0.0.1",
										"port", "2000",
										"flags", "slave",
									},
									[]interface{}{
										"name", "127.0.0.1:3000",
										"ip", "127.0.0.1",
										"port", "3000",
										"flags", "slave",
									},
								}
							},
							InjectError: func() error { return nil },
						},
						client.FakeResponse{
							// cmd: SentinelMaster (shard1)
							InjectResponse: func() interface{} {
								return &client.SentinelMasterCmdResult{Name: "shard1", IP: "127.0.0.1", Port: 5000, Flags: "master"}
							},
							InjectError: func() error { return nil },
						},
						client.FakeResponse{
							// cmd: SentinelGetMasterAddrByName (shard1)
							InjectResponse: func() interface{} {
								return []string{"127.0.0.1", "5000"}
							},
							InjectError: func() error { return nil },
						},
						client.FakeResponse{
							// cmd: SentinelSlaves (shard1)
							InjectResponse: func() interface{} {
								return []interface{}{
									[]interface{}{
										"name", "127.0.0.1:4000",
										"ip", "127.0.0.1",
										"port", "4000",
										"flags", "slave",
									},
									[]interface{}{
										"name", "127.0.0.1:6000",
										"ip", "127.0.0.1",
										"port", "6000",
										"flags", "slave",
									},
								}
							},
							InjectError: func() error { return nil },
						},
					)),
					NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("sentinel-1", "2000",
						client.FakeResponse{
							// cmd: Ping
							InjectResponse: func() interface{} { return nil },
							InjectError:    func() error { return errors.New("ping failed") },
						},
					)),
				},
				pool: &redis.ServerPool{},
			},
			args: args{
				ctx:  context.TODO(),
				opts: []DiscoveryOption{},
			},
			wantErr: false,
		},
		{
			name: "Discovers shards when not provided",
			fields: fields{
				Shards: nil,
				Sentinels: []*SentinelServer{
					NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("sentinel-0", "1000",
						client.FakeResponse{
							// cmd: Ping
							InjectResponse: func() interface{} { return nil },
							InjectError:    func() error { return nil },
						},
						client.FakeResponse{
							// cmd: SentinelMasters()
							InjectResponse: func() interface{} {
								return []interface{}{
									[]interface{}{"name", "shard0", "ip", "127.0.0.1", "port", "1000"},
									[]interface{}{"name", "shard1", "ip", "127.0.0.1", "port", "5000"},
								}
							},
							InjectError: func() error { return nil },
						},
						client.FakeResponse{
							// cmd: SentinelMaster (shard0)
							InjectResponse: func() interface{} {
								return &client.SentinelMasterCmdResult{Name: "shard0", IP: "127.0.0.1", Port: 1000, Flags: "master"}
							},
							InjectError: func() error { return nil },
						},
						client.FakeResponse{
							// cmd: SentinelGetMasterAddrByName (shard0)
							InjectResponse: func() interface{} {
								return []string{"127.0.0.1", "1000"}
							},
							InjectError: func() error { return nil },
						},
						client.FakeResponse{
							// cmd: SentinelSlaves (shard0)
							InjectResponse: func() interface{} {
								return []interface{}{
									[]interface{}{
										"name", "127.0.0.1:2000",
										"ip", "127.0.0.1",
										"port", "2000",
										"flags", "slave",
									},
									[]interface{}{
										"name", "127.0.0.1:3000",
										"ip", "127.0.0.1",
										"port", "3000",
										"flags", "slave",
									},
								}
							},
							InjectError: func() error { return nil },
						},
						client.FakeResponse{
							// cmd: SentinelMaster (shard1)
							InjectResponse: func() interface{} {
								return &client.SentinelMasterCmdResult{Name: "shard1", IP: "127.0.0.1", Port: 5000, Flags: "master"}
							},
							InjectError: func() error { return nil },
						},
						client.FakeResponse{
							// cmd: SentinelGetMasterAddrByName (shard1)
							InjectResponse: func() interface{} {
								return []string{"127.0.0.1", "5000"}
							},
							InjectError: func() error { return nil },
						},
						client.FakeResponse{
							// cmd: SentinelSlaves (shard1)
							InjectResponse: func() interface{} {
								return []interface{}{
									[]interface{}{
										"name", "127.0.0.1:4000",
										"ip", "127.0.0.1",
										"port", "4000",
										"flags", "slave",
									},
									[]interface{}{
										"name", "127.0.0.1:6000",
										"ip", "127.0.0.1",
										"port", "6000",
										"flags", "slave",
									},
								}
							},
							InjectError: func() error { return nil },
						},
					)),
					NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("sentinel-1", "2000",
						client.FakeResponse{
							// cmd: Ping
							InjectResponse: func() interface{} { return nil },
							InjectError:    func() error { return errors.New("ping failed") },
						},
					)),
				},
				pool: redis.NewServerPool(
					redis.NewFakeServerWithFakeClient("127.0.0.1", "1000", client.NewPredefinedRedisFakeResponse("role-master", nil)),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "2000", client.NewPredefinedRedisFakeResponse("role-slave", nil)),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "3000", client.NewPredefinedRedisFakeResponse("role-slave", nil)),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "4000", client.NewPredefinedRedisFakeResponse("role-slave", nil)),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "5000", client.NewPredefinedRedisFakeResponse("role-master", nil)),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "6000", client.NewPredefinedRedisFakeResponse("role-slave", nil)),
				),
			},
			args: args{
				ctx:  context.TODO(),
				opts: []DiscoveryOption{},
			},
			wantErr: false,
		},
		{
			name: "Error discovering server",
			fields: fields{
				Shards: []*Shard{
					{
						Name: "shard0",
						Servers: []*RedisServer{
							{
								Server: redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
									client.NewPredefinedRedisFakeResponse("role-master", nil),
								),
								Role:   client.Unknown,
								Config: map[string]string{},
							},
							{
								Server: redis.NewFakeServerWithFakeClient("127.0.0.1", "2000",
									client.NewPredefinedRedisFakeResponse("role-slave", nil),
								),
								Role:   client.Unknown,
								Config: map[string]string{},
							},
							{
								Server: redis.NewFakeServerWithFakeClient("127.0.0.1", "3000",
									client.NewPredefinedRedisFakeResponse("role-slave", errors.New("error")),
								),
								Role:   client.Unknown,
								Config: map[string]string{},
							},
						},
						pool: redis.NewServerPool(),
					},
				},
				Sentinels: []*SentinelServer{
					NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("sentinel-0", "1000",
						client.FakeResponse{
							// cmd: Ping
							InjectResponse: func() interface{} { return nil },
							InjectError:    func() error { return nil },
						},
						client.FakeResponse{
							// cmd: SentinelMasters()
							InjectResponse: func() interface{} {
								return []interface{}{
									[]interface{}{"name", "shard0", "ip", "127.0.0.1", "port", "1000"},
								}
							},
							InjectError: func() error { return nil },
						},
						client.FakeResponse{
							// cmd: SentinelMaster (shard0)
							InjectResponse: func() interface{} {
								return &client.SentinelMasterCmdResult{Name: "shard0", IP: "127.0.0.1", Port: 1000, Flags: "master"}
							},
							InjectError: func() error { return nil },
						},
						client.FakeResponse{
							// cmd: SentinelGetMasterAddrByName (shard0)
							InjectResponse: func() interface{} {
								return []string{"127.0.0.1", "1000"}
							},
							InjectError: func() error { return nil },
						},
						client.FakeResponse{
							// cmd: SentinelSlaves (shard0)
							InjectResponse: func() interface{} {
								return []interface{}{
									[]interface{}{
										"name", "127.0.0.1:2000",
										"ip", "127.0.0.1",
										"port", "2000",
										"flags", "slave",
									},
									[]interface{}{
										"name", "127.0.0.1:3000",
										"ip", "127.0.0.1",
										"port", "3000",
										"flags", "slave",
									},
								}
							},
							InjectError: func() error { return nil },
						},
					)),
					NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("sentinel-1", "2000",
						client.FakeResponse{
							// cmd: Ping
							InjectResponse: func() interface{} { return nil },
							InjectError:    func() error { return errors.New("ping failed") },
						},
					)),
				},
				pool: &redis.ServerPool{},
			},
			args: args{
				ctx:  context.TODO(),
				opts: []DiscoveryOption{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cluster := Cluster{
				Shards:    tt.fields.Shards,
				Sentinels: tt.fields.Sentinels,
				pool:      tt.fields.pool,
			}
			if err := cluster.SentinelDiscover(tt.args.ctx, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Cluster.SentinelDiscover() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCluster_GetSentinel(t *testing.T) {
	type fields struct {
		Shards    []*Shard
		Sentinels []*SentinelServer
		pool      *redis.ServerPool
	}
	type args struct {
		pctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *SentinelServer
	}{
		{
			name: "Returns the first sentinel",
			fields: fields{
				Shards: []*Shard{},
				Sentinels: []*SentinelServer{
					NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
						// cmd: ping
						client.FakeResponse{
							InjectResponse: func() interface{} { return nil },
							InjectError:    func() error { return nil },
						},
					)),
					NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("127.0.0.1", "2000",
						// cmd: ping
						client.FakeResponse{
							InjectResponse: func() interface{} { return nil },
							InjectError:    func() error { return errors.New("error") },
						},
					)),
					NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("127.0.0.1", "3000",
						// cmd: ping
						client.FakeResponse{
							InjectResponse: func() interface{} { return nil },
							InjectError:    func() error { return errors.New("error") },
						},
					)),
				},
				pool: &redis.ServerPool{},
			},
			args: args{
				pctx: context.TODO(),
			},
			want: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("127.0.0.1", "1000")),
		},
		{
			name: "Returns the third sentinel",
			fields: fields{
				Shards: []*Shard{},
				Sentinels: []*SentinelServer{
					NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
						// cmd: ping
						client.FakeResponse{
							InjectResponse: func() interface{} { return nil },
							InjectError:    func() error { return errors.New("error") },
						},
					)),
					NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("127.0.0.1", "2000",
						// cmd: ping
						client.FakeResponse{
							InjectResponse: func() interface{} { return nil },
							InjectError:    func() error { return errors.New("error") },
						},
					)),
					NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("127.0.0.1", "3000",
						// cmd: ping
						client.FakeResponse{
							InjectResponse: func() interface{} { return nil },
							InjectError:    func() error { return nil },
						},
					)),
				},
				pool: &redis.ServerPool{},
			},
			args: args{
				pctx: context.TODO(),
			},
			want: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("127.0.0.1", "3000")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cluster := Cluster{
				Shards:    tt.fields.Shards,
				Sentinels: tt.fields.Sentinels,
				pool:      tt.fields.pool,
			}
			got := cluster.GetSentinel(tt.args.pctx)
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("Cluster.GetSentinel() = got diff %v", diff)
			}
		})
	}
}
