package sharded

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/3scale/saas-operator/pkg/redis_v2/client"
	redis "github.com/3scale/saas-operator/pkg/redis_v2/server"
	"github.com/3scale/saas-operator/pkg/util"
	"github.com/go-test/deep"
)

var (
	testShardedCluster *Cluster = &Cluster{
		Shards: []Shard{
			{
				Name: "shard00",
				Servers: []*RedisServer{
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.1:2000", util.Pointer("shard00-0")),
						client.Master,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.1:2001", util.Pointer("shard00-1")),
						client.Slave,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.1:2002", util.Pointer("shard00-2")),
						client.Slave,
						map[string]string{},
					),
				},
			},
			{
				Name: "shard01",
				Servers: []*RedisServer{
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.1:3000", util.Pointer("shard01-0")),
						client.Master,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.1:3001", util.Pointer("shard01-1")),
						client.Slave,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.1:3002", util.Pointer("shard01-2")),
						client.Slave,
						map[string]string{},
					),
				},
			},
			{
				Name: "shard02",
				Servers: []*RedisServer{
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.1:4000", util.Pointer("shard02-0")),
						client.Master,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.1:4001", util.Pointer("shard02-1")),
						client.Slave,
						map[string]string{},
					),
					NewRedisServerFromParams(
						redis.MustNewServer("redis://127.0.0.1:4002", util.Pointer("shard02-2")),
						client.Slave,
						map[string]string{},
					),
				},
			},
		},
	}
)

func init() {
	deep.CompareUnexportedFields = true
}

func TestNewSentinelServerFromPool(t *testing.T) {
	type args struct {
		connectionString string
		alias            *string
		pool             *redis.ServerPool
	}
	tests := []struct {
		name    string
		args    args
		want    *SentinelServer
		wantErr bool
	}{
		{
			name: "Returns a SentinelServer",
			args: args{
				connectionString: "redis://127.0.0.1:1000",
				alias:            util.Pointer("sentinel"),
				pool:             &redis.ServerPool{},
			},
			want: &SentinelServer{
				Server: redis.MustNewServer("redis://127.0.0.1:1000", util.Pointer("sentinel")),
			},
			wantErr: false,
		},
		{
			name: "Gets server from pool",
			args: args{
				connectionString: "redis://127.0.0.1:1000",
				alias:            nil,
				pool:             redis.NewServerPool(redis.MustNewServer("redis://127.0.0.1:1000", util.Pointer("sentinel"))),
			},
			want: &SentinelServer{
				Server: redis.MustNewServer("redis://127.0.0.1:1000", util.Pointer("sentinel")),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSentinelServerFromPool(tt.args.connectionString, tt.args.alias, tt.args.pool)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSentinelServerFromPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("NewSentinelServerFromPool() = got diff %v", diff)
			}
		})
	}
}

func TestNewHighAvailableSentinel(t *testing.T) {
	type args struct {
		servers map[string]string
		pool    *redis.ServerPool
	}
	tests := []struct {
		name    string
		args    args
		want    []*SentinelServer
		wantErr bool
	}{
		{
			name: "Returns a list of sentinels",
			args: args{
				servers: map[string]string{
					"sentinel-0": "redis://127.0.0.1:1000",
					"sentinel-1": "redis://127.0.0.1:2000",
				},
				pool: &redis.ServerPool{},
			},
			want: []*SentinelServer{
				{Server: redis.MustNewServer("redis://127.0.0.1:1000", util.Pointer("sentinel-0"))},
				{Server: redis.MustNewServer("redis://127.0.0.1:2000", util.Pointer("sentinel-1"))},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHighAvailableSentinel(tt.args.servers, tt.args.pool)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHighAvailableSentinel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("NewHighAvailableSentinel() = got diff %v", diff)
			}
		})
	}
}

