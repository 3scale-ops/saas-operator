package sharded

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/3scale-ops/basereconciler/util"
	"github.com/3scale-ops/saas-operator/pkg/redis/client"
	redis "github.com/3scale-ops/saas-operator/pkg/redis/server"
	"github.com/go-test/deep"
)

func DiscoveredServersAreEqual(a, b *Shard) (bool, []string) {
	if len(a.Servers) != len(b.Servers) {
		return false, []string{fmt.Sprintf("different number of servers %d != %d", len(a.Servers), len(b.Servers))}
	}

	for idx := range a.Servers {
		if a.Servers[idx].ID() != b.Servers[idx].ID() {
			return false, []string{fmt.Sprintf("%s != %s", a.Servers[idx].ID(), b.Servers[idx].ID())}
		}
		if a.Servers[idx].Role != b.Servers[idx].Role {
			return false, []string{fmt.Sprintf("%s != %s", a.Servers[idx].Role, b.Servers[idx].Role)}
		}
		if diff := deep.Equal(a.Servers[idx].Config, b.Servers[idx].Config); len(diff) > 0 {
			return false, diff
		}
	}
	return true, []string{}
}

func TestNewShard(t *testing.T) {
	type args struct {
		name    string
		servers map[string]string
		pool    *redis.ServerPool
	}
	tests := []struct {
		name    string
		args    args
		want    *Shard
		wantErr bool
	}{
		{
			name: "Returns a new Shard object when the pool is empty",
			args: args{
				name: "test",
				servers: map[string]string{
					"srv0": "redis://127.0.0.1:1000",
					"srv1": "redis://127.0.0.1:2000",
					"srv2": "redis://127.0.0.1:3000",
				},
				pool: redis.NewServerPool(),
			},
			want: &Shard{
				Name: "test",
				Servers: []*RedisServer{
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.1:1000", util.Pointer("srv0")),
						client.Unknown,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.1:2000", util.Pointer("srv1")),
						client.Unknown,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.1:3000", util.Pointer("srv2")),
						client.Unknown,
						map[string]string{},
					),
				},
				pool: redis.NewServerPool(
					redis.MustNewServer("redis://127.0.0.1:1000", util.Pointer("srv0")),
					redis.MustNewServer("redis://127.0.0.1:2000", util.Pointer("srv1")),
					redis.MustNewServer("redis://127.0.0.1:3000", util.Pointer("srv2")),
				),
			},
			wantErr: false,
		},
		{
			name: "Returns an error (bad connection string)",
			args: args{
				name: "test",
				servers: map[string]string{
					"srv0": "redis://127.0.0.1:1000",
					"srv1": "127.0.0.1:2000",
					"srv2": "redis://127.0.0.1:3000",
				},
				pool: redis.NewServerPool(),
			},
			want: &Shard{
				Name: "test",
				Servers: []*RedisServer{
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.1:1000", util.Pointer("srv0")),
						client.Unknown,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.1:3000", util.Pointer("srv2")),
						client.Unknown,
						map[string]string{},
					),
				},
				pool: redis.NewServerPool(
					redis.MustNewServer("redis://127.0.0.1:1000", util.Pointer("srv0")),
					redis.MustNewServer("redis://127.0.0.1:3000", util.Pointer("srv2")),
				),
			},
			wantErr: true,
		},
		{
			name: "Gets servers from the server pool",
			args: args{
				name: "test",
				servers: map[string]string{
					"redis://127.0.0.1:1000": "redis://127.0.0.1:1000",
					"redis://127.0.0.1:2000": "redis://127.0.0.1:2000",
					"redis://127.0.0.1:3000": "redis://127.0.0.1:3000",
				},
				pool: redis.NewServerPool(
					redis.MustNewServer("redis://127.0.0.1:1000", nil),
					redis.MustNewServer("redis://127.0.0.1:2000", util.Pointer("srv1")),
					redis.MustNewServer("redis://127.0.0.1:3000", util.Pointer("srv2")),
				),
			},
			want: &Shard{
				Name: "test",
				Servers: []*RedisServer{
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.1:1000", nil),
						client.Unknown,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.1:2000", util.Pointer("srv1")),
						client.Unknown,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.1:3000", util.Pointer("srv2")),
						client.Unknown,
						map[string]string{},
					),
				},
				pool: redis.NewServerPool(
					redis.MustNewServer("redis://127.0.0.1:1000", nil),
					redis.MustNewServer("redis://127.0.0.1:2000", util.Pointer("srv1")),
					redis.MustNewServer("redis://127.0.0.1:3000", util.Pointer("srv2")),
				),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewShardFromTopology(tt.args.name, tt.args.servers, tt.args.pool)
			if (error(err) != nil) != tt.wantErr {
				t.Errorf("NewShard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("NewShard() got diff: %v", diff)
			}
		})
	}
}

