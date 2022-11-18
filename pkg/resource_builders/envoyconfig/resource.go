package envoyconfig

import (
	"fmt"
	"reflect"

	marin3rv1alpha1 "github.com/3scale-ops/marin3r/apis/marin3r/v1alpha1"
	"github.com/3scale-ops/marin3r/pkg/envoy"
	envoy_serializer "github.com/3scale-ops/marin3r/pkg/envoy/serializer"
	envoy_serializer_v3 "github.com/3scale-ops/marin3r/pkg/envoy/serializer/v3"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_service_runtime_v3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/yaml"
)

var generator = envoyDynamicConfigFactory{
	"ListenerHttp_v1":       {ListenerHTTP_v1, &envoy_config_listener_v3.Listener{}},
	"Cluster_v1":            {Cluster_v1, &envoy_config_cluster_v3.Cluster{}},
	"RouteConfiguration_v1": {RouteConfiguration_v1, &envoy_config_route_v3.RouteConfiguration{}},
	"Runtime_v1":            {Runtime_v1, &envoy_service_runtime_v3.Runtime{}},
}

type envoyDynamicConfigDescriptor interface {
	GetGeneratorVersion() string
	GetName() string
	GetRawConfig() []byte
}

type envoyDynamicConfigGeneratorFn func(envoyDynamicConfigDescriptor) (envoy.Resource, error)

type envoyDynamicConfigClass struct {
	Function envoyDynamicConfigGeneratorFn
	Produces envoy.Resource
}

type envoyDynamicConfigFactory map[string]envoyDynamicConfigClass

func (erf envoyDynamicConfigFactory) newResource(functionName string, descriptor envoyDynamicConfigDescriptor) (envoy.Resource, error) {

	class, ok := erf[functionName]
	if !ok {
		return nil, fmt.Errorf("unregistered class %s", functionName)
	}

	if raw := descriptor.GetRawConfig(); raw != nil {

		err := envoy_serializer_v3.JSON{}.Unmarshal(string(raw), class.Produces)
		if err != nil {
			return nil, err
		}

		return class.Produces, nil
	}

	resource, err := class.Function(descriptor)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

func inspect(v *saasv1alpha1.EnvoyDynamicConfig) (string, envoyDynamicConfigDescriptor) {
	val := reflect.Indirect(reflect.ValueOf(v))
	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		if !val.Field(i).IsNil() {
			descriptor, ok := val.Field(i).Interface().(envoyDynamicConfigDescriptor)
			if !ok {
				// this error cannot occur at runtime
				panic("not an EnvoyDynamicConfigDescriptor")
			}
			generatorFnName := field.Name + "_" + descriptor.GetGeneratorVersion()
			return generatorFnName, descriptor
		}
	}

	return "", nil
}

func newFromProtos(key types.NamespacedName, nodeID string, resources []envoy.Resource) func() (*marin3rv1alpha1.EnvoyConfig, error) {

	return func() (*marin3rv1alpha1.EnvoyConfig, error) {

		clusters := []marin3rv1alpha1.EnvoyResource{}
		routes := []marin3rv1alpha1.EnvoyResource{}
		listeners := []marin3rv1alpha1.EnvoyResource{}
		runtimes := []marin3rv1alpha1.EnvoyResource{}

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
				// should never reach this code in runtime
				panic(fmt.Errorf("unknown resource type"))
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
				},
			},
		}, nil

	}
}

func New(key types.NamespacedName, nodeID string, resources ...saasv1alpha1.EnvoyDynamicConfig) func() (*marin3rv1alpha1.EnvoyConfig, error) {

	return func() (*marin3rv1alpha1.EnvoyConfig, error) {
		protos := []envoy.Resource{}

		for _, res := range resources {

			fn, descriptor := inspect(&res)
			proto, err := generator.newResource(fn, descriptor)
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
