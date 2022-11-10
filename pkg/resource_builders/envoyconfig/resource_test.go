package envoyconfig

import (
	"reflect"
	"testing"
	"time"

	marin3rv1alpha1 "github.com/3scale-ops/marin3r/apis/marin3r/v1alpha1"
	"github.com/3scale-ops/marin3r/pkg/envoy"
	envoy_serializer "github.com/3scale-ops/marin3r/pkg/envoy/serializer"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
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
		resources []saasv1alpha1.EnvoyResource
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
				key:    types.NamespacedName{Name: "test", Namespace: "default"},
				nodeID: "test",
				resources: []saasv1alpha1.EnvoyResource{
					{
						Name: "my_cluster",
						Cluster: &saasv1alpha1.Cluster{
							EnvoyResourceGeneratorMeta: saasv1alpha1.EnvoyResourceGeneratorMeta{
								GeneratorVersion: pointer.String("v1"),
							},
							Host:    "localhost",
							Port:    8080,
							IsHttp2: pointer.Bool(false),
						},
					},
					{
						Name: "my_listener",
						ListenerHttp: &saasv1alpha1.ListenerHttp{
							EnvoyResourceGeneratorMeta: saasv1alpha1.EnvoyResourceGeneratorMeta{
								GeneratorVersion: pointer.String("v1"),
							}, Port: 0,
							RouteConfigName: "routeconfig",
						},
					},
				},
			},
			want: &marin3rv1alpha1.EnvoyConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
				},
				Spec: marin3rv1alpha1.EnvoyConfigSpec{
					NodeID:        "test",
					Serialization: pointer.String(string(envoy_serializer.YAML)),
					EnvoyAPI:      pointer.String(envoy.APIv3.String()),
					EnvoyResources: &marin3rv1alpha1.EnvoyResources{
						Clusters: []marin3rv1alpha1.EnvoyResource{{
							Value: heredoc.Doc(`
                                connect_timeout: 1s
                                dns_lookup_family: V4_ONLY
                                load_assignment:
                                  cluster_name: my_cluster
                                  endpoints:
                                  - lb_endpoints:
                                    - endpoint:
                                        address:
                                          socket_address:
                                            address: localhost
                                            port_value: 8080
                                name: my_cluster
                                type: STRICT_DNS
							`),
						}},
						Listeners: []marin3rv1alpha1.EnvoyResource{{
							Value: heredoc.Doc(`
                                address:
                                  socket_address:
                                    address: 0.0.0.0
                                    port_value: 0
                                filter_chains:
                                - filters:
                                  - name: envoy.filters.network.http_connection_manager
                                    typed_config:
                                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                                      access_log:
                                      - name: envoy.access_loggers.file
                                        typed_config:
                                          '@type': type.googleapis.com/envoy.extensions.access_loggers.file.v3.FileAccessLog
                                          log_format:
                                            json_format:
                                              authority: '%REQ(:AUTHORITY)%'
                                              bytes_received: '%BYTES_RECEIVED%'
                                              bytes_sent: '%BYTES_SENT%'
                                              client_ip: '%REQ(X-ENVOY-EXTERNAL-ADDRESS)%'
                                              duration: '%DURATION%'
                                              listener: my_listener
                                              method: '%REQ(:METHOD)%'
                                              path: '%REQ(X-ENVOY-ORIGINAL-PATH?:PATH)%'
                                              protocol: '%PROTOCOL%'
                                              response_code: '%RESPONSE_CODE%'
                                              response_code_details: '%RESPONSE_CODE_DETAILS%'
                                              response_flags: '%RESPONSE_FLAGS%'
                                              upstream_cluster: '%UPSTREAM_CLUSTER%'
                                              upstream_service_time: '%RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)%'
                                              user_agent: '%REQ(USER-AGENT)%'
                                          path: /dev/stdout
                                      common_http_protocol_options:
                                        idle_timeout: 3600s
                                        max_connection_duration: 900s
                                      http_filters:
                                      - name: envoy.filters.http.router
                                      http_protocol_options: {}
                                      http2_protocol_options:
                                        initial_connection_window_size: 1048576
                                        initial_stream_window_size: 65536
                                        max_concurrent_streams: 100
                                      rds:
                                        config_source:
                                          ads: {}
                                          resource_api_version: V3
                                        route_config_name: routeconfig
                                      request_timeout: 300s
                                      stat_prefix: my_listener
                                      stream_idle_timeout: 300s
                                      use_remote_address: true
                                listener_filters:
                                - name: envoy.filters.listener.proxy_protocol
                                  typed_config:
                                    '@type': type.googleapis.com/envoy.extensions.filters.listener.proxy_protocol.v3.ProxyProtocol
                                name: my_listener
                                per_connection_buffer_limit_bytes: 32768
							`),
						}},
						Routes:   []marin3rv1alpha1.EnvoyResource{},
						Runtimes: []marin3rv1alpha1.EnvoyResource{},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.key, tt.args.nodeID, tt.args.resources...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			ec := got()
			if diff := deep.Equal(ec, tt.want); len(diff) > 0 {
				t.Errorf("New() = got diff %v", diff)
			}
		})
	}
}