func TestShard_Discover(t *testing.T) {
	type fields struct {
		Name    string
		Servers []*RedisServer
		pool    *redis.ServerPool
	}
	type args struct {
		ctx      context.Context
		sentinel *SentinelServer
		options  DiscoveryOptionSet
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Shard
		wantErr bool
	}{
		{
			name: "No sentinel: discovers just the roles for all servers in the shard",
			fields: fields{
				Name: "test",
				Servers: []*RedisServer{
					NewRedisServerFromParams(
						redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
							client.NewPredefinedRedisFakeResponse("role-master", nil),
						),
						client.Unknown,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.NewFakeServerWithFakeClient("127.0.0.1", "2000",
							client.NewPredefinedRedisFakeResponse("role-slave", nil),
						),
						client.Unknown,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.NewFakeServerWithFakeClient("127.0.0.1", "3000",
							client.NewPredefinedRedisFakeResponse("role-slave", nil),
						),
						client.Unknown,
						map[string]string{},
					),
				},
				pool: redis.NewServerPool(),
			},
			args: args{
				ctx:      context.TODO(),
				sentinel: nil,
				options:  DiscoveryOptionSet{},
			},
			want: &Shard{Name: "test",
				Servers: []*RedisServer{
					{Server: redis.NewServerFromParams("", "127.0.0.1", "1000", nil), Role: client.Master, Config: map[string]string{}},
					{Server: redis.NewServerFromParams("", "127.0.0.1", "2000", nil), Role: client.Slave, Config: map[string]string{}},
					{Server: redis.NewServerFromParams("", "127.0.0.1", "3000", nil), Role: client.Slave, Config: map[string]string{}},
				}},
			wantErr: false,
		},
		{
			name: "No sentinel: second server fails, returns error",
			fields: fields{
				Name: "test",
				Servers: []*RedisServer{
					NewRedisServerFromParams(
						redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
							client.NewPredefinedRedisFakeResponse("role-master", nil),
						),
						client.Unknown,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.NewFakeServerWithFakeClient("127.0.0.1", "2000",
							client.NewPredefinedRedisFakeResponse("role-slave", errors.New("error")),
						),
						client.Unknown,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.NewFakeServerWithFakeClient("127.0.0.1", "3000",
							client.NewPredefinedRedisFakeResponse("role-slave", nil),
						),
						client.Unknown,
						map[string]string{},
					),
				},
				pool: redis.NewServerPool(),
			},
			args: args{
				ctx:      context.TODO(),
				sentinel: nil,
				options:  DiscoveryOptionSet{},
			},
			want: &Shard{Name: "test",
				Servers: []*RedisServer{
					{Server: redis.NewServerFromParams("", "127.0.0.1", "1000", nil), Role: client.Master, Config: map[string]string{}},
					{Server: redis.NewServerFromParams("", "127.0.0.1", "2000", nil), Role: client.Unknown, Config: map[string]string{}},
					{Server: redis.NewServerFromParams("", "127.0.0.1", "3000", nil), Role: client.Slave, Config: map[string]string{}},
				}},
			wantErr: true,
		},
		{
			name: "Sentinel: discovers roles and config options within a shard (all available options)",
			fields: fields{
				Name:    "test",
				Servers: []*RedisServer{},
				pool: redis.NewServerPool(
					redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
						client.NewPredefinedRedisFakeResponse("role-master", nil),
						client.NewPredefinedRedisFakeResponse("no-save", nil),
					),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "2000",
						client.NewPredefinedRedisFakeResponse("role-slave", nil),
						client.NewPredefinedRedisFakeResponse("save", nil),
						client.NewPredefinedRedisFakeResponse("slave-read-only-yes", nil),
					),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "3000",
						client.NewPredefinedRedisFakeResponse("role-slave", nil),
						client.NewPredefinedRedisFakeResponse("no-save", nil),
						client.NewPredefinedRedisFakeResponse("slave-read-only-no", nil),
					),
				),
			},
			args: args{
				ctx: context.TODO(),
				sentinel: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
					client.FakeResponse{
						// cmd: SentinelMaster
						InjectResponse: func() interface{} {
							return &client.SentinelMasterCmdResult{Name: "test", IP: "127.0.0.1", Port: 1000, Flags: "master"}
						},
						InjectError: func() error { return nil },
					},
					client.FakeResponse{
						// cmd: SentinelGetMasterAddrByName
						InjectResponse: func() interface{} {
							return []string{"127.0.0.1", "1000"}
						},
						InjectError: func() error { return nil },
					},
					client.FakeResponse{
						// cmd: SentinelSlaves
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
				options: []DiscoveryOption{SaveConfigDiscoveryOpt, SlaveReadOnlyDiscoveryOpt},
			},
			want: &Shard{Name: "test",
				Servers: []*RedisServer{
					{Server: redis.NewServerFromParams("", "127.0.0.1", "1000", nil), Role: client.Master, Config: map[string]string{"save": ""}},
					{Server: redis.NewServerFromParams("", "127.0.0.1", "2000", nil), Role: client.Slave, Config: map[string]string{"save": "900 1 300 10", "slave-read-only": "yes"}},
					{Server: redis.NewServerFromParams("", "127.0.0.1", "3000", nil), Role: client.Slave, Config: map[string]string{"save": "", "slave-read-only": "no"}},
				}},
			wantErr: false,
		},
		{
			name: "Sentinel: 'sentinel master' command fails ",
			fields: fields{
				Name:    "test",
				Servers: []*RedisServer{},
				pool: redis.NewServerPool(
					redis.NewFakeServerWithFakeClient("127.0.0.1", "1000"),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "2000"),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "3000"),
				),
			},
			args: args{
				ctx: context.TODO(),
				sentinel: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
					client.FakeResponse{
						// cmd: SentinelMaster
						InjectResponse: func() interface{} {
							return &client.SentinelMasterCmdResult{Name: "test", IP: "127.0.0.1", Port: 1000, Flags: "master"}
						},
						InjectError: func() error { return errors.New("error") },
					},
				)),
				options: []DiscoveryOption{SaveConfigDiscoveryOpt, SlaveReadOnlyDiscoveryOpt},
			},
			want:    &Shard{Name: "test", Servers: []*RedisServer{}},
			wantErr: true,
		},
		{
			name: "Sentinel: master is down",
			fields: fields{
				Name:    "test",
				Servers: []*RedisServer{},
				pool: redis.NewServerPool(
					redis.NewFakeServerWithFakeClient("127.0.0.1", "1000"),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "2000",
						client.NewPredefinedRedisFakeResponse("role-slave", nil),
						client.NewPredefinedRedisFakeResponse("save", nil),
						client.NewPredefinedRedisFakeResponse("slave-read-only-yes", nil),
					),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "3000",
						client.NewPredefinedRedisFakeResponse("role-slave", nil),
						client.NewPredefinedRedisFakeResponse("no-save", nil),
						client.NewPredefinedRedisFakeResponse("slave-read-only-no", nil),
					),
				),
			},
			args: args{
				ctx: context.TODO(),
				sentinel: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
					client.FakeResponse{
						// cmd: SentinelMaster
						InjectResponse: func() interface{} {
							return &client.SentinelMasterCmdResult{Name: "test", IP: "127.0.0.1", Port: 1000, Flags: "master,o_down"}
						},
						InjectError: func() error { return nil },
					},
					client.FakeResponse{
						// cmd: SentinelGetMasterAddrByName
						InjectResponse: func() interface{} {
							return []string{"127.0.0.1", "1000"}
						},
						InjectError: func() error { return nil },
					},
				)),
				options: []DiscoveryOption{SaveConfigDiscoveryOpt, SlaveReadOnlyDiscoveryOpt},
			},
			want: &Shard{Name: "test",
				Servers: []*RedisServer{
					{Server: redis.NewServerFromParams("", "127.0.0.1", "1000", nil), Role: client.Unknown, Config: map[string]string{}},
				}},
			wantErr: true,
		},
		{
			name: "Sentinel: master's Discover() fails",
			fields: fields{
				Name:    "test",
				Servers: []*RedisServer{},
				pool: redis.NewServerPool(
					redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
						client.NewPredefinedRedisFakeResponse("role-master", errors.New("error")),
					),
				),
			},
			args: args{
				ctx: context.TODO(),
				sentinel: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
					client.FakeResponse{
						// cmd: SentinelMaster
						InjectResponse: func() interface{} {
							return &client.SentinelMasterCmdResult{Name: "test", IP: "127.0.0.1", Port: 1000, Flags: "master"}
						},
						InjectError: func() error { return nil },
					},
					client.FakeResponse{
						// cmd: SentinelGetMasterAddrByName
						InjectResponse: func() interface{} {
							return []string{"127.0.0.1", "1000"}
						},
						InjectError: func() error { return nil },
					},
				)),
				options: []DiscoveryOption{},
			},
			want: &Shard{Name: "test", Servers: []*RedisServer{
				{Server: redis.NewServerFromParams("", "127.0.0.1", "1000", nil), Role: client.Unknown, Config: map[string]string{}},
			}},
			wantErr: true,
		},
		{
			name: "Sentinel: master role reported by sentinel differs from role reported by redis",
			fields: fields{
				Name:    "test",
				Servers: []*RedisServer{},
				pool: redis.NewServerPool(
					redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
						client.NewPredefinedRedisFakeResponse("role-slave", nil),
					),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "2000",
						client.NewPredefinedRedisFakeResponse("role-slave", nil),
					),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "3000",
						client.NewPredefinedRedisFakeResponse("role-slave", nil),
					),
				),
			},
			args: args{
				ctx: context.TODO(),
				sentinel: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
					client.FakeResponse{
						// cmd: SentinelMaster
						InjectResponse: func() interface{} {
							return &client.SentinelMasterCmdResult{Name: "test", IP: "127.0.0.1", Port: 1000, Flags: "master"}
						},
						InjectError: func() error { return nil },
					},
					client.FakeResponse{
						// cmd: SentinelGetMasterAddrByName
						InjectResponse: func() interface{} {
							return []string{"127.0.0.1", "1000"}
						},
						InjectError: func() error { return nil },
					},
				)),
				options: []DiscoveryOption{},
			},
			want: &Shard{Name: "test", Servers: []*RedisServer{
				{Server: redis.NewServerFromParams("", "127.0.0.1", "1000", nil), Role: client.Unknown, Config: map[string]string{}},
			}},
			wantErr: true,
		},
		{
			name: "Sentinel: 'sentinel slaves' command fails",
			fields: fields{
				Name:    "test",
				Servers: []*RedisServer{},
				pool: redis.NewServerPool(
					redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
						client.NewPredefinedRedisFakeResponse("role-master", nil),
					),
				),
			},
			args: args{
				ctx: context.TODO(),
				sentinel: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
					client.FakeResponse{
						// cmd: SentinelMaster
						InjectResponse: func() interface{} {
							return &client.SentinelMasterCmdResult{Name: "test", IP: "127.0.0.1", Port: 1000, Flags: "master"}
						},
						InjectError: func() error { return nil },
					},
					client.FakeResponse{
						// cmd: SentinelGetMasterAddrByName
						InjectResponse: func() interface{} {
							return []string{"127.0.0.1", "1000"}
						},
						InjectError: func() error { return nil },
					},
					client.FakeResponse{
						// cmd: SentinelSlaves
						InjectResponse: func() interface{} {
							return []interface{}{}
						},
						InjectError: func() error { return errors.New("error") },
					},
				)),
				options: []DiscoveryOption{},
			},
			want: &Shard{Name: "test",
				Servers: []*RedisServer{
					{Server: redis.NewServerFromParams("", "127.0.0.1", "1000", nil), Role: client.Master, Config: map[string]string{}},
				}},
			wantErr: true,
		},
		{
			name: "Sentinel: a slave is down",
			fields: fields{
				Name:    "test",
				Servers: []*RedisServer{},
				pool: redis.NewServerPool(
					redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
						client.NewPredefinedRedisFakeResponse("role-master", nil),
						client.NewPredefinedRedisFakeResponse("no-save", nil),
					),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "2000"),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "3000",
						client.NewPredefinedRedisFakeResponse("role-slave", nil),
						client.NewPredefinedRedisFakeResponse("no-save", nil),
						client.NewPredefinedRedisFakeResponse("slave-read-only-no", nil),
					),
				),
			},
			args: args{
				ctx: context.TODO(),
				sentinel: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
					client.FakeResponse{
						// cmd: SentinelMaster
						InjectResponse: func() interface{} {
							return &client.SentinelMasterCmdResult{Name: "test", IP: "127.0.0.1", Port: 1000, Flags: "master"}
						},
						InjectError: func() error { return nil },
					},
					client.FakeResponse{
						// cmd: SentinelGetMasterAddrByName
						InjectResponse: func() interface{} {
							return []string{"127.0.0.1", "1000"}
						},
						InjectError: func() error { return nil },
					},
					client.FakeResponse{
						// cmd: SentinelSlaves
						InjectResponse: func() interface{} {
							return []interface{}{
								[]interface{}{
									"name", "127.0.0.1:2000",
									"ip", "127.0.0.1",
									"port", "2000",
									"flags", "slave,s_down",
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
				options: []DiscoveryOption{SaveConfigDiscoveryOpt, SlaveReadOnlyDiscoveryOpt},
			},
			want: &Shard{Name: "test",
				Servers: []*RedisServer{
					{Server: redis.NewServerFromParams("", "127.0.0.1", "1000", nil), Role: client.Master, Config: map[string]string{"save": ""}},
					{Server: redis.NewServerFromParams("", "127.0.0.1", "2000", nil), Role: client.Unknown, Config: map[string]string{}},
					{Server: redis.NewServerFromParams("", "127.0.0.1", "3000", nil), Role: client.Slave, Config: map[string]string{"save": "", "slave-read-only": "no"}},
				}},
			wantErr: true,
		},
		{
			name: "Sentinel: slave role reported by sentinel differs from role reported by redis",
			fields: fields{
				Name:    "test",
				Servers: []*RedisServer{},
				pool: redis.NewServerPool(
					redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
						client.NewPredefinedRedisFakeResponse("role-master", nil),
					),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "2000",
						client.NewPredefinedRedisFakeResponse("role-slave", nil),
					),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "3000",
						client.NewPredefinedRedisFakeResponse("role-master", nil),
					),
				),
			},
			args: args{
				ctx: context.TODO(),
				sentinel: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
					client.FakeResponse{
						// cmd: SentinelMaster
						InjectResponse: func() interface{} {
							return &client.SentinelMasterCmdResult{Name: "test", IP: "127.0.0.1", Port: 1000, Flags: "master"}
						},
						InjectError: func() error { return nil },
					},
					client.FakeResponse{
						// cmd: SentinelGetMasterAddrByName
						InjectResponse: func() interface{} {
							return []string{"127.0.0.1", "1000"}
						},
						InjectError: func() error { return nil },
					},
					client.FakeResponse{
						// cmd: SentinelSlaves
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
				options: []DiscoveryOption{},
			},
			want: &Shard{Name: "test",
				Servers: []*RedisServer{
					{Server: redis.NewServerFromParams("", "127.0.0.1", "1000", nil), Role: client.Master, Config: map[string]string{}},
					{Server: redis.NewServerFromParams("", "127.0.0.1", "2000", nil), Role: client.Slave, Config: map[string]string{}},
					{Server: redis.NewServerFromParams("", "127.0.0.1", "3000", nil), Role: client.Unknown, Config: map[string]string{}},
				}},
			wantErr: true,
		},
		{
			name: "Sentinel: Discover() fails for a slave",
			fields: fields{
				Name:    "test",
				Servers: []*RedisServer{},
				pool: redis.NewServerPool(
					redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
						client.NewPredefinedRedisFakeResponse("role-master", nil),
						client.NewPredefinedRedisFakeResponse("no-save", nil),
					),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "2000",
						client.NewPredefinedRedisFakeResponse("role-slave", nil),
						client.NewPredefinedRedisFakeResponse("no-save", errors.New("error")),
					),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "3000",
						client.NewPredefinedRedisFakeResponse("role-slave", nil),
						client.NewPredefinedRedisFakeResponse("no-save", nil),
					),
				),
			},
			args: args{
				ctx: context.TODO(),
				sentinel: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
					client.FakeResponse{
						// cmd: SentinelMaster
						InjectResponse: func() interface{} {
							return &client.SentinelMasterCmdResult{Name: "test", IP: "127.0.0.1", Port: 1000, Flags: "master"}
						},
						InjectError: func() error { return nil },
					},
					client.FakeResponse{
						// cmd: SentinelGetMasterAddrByName
						InjectResponse: func() interface{} {
							return []string{"127.0.0.1", "1000"}
						},
						InjectError: func() error { return nil },
					},
					client.FakeResponse{
						// cmd: SentinelSlaves
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
				options: []DiscoveryOption{SaveConfigDiscoveryOpt},
			},
			want: &Shard{Name: "test",
				Servers: []*RedisServer{
					{Server: redis.NewServerFromParams("", "127.0.0.1", "1000", nil), Role: client.Master, Config: map[string]string{"save": ""}},
					{Server: redis.NewServerFromParams("", "127.0.0.1", "2000", nil), Role: client.Unknown, Config: map[string]string{}},
					{Server: redis.NewServerFromParams("", "127.0.0.1", "3000", nil), Role: client.Slave, Config: map[string]string{"save": ""}},
				}},
			wantErr: true,
		},
		{
			name: "Sentinel: Discover() only masters",
			fields: fields{
				Name:    "test",
				Servers: []*RedisServer{},
				pool: redis.NewServerPool(
					redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
						client.NewPredefinedRedisFakeResponse("role-master", nil),
					),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "2000"),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "3000"),
				),
			},
			args: args{
				ctx: context.TODO(),
				sentinel: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
					client.FakeResponse{
						// cmd: SentinelMaster
						InjectResponse: func() interface{} {
							return &client.SentinelMasterCmdResult{Name: "test", IP: "127.0.0.1", Port: 1000, Flags: "master"}
						},
						InjectError: func() error { return nil },
					},
					client.FakeResponse{
						// cmd: SentinelGetMasterAddrByName
						InjectResponse: func() interface{} {
							return []string{"127.0.0.1", "1000"}
						},
						InjectError: func() error { return nil },
					},
				)),
				options: []DiscoveryOption{OnlyMasterDiscoveryOpt},
			},
			want: &Shard{Name: "test",
				Servers: []*RedisServer{
					{Server: redis.NewServerFromParams("", "127.0.0.1", "1000", nil), Role: client.Master, Config: map[string]string{}},
				}},
			wantErr: false,
		},
		{
			name: "Sentinel: failover in progress, refuse to discover slaves",
			fields: fields{
				Name:    "test",
				Servers: []*RedisServer{},
				pool: redis.NewServerPool(
					redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
						client.NewPredefinedRedisFakeResponse("role-master", nil),
					),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "2000",
						client.NewPredefinedRedisFakeResponse("role-master", nil),
					),
					redis.NewFakeServerWithFakeClient("127.0.0.1", "3000"),
				),
			},
			args: args{
				ctx: context.TODO(),
				sentinel: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
					client.FakeResponse{
						// cmd: SentinelMaster
						InjectResponse: func() interface{} {
							return &client.SentinelMasterCmdResult{Name: "test", IP: "127.0.0.1", Port: 1000, Flags: "master"}
						},
						InjectError: func() error { return nil },
					},
					client.FakeResponse{
						// cmd: SentinelGetMasterAddrByName
						InjectResponse: func() interface{} {
							return []string{"127.0.0.1", "2000"}
						},
						InjectError: func() error { return nil },
					},
				)),
				options: []DiscoveryOption{OnlyMasterDiscoveryOpt},
			},
			want: &Shard{Name: "test",
				Servers: []*RedisServer{
					{Server: redis.NewServerFromParams("", "127.0.0.1", "2000", nil), Role: client.Master, Config: map[string]string{}},
				}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shard{
				Name:    tt.fields.Name,
				Servers: tt.fields.Servers,
				pool:    tt.fields.pool,
			}
			if err := s.Discover(tt.args.ctx, tt.args.sentinel, tt.args.options...); (err != nil) != tt.wantErr {
				t.Errorf("Shard.Discover() error = %v, wantErr %v", err, tt.wantErr)
			}
			if equal, diff := DiscoveredServersAreEqual(s, tt.want); !equal {
				t.Errorf("Shard.Discover() got diff = %v", diff)
			}
		})
	}
}

