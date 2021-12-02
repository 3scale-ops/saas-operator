package crud

import (
	"context"
	"reflect"
	"testing"

	"github.com/3scale/saas-operator/pkg/redis/crud/client"
	"github.com/go-redis/redis/v8"
)

func TestClient_GetIP(t *testing.T) {
	type fields struct {
		client Client
		ip     string
		port   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &CRUD{
				client: tt.fields.client,
				ip:     tt.fields.ip,
				port:   tt.fields.port,
			}
			if got := sc.GetIP(); got != tt.want {
				t.Errorf("Client.GetIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetPort(t *testing.T) {
	type fields struct {
		client Client
		ip     string
		port   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &CRUD{
				client: tt.fields.client,
				ip:     tt.fields.ip,
				port:   tt.fields.port,
			}
			if got := sc.GetPort(); got != tt.want {
				t.Errorf("Client.GetPort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_SentinelMaster(t *testing.T) {
	type fields struct {
		client Client
		ip     string
		port   string
	}
	type args struct {
		ctx   context.Context
		shard string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *client.SentinelMasterCmdResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &CRUD{
				client: tt.fields.client,
				ip:     tt.fields.ip,
				port:   tt.fields.port,
			}
			got, err := sc.SentinelMaster(tt.args.ctx, tt.args.shard)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SentinelMaster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.SentinelMaster() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_SentinelMasters(t *testing.T) {
	type fields struct {
		client Client
		ip     string
		port   string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []client.SentinelMasterCmdResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &CRUD{
				client: tt.fields.client,
				ip:     tt.fields.ip,
				port:   tt.fields.port,
			}
			got, err := sc.SentinelMasters(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SentinelMasters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.SentinelMasters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_SentinelSlaves(t *testing.T) {
	type fields struct {
		client Client
		ip     string
		port   string
	}
	type args struct {
		ctx   context.Context
		shard string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []client.SentinelSlaveCmdResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &CRUD{
				client: tt.fields.client,
				ip:     tt.fields.ip,
				port:   tt.fields.port,
			}
			got, err := sc.SentinelSlaves(tt.args.ctx, tt.args.shard)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SentinelSlaves() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.SentinelSlaves() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_SentinelMonitor(t *testing.T) {
	type fields struct {
		client Client
		ip     string
		port   string
	}
	type args struct {
		ctx    context.Context
		name   string
		host   string
		port   string
		quorum int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &CRUD{
				client: tt.fields.client,
				ip:     tt.fields.ip,
				port:   tt.fields.port,
			}
			if err := sc.SentinelMonitor(tt.args.ctx, tt.args.name, tt.args.host, tt.args.port, tt.args.quorum); (err != nil) != tt.wantErr {
				t.Errorf("Client.SentinelMonitor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_SentinelSet(t *testing.T) {
	type fields struct {
		client Client
		ip     string
		port   string
	}
	type args struct {
		ctx       context.Context
		shard     string
		parameter string
		value     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &CRUD{
				client: tt.fields.client,
				ip:     tt.fields.ip,
				port:   tt.fields.port,
			}
			if err := sc.SentinelSet(tt.args.ctx, tt.args.shard, tt.args.parameter, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Client.SentinelSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_SentinelPSubscribe(t *testing.T) {
	type fields struct {
		client Client
		ip     string
		port   string
	}
	type args struct {
		ctx    context.Context
		events []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   <-chan *redis.Message
		want1  func() error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &CRUD{
				client: tt.fields.client,
				ip:     tt.fields.ip,
				port:   tt.fields.port,
			}
			got, got1 := sc.SentinelPSubscribe(tt.args.ctx, tt.args.events...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.SentinelPSubscribe() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Client.SentinelPSubscribe() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestClient_RedisRole(t *testing.T) {
	type fields struct {
		client Client
		ip     string
		port   string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    client.Role
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &CRUD{
				client: tt.fields.client,
				ip:     tt.fields.ip,
				port:   tt.fields.port,
			}
			got, got1, err := sc.RedisRole(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.RedisRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.RedisRole() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Client.RedisRole() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestClient_RedisConfigGet(t *testing.T) {
	type fields struct {
		client Client
		ip     string
		port   string
	}
	type args struct {
		ctx       context.Context
		parameter string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &CRUD{
				client: tt.fields.client,
				ip:     tt.fields.ip,
				port:   tt.fields.port,
			}
			got, err := sc.RedisConfigGet(tt.args.ctx, tt.args.parameter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.RedisConfigGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.RedisConfigGet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_RedisSlaveOf(t *testing.T) {
	type fields struct {
		client Client
		ip     string
		port   string
	}
	type args struct {
		ctx  context.Context
		host string
		port string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &CRUD{
				client: tt.fields.client,
				ip:     tt.fields.ip,
				port:   tt.fields.port,
			}
			if err := sc.RedisSlaveOf(tt.args.ctx, tt.args.host, tt.args.port); (err != nil) != tt.wantErr {
				t.Errorf("Client.RedisSlaveOf() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_sliceCmdToStruct(t *testing.T) {
	type args struct {
		in  []interface{}
		out interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "holy crap",
			args: args{
				in: []interface{}{
					"name",
					"10.244.0.7:6379",
					"ip",
					"10.244.0.7",
					"port",
					"6379",
					"runid",
					"4009bfff0cb1115e8b5d06e00730808a051fa2c5",
					"flags",
					"slave",
					"link-pending-commands",
					"0",
					"link-refcount",
					"1",
					"last-ping-sent",
					"0",
					"last-ok-ping-reply",
					"68",
					"last-ping-reply",
					"68",
					"down-after-milliseconds",
					"5000",
					"info-refresh",
					"2035",
					"role-reported",
					"slave",
					"role-reported-time",
					"91031035",
					"master-link-down-time",
					"0",
					"master-link-status",
					"ok",
					"master-host",
					"10.244.0.8",
					"master-port",
					"6379",
					"slave-priority",
					"100",
					"slave-repl-offset",
					"12707378",
				},
				out: &client.SentinelSlaveCmdResult{},
			},
			want: &client.SentinelSlaveCmdResult{
				Name:                  "10.244.0.7:6379",
				IP:                    "10.244.0.7",
				Port:                  6379,
				RunID:                 "4009bfff0cb1115e8b5d06e00730808a051fa2c5",
				Flags:                 "slave",
				LinkPendingCommands:   0,
				LinkRefcount:          1,
				LastPingSet:           0,
				LastOkPingReply:       68,
				LastPingReply:         68,
				DownAfterMilliseconds: 5000,
				InfoRefresh:           2035,
				RoleReported:          "slave",
				RoleReportedTime:      91031035,
				MasterLinkDownTime:    0,
				MasterLinkStatus:      "ok",
				MasterHost:            "10.244.0.8",
				MasterPort:            "6379",
				SlavePriority:         "100",
				SlaveReplOffset:       12707378,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sliceCmdToStruct(tt.args.in, tt.args.out)
			if (err != nil) != tt.wantErr {
				t.Errorf("sliceCmdToStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.args.out, tt.want) {
				t.Errorf("sliceCmdToStruct() = %v, want %v", tt.args.out, tt.want)
			}
		})
	}
}