func Test_newFromProtos(t *testing.T) {
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
			fn, err := newFromProtos(tt.args.key, tt.args.nodeID, tt.args.resources)
			if (err != nil) != tt.wantErr {
				t.Errorf("newFromProtos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := fn()
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("newFromProtos() = diff %v", diff)
			}
		})
	}
}

func Test_inspect(t *testing.T) {
	type args struct {
		v *saasv1alpha1.EnvoyResource
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 interface{}
	}{
		{
			name: "",
			args: args{
				v: &saasv1alpha1.EnvoyResource{
					Name: "test",
					ListenerHttp: &saasv1alpha1.ListenerHttp{
						EnvoyResourceGeneratorMeta: saasv1alpha1.EnvoyResourceGeneratorMeta{
							GeneratorVersion: pointer.String("v1"),
						},
					},
				},
			},
			want: "ListenerHttp_v1",
			want1: &saasv1alpha1.ListenerHttp{
				EnvoyResourceGeneratorMeta: saasv1alpha1.EnvoyResourceGeneratorMeta{
					GeneratorVersion: pointer.String("v1"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := inspect(tt.args.v)
			if got != tt.want {
				t.Errorf("inspect() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("inspect() got1 = %+v, want %+v", got1, tt.want1)
			}
		})
	}
}
func Test_envoyResourceFactory_newResource(t *testing.T) {
	type args struct {
		resourceName string
		functionName string
		descriptor   envoyResourceDescriptor
	}
	tests := []struct {
		name    string
		erf     envoyResourceFactory
		args    args
		want    envoy.Resource
		wantErr bool
	}{
		{
			name: "Generates a resource proto",
			erf:  generator,
			args: args{
				resourceName: "test",
				functionName: "Runtime_v1",
				descriptor: &saasv1alpha1.Runtime{
					ListenerNames: []string{"http", "https"},
				},
			},
			want: func() envoy.Resource {
				l, _ := structpb.NewStruct(map[string]interface{}{
					"envoy": map[string]interface{}{
						"resource_limits": map[string]interface{}{
							"listener": map[string]interface{}{
								"http": map[string]interface{}{
									"connection_limit": 10000,
								},
								"https": map[string]interface{}{
									"connection_limit": 10000,
								},
							},
						},
					},
					"overload": map[string]interface{}{
						"global_downstream_max_connections": 50000,
					},
				})
				return &envoy_service_runtime_v3.Runtime{
					Name:  "test",
					Layer: l,
				}
			}(),
			wantErr: false,
		},
		{
			name: "Unregistered function",
			erf:  generator,
			args: args{
				resourceName: "test",
				functionName: "Runtime_xx",
				descriptor:   nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.erf.newResource(tt.args.resourceName, tt.args.functionName, tt.args.descriptor)
			if (err != nil) != tt.wantErr {
				t.Errorf("envoyResourceFactory.newResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("envoyResourceFactory.newResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