func TestShard_Init(t *testing.T) {
	type fields struct {
		Name    string
		Servers []*RedisServer
	}
	type args struct {
		ctx            context.Context
		masterHostPort string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "All redis servers configured",
			fields: fields{
				Name: "test",
				Servers: []*RedisServer{
					NewRedisServerFromParams(
						redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"slave", "127.0.0.1"}
								},
								InjectError: func() error { return nil },
							},
							client.FakeResponse{
								InjectResponse: func() interface{} { return nil },
								InjectError:    func() error { return nil },
							},
							client.NewPredefinedRedisFakeResponse("role-master", nil),
							client.NewPredefinedRedisFakeResponse("role-master", nil),
						),
						client.Unknown,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.NewFakeServerWithFakeClient("127.0.0.1", "2000",
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"slave", "127.0.0.1"}
								},
								InjectError: func() error { return nil },
							},
							client.FakeResponse{
								InjectResponse: func() interface{} { return nil },
								InjectError:    func() error { return nil },
							},
						),
						client.Unknown,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.NewFakeServerWithFakeClient("127.0.0.1", "3000",
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"slave", "127.0.0.1"}
								},
								InjectError: func() error { return nil },
							},
							client.FakeResponse{
								InjectResponse: func() interface{} { return nil },
								InjectError:    func() error { return nil },
							},
						),
						client.Unknown,
						map[string]string{},
					),
				},
			},
			args:    args{ctx: context.TODO(), masterHostPort: "127.0.0.1:1000"},
			want:    []string{"127.0.0.1:1000", "127.0.0.1:2000", "127.0.0.1:3000"},
			wantErr: false,
		},
		{
			name: "No configuration needed",
			fields: fields{
				Name: "test",
				Servers: []*RedisServer{
					NewRedisServerFromParams(
						redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
							client.NewPredefinedRedisFakeResponse("role-master", nil),
						),
						client.Unknown,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.NewFakeServerWithFakeClient("127.0.0.1", "2000",
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"slave", "10.0.0.1"}
								},
								InjectError: func() error { return nil },
							},
						),
						client.Unknown,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.NewFakeServerWithFakeClient("127.0.0.1", "3000",
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"slave", "10.0.0.1"}
								},
								InjectError: func() error { return nil },
							},
						),
						client.Unknown,
						map[string]string{},
					),
				},
			},
			args:    args{ctx: context.TODO(), masterHostPort: "127.0.0.1:1000"},
			want:    []string{},
			wantErr: false,
		},
		{
			name: "Returns error",
			fields: fields{
				Name: "test",
				Servers: []*RedisServer{
					NewRedisServerFromParams(
						redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
							client.FakeResponse{
								InjectResponse: func() interface{} { return []interface{}{} },
								InjectError:    func() error { return errors.New("error") },
							},
						),
						client.Unknown,
						map[string]string{},
					),
				},
			},
			args:    args{ctx: context.TODO(), masterHostPort: "127.0.0.1:1000"},
			want:    []string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shard{
				Name:    tt.fields.Name,
				Servers: tt.fields.Servers,
			}
			got, err := s.Init(tt.args.ctx, tt.args.masterHostPort)
			if (err != nil) != tt.wantErr {
				t.Errorf("Shard.Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Shard.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShard_GetMasterAddr(t *testing.T) {
	type fields struct {
		Name    string
		Servers []*RedisServer
		pool    *redis.ServerPool
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Returns the master's address",
			fields: fields{
				Name: "test",
				Servers: []*RedisServer{
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.1:1000", nil),
						client.Master,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.2:2000", util.Pointer("srv1")),
						client.Unknown,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.3:3000", util.Pointer("srv2")),
						client.Unknown,
						map[string]string{},
					),
				},
				pool: redis.NewServerPool(),
			},
			want:    "127.0.0.1:1000",
			wantErr: false,
		},
		{
			name: "Error, no master",
			fields: fields{
				Name: "test",
				Servers: []*RedisServer{
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.1:1000", nil),
						client.Unknown,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.2:2000", util.Pointer("srv1")),
						client.Unknown,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.3:3000", util.Pointer("srv2")),
						client.Unknown,
						map[string]string{},
					),
				},
				pool: redis.NewServerPool(),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Error, more than one master",
			fields: fields{
				Name: "test",
				Servers: []*RedisServer{
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.1:1000", nil),
						client.Master,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.2:2000", util.Pointer("srv1")),
						client.Slave,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.3:3000", util.Pointer("srv2")),
						client.Master,
						map[string]string{},
					),
				},
				pool: redis.NewServerPool(),
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shard := &Shard{
				Name:    tt.fields.Name,
				Servers: tt.fields.Servers,
				pool:    tt.fields.pool,
			}
			got, err := shard.GetMaster()
			if (err != nil) != tt.wantErr {
				t.Errorf("Shard.GetMasterAddr() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got.ID() != tt.want {
				t.Errorf("Shard.GetMasterAddr() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShard_GetServerByID(t *testing.T) {
	type args struct {
		hostport string
	}
	tests := []struct {
		name      string
		servers   map[string]string
		args      args
		wantIndex int
		wantErr   bool
	}{
		{
			name: "Resturns a server",
			servers: map[string]string{
				"host1": "redis://127.0.0.1:1000",
				"host2": "redis://127.0.0.1:2000",
			},
			args: args{
				hostport: "127.0.0.1:1000",
			},
			wantIndex: 0,
			wantErr:   false,
		},
		{
			name: "Adds a server",
			servers: map[string]string{
				"host1": "redis://127.0.0.1:1000",
				"host2": "redis://127.0.0.1:2000",
			},
			args: args{
				hostport: "127.0.0.1:3000",
			},
			wantIndex: 2,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shard, _ := NewShardFromTopology("test", tt.servers, redis.NewServerPool())
			got, err := shard.GetServerByID(tt.args.hostport)
			if (err != nil) != tt.wantErr {
				t.Errorf("Shard.GetServerByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != shard.Servers[tt.wantIndex] {
				t.Errorf("Shard.GetServerByID() = %v, want %v", got, shard.Servers[tt.wantIndex])
			}
		})
	}
}
