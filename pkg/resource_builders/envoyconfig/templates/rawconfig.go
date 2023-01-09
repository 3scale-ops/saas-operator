package templates

import (
	"github.com/3scale-ops/marin3r/pkg/envoy"
	envoy_serializer_v3 "github.com/3scale-ops/marin3r/pkg/envoy/serializer/v3"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_service_runtime_v3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
)

func RawConfig_v1(name string, opts interface{}) (envoy.Resource, error) {
	o := opts.(*saasv1alpha1.RawConfig)

	switch o.Type {
	case "listener":
		return unmarshal(o.Value.Raw, &envoy_config_listener_v3.Listener{})
	case "routeConfiguration":
		return unmarshal(o.Value.Raw, &envoy_config_route_v3.RouteConfiguration{})
	case "cluster":
		return unmarshal(o.Value.Raw, &envoy_config_cluster_v3.Cluster{})
	case "runtime":
		return unmarshal(o.Value.Raw, &envoy_service_runtime_v3.Runtime{})
	}

	return nil, nil
}

func unmarshal(b []byte, proto envoy.Resource) (envoy.Resource, error) {
	err := envoy_serializer_v3.JSON{}.Unmarshal(string(b), proto)
	if err != nil {
		return nil, err
	}
	return proto, nil
}
