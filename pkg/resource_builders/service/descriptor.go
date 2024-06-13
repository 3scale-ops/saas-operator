package service

import (
	"fmt"
	"reflect"

	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/imdario/mergo"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
)

type ServiceDescriptor struct {
	saasv1alpha1.PublishingStrategy
	PortDef corev1.ServicePort
}

type nullTransformer struct {
}

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

	out := make([]ServiceDescriptor, 0, len(def))

	for _, indesc := range in {
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
		case saasv1alpha1.Marin3rStrategy:
			i.Simple = nil
		}
		out = append(out, i)
	})

	return out, nil
}

// func (svc *ProvidedService) SvcPort() corev1.ServicePort {
// 	if svc.Protocol != corev1.ProtocolTCP {
// 		panic(fmt.Sprintf("unsupported Service protocol %s", svc.Protocol))
// 	}
// 	return corev1.ServicePort{
// 		Name:       svc.Name,
// 		Port:       svc.Port,
// 		TargetPort: svc.TargetPort,
// 		Protocol:   corev1.ProtocolTCP,
// 	}
// }
