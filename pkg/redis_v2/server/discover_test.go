package server

import (
	"context"
	"errors"
	"testing"

	"github.com/3scale/saas-operator/pkg/redis_v2/client"
	"github.com/google/go-cmp/cmp"
)

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

func TestServer_Discover(t *testing.T) {
	type fields struct {
		Alias  string
		Client client.TestableInterface
		Host   string
		Port   string
		Role   client.Role
		Config map[string]string
	}
	type args struct {
		ctx  context.Context
		opts DiscoveryOptionSet
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
				Client: client.NewFakeClient(
					client.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{"master", ""}
						},
						InjectError: func() error { return nil },
					},
				),
			},
			args:       args{ctx: context.TODO(), opts: DiscoveryOptionSet{RoleDiscoveryOpt}},
			wantRole:   client.Master,
			wantConfig: map[string]string{},
			wantErr:    false,
		},
		{
			name: "Discovers a master with save config",
			fields: fields{
				Client: client.NewFakeClient(
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
			args:       args{ctx: context.TODO(), opts: DiscoveryOptionSet{RoleDiscoveryOpt, SaveConfigDiscoveryOpt}},
			wantRole:   client.Master,
			wantConfig: map[string]string{"save": "900 1 300 10"},
			wantErr:    false,
		},
		{
			name: "Discovers a ro slave",
			fields: fields{
				Client: client.NewFakeClient(
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
			args:       args{ctx: context.TODO(), opts: DiscoveryOptionSet{RoleDiscoveryOpt, SlaveReadOnlyDiscoveryOpt}},
			wantRole:   client.Slave,
			wantConfig: map[string]string{"slave-read-only": "yes"},
			wantErr:    false,
		},
		{
			name: "Discovers a rw slave",
			fields: fields{
				Client: client.NewFakeClient(
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
			args:       args{ctx: context.TODO(), opts: DiscoveryOptionSet{RoleDiscoveryOpt, SlaveReadOnlyDiscoveryOpt}},
			wantRole:   client.Slave,
			wantConfig: map[string]string{"slave-read-only": "no"},
			wantErr:    false,
		},
		{
			name: "'role' command fails, returns an error",
			fields: fields{
				Client: client.NewFakeClient(
					client.FakeResponse{
						InjectResponse: func() interface{} { return []interface{}{} },
						InjectError:    func() error { return errors.New("error") },
					},
				),
			},
			args:       args{ctx: context.TODO(), opts: DiscoveryOptionSet{RoleDiscoveryOpt}},
			wantRole:   client.Unknown,
			wantConfig: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Server{
				Alias:  tt.fields.Alias,
				Host:   tt.fields.Host,
				Port:   tt.fields.Port,
				Client: tt.fields.Client,
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
			if diff := cmp.Diff(srv.Config, tt.wantConfig); len(diff) > 0 {
				t.Errorf("RedisServer.Discover() got diff: %v", diff)
			}
		})
	}
}
