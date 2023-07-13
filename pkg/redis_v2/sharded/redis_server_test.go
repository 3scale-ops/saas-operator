package sharded

import (
	"context"
	"errors"
	"testing"

	"github.com/3scale/saas-operator/pkg/redis_v2/client"
	redis "github.com/3scale/saas-operator/pkg/redis_v2/server"
	"github.com/3scale/saas-operator/pkg/util"
	"github.com/go-test/deep"
)

func init() {
	deep.CompareUnexportedFields = true
}

func TestNewRedisServerFromParams(t *testing.T) {
	type args struct {
		connectionString string
		alias            *string
		pool             *redis.ServerPool
	}
	tests := []struct {
		name    string
		args    args
		want    *RedisServer
		wantErr bool
	}{
		{
			name: "Retuns a RedisServer",
			args: args{
				connectionString: "redis://127.0.0.1:1000",
				alias:            util.Pointer("host1"),
				pool:             redis.NewServerPool(redis.NewServerFromParams("host1", "127.0.0.1", "1000", client.MustNewFromConnectionString("redis://127.0.0.1:1000"))),
			},
			want: &RedisServer{
				Server: redis.NewServerFromParams("host1", "127.0.0.1", "1000", client.MustNewFromConnectionString("redis://127.0.0.1:1000")),
				Role:   client.Unknown,
				Config: map[string]string{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRedisServerFromPool(tt.args.connectionString, tt.args.alias, tt.args.pool)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServerFromParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("NewServerFromParams() = got diff %v", diff)
			}
		})
	}
}

func TestDiscoveryOptionSet_Has(t *testing.T) {
	type args struct {
		opt DiscoveryOption
	}
	tests := []struct {
		name string
		set  DiscoveryOptionSet
		args args
		want bool
	}{
		{
			name: "Returns true if option in slice",
			set:  DiscoveryOptionSet{SaveConfigDiscoveryOpt, SlaveReadOnlyDiscoveryOpt},
			args: args{opt: SlaveReadOnlyDiscoveryOpt},
			want: true,
		},
		{
			name: "Returns false if option not in slice",
			set:  DiscoveryOptionSet{SlaveReadOnlyDiscoveryOpt},
			args: args{opt: SaveConfigDiscoveryOpt},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.Has(tt.args.opt); got != tt.want {
				t.Errorf("DiscoveryOptions.Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisServer_Discover(t *testing.T) {
	type fields struct {
		Server *redis.Server
		Role   client.Role
		Config map[string]string
	}
	type args struct {
		ctx  context.Context
		opts []DiscoveryOption
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantRole   client.Role
		wantConfig map[string]string
		wantErr    bool
	}{
		{
			name: "Discovers a master",
			fields: fields{
				Server: redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
					client.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{"master", ""}
						},
						InjectError: func() error { return nil },
					},
				)},
			args:       args{ctx: context.TODO(), opts: DiscoveryOptionSet{}},
			wantRole:   client.Master,
			wantConfig: map[string]string{},
			wantErr:    false,
		},
		{
			name: "Discovers a master with save config",
			fields: fields{
				Server: redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
					client.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{"master", ""}
						},
						InjectError: func() error { return nil },
					},
					client.FakeResponse{
						// cmd: RedisConfigGet("save")
						InjectResponse: func() interface{} {
							return []interface{}{"save", "900 1 300 10"}
						},
						InjectError: func() error { return nil },
					},
				),
			},
			args:       args{ctx: context.TODO(), opts: DiscoveryOptionSet{SaveConfigDiscoveryOpt}},
			wantRole:   client.Master,
			wantConfig: map[string]string{"save": "900 1 300 10"},
			wantErr:    false,
		},
		{
			name: "Discovers a ro slave",
			fields: fields{
				Server: redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
					client.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{"slave", "127.0.0.1:3333"}
						},
						InjectError: func() error { return nil },
					},
					client.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{"read-only", "yes"}
						},
						InjectError: func() error { return nil },
					},
				),
			},
			args:       args{ctx: context.TODO(), opts: DiscoveryOptionSet{SlaveReadOnlyDiscoveryOpt}},
			wantRole:   client.Slave,
			wantConfig: map[string]string{"slave-read-only": "yes"},
			wantErr:    false,
		},
		{
			name: "Discovers a rw slave",
			fields: fields{
				Server: redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
					client.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{"slave", "127.0.0.1:3333"}
						},
						InjectError: func() error { return nil },
					},
					client.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{"read-only", "no"}
						},
						InjectError: func() error { return nil },
					},
				),
			},
			args:       args{ctx: context.TODO(), opts: DiscoveryOptionSet{SlaveReadOnlyDiscoveryOpt}},
			wantRole:   client.Slave,
			wantConfig: map[string]string{"slave-read-only": "no"},
			wantErr:    false,
		},
		{
			name: "'role' command fails, returns an error",
			fields: fields{
				Server: redis.NewFakeServerWithFakeClient("127.0.0.1", "1000",
					client.FakeResponse{
						InjectResponse: func() interface{} { return []interface{}{} },
						InjectError:    func() error { return errors.New("error") },
					},
				),
			},
			args:       args{ctx: context.TODO(), opts: DiscoveryOptionSet{}},
			wantRole:   client.Unknown,
			wantConfig: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &RedisServer{
				Server: tt.fields.Server,
				Role:   tt.fields.Role,
				Config: tt.fields.Config,
			}

			if err := srv.Discover(tt.args.ctx, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("RedisServer.Discover() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantRole != srv.Role {
				t.Errorf("RedisServer.Discover() got = %v, want %v", srv.Role, tt.wantRole)
			}
			if diff := deep.Equal(srv.Config, tt.wantConfig); len(diff) > 0 {
				t.Errorf("RedisServer.Discover() got diff: %v", diff)
			}
		})
	}
}
