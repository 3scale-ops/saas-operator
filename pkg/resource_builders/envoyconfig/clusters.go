package envoyconfig

import (
	"time"

	"github.com/3scale-ops/marin3r/pkg/envoy"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	envoy_extensions_upstreams_http_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/upstreams/http/v3"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func Cluster_v1(desc envoyDynamicConfigDescriptor) (envoy.Resource, error) {
	opts := desc.(*saasv1alpha1.Cluster)

	cluster := &envoy_config_cluster_v3.Cluster{
		Name:           desc.GetName(),
		ConnectTimeout: durationpb.New(1 * time.Second),
		ClusterDiscoveryType: &envoy_config_cluster_v3.Cluster_Type{
			Type: envoy_config_cluster_v3.Cluster_STRICT_DNS,
		},
		DnsLookupFamily: envoy_config_cluster_v3.Cluster_V4_ONLY,
		LbPolicy:        envoy_config_cluster_v3.Cluster_ROUND_ROBIN,
		LoadAssignment: &envoy_config_endpoint_v3.ClusterLoadAssignment{
			ClusterName: desc.GetName(),
			Endpoints: []*envoy_config_endpoint_v3.LocalityLbEndpoints{
				{
					LbEndpoints: []*envoy_config_endpoint_v3.LbEndpoint{
						{
							HostIdentifier: &envoy_config_endpoint_v3.LbEndpoint_Endpoint{
								Endpoint: &envoy_config_endpoint_v3.Endpoint{
									Address: Address_v1(opts.Host, opts.Port),
								},
							},
						},
					},
				},
			},
		},
	}

	if *opts.IsHttp2 {
		any, err := anypb.New(&envoy_extensions_upstreams_http_v3.HttpProtocolOptions{
			UpstreamProtocolOptions: &envoy_extensions_upstreams_http_v3.HttpProtocolOptions_ExplicitHttpConfig_{
				ExplicitHttpConfig: &envoy_extensions_upstreams_http_v3.HttpProtocolOptions_ExplicitHttpConfig{
					ProtocolConfig: &envoy_extensions_upstreams_http_v3.HttpProtocolOptions_ExplicitHttpConfig_Http2ProtocolOptions{
						Http2ProtocolOptions: &envoy_config_core_v3.Http2ProtocolOptions{
							InitialStreamWindowSize:     wrapperspb.UInt32(65536),   // 64 KiB
							InitialConnectionWindowSize: wrapperspb.UInt32(1048576), // 1 MiB
						},
					},
				},
			},
		})
		if err != nil {
			panic(err)
		}
		cluster.TypedExtensionProtocolOptions = map[string]*anypb.Any{
			"envoy.extensions.upstreams.http.v3.HttpProtocolOptions": any,
		}
	}

	return cluster, nil
}
