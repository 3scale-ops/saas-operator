package envoyconfig

import (
	"fmt"

	marin3rv1alpha1 "github.com/3scale-ops/marin3r/apis/marin3r/v1alpha1"
	"github.com/3scale-ops/marin3r/pkg/envoy"
	envoy_serializer "github.com/3scale-ops/marin3r/pkg/envoy/serializer"
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_service_runtime_v3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/yaml"
)

func New(key types.NamespacedName, nodeID string, resources ...envoy.Resource) (func() *marin3rv1alpha1.EnvoyConfig, error) {

	clusters := []marin3rv1alpha1.EnvoyResource{}
	routes := []marin3rv1alpha1.EnvoyResource{}
	listeners := []marin3rv1alpha1.EnvoyResource{}
	runtimes := []marin3rv1alpha1.EnvoyResource{}

	s := envoy_serializer.NewResourceMarshaller(envoy_serializer.JSON, envoy.APIv3)

	for i := range resources {

		j, err := s.Marshal(resources[i])
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
			// should never reach this code in runtime
			panic(fmt.Errorf("unknown resource type"))
		}

		if err != nil {
			return nil, err
		}

	}

	return func() *marin3rv1alpha1.EnvoyConfig {
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
				},
			},
		}

	}, nil
}
