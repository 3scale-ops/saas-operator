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

func MergeWithDefaultPublishingStrategy(def []ServiceDescriptor, in saasv1alpha1.PublishingStrategies) ([]ServiceDescriptor, error) {

	var out []ServiceDescriptor

	// default mode is 'Merge'
	mode := saasv1alpha1.PublishingStrategiesReconcileModeMerge
	if in.Mode != nil {
		mode = *in.Mode
	}

	switch mode {

	case saasv1alpha1.PublishingStrategiesReconcileModeReplace:
		var merr util.MultiError

		out = make([]ServiceDescriptor, 0, len(in.Endpoints))
		lo.ForEach(in.Endpoints, func(item saasv1alpha1.PublishingStrategy, index int) {
			indesc := ServiceDescriptor{
				PublishingStrategy: item,
			}
			if defdesc, ok := lo.Find(def, func(i ServiceDescriptor) bool {
				return indesc.EndpointName == i.EndpointName
			}); ok {
				indesc.PortDefinitions = defdesc.PortDefinitions
			} else {
				merr = append(merr, fmt.Errorf("workload has no endpoint named %s", indesc.EndpointName))
			}
			out = append(out, indesc)
		})

		if merr.ErrorOrNil() != nil {
			return nil, merr
		}

	case saasv1alpha1.PublishingStrategiesReconcileModeMerge:

		out = make([]ServiceDescriptor, 0, len(def))
		for _, indesc := range in.Endpoints {
			var defdesc ServiceDescriptor
			var ok bool
			var index int

			if defdesc, index, ok = lo.FindIndexOf(def, func(i ServiceDescriptor) bool {
				return indesc.EndpointName == i.EndpointName
			}); !ok {
				return nil, fmt.Errorf("workload has no endpoint named %s", indesc.EndpointName)
			}

			if err := mergo.Merge(&defdesc.PublishingStrategy, indesc, mergo.WithOverride, mergo.WithTransformers(&nullTransformer{})); err != nil {
				return nil, err
			}
			def[index] = defdesc
		}

		// cleanup collection after merging
		lo.ForEach(def, func(i ServiceDescriptor, _ int) {
			switch i.PublishingStrategy.Strategy {
			case saasv1alpha1.SimpleStrategy:
				i.Marin3rSidecar = nil
			case saasv1alpha1.Marin3rSidecarStrategy:
				i.Simple = nil
			}
			out = append(out, i)
		})
	}

	return out, nil
}
