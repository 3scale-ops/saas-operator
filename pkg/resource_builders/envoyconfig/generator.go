package envoyconfig

import (
	"fmt"

	marin3rv1alpha1 "github.com/3scale-ops/marin3r/apis/marin3r/v1alpha1"
	"github.com/3scale-ops/marin3r/pkg/envoy"
	envoy_serializer "github.com/3scale-ops/marin3r/pkg/envoy/serializer"
	envoy_serializer_v3 "github.com/3scale-ops/marin3r/pkg/envoy/serializer/v3"
	"github.com/3scale/saas-operator/pkg/resource_builders/envoyconfig/dynamic_config"
	"github.com/3scale/saas-operator/pkg/resource_builders/envoyconfig/templates"
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_service_runtime_v3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/yaml"
)

var generator = dynamic_config.EnvoyDynamicConfigFactory{
	"ListenerHttp_v1":       dynamic_config.RegisterTemplate(templates.ListenerHTTP_v1, &envoy_config_listener_v3.Listener{}),
	"Cluster_v1":            dynamic_config.RegisterTemplate(templates.Cluster_v1, &envoy_config_cluster_v3.Cluster{}),
	"RouteConfiguration_v1": dynamic_config.RegisterTemplate(templates.RouteConfiguration_v1, &envoy_config_route_v3.RouteConfiguration{}),
	"Runtime_v1":            dynamic_config.RegisterTemplate(templates.Runtime_v1, &envoy_service_runtime_v3.Runtime{}),
	"RawConfig_v1":          dynamic_config.RegisterTemplate(templates.RawConfig_v1, nil),
}

func newFromProtos(key types.NamespacedName, nodeID string, resources []envoy.Resource) func() (*marin3rv1alpha1.EnvoyConfig, error) {

	return func() (*marin3rv1alpha1.EnvoyConfig, error) {

		clusters := []marin3rv1alpha1.EnvoyResource{}
		routes := []marin3rv1alpha1.EnvoyResource{}
		listeners := []marin3rv1alpha1.EnvoyResource{}
		runtimes := []marin3rv1alpha1.EnvoyResource{}
		secrets, err := templates.GenerateSecrets(resources)
		if err != nil {
			return nil, err
		}

		for i := range resources {

			j, err := envoy_serializer_v3.JSON{}.Marshal(resources[i])
			if err != nil {
				return nil, err
			}
			y, err := yaml.JSONToYAML([]byte(j))
			if err != nil {
				return nil, err
			}

			switch resources[i].(type) {

			case *envoy_config_cluster_v3.Cluster:
				clusters = append(clusters, marin3rv1alpha1.EnvoyResource{Value: string(y)})

			case *envoy_config_route_v3.RouteConfiguration:
				routes = append(routes, marin3rv1alpha1.EnvoyResource{Value: string(y)})

			case *envoy_config_listener_v3.Listener:
				listeners = append(listeners, marin3rv1alpha1.EnvoyResource{Value: string(y)})

			case *envoy_service_runtime_v3.Runtime:
				runtimes = append(runtimes, marin3rv1alpha1.EnvoyResource{Value: string(y)})

			default:
				return nil, fmt.Errorf("unknown dynamic configuration type")
			}
		}

		return &marin3rv1alpha1.EnvoyConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
			Spec: marin3rv1alpha1.EnvoyConfigSpec{
				EnvoyAPI:      pointer.StringPtr(envoy.APIv3.String()),
				NodeID:        nodeID,
				Serialization: pointer.String(string(envoy_serializer.YAML)),
				EnvoyResources: &marin3rv1alpha1.EnvoyResources{
					Clusters:  clusters,
					Routes:    routes,
					Listeners: listeners,
					Runtimes:  runtimes,
					Secrets:   secrets,
				},
			},
		}, nil

	}
}

func New(key types.NamespacedName, nodeID string, resources ...dynamic_config.EnvoyDynamicConfigDescriptor) func() (*marin3rv1alpha1.EnvoyConfig, error) {

	return func() (*marin3rv1alpha1.EnvoyConfig, error) {
		protos := []envoy.Resource{}

		for _, res := range resources {

			proto, err := generator.NewResource(res)
			if err != nil {
				return nil, err
			}
			protos = append(protos, proto)
		}

		ec, err := newFromProtos(key, nodeID, protos)()
		if err != nil {
			return nil, err
		}

		return ec, nil
	}
}
