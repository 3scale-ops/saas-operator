package envoyconfig

import (
	"testing"
	"time"

	envoy_serializer_v3 "github.com/3scale-ops/marin3r/pkg/envoy/serializer/v3"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/MakeNowJust/heredoc"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/yaml"
)

func TestListenerHTTP_v1(t *testing.T) {
	type args struct {
		opts *saasv1alpha1.ListenerHttp
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Generates https listener",
			args: args{
				opts: &saasv1alpha1.ListenerHttp{
					EnvoyDynamicConfigMeta: saasv1alpha1.EnvoyDynamicConfigMeta{Name: "test"},
					Port:                   8080,
					RouteConfigName:        "my_route",
					CertificateSecretName:  pointer.String("my_certificate"),
					RateLimitOptions: &saasv1alpha1.RateLimitOptions{
						Domain:           "test_domain",
						FailureModeDeny:  pointer.Bool(true),
						Timeout:          metav1.Duration{Duration: 10 * time.Millisecond},
						RateLimitCluster: "ratelimit",
					},
					DefaultHostForHttp10: pointer.String("example.com"),
					EnableHttp2:          pointer.Bool(false),
				},
			},
			want: heredoc.Doc(`
                address:
                  socket_address:
                    address: 0.0.0.0
                    port_value: 8080
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
                              downstream_tls_cipher: '%DOWNSTREAM_TLS_CIPHER%'
                              downstream_tls_version: '%DOWNSTREAM_TLS_VERSION%'
                              duration: '%DURATION%'
                              listener: test
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
                      - name: envoy.filters.http.ratelimit
                        typed_config:
                          '@type': type.googleapis.com/envoy.extensions.filters.http.ratelimit.v3.RateLimit
                          domain: test_domain
                          failure_mode_deny: true
                          rate_limit_service:
                            grpc_service:
                              envoy_grpc:
                                cluster_name: ratelimit
                            transport_api_version: V3
                          timeout: 0.010s
                      - name: envoy.filters.http.router
                      http_protocol_options:
                        accept_http_10: true
                        default_host_for_http_10: example.com
                      http2_protocol_options:
                        initial_connection_window_size: 1048576
                        initial_stream_window_size: 65536
                        max_concurrent_streams: 100
                      rds:
                        config_source:
                          ads: {}
                          resource_api_version: V3
                        route_config_name: my_route
                      request_timeout: 300s
                      stat_prefix: test
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
                        - name: my_certificate
                          sds_config:
                            ads: {}
                            resource_api_version: V3
                        tls_params:
                          tls_minimum_protocol_version: TLSv1_2
                listener_filters:
                - name: envoy.filters.listener.tls_inspector
                - name: envoy.filters.listener.proxy_protocol
                  typed_config:
                    '@type': type.googleapis.com/envoy.extensions.filters.listener.proxy_protocol.v3.ProxyProtocol
                name: test
                per_connection_buffer_limit_bytes: 32768
			`),
		},
		{
			name: "Generates http listener",
			args: args{
				opts: &saasv1alpha1.ListenerHttp{
					EnvoyDynamicConfigMeta: saasv1alpha1.EnvoyDynamicConfigMeta{Name: "test"},
					Port:                   8080,
					RouteConfigName:        "my_route",
					RateLimitOptions: &saasv1alpha1.RateLimitOptions{
						Domain:           "test_domain",
						FailureModeDeny:  pointer.Bool(false),
						Timeout:          metav1.Duration{Duration: 10 * time.Millisecond},
						RateLimitCluster: "ratelimit",
					},
					DefaultHostForHttp10: pointer.String("example.com"),
					EnableHttp2:          pointer.Bool(false),
				},
			},
			want: heredoc.Doc(`
                address:
                  socket_address:
                    address: 0.0.0.0
                    port_value: 8080
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
                              listener: test
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
                      - name: envoy.filters.http.ratelimit
                        typed_config:
                          '@type': type.googleapis.com/envoy.extensions.filters.http.ratelimit.v3.RateLimit
                          domain: test_domain
                          rate_limit_service:
                            grpc_service:
                              envoy_grpc:
                                cluster_name: ratelimit
                            transport_api_version: V3
                          timeout: 0.010s
                      - name: envoy.filters.http.router
                      http_protocol_options:
                        accept_http_10: true
                        default_host_for_http_10: example.com
                      http2_protocol_options:
                        initial_connection_window_size: 1048576
                        initial_stream_window_size: 65536
                        max_concurrent_streams: 100
                      rds:
                        config_source:
                          ads: {}
                          resource_api_version: V3
                        route_config_name: my_route
                      request_timeout: 300s
                      stat_prefix: test
                      stream_idle_timeout: 300s
                      use_remote_address: true
                listener_filters:
                - name: envoy.filters.listener.proxy_protocol
                  typed_config:
                    '@type': type.googleapis.com/envoy.extensions.filters.listener.proxy_protocol.v3.ProxyProtocol
                name: test
                per_connection_buffer_limit_bytes: 32768
			`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := ListenerHTTP_v1(tt.args.opts)
			j, err := envoy_serializer_v3.JSON{}.Marshal(got)
			if err != nil {
				t.Error(err)
			}
			y, err := yaml.JSONToYAML([]byte(j))
			if err != nil {
				t.Error(err)
			}
			if string(y) != tt.want {
				t.Errorf("ListenerHTTP_v1():\n# got:\n%v\n# want:\n%v", string(y), tt.want)
			}
		})
	}
}
