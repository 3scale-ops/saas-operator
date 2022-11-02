package envoyconfig

import (
	"testing"
	"time"

	marin3rv1alpha1 "github.com/3scale-ops/marin3r/apis/marin3r/v1alpha1"
	"github.com/3scale-ops/marin3r/pkg/envoy"
	envoy_serializer "github.com/3scale-ops/marin3r/pkg/envoy/serializer"
	"github.com/MakeNowJust/heredoc"
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_service_runtime_v3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	"github.com/go-test/deep"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
)

func TestNew(t *testing.T) {
	type args struct {
		key       types.NamespacedName
		nodeID    string
		resources []envoy.Resource
	}
	tests := []struct {
		name    string
		args    args
		want    *marin3rv1alpha1.EnvoyConfig
		wantErr bool
	}{
		{
			name: "Generates an EnvoyConfig",
			args: args{
				key:    types.NamespacedName{Name: "test", Namespace: "ns"},
				nodeID: "test",
				resources: []envoy.Resource{
					&envoy_config_cluster_v3.Cluster{
						Name:           "cluster1",
						ConnectTimeout: durationpb.New(2 * time.Second),
						ClusterDiscoveryType: &envoy_config_cluster_v3.Cluster_Type{
							Type: envoy_config_cluster_v3.Cluster_STRICT_DNS,
						},
						LbPolicy: envoy_config_cluster_v3.Cluster_ROUND_ROBIN,
						LoadAssignment: &envoy_config_endpoint_v3.ClusterLoadAssignment{
							ClusterName: "cluster1",
						},
					},
					&envoy_config_route_v3.RouteConfiguration{
						Name: "route1",
						VirtualHosts: []*envoy_config_route_v3.VirtualHost{{
							Name:    "vhost",
							Domains: []string{"*"},
							Routes: []*envoy_config_route_v3.Route{{
								Match: &envoy_config_route_v3.RouteMatch{
									PathSpecifier: &envoy_config_route_v3.RouteMatch_Prefix{Prefix: "/"}},
								Action: &envoy_config_route_v3.Route_DirectResponse{
									DirectResponse: &envoy_config_route_v3.DirectResponseAction{Status: 200}},
							}},
						}},
					},
					&envoy_service_runtime_v3.Runtime{
						Name: "runtime1",
						Layer: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"static_layer_0": {Kind: &structpb.Value_StringValue{StringValue: "value"}},
							}},
					},
					&envoy_config_listener_v3.Listener{
						Name: "listener1",
						Address: &envoy_config_core_v3.Address{
							Address: &envoy_config_core_v3.Address_SocketAddress{
								SocketAddress: &envoy_config_core_v3.SocketAddress{
									Address: "0.0.0.0",
									PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
										PortValue: 8443,
									}}}},
						FilterChains: []*envoy_config_listener_v3.FilterChain{},
					},
				},
			},
			want: &marin3rv1alpha1.EnvoyConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "ns"},
				Spec: marin3rv1alpha1.EnvoyConfigSpec{
					NodeID:        "test",
					Serialization: pointer.String(string(envoy_serializer.YAML)),
					EnvoyAPI:      pointer.StringPtr(envoy.APIv3.String()),
					EnvoyResources: &marin3rv1alpha1.EnvoyResources{
						Clusters: []marin3rv1alpha1.EnvoyResource{{
							Value: heredoc.Doc(`
								connect_timeout: 2s
								load_assignment:
								  cluster_name: cluster1
								name: cluster1
								type: STRICT_DNS
							`)}},
						Routes: []marin3rv1alpha1.EnvoyResource{{
							Value: heredoc.Doc(`
								name: route1
								virtual_hosts:
								- domains:
								  - '*'
								  name: vhost
								  routes:
								  - direct_response:
								      status: 200
								    match:
								      prefix: /
						`)}},
						Listeners: []marin3rv1alpha1.EnvoyResource{{
							Value: heredoc.Doc(`
								address:
								  socket_address:
								    address: 0.0.0.0
								    port_value: 8443
								name: listener1
						`)}},
						Runtimes: []marin3rv1alpha1.EnvoyResource{{
							Value: heredoc.Doc(`
								layer:
								  static_layer_0: value
								name: runtime1
						`)}},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, err := New(tt.args.key, tt.args.nodeID, tt.args.resources...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := fn()
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("New() = diff %v", diff)
			}
		})
	}
}
