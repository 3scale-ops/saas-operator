package factory

import (
	"fmt"
	"reflect"

	"github.com/3scale-ops/marin3r/pkg/envoy"
	descriptor "github.com/3scale/saas-operator/pkg/resource_builders/envoyconfig/descriptor"
)

// EnvoyDynamicConfigClass contains properties to generate specific types
// of Envoy dynamic configurations
type EnvoyDynamicConfigClass struct {
	Function func(name string, opts interface{}) (envoy.Resource, error)
	Produces envoy.Resource
}

func RegisterTemplate(f func(name string, opts interface{}) (envoy.Resource, error), p envoy.Resource) *EnvoyDynamicConfigClass {
	return &EnvoyDynamicConfigClass{
		Function: f,
		Produces: p,
	}
}

// EnvoyDynamicConfigFactory has methods to produce different types of
// Envoy dynamic resources
type EnvoyDynamicConfigFactory map[string]*EnvoyDynamicConfigClass

// GetClass translates from the external saas-operator API to the internal
// EnvoyDynamicConfigClass that can generate the envoy dynamic resource described
// by the external API
func (factory EnvoyDynamicConfigFactory) GetClass(v descriptor.EnvoyDynamicConfigDescriptor) (*EnvoyDynamicConfigClass, error) {
	opts := v.GetOptions()
	name := reflect.TypeOf(opts).Elem().Name() + "_" + v.GetGeneratorVersion()
	class, ok := factory[name]
	if !ok {
		return nil, fmt.Errorf("unregistered function for '%s'", name)
	}

	return class, nil
}

func (factory EnvoyDynamicConfigFactory) NewResource(desc descriptor.EnvoyDynamicConfigDescriptor) (envoy.Resource, error) {

	class, err := factory.GetClass(desc)
	if err != nil {
		return nil, err
	}

	resource, err := class.Function(desc.GetName(), desc.GetOptions())
	if err != nil {
		return nil, err
	}
	return resource, nil
}
