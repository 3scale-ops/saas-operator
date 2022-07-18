package redis

import (
	"context"
	"errors"
	"reflect"
	"testing"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/redis/crud"
	redis "github.com/3scale/saas-operator/pkg/redis/crud/client"
	"github.com/go-test/deep"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8s "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

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
