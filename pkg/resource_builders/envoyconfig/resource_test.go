package envoyconfig

import (
	"testing"
	"time"

	marin3rv1alpha1 "github.com/3scale-ops/marin3r/apis/marin3r/v1alpha1"
	"github.com/3scale-ops/marin3r/pkg/envoy"
	envoy_serializer "github.com/3scale-ops/marin3r/pkg/envoy/serializer"
	marin3r_pointer "github.com/3scale-ops/marin3r/pkg/util/pointer"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	descriptor "github.com/3scale/saas-operator/pkg/resource_builders/envoyconfig/descriptor"
	"github.com/3scale/saas-operator/pkg/resource_builders/envoyconfig/factory"
	"github.com/3scale/saas-operator/pkg/util"
	"github.com/MakeNowJust/heredoc"
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_service_runtime_v3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	"github.com/go-test/deep"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/protobuf/types/known/durationpb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
)

func TestNew(t *testing.T) {
	type args struct {
		key       types.NamespacedName
		nodeID    string
		factory   factory.EnvoyDynamicConfigFactory
		resources []descriptor.EnvoyDynamicConfigDescriptor
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
				key:     types.NamespacedName{Name: "test", Namespace: "default"},
				nodeID:  "test",
				factory: factory.Default(),
				resources: saasv1alpha1.MapOfEnvoyDynamicConfig{
					"my_cluster": {
						GeneratorVersion: pointer.String("v1"),
						Cluster: &saasv1alpha1.Cluster{
							Host:    "localhost",
							Port:    8080,
							IsHttp2: pointer.Bool(false),
						},
					},
					"my_listener": {
						GeneratorVersion: pointer.String("v1"),
						ListenerHttp: &saasv1alpha1.ListenerHttp{
							Port:                        0,
							RouteConfigName:             "routeconfig",
							CertificateSecretName:       pointer.String("certificate"),
							EnableHttp2:                 pointer.Bool(false),
							AllowHeadersWithUnderscores: pointer.Bool(true),
							MaxConnectionDuration:       util.Metav1DurationPtr(900 * time.Second),
							ProxyProtocol:               pointer.Bool(true),
						},
					},
				}.AsList(),
			},
			want: &marin3rv1alpha1.EnvoyConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
				},
				Spec: marin3rv1alpha1.EnvoyConfigSpec{
					NodeID:        "test",
					Serialization: marin3r_pointer.New(envoy_serializer.YAML),
					EnvoyAPI:      marin3r_pointer.New(envoy.APIv3),
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
                                            client_ip: '%DOWNSTREAM_REMOTE_ADDRESS_WITHOUT_PORT%'
                                            downstream_tls_cipher: '%DOWNSTREAM_TLS_CIPHER%'
                                            downstream_tls_version: '%DOWNSTREAM_TLS_VERSION%'
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
                                transport_socket:
                                  name: envoy.transport_sockets.tls
                                  typed_config:
                                    '@type': type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.DownstreamTlsContext
                                    common_tls_context:
                                      alpn_protocols:
                                      - http/1.1
                                      tls_certificate_sds_secret_configs:
                                      - name: certificate
                                        sds_config:
                                          ads: {}
                                          resource_api_version: V3
                                      tls_params:
                                        tls_minimum_protocol_version: TLSv1_2
                              listener_filters:
                              - name: envoy.filters.listener.tls_inspector
                              - name: envoy.filters.listener.proxy_protocol
                              name: my_listener
                              per_connection_buffer_limit_bytes: 32768
							`),
						}},
						Routes:   []marin3rv1alpha1.EnvoyResource{},
						Runtimes: []marin3rv1alpha1.EnvoyResource{},
						Secrets:  []marin3rv1alpha1.EnvoySecretResource{{Name: "certificate"}},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.key, tt.args.nodeID, tt.args.factory, tt.args.resources...)(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
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
						FilterChains: []*envoy_config_listener_v3.FilterChain{{}},
					},
				},
			},
			want: &marin3rv1alpha1.EnvoyConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "ns"},
				Spec: marin3rv1alpha1.EnvoyConfigSpec{
					NodeID:        "test",
					Serialization: marin3r_pointer.New(envoy_serializer.YAML),
					EnvoyAPI:      marin3r_pointer.New(envoy.APIv3),
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
								filter_chains:
								- {}
								name: listener1
						`)}},
						Runtimes: []marin3rv1alpha1.EnvoyResource{{
							Value: heredoc.Doc(`
								layer:
								  static_layer_0: value
								name: runtime1
						`)}},
						Secrets: []marin3rv1alpha1.EnvoySecretResource{},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newFromProtos(tt.args.key, tt.args.nodeID, tt.args.resources)()
			if (err != nil) != tt.wantErr {
				t.Errorf("newFromProtos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("newFromProtos() = diff %v", diff)
			}
		})
	}
}