func TestSentinelServer_IsMonitoringShards(t *testing.T) {
	type args struct {
		ctx    context.Context
		shards []string
	}
	tests := []struct {
		name    string
		ss      *SentinelServer
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "All shards monitored by SentinelServer",
			ss: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
				client.FakeResponse{
					InjectResponse: func() interface{} {
						return []interface{}{
							[]interface{}{"name", "shard01"},
							[]interface{}{"name", "shard02"},
						}
					},
					InjectError: func() error { return nil },
				})),
			args: args{
				ctx:    context.TODO(),
				shards: []string{"shard01", "shard02"},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "No shard monitored",
			ss: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
				client.FakeResponse{
					InjectResponse: func() interface{} { return []interface{}{} },
					InjectError:    func() error { return nil },
				})),
			args: args{
				ctx:    context.TODO(),
				shards: []string{"shard01", "shard02"},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "One shard is not monitored",
			ss: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
				client.FakeResponse{
					InjectResponse: func() interface{} {
						return []interface{}{
							[]interface{}{"name", "shard01"},
						}
					},
					InjectError: func() error { return nil },
				})),
			args: args{
				ctx:    context.TODO(),
				shards: []string{"shard01", "shard02"},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Returns an error",
			ss: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
				client.FakeResponse{
					InjectResponse: func() interface{} { return []interface{}{} },
					InjectError:    func() error { return errors.New("error") },
				})),
			args: args{
				ctx:    context.TODO(),
				shards: []string{"shard01", "shard02"},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ss.IsMonitoringShards(tt.args.ctx, tt.args.shards)
			if (err != nil) != tt.wantErr {
				t.Errorf("SentinelServer.IsMonitoringShards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SentinelServer.IsMonitoringShards() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSentinelServer_Monitor(t *testing.T) {
	type args struct {
		ctx    context.Context
		shards *Cluster
	}
	tests := []struct {
		name    string
		ss      *SentinelServer
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "All shards monitored",
			ss: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
				// SentinelMaster response for shard00
				client.FakeResponse{
					InjectResponse: func() interface{} {
						return &client.SentinelMasterCmdResult{
							Name: "shard00",
							IP:   "127.0.0.1",
							Port: 2000,
						}
					},
					InjectError: func() error { return nil },
				},
				// SentinelMaster response for shard01
				client.FakeResponse{
					InjectResponse: func() interface{} {
						return &client.SentinelMasterCmdResult{
							Name: "shard01",
							IP:   "127.0.0.1",
							Port: 3000,
						}
					},
					InjectError: func() error { return nil },
				},
				// SentinelMaster response for shard02
				client.FakeResponse{
					InjectResponse: func() interface{} {
						return &client.SentinelMasterCmdResult{
							Name: "shard02",
							IP:   "127.0.0.1",
							Port: 4000,
						}
					},
					InjectError: func() error { return nil },
				},
			)),
			args: args{
				ctx:    context.TODO(),
				shards: testShardedCluster,
			},
			want:    []string{},
			wantErr: false,
		},
		{
			name: "shard01 is not monitored",
			ss: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
				// SentinelMaster response for shard00
				client.FakeResponse{
					InjectResponse: func() interface{} {
						return &client.SentinelMasterCmdResult{
							Name: "shard00",
							IP:   "127.0.0.1",
							Port: 2000,
						}
					},
					InjectError: func() error { return nil },
				},
				// SentinelMaster response for shard01 (returns error as it is unmonitored)
				client.FakeResponse{
					InjectResponse: func() interface{} { return &client.SentinelMasterCmdResult{} },
					InjectError:    func() error { return errors.New(shardNotInitializedError) },
				},
				// SentinelMonitor response for shard01
				client.FakeResponse{
					InjectResponse: nil,
					InjectError:    func() error { return nil },
				},
				// SentinelSet response for shard01
				client.FakeResponse{
					InjectResponse: nil,
					InjectError:    func() error { return nil },
				},
				// SentinelMaster response for shard02
				client.FakeResponse{
					InjectResponse: func() interface{} {
						return &client.SentinelMasterCmdResult{
							Name: "shard02",
							IP:   "127.0.0.1",
							Port: 4000,
						}
					},
					InjectError: func() error { return nil },
				},
			)),
			args: args{
				ctx:    context.TODO(),
				shards: testShardedCluster,
			},
			want:    []string{"shard01"},
			wantErr: false,
		},
		{
			name: "all shards are unmonitored",
			ss: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
				// SentinelMaster response for shard00 (returns error as it is unmonitored)
				client.FakeResponse{
					InjectResponse: func() interface{} { return &client.SentinelMasterCmdResult{} },
					InjectError:    func() error { return errors.New(shardNotInitializedError) },
				},
				// SentinelMonitor response for shard00
				client.FakeResponse{
					InjectResponse: nil,
					InjectError:    func() error { return nil },
				},
				// SentinelSet response for shard00
				client.FakeResponse{
					InjectResponse: nil,
					InjectError:    func() error { return nil },
				},
				// SentinelMaster response for shard01 (returns error as it is unmonitored)
				client.FakeResponse{
					InjectResponse: func() interface{} { return &client.SentinelMasterCmdResult{} },
					InjectError:    func() error { return errors.New(shardNotInitializedError) },
				},
				// SentinelMonitor response for shard01
				client.FakeResponse{
					InjectResponse: nil,
					InjectError:    func() error { return nil },
				},
				// SentinelSet response for shard01
				client.FakeResponse{
					InjectResponse: nil,
					InjectError:    func() error { return nil },
				},
				// SentinelMaster response for shard02 (returns error as it is unmonitored)
				client.FakeResponse{
					InjectResponse: func() interface{} { return &client.SentinelMasterCmdResult{} },
					InjectError:    func() error { return errors.New(shardNotInitializedError) },
				},
				// SentinelMonitor response for shard02
				client.FakeResponse{
					InjectResponse: nil,
					InjectError:    func() error { return nil },
				},
				// SentinelSet response for shard02
				client.FakeResponse{
					InjectResponse: nil,
					InjectError:    func() error { return nil },
				},
			)),
			args: args{
				ctx:    context.TODO(),
				shards: testShardedCluster,
			},
			want:    []string{"shard00", "shard01", "shard02"},
			wantErr: false,
		},
		{
			name: "All shards unmonitored, failure on the 2nd one",
			ss: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
				client.FakeResponse{
					InjectResponse: func() interface{} { return &client.SentinelMasterCmdResult{} },
					InjectError:    func() error { return errors.New(shardNotInitializedError) },
				},
				// SentinelMonitor response for shard00
				client.FakeResponse{
					InjectResponse: nil,
					InjectError:    func() error { return nil },
				},
				// SentinelSet response for shard00
				client.FakeResponse{
					InjectResponse: nil,
					InjectError:    func() error { return nil },
				},
				// SentinelMaster response for shard01 (returns error as it is unmonitored)
				client.FakeResponse{
					InjectResponse: func() interface{} { return &client.SentinelMasterCmdResult{} },
					InjectError:    func() error { return errors.New("error") },
				},
				// SentinelMaster response for shard02 (returns error as it is unmonitored)
				client.FakeResponse{
					InjectResponse: func() interface{} { return &client.SentinelMasterCmdResult{} },
					InjectError:    func() error { return errors.New(shardNotInitializedError) },
				},
				// SentinelMonitor response for shard02
				client.FakeResponse{
					InjectResponse: nil,
					InjectError:    func() error { return nil },
				},
				// SentinelSet response for shard02
				client.FakeResponse{
					InjectResponse: nil,
					InjectError:    func() error { return nil },
				},
			)),
			args: args{
				ctx:    context.TODO(),
				shards: testShardedCluster,
			},
			want:    []string{"shard00"},
			wantErr: true,
		},
		{
			name: "All shards monitored, failure on the 2nd one",
			ss: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
				// SentinelMaster response for shard00 (returns error as it is unmonitored)
				client.FakeResponse{
					InjectResponse: func() interface{} {
						return &client.SentinelMasterCmdResult{
							Name: "shard00",
							IP:   "127.0.0.1",
							Port: 2000,
						}
					},
					InjectError: func() error { return nil },
				},
				// SentinelMaster response for shard01 (returns error as it is unmonitored)
				client.FakeResponse{
					InjectResponse: func() interface{} { return &client.SentinelMasterCmdResult{} },
					InjectError:    func() error { return errors.New("error") },
				},
			)),
			args: args{
				ctx:    context.TODO(),
				shards: testShardedCluster,
			},
			want:    []string{},
			wantErr: true,
		},
		{
			name: "'sentinel monitor' fails for shard00, returns no shards changed",
			ss: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
				// SentinelMaster response for shard00 (returns error as it is unmonitored)
				client.FakeResponse{
					InjectResponse: func() interface{} { return &client.SentinelMasterCmdResult{} },
					InjectError:    func() error { return errors.New(shardNotInitializedError) },
				},
				// SentinelMonitor response for shard00
				client.FakeResponse{
					InjectResponse: nil,
					InjectError:    func() error { return errors.New("error") },
				},
			)),
			args: args{
				ctx:    context.TODO(),
				shards: testShardedCluster,
			},
			want:    []string{},
			wantErr: true,
		},
		{
			name: "Error writing config param, returns shard00 changed",
			ss: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
				// SentinelMaster response for shard00 (returns error as it is unmonitored)
				client.FakeResponse{
					InjectResponse: func() interface{} { return &client.SentinelMasterCmdResult{} },
					InjectError:    func() error { return errors.New(shardNotInitializedError) },
				},
				// SentinelMonitor response for shard00
				client.FakeResponse{
					InjectResponse: nil,
					InjectError:    func() error { return nil },
				},
				// SentinelSet response for shard01
				client.FakeResponse{
					InjectResponse: nil,
					InjectError:    func() error { return errors.New("error") },
				},
			)),
			args: args{
				ctx:    context.TODO(),
				shards: testShardedCluster,
			},
			want:    []string{"shard00"},
			wantErr: true,
		},
		{
			name: "No master found, returns error, no shards changed",
			ss: NewSentinelServerFromParams(redis.NewFakeServerWithFakeClient("host", "port",
				// SentinelMaster response for shard00 (returns error as it is unmonitored)
				client.FakeResponse{
					InjectResponse: func() interface{} { return &client.SentinelMasterCmdResult{} },
					InjectError:    func() error { return errors.New(shardNotInitializedError) },
				},
			)),
			args: args{
				ctx: context.TODO(),
				shards: &Cluster{
					Shards: []Shard{
						{
							Name: "shard00",
							Servers: []*RedisServer{
								NewRedisServerFromParams(
									redis.MustNewServer("redis://127.0.0.1:2000", util.Pointer("shard00-0")),
									client.Slave,
									map[string]string{},
								),
								NewRedisServerFromParams(
									redis.MustNewServer("redis://127.0.0.1:2001", util.Pointer("shard00-1")),
									client.Slave,
									map[string]string{},
								),
								NewRedisServerFromParams(
									redis.MustNewServer("redis://127.0.0.1:2002", util.Pointer("shard00-2")),
									client.Slave,
									map[string]string{},
								),
							},
						},
					},
				},
			},
			want:    []string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ss.Monitor(tt.args.ctx, tt.args.shards)
			if (err != nil) != tt.wantErr {
				t.Errorf("SentinelServer.Monitor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SentinelServer.Monitor() = %v, want %v", got, tt.want)
			}
		})
	}
}
