package factory

import (
	"github.com/3scale/saas-operator/pkg/resource_builders/envoyconfig/templates"
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_service_runtime_v3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
)

var f = EnvoyDynamicConfigFactory{
	"ListenerHttp_v1":       RegisterTemplate(templates.ListenerHTTP_v1, &envoy_config_listener_v3.Listener{}),
	"Cluster_v1":            RegisterTemplate(templates.Cluster_v1, &envoy_config_cluster_v3.Cluster{}),
	"RouteConfiguration_v1": RegisterTemplate(templates.RouteConfiguration_v1, &envoy_config_route_v3.RouteConfiguration{}),
	"Runtime_v1":            RegisterTemplate(templates.Runtime_v1, &envoy_service_runtime_v3.Runtime{}),
	"RawConfig_v1":          RegisterTemplate(templates.RawConfig_v1, nil),
}

func Default() EnvoyDynamicConfigFactory {
	return f
}
