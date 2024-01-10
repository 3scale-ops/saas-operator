package templates

import (
	"time"

	"github.com/3scale-ops/marin3r/pkg/envoy"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	envoy_config_accesslog_v3 "github.com/envoyproxy/go-control-plane/envoy/config/accesslog/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_ratelimit_v3 "github.com/envoyproxy/go-control-plane/envoy/config/ratelimit/v3"
	envoy_extensions_access_loggers_file_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/access_loggers/file/v3"
	envoy_extensions_filters_http_ratelimit_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ratelimit/v3"
	envoy_extensions_filters_http_router_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	envoy_extensions_filters_listener_proxy_protocol_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/listener/proxy_protocol/v3"
	envoy_extensions_filters_listener_tls_inspector_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/listener/tls_inspector/v3"
	http_connection_manager_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	envoy_extensions_transport_sockets_tls_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func ListenerHTTP_v1(name string, opts interface{}) (envoy.Resource, error) {
	o := opts.(*saasv1alpha1.ListenerHttp)

	listener := &envoy_config_listener_v3.Listener{
		Name:            name,
		Address:         Address_v1("0.0.0.0", o.Port),
		ListenerFilters: ListenerFilters_v1(o.CertificateSecretName != nil, *o.ProxyProtocol),
		FilterChains: []*envoy_config_listener_v3.FilterChain{{
			Filters: []*envoy_config_listener_v3.Filter{{
				Name: "envoy.filters.network.http_connection_manager",
				ConfigType: &envoy_config_listener_v3.Filter_TypedConfig{
					TypedConfig: func() *anypb.Any {
						any, err := anypb.New(
							&http_connection_manager_v3.HttpConnectionManager{
								AccessLog: AccessLogConfig_v1(name, o.CertificateSecretName != nil),
								CommonHttpProtocolOptions: func() *envoy_config_core_v3.HttpProtocolOptions {
									po := &envoy_config_core_v3.HttpProtocolOptions{
										IdleTimeout: durationpb.New(3600 * time.Second),
									}
									if o.MaxConnectionDuration != nil {
										po.MaxConnectionDuration = durationpb.New(o.MaxConnectionDuration.Duration)
									}
									if o.AllowHeadersWithUnderscores != nil && *o.AllowHeadersWithUnderscores {
										po.HeadersWithUnderscoresAction = envoy_config_core_v3.HttpProtocolOptions_ALLOW
									} else {
										po.HeadersWithUnderscoresAction = envoy_config_core_v3.HttpProtocolOptions_REJECT_REQUEST
									}
									return po
								}(),

								HttpFilters: HttpFilters_v1(o.RateLimitOptions),
								HttpProtocolOptions: func() *envoy_config_core_v3.Http1ProtocolOptions {
									if o.DefaultHostForHttp10 != nil {
										return &envoy_config_core_v3.Http1ProtocolOptions{
											AcceptHttp_10:         true,
											DefaultHostForHttp_10: *o.DefaultHostForHttp10,
										}
									}
									return &envoy_config_core_v3.Http1ProtocolOptions{}
								}(),
								Http2ProtocolOptions: &envoy_config_core_v3.Http2ProtocolOptions{
									MaxConcurrentStreams:        wrapperspb.UInt32(100),
									InitialStreamWindowSize:     wrapperspb.UInt32(65536),   // 64 KiB
									InitialConnectionWindowSize: wrapperspb.UInt32(1048576), // 1 MiB
								},
								RequestTimeout:    durationpb.New(300 * time.Second),
								RouteSpecifier:    RouteConfigFromAds_v1(o.RouteConfigName),
								StatPrefix:        name,
								StreamIdleTimeout: durationpb.New(300 * time.Second),
								UseRemoteAddress:  wrapperspb.Bool(*o.ProxyProtocol),
							})
						if err != nil {
							panic(err)
						}
						return any
					}(),
				},
			}},
		}},
		PerConnectionBufferLimitBytes: wrapperspb.UInt32(32768), // 32 KiB
	}

	// Apply TLS config if this is a HTTPS listener
	if o.CertificateSecretName != nil {
		listener.FilterChains[0].TransportSocket = TransportSocket_v1(*o.CertificateSecretName, *o.EnableHttp2)
	}

	return listener, nil
}

func ListenerFilters_v1(tls, proxyProtocol bool) []*envoy_config_listener_v3.ListenerFilter {
	filters := []*envoy_config_listener_v3.ListenerFilter{}
	if tls {
		filters = append(filters, &envoy_config_listener_v3.ListenerFilter{
			Name: "envoy.filters.listener.tls_inspector",
			ConfigType: &envoy_config_listener_v3.ListenerFilter_TypedConfig{
				TypedConfig: func() *anypb.Any {
					any, err := anypb.New(
						&envoy_extensions_filters_listener_tls_inspector_v3.TlsInspector{},
					)
					if err != nil {
						panic(err)
					}
					return any
				}(),
			},
		})
	}
	if proxyProtocol {
		filters = append(filters, &envoy_config_listener_v3.ListenerFilter{
			Name: "envoy.filters.listener.proxy_protocol",
			ConfigType: &envoy_config_listener_v3.ListenerFilter_TypedConfig{
				TypedConfig: func() *anypb.Any {
					any, err := anypb.New(
						&envoy_extensions_filters_listener_proxy_protocol_v3.ProxyProtocol{},
					)
					if err != nil {
						panic(err)
					}
					return any
				}(),
			},
		})
	}
	return filters
}

func RouteConfigFromAds_v1(name string) *http_connection_manager_v3.HttpConnectionManager_Rds {
	return &http_connection_manager_v3.HttpConnectionManager_Rds{
		Rds: &http_connection_manager_v3.Rds{
			ConfigSource: &envoy_config_core_v3.ConfigSource{
				ConfigSourceSpecifier: &envoy_config_core_v3.ConfigSource_Ads{
					Ads: &envoy_config_core_v3.AggregatedConfigSource{},
				},
				ResourceApiVersion: envoy_config_core_v3.ApiVersion_V3,
			},
			RouteConfigName: name,
		},
	}
}

func HttpFilters_v1(rlOpts *saasv1alpha1.RateLimitOptions) []*http_connection_manager_v3.HttpFilter {

	filters := []*http_connection_manager_v3.HttpFilter{}
	if rlOpts != nil {
		filters = append(filters, HttpFilterRateLimit_v1(rlOpts))
	}
	filters = append(filters, &http_connection_manager_v3.HttpFilter{
		Name: "envoy.filters.http.router",
		ConfigType: &http_connection_manager_v3.HttpFilter_TypedConfig{
			TypedConfig: func() *anypb.Any {
				any, err := anypb.New(
					&envoy_extensions_filters_http_router_v3.Router{},
				)
				if err != nil {
					panic(err)
				}
				return any
			}(),
		},
	})
	return filters
}

func HttpFilterRateLimit_v1(opts *saasv1alpha1.RateLimitOptions) *http_connection_manager_v3.HttpFilter {
	return &http_connection_manager_v3.HttpFilter{
		Name: "envoy.filters.http.ratelimit",
		ConfigType: &http_connection_manager_v3.HttpFilter_TypedConfig{
			TypedConfig: func() *anypb.Any {
				any, err := anypb.New(
					&envoy_extensions_filters_http_ratelimit_v3.RateLimit{
						Domain:          opts.Domain,
						Timeout:         durationpb.New(opts.Timeout.Duration),
						FailureModeDeny: *opts.FailureModeDeny,
						RateLimitService: &envoy_config_ratelimit_v3.RateLimitServiceConfig{
							GrpcService: &envoy_config_core_v3.GrpcService{
								TargetSpecifier: &envoy_config_core_v3.GrpcService_EnvoyGrpc_{
									EnvoyGrpc: &envoy_config_core_v3.GrpcService_EnvoyGrpc{
										ClusterName: opts.RateLimitCluster,
									},
								},
							},
							TransportApiVersion: envoy_config_core_v3.ApiVersion_V3,
						},
					},
				)
				if err != nil {
					panic(err)
				}
				return any
			}(),
		},
	}
}

func AccessLogConfig_v1(name string, tls bool) []*envoy_config_accesslog_v3.AccessLog {
	return []*envoy_config_accesslog_v3.AccessLog{{
		Name: "envoy.access_loggers.file",
		ConfigType: &envoy_config_accesslog_v3.AccessLog_TypedConfig{
			TypedConfig: func() *anypb.Any {
				logfmt := &envoy_extensions_access_loggers_file_v3.FileAccessLog{
					Path: "/dev/stdout",
					AccessLogFormat: &envoy_extensions_access_loggers_file_v3.FileAccessLog_LogFormat{
						LogFormat: &envoy_config_core_v3.SubstitutionFormatString{
							Format: &envoy_config_core_v3.SubstitutionFormatString_JsonFormat{
								JsonFormat: &structpb.Struct{
									Fields: func() map[string]*structpb.Value {
										m := map[string]*structpb.Value{
											"authority":             structpb.NewStringValue("%REQ(:AUTHORITY)%"),
											"bytes_received":        structpb.NewStringValue("%BYTES_RECEIVED%"),
											"bytes_sent":            structpb.NewStringValue("%BYTES_SENT%"),
											"duration":              structpb.NewStringValue("%DURATION%"),
											"method":                structpb.NewStringValue("%REQ(:METHOD)%"),
											"path":                  structpb.NewStringValue("%REQ(X-ENVOY-ORIGINAL-PATH?:PATH)%"),
											"protocol":              structpb.NewStringValue("%PROTOCOL%"),
											"response_code":         structpb.NewStringValue("%RESPONSE_CODE%"),
											"response_code_details": structpb.NewStringValue("%RESPONSE_CODE_DETAILS%"),
											"response_flags":        structpb.NewStringValue("%RESPONSE_FLAGS%"),
											"listener":              structpb.NewStringValue(name),
											"upstream_cluster":      structpb.NewStringValue("%UPSTREAM_CLUSTER%"),
											"upstream_service_time": structpb.NewStringValue("%RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)%"),
											"user_agent":            structpb.NewStringValue("%REQ(USER-AGENT)%"),
											"client_ip":             structpb.NewStringValue("%DOWNSTREAM_REMOTE_ADDRESS_WITHOUT_PORT%"),
										}
										if tls {
											m["downstream_tls_cipher"] = structpb.NewStringValue("%DOWNSTREAM_TLS_CIPHER%")
											m["downstream_tls_version"] = structpb.NewStringValue("%DOWNSTREAM_TLS_VERSION%")
										}
										return m
									}(),
								},
							},
						},
					},
				}

				any, err := anypb.New(logfmt)
				if err != nil {
					panic(err)
				}
				return any
			}(),
		},
	}}
}

func TransportSocket_v1(secretName string, http2 bool) *envoy_config_core_v3.TransportSocket {
	return &envoy_config_core_v3.TransportSocket{
		Name: "envoy.transport_sockets.tls",
		ConfigType: &envoy_config_core_v3.TransportSocket_TypedConfig{
			TypedConfig: func() *anypb.Any {
				any, err := anypb.New(&envoy_extensions_transport_sockets_tls_v3.DownstreamTlsContext{
					CommonTlsContext: &envoy_extensions_transport_sockets_tls_v3.CommonTlsContext{
						TlsCertificateSdsSecretConfigs: []*envoy_extensions_transport_sockets_tls_v3.SdsSecretConfig{
							{
								Name: secretName,
								SdsConfig: &envoy_config_core_v3.ConfigSource{
									ConfigSourceSpecifier: &envoy_config_core_v3.ConfigSource_Ads{
										Ads: &envoy_config_core_v3.AggregatedConfigSource{},
									},
									ResourceApiVersion: envoy_config_core_v3.ApiVersion_V3,
								},
							},
						},
						TlsParams: &envoy_extensions_transport_sockets_tls_v3.TlsParameters{
							TlsMinimumProtocolVersion: envoy_extensions_transport_sockets_tls_v3.TlsParameters_TLSv1_2,
						},
						AlpnProtocols: func() []string {
							if http2 {
								return []string{"h2,http/1.1"}
							} else {
								return []string{"http/1.1"}
							}
						}(),
					},
				})
				if err != nil {
					panic(err)
				}
				return any
			}(),
		},
	}
}
