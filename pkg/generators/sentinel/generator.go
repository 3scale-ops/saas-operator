package sentinel

import (
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pdb"
	"github.com/3scale/saas-operator/pkg/generators/sentinel/config"
	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	component string = "redis-sentinel"
)

// Generator configures the generators for Sentinel
type Generator struct {
	generators.BaseOptions
	Spec    saasv1alpha1.SentinelSpec
	Options config.Options
}

// NewGenerator returns a new Options struct
func NewGenerator(instance, namespace string, spec saasv1alpha1.SentinelSpec) Generator {
	return Generator{
		BaseOptions: generators.BaseOptions{
			Component:    component,
			InstanceName: instance,
			Namespace:    namespace,
			Labels: map[string]string{
				"app":     component,
				"part-of": "3scale-saas",
			},
		},
		Spec:    spec,
		Options: config.NewOptions(spec),
	}
}

// PDB returns a basereconciler.GeneratorFunction
func (gen *Generator) PDB() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return pdb.New(key, gen.GetLabels(), gen.Selector().MatchLabels, *gen.Spec.PDB)
}
