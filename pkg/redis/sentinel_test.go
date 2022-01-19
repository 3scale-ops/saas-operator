package redis

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/redis/crud"
	redis "github.com/3scale/saas-operator/pkg/redis/crud/client"
	"github.com/go-test/deep"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8s "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	s                  *runtime.Scheme = scheme.Scheme
	testShardedCluster ShardedCluster  = ShardedCluster{
		{
			Name: "shard00",
			Servers: []RedisServer{
				{
					Name:     "shard00-0",
					IP:       "127.0.0.1",
					Port:     "2000",
					Role:     redis.Master,
					ReadOnly: false,
					CRUD:     nil,
				},
				{
					Name:     "shard00-1",
					IP:       "127.0.0.1",
					Port:     "2001",
					Role:     redis.Slave,
					ReadOnly: true,
					CRUD:     nil,
				},
				{
					Name:     "shard00-2",
					IP:       "127.0.0.1",
					Port:     "2002",
					Role:     redis.Slave,
					ReadOnly: true,
					CRUD:     nil,
				},
			},
		},
		{
			Name: "shard01",
			Servers: []RedisServer{
				{
					Name:     "shard01-0",
					IP:       "127.0.0.1",
					Port:     "3000",
					Role:     redis.Master,
					ReadOnly: false,
					CRUD:     nil,
				},
				{
					Name:     "shard01-1",
					IP:       "127.0.0.1",
					Port:     "3001",
					Role:     redis.Slave,
					ReadOnly: true,
					CRUD:     nil,
				},
				{
					Name:     "shard01-2",
					IP:       "127.0.0.1",
					Port:     "3002",
					Role:     redis.Slave,
					ReadOnly: true,
					CRUD:     nil,
				},
			},
		},
		{
			Name: "shard02",
			Servers: []RedisServer{
				{
					Name:     "shard02-0",
					IP:       "127.0.0.1",
					Port:     "4000",
					Role:     redis.Master,
					ReadOnly: false,
					CRUD:     nil,
				},
				{
					Name:     "shard02-1",
					IP:       "127.0.0.1",
					Port:     "4001",
					Role:     redis.Slave,
					ReadOnly: true,
					CRUD:     nil,
				},
				{
					Name:     "shard02-2",
					IP:       "127.0.0.1",
					Port:     "4002",
					Role:     redis.Slave,
					ReadOnly: true,
					CRUD:     nil,
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

func TestSentinelServer_MonitoredShards(t *testing.T) {
	type args struct {
		ctx            context.Context
		discoverSlaves bool
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
				CRUD: crud.NewFakeCRUD(redis.FakeResponse{
					InjectResponse: func() interface{} {
						return []interface{}{
							[]interface{}{"name", "shard01", "ip", "127.0.0.1", "port", "6379"},
							[]interface{}{"name", "shard02", "ip", "127.0.0.2", "port", "6379"},
						}
					},
					InjectError: func() error { return nil },
				}),
			},
			args: args{
				ctx:            context.TODO(),
				discoverSlaves: false,
			},
			want: saasv1alpha1.MonitoredShards{
				{Name: "shard01", Master: "127.0.0.1:6379"},
				{Name: "shard02", Master: "127.0.0.2:6379"},
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
				ctx:            context.TODO(),
				discoverSlaves: false,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ss.MonitoredShards(tt.args.ctx, tt.args.discoverSlaves)
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

func TestSentinelServer_DiscoverSlaves(t *testing.T) {
	type fields struct {
		Name string
		IP   string
		Port string
		CRUD *crud.CRUD
	}
	type args struct {
		ctx             context.Context
		shard           string
		maxInfoCacheAge time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		want1   []string
		wantErr bool
	}{
		{
			name: "Returns the RO and RW slaves of the shard",
			fields: fields{
				Name: "test-server",
				CRUD: crud.NewFakeCRUD(
					redis.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{
								"shard01",
								[]interface{}{
									[]interface{}{
										int64(1000),
										"# Server\r\nrun_id:bce68a863acb3bb1e02c2caae48ce36373c524fc\r\n\r\n# Replication\r\nslave_read_only:1",
									},
									[]interface{}{
										int64(1000),
										"# Server\r\nrun_id:1f67e9246d3017be5d5cb9a1fdc6020c8338da76\r\n\r\n# Replication\r\nslave_read_only:0",
									},
									[]interface{}{
										int64(1000),
										"# Server\r\nrun_id:addb9655385adb19c7944889c0efc3b3cb63602b\r\n\r\n# Replication\r\nslave_read_only:1",
									},
								},
							}
						},
						InjectError: func() error { return nil },
					},
					redis.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{
								[]interface{}{
									"name", "10.244.0.7:6379",
									"ip", "10.244.0.7",
									"port", "6379",
									"runid", "1f67e9246d3017be5d5cb9a1fdc6020c8338da76",
									"flags", "slave",
								},
								[]interface{}{
									"name", "10.244.0.6:6379",
									"ip", "10.244.0.6",
									"port", "6379",
									"runid", "bce68a863acb3bb1e02c2caae48ce36373c524fc",
									"flags", "slave",
								},
								[]interface{}{
									"name", "10.244.0.8:6379",
									"ip", "10.244.0.8",
									"port", "6379",
									"runid", "addb9655385adb19c7944889c0efc3b3cb63602b",
									"flags", "slave",
								},
							}
						},
						InjectError: func() error { return nil },
					},
				),
			},
			args:    args{ctx: nil, shard: "shard01", maxInfoCacheAge: 1 * time.Second},
			want:    []string{"10.244.0.6:6379", "10.244.0.8:6379"},
			want1:   []string{"10.244.0.7:6379"},
			wantErr: false,
		},
		{
			name: "Ignores servers if their cache is bigger than maxInfoCacheAge",
			fields: fields{
				Name: "test-server",
				CRUD: crud.NewFakeCRUD(
					redis.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{
								"shard01",
								[]interface{}{
									[]interface{}{
										int64(1000),
										"# Server\r\nrun_id:bce68a863acb3bb1e02c2caae48ce36373c524fc\r\n\r\n# Replication\r\nslave_read_only:1",
									},
									[]interface{}{
										int64(1000),
										"# Server\r\nrun_id:1f67e9246d3017be5d5cb9a1fdc6020c8338da76\r\n\r\n# Replication\r\nslave_read_only:0",
									},
									[]interface{}{
										int64(2000),
										"# Server\r\nrun_id:addb9655385adb19c7944889c0efc3b3cb63602b\r\n\r\n# Replication\r\nslave_read_only:1",
									},
								},
							}
						},
						InjectError: func() error { return nil },
					},
					redis.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{
								[]interface{}{
									"name", "10.244.0.6:6379",
									"ip", "10.244.0.6",
									"port", "6379",
									"runid", "bce68a863acb3bb1e02c2caae48ce36373c524fc",
									"flags", "slave",
								},
								[]interface{}{
									"name", "10.244.0.7:6379",
									"ip", "10.244.0.7",
									"port", "6379",
									"runid", "1f67e9246d3017be5d5cb9a1fdc6020c8338da76",
									"flags", "slave",
								},
								[]interface{}{
									"name", "10.244.0.8:6379",
									"ip", "10.244.0.8",
									"port", "6379",
									"runid", "addb9655385adb19c7944889c0efc3b3cb63602b",
									"flags", "slave",
								},
							}
						},
						InjectError: func() error { return nil },
					},
				),
			},
			args:    args{ctx: nil, shard: "shard01", maxInfoCacheAge: 1 * time.Second},
			want:    []string{"10.244.0.6:6379"},
			want1:   []string{"10.244.0.7:6379"},
			wantErr: false,
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
			got, got1, err := ss.DiscoverSlaves(tt.args.ctx, tt.args.shard, tt.args.maxInfoCacheAge)
			if (err != nil) != tt.wantErr {
				t.Errorf("SentinelServer.DiscoverSlaves() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SentinelServer.DiscoverSlaves() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SentinelServer.DiscoverSlaves() got1 = %v, want %v", got1, tt.want1)
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
							Name:     "shard00-0",
							IP:       "127.0.0.1",
							Port:     "2000",
							Role:     redis.Slave,
							ReadOnly: true,
							CRUD:     nil,
						},
						{
							Name:     "shard00-1",
							IP:       "127.0.0.1",
							Port:     "2001",
							Role:     redis.Slave,
							ReadOnly: true,
							CRUD:     nil,
						},
						{
							Name:     "shard00-2",
							IP:       "127.0.0.1",
							Port:     "2002",
							Role:     redis.Slave,
							ReadOnly: true,
							CRUD:     nil,
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

func TestNewSentinelPool(t *testing.T) {
	type args struct {
		ctx      context.Context
		cl       client.Client
		key      types.NamespacedName
		replicas int
	}
	tests := []struct {
		name    string
		args    args
		want    SentinelPool
		wantErr bool
	}{
		{
			name: "Returns a SentinelPool object",
			args: args{
				ctx: context.TODO(),
				cl: k8s.NewClientBuilder().WithScheme(s).WithObjects(
					&corev1.Pod{
						ObjectMeta: metav1.ObjectMeta{Name: "sentinel-0", Namespace: "test"},
						Status:     corev1.PodStatus{PodIP: "127.0.0.1"},
					},
					&corev1.Pod{
						ObjectMeta: metav1.ObjectMeta{Name: "sentinel-1", Namespace: "test"},
						Status:     corev1.PodStatus{PodIP: "127.0.0.2"},
					},
					&corev1.Pod{
						ObjectMeta: metav1.ObjectMeta{Name: "sentinel-2", Namespace: "test"},
						Status:     corev1.PodStatus{PodIP: "127.0.0.3"},
					},
				).Build(),
				key:      types.NamespacedName{Name: "sentinel", Namespace: "test"},
				replicas: 3,
			},
			want: []SentinelServer{
				{
					Name: "sentinel-0",
					IP:   "127.0.0.1",
					Port: "26379",
					CRUD: func() *crud.CRUD { c, _ := crud.NewRedisCRUDFromConnectionString("redis://127.0.0.1:26379"); return c }(),
				},
				{
					Name: "sentinel-1",
					IP:   "127.0.0.2",
					Port: "26379",
					CRUD: func() *crud.CRUD { c, _ := crud.NewRedisCRUDFromConnectionString("redis://127.0.0.2:26379"); return c }(),
				},
				{
					Name: "sentinel-2",
					IP:   "127.0.0.3",
					Port: "26379",
					CRUD: func() *crud.CRUD { c, _ := crud.NewRedisCRUDFromConnectionString("redis://127.0.0.3:26379"); return c }(),
				},
			},
			wantErr: false,
		},
		{
			name: "Pod not found, returns error",
			args: args{
				ctx: context.TODO(),
				cl: k8s.NewClientBuilder().WithScheme(s).WithObjects(
					&corev1.Pod{
						ObjectMeta: metav1.ObjectMeta{Name: "sentinel-0", Namespace: "test"},
						Status:     corev1.PodStatus{PodIP: "127.0.0.1"},
					},
					&corev1.Pod{
						ObjectMeta: metav1.ObjectMeta{Name: "sentinel-1", Namespace: "test"},
						Status:     corev1.PodStatus{PodIP: "127.0.0.2"},
					},
				).Build(),
				key:      types.NamespacedName{Name: "sentinel", Namespace: "test"},
				replicas: 3,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Pod not found, returns error",
			args: args{
				ctx: context.TODO(),
				cl: k8s.NewClientBuilder().WithScheme(s).WithObjects(
					&corev1.Pod{
						ObjectMeta: metav1.ObjectMeta{Name: "sentinel-0", Namespace: "test"},
						Status:     corev1.PodStatus{PodIP: "127.0.0.1:wrong"},
					},
				).Build(),
				key:      types.NamespacedName{Name: "sentinel", Namespace: "test"},
				replicas: 1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSentinelPool(tt.args.ctx, tt.args.cl, tt.args.key, tt.args.replicas)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSentinelPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("NewSentinelServer() got diff: %v", diff)
			}
		})
	}
}

func TestSentinelPool_IsMonitoringShards(t *testing.T) {
	type args struct {
		ctx    context.Context
		shards []string
	}
	tests := []struct {
		name    string
		sp      SentinelPool
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Returns true",
			sp: []SentinelServer{
				{
					Name: "sentinel-0",
					IP:   "127.0.0.1",
					Port: "26379",
					CRUD: crud.NewFakeCRUD(redis.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{
								[]interface{}{"name", "shard00"},
								[]interface{}{"name", "shard01"},
							}
						},
						InjectError: func() error { return nil },
					}),
				},
				{
					Name: "sentinel-1",
					IP:   "127.0.0.2",
					Port: "26379",
					CRUD: crud.NewFakeCRUD(redis.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{
								[]interface{}{"name", "shard00"},
								[]interface{}{"name", "shard01"},
							}
						},
						InjectError: func() error { return nil },
					}),
				},
				{
					Name: "sentinel-2",
					IP:   "127.0.0.3",
					Port: "26379",
					CRUD: crud.NewFakeCRUD(redis.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{
								[]interface{}{"name", "shard00"},
								[]interface{}{"name", "shard01"},
							}
						},
						InjectError: func() error { return nil },
					}),
				},
			},
			args: args{
				ctx:    context.TODO(),
				shards: []string{"shard00", "shard01"},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Returns false",
			sp: []SentinelServer{
				{
					Name: "sentinel-0",
					IP:   "127.0.0.1",
					Port: "26379",
					CRUD: crud.NewFakeCRUD(redis.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{
								[]interface{}{"name", "shard00"},
							}
						},
						InjectError: func() error { return nil },
					}),
				},
				{
					Name: "sentinel-1",
					IP:   "127.0.0.2",
					Port: "26379",
					CRUD: crud.NewFakeCRUD(redis.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{
								[]interface{}{"name", "shard00"},
								[]interface{}{"name", "shard01"},
							}
						},
						InjectError: func() error { return nil },
					}),
				},
				{
					Name: "sentinel-2",
					IP:   "127.0.0.3",
					Port: "26379",
					CRUD: crud.NewFakeCRUD(redis.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{
								[]interface{}{"name", "shard00"},
								[]interface{}{"name", "shard01"},
							}
						},
						InjectError: func() error { return nil },
					}),
				},
			},
			args: args{
				ctx:    context.TODO(),
				shards: []string{"shard00", "shard01"},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Returns false",
			sp: []SentinelServer{
				{
					Name: "sentinel-0",
					IP:   "127.0.0.1",
					Port: "26379",
					CRUD: crud.NewFakeCRUD(redis.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{}
						},
						InjectError: func() error { return nil },
					}),
				},
				{
					Name: "sentinel-1",
					IP:   "127.0.0.2",
					Port: "26379",
					CRUD: crud.NewFakeCRUD(redis.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{}
						},
						InjectError: func() error { return nil },
					}),
				},
				{
					Name: "sentinel-2",
					IP:   "127.0.0.3",
					Port: "26379",
					CRUD: crud.NewFakeCRUD(redis.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{}
						},
						InjectError: func() error { return nil },
					}),
				},
			},
			args: args{
				ctx:    context.TODO(),
				shards: []string{"shard00", "shard01"},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Returns error",
			sp: []SentinelServer{
				{
					Name: "sentinel-0",
					IP:   "127.0.0.1",
					Port: "26379",
					CRUD: crud.NewFakeCRUD(redis.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{}
						},
						InjectError: func() error { return errors.New("error") },
					}),
				},
			},
			args: args{
				ctx:    context.TODO(),
				shards: []string{"shard00", "shard01"},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.sp.IsMonitoringShards(tt.args.ctx, tt.args.shards)
			if (err != nil) != tt.wantErr {
				t.Errorf("SentinelPool.IsMonitoringShards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SentinelPool.IsMonitoringShards() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSentinelPool_Monitor(t *testing.T) {
	type args struct {
		ctx    context.Context
		shards ShardedCluster
	}
	tests := []struct {
		name    string
		sp      SentinelPool
		args    args
		want    map[string][]string
		wantErr bool
	}{
		{
			name: "No changes",
			sp: []SentinelServer{
				{
					Name: "sentinel-0",
					IP:   "127.0.0.1",
					Port: "26379",
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
					),
				},
				{
					Name: "sentinel-1",
					IP:   "127.0.0.2",
					Port: "26379",
					CRUD: crud.NewFakeCRUD(
						// SentinelMaster response for shard00
						redis.FakeResponse{
							InjectResponse: func() interface{} {
								return &redis.SentinelMasterCmdResult{
									Name: "shard00",
									IP:   "127.0.0.2",
									Port: 3000,
								}
							},
							InjectError: func() error { return nil },
						},
					),
				},
				{
					Name: "sentinel-2",
					IP:   "127.0.0.3",
					Port: "26379",
					CRUD: crud.NewFakeCRUD(
						// SentinelMaster response for shard00
						redis.FakeResponse{
							InjectResponse: func() interface{} {
								return &redis.SentinelMasterCmdResult{
									Name: "shard00",
									IP:   "127.0.0.3",
									Port: 4000,
								}
							},
							InjectError: func() error { return nil },
						},
					),
				},
			},
			args: args{
				ctx: context.TODO(),
				shards: ShardedCluster{{
					Name: "shard00",
					Servers: []RedisServer{
						{
							Name:     "shard00-0",
							IP:       "127.0.0.1",
							Port:     "2000",
							Role:     redis.Master,
							ReadOnly: false,
							CRUD:     nil,
						},
						{
							Name:     "shard00-1",
							IP:       "127.0.0.1",
							Port:     "2001",
							Role:     redis.Slave,
							ReadOnly: true,
							CRUD:     nil,
						},
						{
							Name:     "shard00-2",
							IP:       "127.0.0.1",
							Port:     "2002",
							Role:     redis.Slave,
							ReadOnly: true,
							CRUD:     nil,
						},
					}},
				},
			},
			want:    map[string][]string{},
			wantErr: false,
		},
		{
			name: "Returns changes for all sentinel servers",
			sp: []SentinelServer{
				{
					Name: "sentinel-0",
					IP:   "127.0.0.1",
					Port: "26379",
					CRUD: crud.NewFakeCRUD(
						// SentinelMaster response for shard00
						redis.FakeResponse{
							InjectResponse: func() interface{} {
								return &redis.SentinelMasterCmdResult{}
							},
							InjectError: func() error { return errors.New(shardNotInitializedError) },
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
					),
				},
				{
					Name: "sentinel-1",
					IP:   "127.0.0.2",
					Port: "26379",
					CRUD: crud.NewFakeCRUD(
						// SentinelMaster response for shard00
						redis.FakeResponse{
							InjectResponse: func() interface{} {
								return &redis.SentinelMasterCmdResult{}
							},
							InjectError: func() error { return errors.New(shardNotInitializedError) },
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
					),
				},
				{
					Name: "sentinel-2",
					IP:   "127.0.0.3",
					Port: "26379",
					CRUD: crud.NewFakeCRUD(
						// SentinelMaster response for shard00
						redis.FakeResponse{
							InjectResponse: func() interface{} {
								return &redis.SentinelMasterCmdResult{}
							},
							InjectError: func() error { return errors.New(shardNotInitializedError) },
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
					),
				},
			},
			args: args{
				ctx: context.TODO(),
				shards: ShardedCluster{{
					Name: "shard00",
					Servers: []RedisServer{
						{
							Name:     "shard00-0",
							IP:       "127.0.0.1",
							Port:     "2000",
							Role:     redis.Master,
							ReadOnly: false,
							CRUD:     nil,
						},
						{
							Name:     "shard00-1",
							IP:       "127.0.0.1",
							Port:     "2001",
							Role:     redis.Slave,
							ReadOnly: true,
							CRUD:     nil,
						},
						{
							Name:     "shard00-2",
							IP:       "127.0.0.1",
							Port:     "2002",
							Role:     redis.Slave,
							ReadOnly: true,
							CRUD:     nil,
						},
					}},
				},
			},
			want: map[string][]string{
				"sentinel-0": {"shard00"},
				"sentinel-1": {"shard00"},
				"sentinel-2": {"shard00"},
			},
			wantErr: false,
		},
		{
			name: "Error returned by sentinel-1, sentinel-0 changed",
			sp: []SentinelServer{
				{
					Name: "sentinel-0",
					IP:   "127.0.0.1",
					Port: "26379",
					CRUD: crud.NewFakeCRUD(
						// SentinelMaster response for shard00
						redis.FakeResponse{
							InjectResponse: func() interface{} {
								return &redis.SentinelMasterCmdResult{}
							},
							InjectError: func() error { return errors.New(shardNotInitializedError) },
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
					),
				},
				{
					Name: "sentinel-1",
					IP:   "127.0.0.2",
					Port: "26379",
					CRUD: crud.NewFakeCRUD(
						// SentinelMaster response for shard00
						redis.FakeResponse{
							InjectResponse: func() interface{} {
								return &redis.SentinelMasterCmdResult{}
							},
							InjectError: func() error { return errors.New("error") },
						},
					),
				},
				{
					Name: "sentinel-2",
					IP:   "127.0.0.3",
					Port: "26379",
					// function code should error before reaching sentinel-2
					CRUD: crud.NewFakeCRUD(),
				},
			},
			args: args{
				ctx: context.TODO(),
				shards: ShardedCluster{{
					Name: "shard00",
					Servers: []RedisServer{
						{
							Name:     "shard00-0",
							IP:       "127.0.0.1",
							Port:     "2000",
							Role:     redis.Master,
							ReadOnly: false,
							CRUD:     nil,
						},
						{
							Name:     "shard00-1",
							IP:       "127.0.0.1",
							Port:     "2001",
							Role:     redis.Slave,
							ReadOnly: true,
							CRUD:     nil,
						},
						{
							Name:     "shard00-2",
							IP:       "127.0.0.1",
							Port:     "2002",
							Role:     redis.Slave,
							ReadOnly: true,
							CRUD:     nil,
						},
					}},
				},
			},
			want: map[string][]string{
				"sentinel-0": {"shard00"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.sp.Monitor(tt.args.ctx, tt.args.shards)
			if (err != nil) != tt.wantErr {
				t.Errorf("SentinelPool.Monitor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SentinelPool.Monitor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_applyQuorum(t *testing.T) {
	type args struct {
		responses []saasv1alpha1.MonitoredShards
		quorum    int
	}
	tests := []struct {
		name    string
		args    args
		want    saasv1alpha1.MonitoredShards
		wantErr bool
	}{
		{
			name: "Should return the accepted response",
			args: args{
				responses: []saasv1alpha1.MonitoredShards{
					{
						{Name: "shard01", Master: "127.0.0.1:1111"},
						{Name: "shard02", Master: "127.0.0.2:2222"},
						{Name: "shard03", Master: "127.0.0.3:3333"},
					},
					{
						{Name: "shard03", Master: "127.0.0.3:3333"},
						{Name: "shard02", Master: "127.0.0.2:2222"},
						{Name: "shard01", Master: "127.0.0.1:1111"},
					},
				},
				quorum: 2,
			},
			want: []saasv1alpha1.MonitoredShard{
				{Name: "shard01", Master: "127.0.0.1:1111"},
				{Name: "shard02", Master: "127.0.0.2:2222"},
				{Name: "shard03", Master: "127.0.0.3:3333"},
			},
			wantErr: false,
		},
		{
			name: "Should fail, no quorum",
			args: args{
				responses: []saasv1alpha1.MonitoredShards{
					{
						{Name: "shard01", Master: "127.0.0.1:1111"},
						{Name: "shard02", Master: "127.0.0.2:2222"},
						{Name: "shard03", Master: "127.0.0.3:3333"},
					},
					{
						{Name: "shard03", Master: "127.0.0.3:3333"},
						{Name: "shard02", Master: "127.0.0.2:2222"},
						{Name: "shard01", Master: "127.0.0.4:4444"},
					},
				},
				quorum: 2,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := applyQuorum(tt.args.responses, tt.args.quorum)
			if (err != nil) != tt.wantErr {
				t.Errorf("applyQuorum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("applyQuorum() = %v, want %v", got, tt.want)
			}
		})
	}
}
