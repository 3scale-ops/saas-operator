package redis

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/redis/crud"
	"github.com/3scale/saas-operator/pkg/redis/crud/client"
	redis "github.com/3scale/saas-operator/pkg/redis/crud/client"
	redis_client "github.com/3scale/saas-operator/pkg/redis/crud/client"
	"github.com/go-test/deep"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

var (
	s                  *runtime.Scheme = scheme.Scheme
	testShardedCluster ShardedCluster  = ShardedCluster{
		{
			Name: "shard00",
			Servers: []RedisServer{
				{
					Name: "shard00-0",
					Host: "127.0.0.1",
					Port: "2000",
					Role: redis.Master,
					CRUD: nil,
				},
				{
					Name: "shard00-1",
					Host: "127.0.0.1",
					Port: "2001",
					Role: redis.Slave,
					CRUD: nil,
				},
				{
					Name: "shard00-2",
					Host: "127.0.0.1",
					Port: "2002",
					Role: redis.Slave,
					CRUD: nil,
				},
			},
		},
		{
			Name: "shard01",
			Servers: []RedisServer{
				{
					Name: "shard01-0",
					Host: "127.0.0.1",
					Port: "3000",
					Role: redis.Master,
					CRUD: nil,
				},
				{
					Name: "shard01-1",
					Host: "127.0.0.1",
					Port: "3001",
					Role: redis.Slave,
					CRUD: nil,
				},
				{
					Name: "shard01-2",
					Host: "127.0.0.1",
					Port: "3002",
					Role: redis.Slave,
					CRUD: nil,
				},
			},
		},
		{
			Name: "shard02",
			Servers: []RedisServer{
				{
					Name: "shard02-0",
					Host: "127.0.0.1",
					Port: "4000",
					Role: redis.Master,
					CRUD: nil,
				},
				{
					Name: "shard02-1",
					Host: "127.0.0.1",
					Port: "4001",
					Role: redis.Slave,
					CRUD: nil,
				},
				{
					Name: "shard02-2",
					Host: "127.0.0.1",
					Port: "4002",
					Role: redis.Slave,
					CRUD: nil,
				},
			},
		},
	}
)

func init() {
	deep.CompareUnexportedFields = true
	s.AddKnownTypes(saasv1alpha1.GroupVersion)
}

