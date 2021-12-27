package mappingservice

import (
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	basereconciler_types "github.com/3scale/saas-operator/pkg/basereconciler/types"
	"github.com/3scale/saas-operator/pkg/generators"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/grafanadashboard"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/hpa"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pdb"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pod"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/podmonitor"
	"github.com/3scale/saas-operator/pkg/generators/mappingservice/config"

	"k8s.io/apimachinery/pkg/types"
)

const (
	component string = "mapping-service"
)

// Generator configures the generators for MappingService
type Generator struct {
	generators.BaseOptions
	Spec    saasv1alpha1.MappingServiceSpec
	Options config.Options
}

// NewGenerator returns a new Options struct
func NewGenerator(instance, namespace string, spec saasv1alpha1.MappingServiceSpec) Generator {
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

// HPA returns a basereconciler_types.GeneratorFunction
func (gen *Generator) HPA() basereconciler_types.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return hpa.New(key, gen.GetLabels(), *gen.Spec.HPA)
}

// PDB returns a basereconciler_types.GeneratorFunction
func (gen *Generator) PDB() basereconciler_types.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return pdb.New(key, gen.GetLabels(), gen.Selector().MatchLabels, *gen.Spec.PDB)
}

// PodMonitor returns a basereconciler_types.GeneratorFunction
func (gen *Generator) PodMonitor() basereconciler_types.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return podmonitor.New(key, gen.GetLabels(), gen.Selector().MatchLabels, podmonitor.PodMetricsEndpoint("/metrics", "metrics", 30))
}

// GrafanaDashboard returns a basereconciler_types.GeneratorFunction
func (gen *Generator) GrafanaDashboard() basereconciler_types.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return grafanadashboard.New(key, gen.GetLabels(), *gen.Spec.GrafanaDashboard, "dashboards/mapping-service.json.gtpl")
}

// SecretDefinition returns a basereconciler_types.GeneratorFunction
func (gen *Generator) SecretDefinition() basereconciler_types.GeneratorFunction {
	return pod.GenerateSecretDefinitionFn("mapping-service-system-master-access-token", gen.GetNamespace(), gen.GetLabels(), gen.Options)
}
