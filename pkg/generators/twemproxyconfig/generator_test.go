package twemproxyconfig

import (
	"context"
	"fmt"
	"testing"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators"
	redis_client "github.com/3scale/saas-operator/pkg/redis/client"
	"github.com/3scale/saas-operator/pkg/redis/server"
	"github.com/3scale/saas-operator/pkg/resource_builders/twemproxy"
	"github.com/3scale/saas-operator/pkg/util"
	"github.com/go-logr/logr"
	"github.com/go-test/deep"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestNewGenerator(t *testing.T) {
	type args struct {
		ctx      context.Context
		instance *saasv1alpha1.TwemproxyConfig
		cl       client.Client
		pool     *server.ServerPool
		log      logr.Logger
	}
	tests := []struct {
		name    string
		args    args
		want    Generator
		wantErr bool
	}{
		{
			name: "Populates the generation with the current cluster topology (target masters)",
			args: args{
				ctx: context.TODO(),
				instance: &saasv1alpha1.TwemproxyConfig{
					ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"},
					Spec: saasv1alpha1.TwemproxyConfigSpec{
						SentinelURIs: []string{"redis://127.0.0.1:26379"},
						ServerPools: []saasv1alpha1.TwemproxyServerPool{{
							Name:   "test-pool",
							Target: util.Pointer(saasv1alpha1.Masters),
							Topology: []saasv1alpha1.ShardedRedisTopology{
								{ShardName: "l-shard00", PhysicalShard: "shard0"},
								{ShardName: "l-shard01", PhysicalShard: "shard0"},
								{ShardName: "l-shard02", PhysicalShard: "shard0"},
								{ShardName: "l-shard03", PhysicalShard: "shard1"},
								{ShardName: "l-shard04", PhysicalShard: "shard1"},
							},
							BindAddress: "0.0.0.0:22121",
							Timeout:     5000,
							TCPBacklog:  512,
							PreConnect:  false,
						}},
						ReconcileServerPools: util.Pointer(true),
					},
				},
				cl: nil,
				pool: server.NewServerPool(
					// redis servers
					server.NewFakeServerWithFakeClient("127.0.0.1", "1000", redis_client.NewPredefinedRedisFakeResponse("role-master", nil)),
					server.NewFakeServerWithFakeClient("127.0.0.1", "2000", redis_client.NewPredefinedRedisFakeResponse("role-slave", nil)),
					server.NewFakeServerWithFakeClient("127.0.0.1", "3000", redis_client.NewPredefinedRedisFakeResponse("role-slave", nil)),
					server.NewFakeServerWithFakeClient("127.0.0.1", "4000", redis_client.NewPredefinedRedisFakeResponse("role-slave", nil)),
					server.NewFakeServerWithFakeClient("127.0.0.1", "5000", redis_client.NewPredefinedRedisFakeResponse("role-master", nil)),
					server.NewFakeServerWithFakeClient("127.0.0.1", "6000", redis_client.NewPredefinedRedisFakeResponse("role-slave", nil)),
					// sentinel
					server.NewFakeServerWithFakeClient("127.0.0.1", "26379",
						redis_client.FakeResponse{
							// cmd: Ping
							InjectResponse: func() interface{} { return nil },
							InjectError:    func() error { return nil },
						},
						redis_client.FakeResponse{
							// cmd: SentinelMasters()
							InjectResponse: func() interface{} {
								return []interface{}{
									[]interface{}{"name", "shard0", "ip", "127.0.0.1", "port", "1000"},
									[]interface{}{"name", "shard1", "ip", "127.0.0.1", "port", "5000"},
								}
							},
							InjectError: func() error { return nil },
						},
						redis_client.FakeResponse{
							// cmd: SentinelMaster (shard0)
							InjectResponse: func() interface{} {
								return &redis_client.SentinelMasterCmdResult{Name: "shard0", IP: "127.0.0.1", Port: 1000, Flags: "master"}
							},
							InjectError: func() error { return nil },
						},
						redis_client.FakeResponse{
							// cmd: SentinelGetMasterAddrByName (shard0)
							InjectResponse: func() interface{} {
								return []string{"127.0.0.1", "1000"}
							},
							InjectError: func() error { return nil },
						},
						redis_client.FakeResponse{
							// cmd: SentinelMaster (shard1)
							InjectResponse: func() interface{} {
								return &redis_client.SentinelMasterCmdResult{Name: "shard1", IP: "127.0.0.1", Port: 5000, Flags: "master"}
							},
							InjectError: func() error { return nil },
						},
						redis_client.FakeResponse{
							// cmd: SentinelGetMasterAddrByName (shard1)
							InjectResponse: func() interface{} {
								return []string{"127.0.0.1", "5000"}
							},
							InjectError: func() error { return nil },
						},
					),
				),
				log: logr.Discard(),
			},
			want: Generator{
				BaseOptionsV2: generators.BaseOptionsV2{
					Component:    "twemproxy",
					InstanceName: "test",
					Namespace:    "test",
					Labels: map[string]string{
						"app":     component,
						"part-of": "3scale-saas",
					}},
				Spec: saasv1alpha1.TwemproxyConfigSpec{
					SentinelURIs: []string{"redis://127.0.0.1:26379"},
					ServerPools: []saasv1alpha1.TwemproxyServerPool{{
						Name:   "test-pool",
						Target: util.Pointer(saasv1alpha1.Masters),
						Topology: []saasv1alpha1.ShardedRedisTopology{
							{ShardName: "l-shard00", PhysicalShard: "shard0"},
							{ShardName: "l-shard01", PhysicalShard: "shard0"},
							{ShardName: "l-shard02", PhysicalShard: "shard0"},
							{ShardName: "l-shard03", PhysicalShard: "shard1"},
							{ShardName: "l-shard04", PhysicalShard: "shard1"},
						},
						BindAddress: "0.0.0.0:22121",
						Timeout:     5000,
						TCPBacklog:  512,
						PreConnect:  false,
					}},
					ReconcileServerPools: util.Pointer(true),
				},
				masterTargets: map[string]twemproxy.Server{
					"shard0": {
						Address:  "127.0.0.1:1000",
						Priority: 1,
					},
					"shard1": {
						Address:  "127.0.0.1:5000",
						Priority: 1,
					},
				},
				slaverwTargets: nil,
			},
			wantErr: false,
		},
		{
			name: "Returns error (target masters)",
			args: args{
				ctx: context.TODO(),
				instance: &saasv1alpha1.TwemproxyConfig{
					ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"},
					Spec: saasv1alpha1.TwemproxyConfigSpec{
						SentinelURIs: []string{"redis://127.0.0.1:26379"},
						ServerPools: []saasv1alpha1.TwemproxyServerPool{{
							Name:   "test-pool",
							Target: util.Pointer(saasv1alpha1.Masters),
							Topology: []saasv1alpha1.ShardedRedisTopology{
								{ShardName: "l-shard00", PhysicalShard: "shard0"},
								{ShardName: "l-shard01", PhysicalShard: "shard0"},
								{ShardName: "l-shard02", PhysicalShard: "shard0"},
								{ShardName: "l-shard03", PhysicalShard: "shard1"},
								{ShardName: "l-shard04", PhysicalShard: "shard1"},
							},
							BindAddress: "0.0.0.0:22121",
							Timeout:     5000,
							TCPBacklog:  512,
							PreConnect:  false,
						}},
						ReconcileServerPools: util.Pointer(true),
					},
				},
				cl: nil,
				pool: server.NewServerPool(
					// sentinel
					server.NewFakeServerWithFakeClient("127.0.0.1", "26379",
						redis_client.FakeResponse{
							// cmd: Ping
							InjectResponse: func() interface{} { return nil },
							InjectError:    func() error { return fmt.Errorf("error") },
						},
					),
				),
				log: logr.Discard(),
			},
			want:    Generator{},
			wantErr: true,
		},
		{
			name: "Populates the generation with the current cluster topology (target rw slaves)",
			args: args{
				ctx: context.TODO(),
				instance: &saasv1alpha1.TwemproxyConfig{
					ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"},
					Spec: saasv1alpha1.TwemproxyConfigSpec{
						SentinelURIs: []string{"redis://127.0.0.1:26379"},
						ServerPools: []saasv1alpha1.TwemproxyServerPool{{
							Name:   "test-pool",
							Target: util.Pointer(saasv1alpha1.SlavesRW),
							Topology: []saasv1alpha1.ShardedRedisTopology{
								{ShardName: "l-shard00", PhysicalShard: "shard0"},
								{ShardName: "l-shard01", PhysicalShard: "shard0"},
								{ShardName: "l-shard02", PhysicalShard: "shard0"},
								{ShardName: "l-shard03", PhysicalShard: "shard1"},
								{ShardName: "l-shard04", PhysicalShard: "shard1"},
							},
							BindAddress: "0.0.0.0:22121",
							Timeout:     5000,
							TCPBacklog:  512,
							PreConnect:  false,
						}},
						ReconcileServerPools: util.Pointer(true),
					},
				},
				cl: nil,
				pool: server.NewServerPool(
					// redis servers
					server.NewFakeServerWithFakeClient("127.0.0.1", "1000", redis_client.NewPredefinedRedisFakeResponse("role-master", nil)),
					server.NewFakeServerWithFakeClient("127.0.0.1", "2000",
						redis_client.NewPredefinedRedisFakeResponse("role-slave", nil),
						redis_client.NewPredefinedRedisFakeResponse("slave-read-only-no", nil),
					),
					server.NewFakeServerWithFakeClient("127.0.0.1", "3000",
						redis_client.NewPredefinedRedisFakeResponse("role-slave", nil),
						redis_client.NewPredefinedRedisFakeResponse("slave-read-only-yes", nil),
					),
					server.NewFakeServerWithFakeClient("127.0.0.1", "4000",
						redis_client.NewPredefinedRedisFakeResponse("role-slave", nil),
						redis_client.NewPredefinedRedisFakeResponse("slave-read-only-no", nil),
					),
					server.NewFakeServerWithFakeClient("127.0.0.1", "5000", redis_client.NewPredefinedRedisFakeResponse("role-master", nil)),
					server.NewFakeServerWithFakeClient("127.0.0.1", "6000",
						redis_client.NewPredefinedRedisFakeResponse("role-slave", nil),
						redis_client.NewPredefinedRedisFakeResponse("slave-read-only-yes", nil),
					),
					// sentinel
					server.NewFakeServerWithFakeClient("127.0.0.1", "26379",
						redis_client.FakeResponse{
							// cmd: Ping
							InjectResponse: func() interface{} { return nil },
							InjectError:    func() error { return nil },
						},
						redis_client.FakeResponse{
							// cmd: SentinelMasters()
							InjectResponse: func() interface{} {
								return []interface{}{
									[]interface{}{"name", "shard0", "ip", "127.0.0.1", "port", "1000"},
									[]interface{}{"name", "shard1", "ip", "127.0.0.1", "port", "5000"},
								}
							},
							InjectError: func() error { return nil },
						},
						redis_client.FakeResponse{
							// cmd: SentinelMaster (shard0)
							InjectResponse: func() interface{} {
								return &redis_client.SentinelMasterCmdResult{Name: "shard0", IP: "127.0.0.1", Port: 1000, Flags: "master"}
							},
							InjectError: func() error { return nil },
						},
						redis_client.FakeResponse{
							// cmd: SentinelGetMasterAddrByName (shard0)
							InjectResponse: func() interface{} {
								return []string{"127.0.0.1", "1000"}
							},
							InjectError: func() error { return nil },
						},
						redis_client.FakeResponse{
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
						redis_client.FakeResponse{
							// cmd: SentinelMaster (shard1)
							InjectResponse: func() interface{} {
								return &redis_client.SentinelMasterCmdResult{Name: "shard1", IP: "127.0.0.1", Port: 5000, Flags: "master"}
							},
							InjectError: func() error { return nil },
						},
						redis_client.FakeResponse{
							// cmd: SentinelGetMasterAddrByName (shard1)
							InjectResponse: func() interface{} {
								return []string{"127.0.0.1", "5000"}
							},
							InjectError: func() error { return nil },
						},
						redis_client.FakeResponse{
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
					),
				),
				log: logr.Discard(),
			},
			want: Generator{
				BaseOptionsV2: generators.BaseOptionsV2{
					Component:    "twemproxy",
					InstanceName: "test",
					Namespace:    "test",
					Labels: map[string]string{
						"app":     component,
						"part-of": "3scale-saas",
					}},
				Spec: saasv1alpha1.TwemproxyConfigSpec{
					SentinelURIs: []string{"redis://127.0.0.1:26379"},
					ServerPools: []saasv1alpha1.TwemproxyServerPool{{
						Name:   "test-pool",
						Target: util.Pointer(saasv1alpha1.SlavesRW),
						Topology: []saasv1alpha1.ShardedRedisTopology{
							{ShardName: "l-shard00", PhysicalShard: "shard0"},
							{ShardName: "l-shard01", PhysicalShard: "shard0"},
							{ShardName: "l-shard02", PhysicalShard: "shard0"},
							{ShardName: "l-shard03", PhysicalShard: "shard1"},
							{ShardName: "l-shard04", PhysicalShard: "shard1"},
						},
						BindAddress: "0.0.0.0:22121",
						Timeout:     5000,
						TCPBacklog:  512,
						PreConnect:  false,
					}},
					ReconcileServerPools: util.Pointer(true),
				},
				masterTargets: map[string]twemproxy.Server{
					"shard0": {
						Address:  "127.0.0.1:1000",
						Priority: 1,
					},
					"shard1": {
						Address:  "127.0.0.1:5000",
						Priority: 1,
					},
				},
				slaverwTargets: map[string]twemproxy.Server{
					"shard0": {
						Address:  "127.0.0.1:2000",
						Priority: 1,
					},
					"shard1": {
						Address:  "127.0.0.1:4000",
						Priority: 1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Populates the generation with the current cluster topology (target rw slaves with failover to master)",
			args: args{
				ctx: context.TODO(),
				instance: &saasv1alpha1.TwemproxyConfig{
					ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"},
					Spec: saasv1alpha1.TwemproxyConfigSpec{
						SentinelURIs: []string{"redis://127.0.0.1:26379"},
						ServerPools: []saasv1alpha1.TwemproxyServerPool{{
							Name:   "test-pool",
							Target: util.Pointer(saasv1alpha1.SlavesRW),
							Topology: []saasv1alpha1.ShardedRedisTopology{
								{ShardName: "l-shard00", PhysicalShard: "shard0"},
								{ShardName: "l-shard01", PhysicalShard: "shard0"},
								{ShardName: "l-shard02", PhysicalShard: "shard0"},
								{ShardName: "l-shard03", PhysicalShard: "shard1"},
								{ShardName: "l-shard04", PhysicalShard: "shard1"},
							},
							BindAddress: "0.0.0.0:22121",
							Timeout:     5000,
							TCPBacklog:  512,
							PreConnect:  false,
						}},
						ReconcileServerPools: util.Pointer(true),
					},
				},
				cl: nil,
				pool: server.NewServerPool(
					// redis servers
					server.NewFakeServerWithFakeClient("127.0.0.1", "1000", redis_client.NewPredefinedRedisFakeResponse("role-master", nil)),
					server.NewFakeServerWithFakeClient("127.0.0.1", "2000"), // is down
					server.NewFakeServerWithFakeClient("127.0.0.1", "3000",
						redis_client.NewPredefinedRedisFakeResponse("role-slave", nil),
						redis_client.NewPredefinedRedisFakeResponse("slave-read-only-yes", nil),
					),
					server.NewFakeServerWithFakeClient("127.0.0.1", "4000",
						redis_client.NewPredefinedRedisFakeResponse("role-slave", nil),
						redis_client.NewPredefinedRedisFakeResponse("slave-read-only-no", nil),
					),
					server.NewFakeServerWithFakeClient("127.0.0.1", "5000", redis_client.NewPredefinedRedisFakeResponse("role-master", nil)),
					server.NewFakeServerWithFakeClient("127.0.0.1", "6000",
						redis_client.NewPredefinedRedisFakeResponse("role-slave", nil),
						redis_client.NewPredefinedRedisFakeResponse("slave-read-only-yes", nil),
					),
					// sentinel
					server.NewFakeServerWithFakeClient("127.0.0.1", "26379",
						redis_client.FakeResponse{
							// cmd: Ping
							InjectResponse: func() interface{} { return nil },
							InjectError:    func() error { return nil },
						},
						redis_client.FakeResponse{
							// cmd: SentinelMasters()
							InjectResponse: func() interface{} {
								return []interface{}{
									[]interface{}{"name", "shard0", "ip", "127.0.0.1", "port", "1000"},
									[]interface{}{"name", "shard1", "ip", "127.0.0.1", "port", "5000"},
								}
							},
							InjectError: func() error { return nil },
						},
						redis_client.FakeResponse{
							// cmd: SentinelMaster (shard0)
							InjectResponse: func() interface{} {
								return &redis_client.SentinelMasterCmdResult{Name: "shard0", IP: "127.0.0.1", Port: 1000, Flags: "master"}
							},
							InjectError: func() error { return nil },
						},
						redis_client.FakeResponse{
							// cmd: SentinelGetMasterAddrByName (shard0)
							InjectResponse: func() interface{} {
								return []string{"127.0.0.1", "1000"}
							},
							InjectError: func() error { return nil },
						},
						redis_client.FakeResponse{
							// cmd: SentinelSlaves (shard0)
							InjectResponse: func() interface{} {
								return []interface{}{
									[]interface{}{
										"name", "127.0.0.1:2000",
										"ip", "127.0.0.1",
										"port", "2000",
										"flags", "slave,s_down", // slave is down
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
						redis_client.FakeResponse{
							// cmd: SentinelMaster (shard1)
							InjectResponse: func() interface{} {
								return &redis_client.SentinelMasterCmdResult{Name: "shard1", IP: "127.0.0.1", Port: 5000, Flags: "master"}
							},
							InjectError: func() error { return nil },
						},
						redis_client.FakeResponse{
							// cmd: SentinelGetMasterAddrByName (shard1)
							InjectResponse: func() interface{} {
								return []string{"127.0.0.1", "5000"}
							},
							InjectError: func() error { return nil },
						},
						redis_client.FakeResponse{
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
					),
				),
				log: logr.Discard(),
			},
			want: Generator{
				BaseOptionsV2: generators.BaseOptionsV2{
					Component:    "twemproxy",
					InstanceName: "test",
					Namespace:    "test",
					Labels: map[string]string{
						"app":     component,
						"part-of": "3scale-saas",
					}},
				Spec: saasv1alpha1.TwemproxyConfigSpec{
					SentinelURIs: []string{"redis://127.0.0.1:26379"},
					ServerPools: []saasv1alpha1.TwemproxyServerPool{{
						Name:   "test-pool",
						Target: util.Pointer(saasv1alpha1.SlavesRW),
						Topology: []saasv1alpha1.ShardedRedisTopology{
							{ShardName: "l-shard00", PhysicalShard: "shard0"},
							{ShardName: "l-shard01", PhysicalShard: "shard0"},
							{ShardName: "l-shard02", PhysicalShard: "shard0"},
							{ShardName: "l-shard03", PhysicalShard: "shard1"},
							{ShardName: "l-shard04", PhysicalShard: "shard1"},
						},
						BindAddress: "0.0.0.0:22121",
						Timeout:     5000,
						TCPBacklog:  512,
						PreConnect:  false,
					}},
					ReconcileServerPools: util.Pointer(true),
				},
				masterTargets: map[string]twemproxy.Server{
					"shard0": {
						Address:  "127.0.0.1:1000",
						Priority: 1,
					},
					"shard1": {
						Address:  "127.0.0.1:5000",
						Priority: 1,
					},
				},
				slaverwTargets: map[string]twemproxy.Server{
					"shard0": {
						Address:  "127.0.0.1:1000",
						Priority: 1,
					},
					"shard1": {
						Address:  "127.0.0.1:4000",
						Priority: 1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Returns error (target rw slaves)",
			args: args{
				ctx: context.TODO(),
				instance: &saasv1alpha1.TwemproxyConfig{
					ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"},
					Spec: saasv1alpha1.TwemproxyConfigSpec{
						SentinelURIs: []string{"redis://127.0.0.1:26379"},
						ServerPools: []saasv1alpha1.TwemproxyServerPool{{
							Name:   "test-pool",
							Target: util.Pointer(saasv1alpha1.SlavesRW),
							Topology: []saasv1alpha1.ShardedRedisTopology{
								{ShardName: "l-shard00", PhysicalShard: "shard0"},
								{ShardName: "l-shard01", PhysicalShard: "shard0"},
								{ShardName: "l-shard02", PhysicalShard: "shard0"},
								{ShardName: "l-shard03", PhysicalShard: "shard1"},
								{ShardName: "l-shard04", PhysicalShard: "shard1"},
							},
							BindAddress: "0.0.0.0:22121",
							Timeout:     5000,
							TCPBacklog:  512,
							PreConnect:  false,
						}},
						ReconcileServerPools: util.Pointer(true),
					},
				},
				cl: nil,
				pool: server.NewServerPool(
					// redis servers
					server.NewFakeServerWithFakeClient("127.0.0.1", "1000"), // is down
					server.NewFakeServerWithFakeClient("127.0.0.1", "4000",
						redis_client.NewPredefinedRedisFakeResponse("role-slave", nil),
						redis_client.NewPredefinedRedisFakeResponse("slave-read-only-no", nil),
					),
					server.NewFakeServerWithFakeClient("127.0.0.1", "5000", redis_client.NewPredefinedRedisFakeResponse("role-master", nil)),
					server.NewFakeServerWithFakeClient("127.0.0.1", "6000",
						redis_client.NewPredefinedRedisFakeResponse("role-slave", nil),
						redis_client.NewPredefinedRedisFakeResponse("slave-read-only-yes", nil),
					),
					// sentinel
					server.NewFakeServerWithFakeClient("127.0.0.1", "26379",
						redis_client.FakeResponse{
							// cmd: Ping
							InjectResponse: func() interface{} { return nil },
							InjectError:    func() error { return nil },
						},
						redis_client.FakeResponse{
							// cmd: SentinelMasters()
							InjectResponse: func() interface{} {
								return []interface{}{
									[]interface{}{"name", "shard0", "ip", "127.0.0.1", "port", "1000"},
									[]interface{}{"name", "shard1", "ip", "127.0.0.1", "port", "5000"},
								}
							},
							InjectError: func() error { return nil },
						},
						redis_client.FakeResponse{
							// cmd: SentinelMaster (shard0)
							InjectResponse: func() interface{} {
								return &redis_client.SentinelMasterCmdResult{Name: "shard0", IP: "127.0.0.1", Port: 1000, Flags: "master,o_down"}
							},
							InjectError: func() error { return nil },
						},
						redis_client.FakeResponse{
							// cmd: SentinelGetMasterAddrByName (shard0)
							InjectResponse: func() interface{} {
								return []string{"127.0.0.1", "1000"}
							},
							InjectError: func() error { return nil },
						},
						redis_client.FakeResponse{
							// cmd: SentinelMaster (shard1)
							InjectResponse: func() interface{} {
								return &redis_client.SentinelMasterCmdResult{Name: "shard1", IP: "127.0.0.1", Port: 5000, Flags: "master"}
							},
							InjectError: func() error { return nil },
						},
						redis_client.FakeResponse{
							// cmd: SentinelGetMasterAddrByName (shard1)
							InjectResponse: func() interface{} {
								return []string{"127.0.0.1", "5000"}
							},
							InjectError: func() error { return nil },
						},
						redis_client.FakeResponse{
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
					),
				),
				log: logr.Discard(),
			},
			want:    Generator{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGenerator(tt.args.ctx, tt.args.instance, tt.args.cl, tt.args.pool, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGenerator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			deep.CompareUnexportedFields = true
			if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(Generator{}), cmpopts.IgnoreUnexported(twemproxy.Server{})); len(diff) != 0 {
				t.Errorf("NewGenerator() = diff %v", diff)
			}
		})
	}
}