func TestNewSentinelServerFromConnectionString(t *testing.T) {
	type args struct {
		name             string
		connectionString string
	}
	tests := []struct {
		name    string
		args    args
		want    *SentinelServer
		wantErr bool
	}{
		{
			name: "Returns a SentinelServer object",
			args: args{
				name:             "redis://127.0.0.1:6379",
				connectionString: "redis://127.0.0.1:6379",
			},
			want: &SentinelServer{
				Name: "redis://127.0.0.1:6379",
				CRUD: func() *crud.CRUD { c, _ := crud.NewRedisCRUDFromConnectionString("redis://127.0.0.1:6379"); return c }(),
				Port: "6379",
				IP:   "127.0.0.1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSentinelServerFromConnectionString(tt.args.name, tt.args.connectionString)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSentinelServerFromConnectionString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("NewSentinelServer() got diff: %v", diff)
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
			ss: &SentinelServer{
				Name: "test-server",
				CRUD: crud.NewFakeCRUD(redis.FakeResponse{
					InjectResponse: func() interface{} {
						return []interface{}{
							[]interface{}{"name", "shard01"},
							[]interface{}{"name", "shard02"},
						}
					},
					InjectError: func() error { return nil },
				}),
			},
			args: args{
				ctx:    context.TODO(),
				shards: []string{"shard01", "shard02"},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "No shard monitored",
			ss: &SentinelServer{
				Name: "test-server",
				CRUD: crud.NewFakeCRUD(redis.FakeResponse{
					InjectResponse: func() interface{} { return []interface{}{} },
					InjectError:    func() error { return nil },
				}),
			},
			args: args{
				ctx:    context.TODO(),
				shards: []string{"shard01", "shard02"},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "One shard is not monitored",
			ss: &SentinelServer{
				Name: "test-server",
				CRUD: crud.NewFakeCRUD(redis.FakeResponse{
					InjectResponse: func() interface{} {
						return []interface{}{
							[]interface{}{"name", "shard01"},
						}
					},
					InjectError: func() error { return nil },
				}),
			},
			args: args{
				ctx:    context.TODO(),
				shards: []string{"shard01", "shard02"},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Returns an error",
			ss: &SentinelServer{
				Name: "test-server",
				CRUD: crud.NewFakeCRUD(redis.FakeResponse{
					InjectResponse: func() interface{} { return []interface{}{} },
					InjectError:    func() error { return errors.New("error") },
				}),
			},
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
	type fields struct {
		Name string
		IP   string
		Port string
		CRUD *crud.CRUD
	}
	type args struct {
		ctx    context.Context
		shards ShardedCluster
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "All shards monitored",
			fields: fields{
				Name: "test-server",
				CRUD: crud.NewFakeCRUD(
					// SentinelMaster response for shard00
					redis.FakeResponse{
						InjectResponse: func() interface{} {
							return &redis.SentinelMasterCmdResult{
								Name: "shard00",
								IP:   "127.0.0.1",
								Port: 2000,
							}
						},
						InjectError: func() error { return nil },
					},
					// SentinelMaster response for shard01
					redis.FakeResponse{
						InjectResponse: func() interface{} {
							return &redis.SentinelMasterCmdResult{
								Name: "shard01",
								IP:   "127.0.0.1",
								Port: 3000,
							}
						},
						InjectError: func() error { return nil },
					},
					// SentinelMaster response for shard02
					redis.FakeResponse{
						InjectResponse: func() interface{} {
							return &redis.SentinelMasterCmdResult{
								Name: "shard02",
								IP:   "127.0.0.1",
								Port: 4000,
							}
						},
						InjectError: func() error { return nil },
					},
				),
			},
			args: args{
				ctx:    context.TODO(),
				shards: testShardedCluster,
			},
			want:    []string{},
			wantErr: false,
		},
		{
			name: "shard01 is not monitored",
			fields: fields{
				Name: "test-server",
				CRUD: crud.NewFakeCRUD(
					// SentinelMaster response for shard00
					redis.FakeResponse{
						InjectResponse: func() interface{} {
							return &redis.SentinelMasterCmdResult{
								Name: "shard00",
								IP:   "127.0.0.1",
								Port: 2000,
							}
						},
						InjectError: func() error { return nil },
					},
					// SentinelMaster response for shard01 (returns error as it is unmonitored)
					redis.FakeResponse{
						InjectResponse: func() interface{} { return &redis.SentinelMasterCmdResult{} },
						InjectError:    func() error { return errors.New(shardNotInitializedError) },
					},
					// SentinelMonitor response for shard01
					redis.FakeResponse{
						InjectResponse: nil,
						InjectError:    func() error { return nil },
					},
					// SentinelSet response for shard01
					redis.FakeResponse{
						InjectResponse: nil,
						InjectError:    func() error { return nil },
					},
					// SentinelMaster response for shard02
					redis.FakeResponse{
						InjectResponse: func() interface{} {
							return &redis.SentinelMasterCmdResult{
								Name: "shard02",
								IP:   "127.0.0.1",
								Port: 4000,
							}
						},
						InjectError: func() error { return nil },
					},
				),
			},
			args: args{
				ctx:    context.TODO(),
				shards: testShardedCluster,
			},
			want:    []string{"shard01"},
			wantErr: false,
		},
		{
			name: "all shards are unmonitored",
			fields: fields{
				Name: "test-server",
				CRUD: crud.NewFakeCRUD(
					// SentinelMaster response for shard00 (returns error as it is unmonitored)
					redis.FakeResponse{
						InjectResponse: func() interface{} { return &redis.SentinelMasterCmdResult{} },
						InjectError:    func() error { return errors.New(shardNotInitializedError) },
					},
					// SentinelMonitor response for shard00
					redis.FakeResponse{
						InjectResponse: nil,
						InjectError:    func() error { return nil },
					},
					// SentinelSet response for shard00
					redis.FakeResponse{
						InjectResponse: nil,
						InjectError:    func() error { return nil },
					},
					// SentinelMaster response for shard01 (returns error as it is unmonitored)
					redis.FakeResponse{
						InjectResponse: func() interface{} { return &redis.SentinelMasterCmdResult{} },
						InjectError:    func() error { return errors.New(shardNotInitializedError) },
					},
					// SentinelMonitor response for shard01
					redis.FakeResponse{
						InjectResponse: nil,
						InjectError:    func() error { return nil },
					},
					// SentinelSet response for shard01
					redis.FakeResponse{
						InjectResponse: nil,
						InjectError:    func() error { return nil },
					},
					// SentinelMaster response for shard02 (returns error as it is unmonitored)
					redis.FakeResponse{
						InjectResponse: func() interface{} { return &redis.SentinelMasterCmdResult{} },
						InjectError:    func() error { return errors.New(shardNotInitializedError) },
					},
					// SentinelMonitor response for shard02
					redis.FakeResponse{
						InjectResponse: nil,
						InjectError:    func() error { return nil },
					},
					// SentinelSet response for shard02
					redis.FakeResponse{
						InjectResponse: nil,
						InjectError:    func() error { return nil },
					},
				),
			},
			args: args{
				ctx:    context.TODO(),
				shards: testShardedCluster,
			},
			want:    []string{"shard00", "shard01", "shard02"},
			wantErr: false,
		},
		{
			name: "All shards unmonitored, failure on the 2nd one",
			fields: fields{
				Name: "test-server",
				CRUD: crud.NewFakeCRUD(
					// SentinelMaster response for shard00 (returns error as it is unmonitored)
					redis.FakeResponse{
						InjectResponse: func() interface{} { return &redis.SentinelMasterCmdResult{} },
						InjectError:    func() error { return errors.New(shardNotInitializedError) },
					},
					// SentinelMonitor response for shard00
					redis.FakeResponse{
						InjectResponse: nil,
						InjectError:    func() error { return nil },
					},
					// SentinelSet response for shard00
					redis.FakeResponse{
						InjectResponse: nil,
						InjectError:    func() error { return nil },
					},
					// SentinelMaster response for shard01 (returns error as it is unmonitored)
					redis.FakeResponse{
						InjectResponse: func() interface{} { return &redis.SentinelMasterCmdResult{} },
						InjectError:    func() error { return errors.New("error") },
					},
					// SentinelMaster response for shard02 (returns error as it is unmonitored)
					redis.FakeResponse{
						InjectResponse: func() interface{} { return &redis.SentinelMasterCmdResult{} },
						InjectError:    func() error { return errors.New(shardNotInitializedError) },
					},
					// SentinelMonitor response for shard02
					redis.FakeResponse{
						InjectResponse: nil,
						InjectError:    func() error { return nil },
					},
					// SentinelSet response for shard02
					redis.FakeResponse{
						InjectResponse: nil,
						InjectError:    func() error { return nil },
					},
				),
			},
			args: args{
				ctx:    context.TODO(),
				shards: testShardedCluster,
			},
			want:    []string{"shard00"},
			wantErr: true,
		},
		{
			name: "All shards monitored, failure on the 2nd one",
			fields: fields{
				Name: "test-server",
				CRUD: crud.NewFakeCRUD(
					// SentinelMaster response for shard00 (returns error as it is unmonitored)
					redis.FakeResponse{
						InjectResponse: func() interface{} {
							return &redis.SentinelMasterCmdResult{
								Name: "shard00",
								IP:   "127.0.0.1",
								Port: 2000,
							}
						},
						InjectError: func() error { return nil },
					},
					// SentinelMaster response for shard01 (returns error as it is unmonitored)
					redis.FakeResponse{
						InjectResponse: func() interface{} { return &redis.SentinelMasterCmdResult{} },
						InjectError:    func() error { return errors.New("error") },
					},
				),
			},
			args: args{
				ctx:    context.TODO(),
				shards: testShardedCluster,
			},
			want:    []string{},
			wantErr: true,
		},
		{
			name: "'sentinel monitor' fails for shard00, returns no shards changed",
			fields: fields{
				Name: "test-server",
				CRUD: crud.NewFakeCRUD(
					// SentinelMaster response for shard00 (returns error as it is unmonitored)
					redis.FakeResponse{
						InjectResponse: func() interface{} { return &redis.SentinelMasterCmdResult{} },
						InjectError:    func() error { return errors.New(shardNotInitializedError) },
					},
					// SentinelMonitor response for shard00
					redis.FakeResponse{
						InjectResponse: nil,
						InjectError:    func() error { return errors.New("error") },
					},
				),
			},
			args: args{
				ctx:    context.TODO(),
				shards: testShardedCluster,
			},
			want:    []string{},
			wantErr: true,
		},
		{
			name: "Error writing config param, returns shard00 changed",
			fields: fields{
				Name: "test-server",
				CRUD: crud.NewFakeCRUD(
					// SentinelMaster response for shard00 (returns error as it is unmonitored)
					redis.FakeResponse{
						InjectResponse: func() interface{} { return &redis.SentinelMasterCmdResult{} },
						InjectError:    func() error { return errors.New(shardNotInitializedError) },
					},
					// SentinelMonitor response for shard00
					redis.FakeResponse{
						InjectResponse: nil,
						InjectError:    func() error { return nil },
					},
					// SentinelSet response for shard01
					redis.FakeResponse{
						InjectResponse: nil,
						InjectError:    func() error { return errors.New("error") },
					},
				),
			},
			args: args{
				ctx:    context.TODO(),
				shards: testShardedCluster,
			},
			want:    []string{"shard00"},
			wantErr: true,
		},
		{
			name: "No master found, returns error, no shards changed",
			fields: fields{
				Name: "test-server",
				CRUD: crud.NewFakeCRUD(
					// SentinelMaster response for shard00 (returns error as it is unmonitored)
					redis.FakeResponse{
						InjectResponse: func() interface{} { return &redis.SentinelMasterCmdResult{} },
						InjectError:    func() error { return errors.New(shardNotInitializedError) },
					},
					// SentinelMonitor response for shard00
					redis.FakeResponse{
						InjectResponse: nil,
						InjectError:    func() error { return nil },
					},
					// SentinelSet response for shard01
					redis.FakeResponse{
						InjectResponse: nil,
						InjectError:    func() error { return errors.New("error") },
					},
				),
			},
			args: args{
				ctx: context.TODO(),
				shards: ShardedCluster{{
					Name: "shard00",
					Servers: []RedisServer{
						{
							Name: "shard00-0",
							Host: "127.0.0.1",
							Port: "2000",
							Role: redis.Slave,
							CRUD: nil,
						},
						{
							Name: "shard00-1",
							Host: "127.0.0.1",
							Port: "2001",
							Role: redis.Slave,
							CRUD: nil,
						},
						{
							Name: "shard00-2",
							Host: "127.0.0.1",
							Port: "2002",
							Role: redis.Slave,
							CRUD: nil,
						},
					},
				}},
			},
			want:    []string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &SentinelServer{
				Name: tt.fields.Name,
				IP:   tt.fields.IP,
				Port: tt.fields.Port,
				CRUD: tt.fields.CRUD,
			}

			got, err := ss.Monitor(tt.args.ctx, tt.args.shards)
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

func TestSentinelServer_MonitoredShards(t *testing.T) {
	type args struct {
		ctx  context.Context
		opts ShardDiscoveryOption
	}
	tests := []struct {
		name    string
		ss      *SentinelServer
		args    args
		want    saasv1alpha1.MonitoredShards
		wantErr bool
	}{
		{
			name: "Returns all shards monitored by sentinel",
			ss: &SentinelServer{
				Name: "test-server",
				CRUD: crud.NewFakeCRUD(
					redis.FakeResponse{
						// cmd: SentinelMasters()
						InjectResponse: func() interface{} {
							return []interface{}{
								[]interface{}{"name", "shard01", "ip", "127.0.0.1", "port", "6379"},
								[]interface{}{"name", "shard02", "ip", "127.0.0.2", "port", "6379"},
							}
						},
						InjectError: func() error { return nil },
					},
					redis.FakeResponse{
						// cmd: SentinelMaster
						InjectResponse: func() interface{} {
							return &redis_client.SentinelMasterCmdResult{Name: "shard01", IP: "127.0.0.1", Port: 6379, Flags: "master"}
						},
						InjectError: func() error { return nil },
					},
					redis.FakeResponse{
						// cmd: SentinelMaster
						InjectResponse: func() interface{} {
							return &redis_client.SentinelMasterCmdResult{Name: "shard02", IP: "127.0.0.2", Port: 6379, Flags: "master"}
						},
						InjectError: func() error { return nil },
					},
				),
			},
			args: args{
				ctx:  context.TODO(),
				opts: OnlyMasterDiscoveryOpt,
			},
			want: saasv1alpha1.MonitoredShards{
				{
					Name: "shard01",
					Servers: map[string]saasv1alpha1.RedisServerDetails{
						"127.0.0.1:6379": {Role: client.Master, Config: map[string]string{}}},
				},
				{
					Name: "shard02",
					Servers: map[string]saasv1alpha1.RedisServerDetails{
						"127.0.0.2:6379": {Role: client.Master, Config: map[string]string{}},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Returns an error",
			ss: &SentinelServer{
				Name: "test-server",
				CRUD: crud.NewFakeCRUD(redis.FakeResponse{
					InjectResponse: func() interface{} { return []interface{}{} },
					InjectError:    func() error { return errors.New("error") },
				}),
			},
			args: args{
				ctx:  context.TODO(),
				opts: OnlyMasterDiscoveryOpt,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ss.MonitoredShards(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("SentinelServer.MonitoredShards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SentinelServer.MonitoredShards() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSentinelServer_DiscoverShard(t *testing.T) {
	type args struct {
		ctx             context.Context
		shard           string
		maxInfoCacheAge time.Duration
		opts            ShardDiscoveryOptions
	}
	tests := []struct {
		name    string
		ss      *SentinelServer
		rss     []RedisServer
		args    args
		want    map[string]saasv1alpha1.RedisServerDetails
		wantErr bool
	}{
		{
			name: "Discovers roles and config options within a shard (all available options)",
			ss: &SentinelServer{
				Name: "tsentinel",
				CRUD: crud.NewFakeCRUD(
					redis.FakeResponse{
						// cmd: SentinelMaster
						InjectResponse: func() interface{} {
							return &redis_client.SentinelMasterCmdResult{Name: "shard0", IP: "127.0.0.1", Port: 1000, Flags: "master"}
						},
						InjectError: func() error { return nil },
					},
					redis.FakeResponse{
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
								[]interface{}{
									"name", "127.0.0.1:4000",
									"ip", "127.0.0.1",
									"port", "4000",
									"flags", "slave",
								},
							}
						},
						InjectError: func() error { return nil },
					},
				),
				MonitoredRedisServers: map[string]*RedisServer{
					// redis master
					"127.0.0.1:1000": {
						CRUD: crud.NewFakeCRUD(
							redis.FakeResponse{
								// cmd: RedisConfigGet("save")
								InjectResponse: func() interface{} {
									return []interface{}{"save", "900 1 300 10"}
								},
								InjectError: func() error { return nil },
							}),
					},
					// redis slaves
					"127.0.0.1:2000": {
						CRUD: crud.NewFakeCRUD(
							redis.FakeResponse{
								// cmd: RedisConfigGet("save")
								InjectResponse: func() interface{} {
									return []interface{}{"save", "900 1 300 10"}
								},
								InjectError: func() error { return nil },
							},
							redis.FakeResponse{
								// cmd: RedisConfigGet("slave-read-only")
								InjectResponse: func() interface{} {
									return []interface{}{"slave-read-only", "yes"}
								},
								InjectError: func() error { return nil },
							},
						),
					},
					"127.0.0.1:3000": {
						CRUD: crud.NewFakeCRUD(
							redis.FakeResponse{
								// cmd: RedisConfigGet("save")
								InjectResponse: func() interface{} {
									return []interface{}{"save", "900 1 300 10"}
								},
								InjectError: func() error { return nil },
							},
							redis.FakeResponse{
								// cmd: RedisConfigGet("slave-read-only")
								InjectResponse: func() interface{} {
									return []interface{}{"slave-read-only", "yes"}
								},
								InjectError: func() error { return nil },
							},
						),
					},
					"127.0.0.1:4000": {
						CRUD: crud.NewFakeCRUD(
							redis.FakeResponse{
								// cmd: RedisConfigGet("save")
								InjectResponse: func() interface{} {
									return []interface{}{"save", "900 1 300 10"}
								},
								InjectError: func() error { return nil },
							},
							redis.FakeResponse{
								// cmd: RedisConfigGet("slave-read-only")
								InjectResponse: func() interface{} {
									return []interface{}{"slave-read-only", "no"}
								},
								InjectError: func() error { return nil },
							},
						),
					},
				},
			},
			args: args{
				opts: ShardDiscoveryOptions{SaveConfigDiscoveryOpt, SlaveReadOnlyDiscoveryOpt},
			},
			want: map[string]saasv1alpha1.RedisServerDetails{
				"127.0.0.1:1000": {
					Role:   client.Master,
					Config: map[string]string{"save": "900 1 300 10"},
				},
				"127.0.0.1:2000": {
					Role:   client.Slave,
					Config: map[string]string{"save": "900 1 300 10", "slave-read-only": "yes"},
				},
				"127.0.0.1:3000": {
					Role:   client.Slave,
					Config: map[string]string{"save": "900 1 300 10", "slave-read-only": "yes"},
				},
				"127.0.0.1:4000": {
					Role:   client.Slave,
					Config: map[string]string{"save": "900 1 300 10", "slave-read-only": "no"},
				},
			},
			wantErr: false,
		},
		{
			name: "Only discovers master, no config",
			ss: &SentinelServer{
				Name: "tsentinel",
				CRUD: crud.NewFakeCRUD(
					redis.FakeResponse{
						// cmd: SentinelMaster
						InjectResponse: func() interface{} {
							return &redis_client.SentinelMasterCmdResult{Name: "shard0", IP: "127.0.0.1", Port: 1000, Flags: "master"}
						},
						InjectError: func() error { return nil },
					},
				),
				MonitoredRedisServers: map[string]*RedisServer{},
			},
			args: args{
				opts: ShardDiscoveryOptions{OnlyMasterDiscoveryOpt},
			},
			want: map[string]saasv1alpha1.RedisServerDetails{
				"127.0.0.1:1000": {
					Role:   client.Master,
					Config: map[string]string{},
				},
			},
			wantErr: false,
		},
		{
			name: "Discovers roles and slave-read-only option",
			ss: &SentinelServer{
				Name: "tsentinel",
				CRUD: crud.NewFakeCRUD(
					redis.FakeResponse{
						// cmd: SentinelMaster
						InjectResponse: func() interface{} {
							return &redis_client.SentinelMasterCmdResult{Name: "shard0", IP: "127.0.0.1", Port: 1000, Flags: "master"}
						},
						InjectError: func() error { return nil },
					},
					redis.FakeResponse{
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
								[]interface{}{
									"name", "127.0.0.1:4000",
									"ip", "127.0.0.1",
									"port", "4000",
									"flags", "slave",
								},
							}
						},
						InjectError: func() error { return nil },
					},
				),
				MonitoredRedisServers: map[string]*RedisServer{
					// redis master
					"127.0.0.1:1000": {
						CRUD: crud.NewFakeCRUD(),
					},
					// redis slaves
					"127.0.0.1:2000": {
						CRUD: crud.NewFakeCRUD(
							redis.FakeResponse{
								// cmd: RedisConfigGet("slave-read-only")
								InjectResponse: func() interface{} {
									return []interface{}{"slave-read-only", "yes"}
								},
								InjectError: func() error { return nil },
							},
						),
					},
					"127.0.0.1:3000": {
						CRUD: crud.NewFakeCRUD(
							redis.FakeResponse{
								// cmd: RedisConfigGet("slave-read-only")
								InjectResponse: func() interface{} {
									return []interface{}{"slave-read-only", "yes"}
								},
								InjectError: func() error { return nil },
							},
						),
					},
					"127.0.0.1:4000": {
						CRUD: crud.NewFakeCRUD(
							redis.FakeResponse{
								// cmd: RedisConfigGet("slave-read-only")
								InjectResponse: func() interface{} {
									return []interface{}{"slave-read-only", "no"}
								},
								InjectError: func() error { return nil },
							},
						),
					},
				},
			},
			args: args{
				opts: ShardDiscoveryOptions{SlaveReadOnlyDiscoveryOpt},
			},
			want: map[string]saasv1alpha1.RedisServerDetails{
				"127.0.0.1:1000": {
					Role:   client.Master,
					Config: map[string]string{},
				},
				"127.0.0.1:2000": {
					Role:   client.Slave,
					Config: map[string]string{"slave-read-only": "yes"},
				},
				"127.0.0.1:3000": {
					Role:   client.Slave,
					Config: map[string]string{"slave-read-only": "yes"},
				},
				"127.0.0.1:4000": {
					Role:   client.Slave,
					Config: map[string]string{"slave-read-only": "no"},
				},
			},
			wantErr: false,
		},
		{
			name: "Discovers roles and save option",
			ss: &SentinelServer{
				Name: "tsentinel",
				CRUD: crud.NewFakeCRUD(
					redis.FakeResponse{
						// cmd: SentinelMaster
						InjectResponse: func() interface{} {
							return &redis_client.SentinelMasterCmdResult{Name: "shard0", IP: "127.0.0.1", Port: 1000, Flags: "master"}
						},
						InjectError: func() error { return nil },
					},
					redis.FakeResponse{
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
								[]interface{}{
									"name", "127.0.0.1:4000",
									"ip", "127.0.0.1",
									"port", "4000",
									"flags", "slave",
								},
							}
						},
						InjectError: func() error { return nil },
					},
				),
				MonitoredRedisServers: map[string]*RedisServer{
					// redis master
					"127.0.0.1:1000": {
						CRUD: crud.NewFakeCRUD(
							redis.FakeResponse{
								// cmd: RedisConfigGet("save")
								InjectResponse: func() interface{} {
									return []interface{}{"save", "900 1 300 10"}
								},
								InjectError: func() error { return nil },
							},
						)},
					// redis slaves
					"127.0.0.1:2000": {
						CRUD: crud.NewFakeCRUD(
							redis.FakeResponse{
								// cmd: RedisConfigGet("save")
								InjectResponse: func() interface{} {
									return []interface{}{"save", "900 1 300 10"}
								},
								InjectError: func() error { return nil },
							},
						),
					},
					"127.0.0.1:3000": {
						CRUD: crud.NewFakeCRUD(
							redis.FakeResponse{
								// cmd: RedisConfigGet("save")
								InjectResponse: func() interface{} {
									return []interface{}{"save", "900 1 300 10"}
								},
								InjectError: func() error { return nil },
							},
						),
					},
					"127.0.0.1:4000": {
						CRUD: crud.NewFakeCRUD(
							redis.FakeResponse{
								// cmd: RedisConfigGet("save")
								InjectResponse: func() interface{} {
									return []interface{}{"save", "900 1 300 10"}
								},
								InjectError: func() error { return nil },
							},
						),
					},
				},
			},
			args: args{
				opts: ShardDiscoveryOptions{SaveConfigDiscoveryOpt},
			},
			want: map[string]saasv1alpha1.RedisServerDetails{
				"127.0.0.1:1000": {
					Role:   client.Master,
					Config: map[string]string{"save": "900 1 300 10"},
				},
				"127.0.0.1:2000": {
					Role:   client.Slave,
					Config: map[string]string{"save": "900 1 300 10"},
				},
				"127.0.0.1:3000": {
					Role:   client.Slave,
					Config: map[string]string{"save": "900 1 300 10"},
				},
				"127.0.0.1:4000": {
					Role:   client.Slave,
					Config: map[string]string{"save": "900 1 300 10"},
				},
			},
			wantErr: false,
		},
		{
			name: "Avoids down slaves",
			ss: &SentinelServer{
				Name: "tsentinel",
				CRUD: crud.NewFakeCRUD(
					redis.FakeResponse{
						// cmd: SentinelMaster
						InjectResponse: func() interface{} {
							return &redis_client.SentinelMasterCmdResult{Name: "shard0", IP: "127.0.0.1", Port: 1000, Flags: "master"}
						},
						InjectError: func() error { return nil },
					},
					redis.FakeResponse{
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
								[]interface{}{
									"name", "127.0.0.1:4000",
									"ip", "127.0.0.1",
									"port", "4000",
									"flags", "slave",
								},
							}
						},
						InjectError: func() error { return nil },
					},
				),
				MonitoredRedisServers: map[string]*RedisServer{
					// redis master
					"127.0.0.1:1000": {
						CRUD: crud.NewFakeCRUD(
							redis.FakeResponse{
								// cmd: RedisConfigGet("save")
								InjectResponse: func() interface{} {
									return []interface{}{"save", "900 1 300 10"}
								},
								InjectError: func() error { return nil },
							}),
					},
					// redis slaves
					"127.0.0.1:2000": {
						CRUD: crud.NewFakeCRUD(
							redis.FakeResponse{
								// cmd: RedisConfigGet("save")
								InjectResponse: func() interface{} {
									return []interface{}{"save", "900 1 300 10"}
								},
								InjectError: func() error { return nil },
							},
							redis.FakeResponse{
								// cmd: RedisConfigGet("slave-read-only")
								InjectResponse: func() interface{} {
									return []interface{}{"slave-read-only", "yes"}
								},
								InjectError: func() error { return nil },
							},
						),
					},
					"127.0.0.1:3000": {
						CRUD: crud.NewFakeCRUD(
							redis.FakeResponse{
								// cmd: RedisConfigGet("save")
								InjectResponse: func() interface{} {
									return []interface{}{"save", "900 1 300 10"}
								},
								InjectError: func() error { return nil },
							},
							redis.FakeResponse{
								// cmd: RedisConfigGet("slave-read-only")
								InjectResponse: func() interface{} {
									return []interface{}{"slave-read-only", "yes"}
								},
								InjectError: func() error { return nil },
							},
						),
					},
					"127.0.0.1:4000": {
						CRUD: crud.NewFakeCRUD(
							redis.FakeResponse{
								// cmd: RedisConfigGet("save")
								InjectResponse: func() interface{} {
									return []interface{}{"save", "900 1 300 10"}
								},
								InjectError: func() error { return nil },
							},
							redis.FakeResponse{
								// cmd: RedisConfigGet("slave-read-only")
								InjectResponse: func() interface{} {
									return []interface{}{"slave-read-only", "no"}
								},
								InjectError: func() error { return nil },
							},
						),
					},
				},
			},
			args: args{
				opts: ShardDiscoveryOptions{SaveConfigDiscoveryOpt, SlaveReadOnlyDiscoveryOpt},
			},
			want: map[string]saasv1alpha1.RedisServerDetails{
				"127.0.0.1:1000": {
					Role:   client.Master,
					Config: map[string]string{"save": "900 1 300 10"},
				},
				"127.0.0.1:3000": {
					Role:   client.Slave,
					Config: map[string]string{"save": "900 1 300 10", "slave-read-only": "yes"},
				},
				"127.0.0.1:4000": {
					Role:   client.Slave,
					Config: map[string]string{"save": "900 1 300 10", "slave-read-only": "no"},
				},
			},
			wantErr: false,
		},
		{
			name: "Fails if master is down",
			ss: &SentinelServer{
				Name: "tsentinel",
				CRUD: crud.NewFakeCRUD(
					redis.FakeResponse{
						// cmd: SentinelMaster
						InjectResponse: func() interface{} {
							return &redis_client.SentinelMasterCmdResult{Name: "shard0", IP: "127.0.0.1", Port: 1000, Flags: "o_down,master"}
						},
						InjectError: func() error { return nil },
					},
				),
				MonitoredRedisServers: map[string]*RedisServer{},
			},
			args: args{
				opts: ShardDiscoveryOptions{SaveConfigDiscoveryOpt, SlaveReadOnlyDiscoveryOpt},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ss.DiscoverShard(tt.args.ctx, tt.args.shard, tt.args.maxInfoCacheAge, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("SentinelServer.DiscoverShard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SentinelServer.DiscoverShard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShardDiscoveryOptions_Has(t *testing.T) {
	type args struct {
		sdo ShardDiscoveryOption
	}
	tests := []struct {
		name string
		sdos ShardDiscoveryOptions
		args args
		want bool
	}{
		{
			name: "Returns true if option in slice",
			sdos: []ShardDiscoveryOption{SaveConfigDiscoveryOpt, SlaveReadOnlyDiscoveryOpt},
			args: args{SlaveReadOnlyDiscoveryOpt},
			want: true,
		},
		{
			name: "Returns false if option not in slice",
			sdos: []ShardDiscoveryOption{SaveConfigDiscoveryOpt, SlaveReadOnlyDiscoveryOpt},
			args: args{OnlyMasterDiscoveryOpt},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sdos.Has(tt.args.sdo); got != tt.want {
				t.Errorf("ShardDiscoveryOptions.Has() = %v, want %v", got, tt.want)
			}
		})
	}
}
