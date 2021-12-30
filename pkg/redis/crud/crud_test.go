package crud

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/3scale/saas-operator/pkg/redis/crud/client"
	"github.com/go-redis/redis/v8"
	"github.com/go-test/deep"
)

func init() {
	deep.CompareUnexportedFields = true
}

func TestNewRedisCRUD(t *testing.T) {
	type args struct {
		connectionString string
	}
	tests := []struct {
		name    string
		args    args
		want    *CRUD
		wantErr bool
	}{
		{
			name: "Returns a CRUD object",
			args: args{
				connectionString: "redis://127.0.0.1:1234",
			},
			want: &CRUD{
				Client: func() Client { c, _ := client.NewFromConnectionString("redis://127.0.0.1:1234"); return c }(),
				IP:     "127.0.0.1",
				Port:   "1234",
			},
			wantErr: false,
		},
		{
			name: "Returns an error",
			args: args{
				connectionString: "127.0.0.1:1234",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRedisCRUDFromConnectionString(tt.args.connectionString)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRedisCRUD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("NewRedisCRUD() got diff: %v", diff)
			}
		})
	}
}

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
		{
			name: "Returns the server IP",
			fields: fields{
				client: nil,
				ip:     "127.0.0.1",
				port:   "2222",
			},
			want: "127.0.0.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &CRUD{
				Client: tt.fields.client,
				IP:     tt.fields.ip,
				Port:   tt.fields.port,
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
		{
			name: "Returns the server port",
			fields: fields{
				client: nil,
				ip:     "127.0.0.1",
				port:   "2222",
			},
			want: "2222",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &CRUD{
				Client: tt.fields.client,
				IP:     tt.fields.ip,
				Port:   tt.fields.port,
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
		{
			name: "Sends the 'master' command to a sentinel",
			fields: fields{
				client: &client.FakeClient{
					Responses: []client.FakeResponse{{
						InjectResponse: func() interface{} {
							return &client.SentinelMasterCmdResult{
								Name:                  "test",
								IP:                    "127.0.0.1",
								Port:                  6379,
								RunID:                 "b346b91b9492b4acd2cb7f04d466055ae1eca9b7",
								Flags:                 "master",
								LinkPendingCommands:   0,
								LinkRefcount:          1,
								LastPingSent:          0,
								LastOkPingReply:       424,
								LastPingReply:         424,
								DownAfterMilliseconds: 5000,
								InfoRefresh:           8877,
								RoleReported:          "master",
								RoleReportedTime:      141754639,
								ConfigEpoch:           0,
								NumSlaves:             2,
								NumOtherSentinels:     2,
								Quorum:                2,
								FailoverTimeout:       180000,
								ParallelSyncs:         1,
							}
						},
						InjectError: func() error { return nil },
					}},
				},
				ip:   "abc",
				port: "abc",
			},
			args: args{ctx: context.TODO(), shard: "x"},
			want: &client.SentinelMasterCmdResult{
				Name:                  "test",
				IP:                    "127.0.0.1",
				Port:                  6379,
				RunID:                 "b346b91b9492b4acd2cb7f04d466055ae1eca9b7",
				Flags:                 "master",
				LinkPendingCommands:   0,
				LinkRefcount:          1,
				LastPingSent:          0,
				LastOkPingReply:       424,
				LastPingReply:         424,
				DownAfterMilliseconds: 5000,
				InfoRefresh:           8877,
				RoleReported:          "master",
				RoleReportedTime:      141754639,
				ConfigEpoch:           0,
				NumSlaves:             2,
				NumOtherSentinels:     2,
				Quorum:                2,
				FailoverTimeout:       180000,
				ParallelSyncs:         1,
			},
			wantErr: false,
		},
		{
			name: "Returns an error",
			fields: fields{
				client: &client.FakeClient{
					Responses: []client.FakeResponse{{
						InjectResponse: func() interface{} { return &client.SentinelMasterCmdResult{} },
						InjectError:    func() error { return errors.New("error") },
					}},
				},
				ip:   "abc",
				port: "abc",
			},
			args:    args{ctx: context.TODO(), shard: "x"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &CRUD{
				Client: tt.fields.client,
				IP:     tt.fields.ip,
				Port:   tt.fields.port,
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
		{
			name: "Send the 'masters' command to a sentinel",
			fields: fields{
				client: &client.FakeClient{
					Responses: []client.FakeResponse{{
						InjectResponse: func() interface{} {
							return []interface{}{
								[]interface{}{
									"name", "shard01",
									"ip", "10.244.0.8",
									"port", "6379",
									"runid", "b346b91b9492b4acd2cb7f04d466055ae1eca9b7",
									"flags", "master",
									"link-pending-commands", "0",
									"link-refcount", "1",
									"last-ping-sent", "0",
									"last-ok-ping-reply", "424",
									"last-ping-reply", "424",
									"down-after-milliseconds", "5000",
									"info-refresh", "8877",
									"role-reported", "master",
									"role-reported-time", "141754639",
									"config-epoch", "0",
									"num-slaves", "2",
									"num-other-sentinels", "2",
									"quorum", "2",
									"failover-timeout", "180000",
									"parallel-syncs", "1",
								},
								[]interface{}{
									"name", "shard02",
									"ip", "10.244.0.10",
									"port", "6379",
									"runid", "0493dfc108d3becd49da2f47695d0b11c442f921",
									"flags", "master",
									"link-pending-commands", "0",
									"link-refcount", "1",
									"last-ping-sent", "0",
									"last-ok-ping-reply", "424",
									"last-ping-reply", "424",
									"down-after-milliseconds", "5000",
									"info-refresh", "8877",
									"role-reported", "master",
									"role-reported-time", "141754633",
									"config-epoch", "0",
									"num-slaves", "2",
									"num-other-sentinels", "2",
									"quorum", "2",
									"failover-timeout", "180000",
									"parallel-syncs", "1",
								},
							}
						},
						InjectError: func() error { return nil },
					}},
				},
				ip:   "abc",
				port: "abc",
			},
			args: args{ctx: context.TODO()},
			want: []client.SentinelMasterCmdResult{
				{
					Name:                  "shard01",
					IP:                    "10.244.0.8",
					Port:                  6379,
					RunID:                 "b346b91b9492b4acd2cb7f04d466055ae1eca9b7",
					Flags:                 "master",
					LinkPendingCommands:   0,
					LinkRefcount:          1,
					LastPingSent:          0,
					LastOkPingReply:       424,
					LastPingReply:         424,
					DownAfterMilliseconds: 5000,
					InfoRefresh:           8877,
					RoleReported:          "master",
					RoleReportedTime:      141754639,
					ConfigEpoch:           0,
					NumSlaves:             2,
					NumOtherSentinels:     2,
					Quorum:                2,
					FailoverTimeout:       180000,
					ParallelSyncs:         1,
				},
				{
					Name:                  "shard02",
					IP:                    "10.244.0.10",
					Port:                  6379,
					RunID:                 "0493dfc108d3becd49da2f47695d0b11c442f921",
					Flags:                 "master",
					LinkPendingCommands:   0,
					LinkRefcount:          1,
					LastPingSent:          0,
					LastOkPingReply:       424,
					LastPingReply:         424,
					DownAfterMilliseconds: 5000,
					InfoRefresh:           8877,
					RoleReported:          "master",
					RoleReportedTime:      141754633,
					ConfigEpoch:           0,
					NumSlaves:             2,
					NumOtherSentinels:     2,
					Quorum:                2,
					FailoverTimeout:       180000,
					ParallelSyncs:         1,
				},
			},
			wantErr: false,
		},
		{
			name: "Returns an error",
			fields: fields{
				client: &client.FakeClient{
					Responses: []client.FakeResponse{{
						InjectResponse: func() interface{} { return []interface{}{} },
						InjectError:    func() error { return errors.New("error") },
					}},
				},
				ip:   "abc",
				port: "abc",
			},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &CRUD{
				Client: tt.fields.client,
				IP:     tt.fields.ip,
				Port:   tt.fields.port,
			}
			got, err := sc.SentinelMasters(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SentinelMasters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.SentinelMasters() = %+v, want %+v", got, tt.want)
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
		{
			name: "Sends the 'slaves' command to sentinel",
			fields: fields{
				client: &client.FakeClient{
					Responses: []client.FakeResponse{{
						InjectResponse: func() interface{} {
							return []interface{}{
								[]interface{}{
									"name", "10.244.0.6:6379",
									"ip", "10.244.0.6",
									"port", "6379",
									"runid", "aa63203ce2d7243c819a5c26d7105ab92b897a0c",
									"flags", "slave",
									"link-pending-commands", "0",
									"link-refcount", "1",
									"last-ping-sent", "0",
									"last-ok-ping-reply", "168",
									"last-ping-reply", "168",
									"down-after-milliseconds", "5000",
									"info-refresh", "9517",
									"role-reported", "slave",
									"role-reported-time", "155699495",
									"master-link-down-time", "0",
									"master-link-status", "ok",
									"master-host", "10.244.0.8",
									"master-port", "6379",
									"slave-priority", "100",
									"slave-repl-offset", "11218922",
								},
								[]interface{}{
									"name", "10.244.0.7:6379",
									"ip", "10.244.0.7",
									"port", "6379",
									"runid", "678106ab45d26fdb254e1e276ecd84999f9f969f",
									"flags", "slave",
									"link-pending-commands", "0",
									"link-refcount", "1",
									"last-ping-sent", "0",
									"last-ok-ping-reply", "168",
									"last-ping-reply", "168",
									"down-after-milliseconds", "5000",
									"info-refresh", "9516",
									"role-reported", "slave",
									"role-reported-time", "155699496",
									"master-link-down-time", "0",
									"master-link-status", "ok",
									"master-host", "10.244.0.8",
									"master-port", "6379",
									"slave-priority", "100",
									"slave-repl-offset", "11218922",
								},
							}
						},
						InjectError: func() error { return nil },
					}},
				},
				ip:   "abc",
				port: "abc",
			}, args: args{},
			want: []client.SentinelSlaveCmdResult{
				{
					Name:                  "10.244.0.6:6379",
					IP:                    "10.244.0.6",
					Port:                  6379,
					RunID:                 "aa63203ce2d7243c819a5c26d7105ab92b897a0c",
					Flags:                 "slave",
					LinkPendingCommands:   0,
					LinkRefcount:          1,
					LastPingSent:          0,
					LastOkPingReply:       168,
					LastPingReply:         168,
					DownAfterMilliseconds: 5000,
					InfoRefresh:           9517,
					RoleReported:          "slave",
					RoleReportedTime:      155699495,
					MasterLinkDownTime:    0,
					MasterLinkStatus:      "ok",
					MasterHost:            "10.244.0.8",
					MasterPort:            6379,
					SlavePriority:         100,
					SlaveReplOffset:       11218922,
				},
				{
					Name:                  "10.244.0.7:6379",
					IP:                    "10.244.0.7",
					Port:                  6379,
					RunID:                 "678106ab45d26fdb254e1e276ecd84999f9f969f",
					Flags:                 "slave",
					LinkPendingCommands:   0,
					LinkRefcount:          1,
					LastPingSent:          0,
					LastOkPingReply:       168,
					LastPingReply:         168,
					DownAfterMilliseconds: 5000,
					InfoRefresh:           9516,
					RoleReported:          "slave",
					RoleReportedTime:      155699496,
					MasterLinkDownTime:    0,
					MasterLinkStatus:      "ok",
					MasterHost:            "10.244.0.8",
					MasterPort:            6379,
					SlavePriority:         100,
					SlaveReplOffset:       11218922,
				},
			},
			wantErr: false,
		},
		{
			name: "Returns an error",
			fields: fields{
				client: &client.FakeClient{
					Responses: []client.FakeResponse{{
						InjectResponse: func() interface{} { return []interface{}{} },
						InjectError:    func() error { return errors.New("error") },
					}},
				},
				ip:   "abc",
				port: "abc",
			},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &CRUD{
				Client: tt.fields.client,
				IP:     tt.fields.ip,
				Port:   tt.fields.port,
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
		{
			name: "Returns ok",
			fields: fields{
				client: &client.FakeClient{
					Responses: []client.FakeResponse{{
						InjectResponse: func() interface{} { return nil },
						InjectError:    func() error { return nil },
					}},
				},
				ip:   "abc",
				port: "abc",
			},
			args: args{
				ctx: context.TODO(), name: "abc", host: "abc", port: "abc", quorum: 1,
			},
			wantErr: false,
		},
		{
			name: "Returns error",
			fields: fields{
				client: &client.FakeClient{
					Responses: []client.FakeResponse{{
						InjectResponse: func() interface{} { return nil },
						InjectError:    func() error { return errors.New("error") },
					}},
				},
				ip:   "abc",
				port: "abc",
			},
			args: args{
				ctx: context.TODO(), name: "abc", host: "abc", port: "abc", quorum: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &CRUD{
				Client: tt.fields.client,
				IP:     tt.fields.ip,
				Port:   tt.fields.port,
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
		{
			name: "Returns ok",
			fields: fields{
				client: &client.FakeClient{
					Responses: []client.FakeResponse{{
						InjectResponse: func() interface{} { return nil },
						InjectError:    func() error { return nil },
					}},
				},
				ip:   "abc",
				port: "abc",
			},
			args:    args{},
			wantErr: false,
		},
		{
			name: "Returns error",
			fields: fields{
				client: &client.FakeClient{
					Responses: []client.FakeResponse{{
						InjectResponse: func() interface{} { return nil },
						InjectError:    func() error { return errors.New("error") },
					}},
				},
				ip:   "abc",
				port: "abc",
			},
			args: args{
				ctx:       nil,
				shard:     "abc",
				parameter: "abc",
				value:     "abc",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &CRUD{
				Client: tt.fields.client,
				IP:     tt.fields.ip,
				Port:   tt.fields.port,
			}
			if err := sc.SentinelSet(tt.args.ctx, tt.args.shard, tt.args.parameter, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Client.SentinelSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCRUD_SentinelPSubscribe(t *testing.T) {
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
		want   string
	}{
		{
			name: "Returns an events channel",
			fields: fields{
				client: &client.FakeClient{
					Responses: []client.FakeResponse{{
						InjectResponse: func() interface{} {
							ch := make(chan *redis.Message)
							go func() {
								ch <- &redis.Message{
									Channel: "events",
									Payload: "this is a test",
								}
							}()
							var roCh <-chan *redis.Message = ch
							return roCh
						},
						InjectError: func() error { return nil },
					}},
				},
				ip:   "abc",
				port: "abc",
			},
			args: args{
				ctx:    context.TODO(),
				events: []string{"event"},
			},
			want: "Message<events: this is a test>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crud := &CRUD{
				Client: tt.fields.client,
				IP:     tt.fields.ip,
				Port:   tt.fields.port,
			}
			timeout := time.After(100 * time.Millisecond)
			done := make(chan bool)
			go func() {
				ch, _ := crud.SentinelPSubscribe(tt.args.ctx, tt.args.events...)
				got := <-ch
				if got.String() != tt.want {
					t.Errorf("CRUD.SentinelPSubscribe() got = %v, want %v", got.String(), tt.want)
				}
				done <- true
			}()

			select {
			case <-timeout:
				t.Fatal("CRUD.SentinelPSubscribe() didn't finish in time")
			case <-done:
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
		{
			name: "Server is a slave",
			fields: fields{
				client: &client.FakeClient{
					Responses: []client.FakeResponse{{
						InjectResponse: func() interface{} {
							return []interface{}{
								"slave", "10.244.0.8", 6379, "connected", 12046989,
							}
						},
						InjectError: func() error { return nil },
					}},
				},
				ip:   "abc",
				port: "abc",
			},
			args:    args{ctx: context.TODO()},
			want:    client.Slave,
			want1:   "10.244.0.8",
			wantErr: false,
		},
		{
			name: "Server is a master",
			fields: fields{
				client: &client.FakeClient{
					Responses: []client.FakeResponse{{
						InjectResponse: func() interface{} {
							return []interface{}{
								"master",
								12204347,
								[]interface{}{
									[]interface{}{
										"10.244.0.9",
										"6379",
										"12204211",
									},
									[]interface{}{
										"10.244.0.11",
										"6379",
										"12204211",
									},
								},
							}
						},
						InjectError: func() error { return nil },
					}},
				},
				ip:   "abc",
				port: "abc",
			},
			args:    args{ctx: context.TODO()},
			want:    client.Master,
			want1:   "",
			wantErr: false,
		},
		{
			name: "Returns an error",
			fields: fields{
				client: &client.FakeClient{
					Responses: []client.FakeResponse{{
						InjectResponse: func() interface{} {
							return []interface{}{}
						},
						InjectError: func() error { return errors.New("error") },
					}},
				},
				ip:   "abc",
				port: "abc",
			},
			args:    args{ctx: context.TODO()},
			want:    client.Unknown,
			want1:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &CRUD{
				Client: tt.fields.client,
				IP:     tt.fields.ip,
				Port:   tt.fields.port,
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
		{
			name: "Returns a redis configuration parameter",
			fields: fields{
				client: &client.FakeClient{
					Responses: []client.FakeResponse{{
						InjectResponse: func() interface{} {
							return []interface{}{
								"param1", "value1",
							}
						},
						InjectError: func() error { return nil },
					}},
				},
				ip:   "abc",
				port: "abc",
			},
			args: args{
				ctx:       context.TODO(),
				parameter: "param1",
			},
			want:    "value1",
			wantErr: false,
		},
		{
			name: "Returns an error",
			fields: fields{
				client: &client.FakeClient{
					Responses: []client.FakeResponse{{
						InjectResponse: func() interface{} { return []interface{}{} },
						InjectError:    func() error { return errors.New("error") },
					}},
				},
				ip:   "abc",
				port: "abc",
			},
			args: args{
				ctx:       context.TODO(),
				parameter: "param1",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &CRUD{
				Client: tt.fields.client,
				IP:     tt.fields.ip,
				Port:   tt.fields.port,
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
		{
			name: "Configures the redis server as a slave of the specified server",
			fields: fields{
				client: &client.FakeClient{
					Responses: []client.FakeResponse{{
						InjectResponse: func() interface{} { return nil },
						InjectError:    func() error { return nil },
					}},
				},
				ip:   "abc",
				port: "abc",
			},
			args:    args{ctx: context.TODO(), host: "abc", port: "abc"},
			wantErr: false,
		},
		{
			name: "Returns an error",
			fields: fields{
				client: &client.FakeClient{
					Responses: []client.FakeResponse{{
						InjectResponse: func() interface{} { return nil },
						InjectError:    func() error { return errors.New("error") },
					}},
				},
				ip:   "abc",
				port: "abc",
			},
			args:    args{ctx: context.TODO(), host: "abc", port: "abc"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &CRUD{
				Client: tt.fields.client,
				IP:     tt.fields.ip,
				Port:   tt.fields.port,
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
			name: "Parses a redis-go administrative command response",
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
				LastPingSent:          0,
				LastOkPingReply:       68,
				LastPingReply:         68,
				DownAfterMilliseconds: 5000,
				InfoRefresh:           2035,
				RoleReported:          "slave",
				RoleReportedTime:      91031035,
				MasterLinkDownTime:    0,
				MasterLinkStatus:      "ok",
				MasterHost:            "10.244.0.8",
				MasterPort:            6379,
				SlavePriority:         100,
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
