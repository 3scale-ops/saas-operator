package service

import (
	"fmt"
	"reflect"

	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/util"
	"github.com/imdario/mergo"
	"github.com/samber/lo"
)

type nullTransformer struct{}

func (t *nullTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() != reflect.Struct {
		return func(dst, src reflect.Value) error {
			// DEBUG:
			// fmt.Printf("\nTYPE %s\n", typ.Elem().Kind())
			// spew.Printf("\tSRC: %+v\n", src.Interface())
			// spew.Printf("\tDST: %+v\n", dst.Interface())
			if dst.CanSet() && !src.IsNil() {
				// DEBUG:
				// fmt.Printf(" ---> WRITE\n")
				dst.Set(src)
			}
			return nil
		}
	}
	return nil
}

func MergeWithDefaultPublishingStrategy(defaults []ServiceDescriptor, in saasv1alpha1.PublishingStrategies) ([]ServiceDescriptor, error) {

	out := []ServiceDescriptor{}

	// NOTE: the mode is always set to Merge by the API defaulter
	// if the user leaves it unset
	switch *in.Mode {

	case saasv1alpha1.PublishingStrategiesReconcileModeReplace:
		var merr util.MultiError

		lo.ForEach(in.Endpoints, func(item saasv1alpha1.PublishingStrategy, index int) {
			indesc := ServiceDescriptor{
				PublishingStrategy: item,
			}

			defdesc, found := lo.Find(defaults, func(i ServiceDescriptor) bool {
				return indesc.EndpointName == i.EndpointName
			})

			if found {
				indesc.PortDefinitions = defdesc.PortDefinitions

			} else {
				// If create is not explicitly set, it's an error
				if indesc.Create == nil || !*indesc.Create {
					merr = append(merr, fmt.Errorf("workload has no endpoint named %s, set 'create=true' if you want to add a new endpoint", indesc.EndpointName))
				}
			}
			out = append(out, indesc)
		})

		if merr.ErrorOrNil() != nil {
			return nil, merr
		}

	case saasv1alpha1.PublishingStrategiesReconcileModeMerge:
		out = defaults
		for _, indesc := range in.Endpoints {

			defdesc, index, found := lo.FindIndexOf(defaults, func(i ServiceDescriptor) bool {
				return indesc.EndpointName == i.EndpointName
			})

			if found {
				// merge with the publishing strategy
				if err := mergo.Merge(&defdesc.PublishingStrategy, indesc, mergo.WithOverride, mergo.WithTransformers(&nullTransformer{})); err != nil {
					return nil, err
				}
				out[index] = defdesc

			} else {
				// If create is not explicitly set, it's an error
				if indesc.Create != nil && *indesc.Create {
					defdesc = ServiceDescriptor{PublishingStrategy: indesc}
				} else {
					return nil, fmt.Errorf("workload has no endpoint named %s, set 'create=true' if you want to add a new endpoint", indesc.EndpointName)
				}
				out = append(out, defdesc)
			}
		}

		// cleanup collection after merging
		lo.ForEach(out, func(i ServiceDescriptor, index int) {
			switch i.PublishingStrategy.Strategy {
			case saasv1alpha1.SimpleStrategy:
				i.Marin3rSidecar = nil
			case saasv1alpha1.Marin3rSidecarStrategy:
				i.Simple = nil
			}
			out[index] = i
		})
	}

	return out, nil
}
