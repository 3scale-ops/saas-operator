package zync

import (
	"strings"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	"github.com/3scale/saas-operator/pkg/generators"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/grafanadashboard"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/hpa"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pdb"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pod"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/podmonitor"
	"github.com/3scale/saas-operator/pkg/generators/zync/config"
	"k8s.io/apimachinery/pkg/types"
)

const (
	component string = "zync"
	api       string = "zync"
	que       string = "que"
)

// Generator configures the generators for Zync
type Generator struct {
	generators.BaseOptions
	API                  APIGenerator
	Que                  QueGenerator
	GrafanaDashboardSpec saasv1alpha1.GrafanaDashboardSpec
	Config               saasv1alpha1.ZyncConfig
}

// NewGenerator returns a new Options struct
func NewGenerator(instance, namespace string, spec saasv1alpha1.ZyncSpec) Generator {
	return Generator{
		BaseOptions: generators.BaseOptions{
			Component:    component,
			InstanceName: instance,
			Namespace:    namespace,
			Labels: map[string]string{
				"app":                  "3scale-api-management",
				"threescale_component": component,
			},
		},
		API: APIGenerator{
			BaseOptions: generators.BaseOptions{
				Component:    api,
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": api,
				},
			},
			APISpec: *spec.API,
			Image:   *spec.Image,
			Options: config.NewAPIOptions(spec),
		},
		Que: QueGenerator{
			BaseOptions: generators.BaseOptions{
				Component:    strings.Join([]string{component, que}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": que,
				},
			},
			QueSpec: *spec.Que,
			Image:   *spec.Image,
			Options: config.NewQueOptions(spec),
		},
		GrafanaDashboardSpec: *spec.GrafanaDashboard,
		Config:               spec.Config,
	}
}

// GrafanaDashboard returns a basereconciler.GeneratorFunction
func (gen *Generator) GrafanaDashboard() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return grafanadashboard.New(key, gen.GetLabels(), gen.GrafanaDashboardSpec, "dashboards/zync.json.tpl")
}

// ZyncSecretDefinition returns a basereconciler.GeneratorFunction
func (gen *Generator) ZyncSecretDefinition() basereconciler.GeneratorFunction {
	return pod.GenerateSecretDefinitionFn("zync", gen.GetNamespace(), gen.GetLabels(), gen.API.Options)
}

// APIGenerator has methods to generate resources for a
// Zync environment
type APIGenerator struct {
	generators.BaseOptions
	Image   saasv1alpha1.ImageSpec
	APISpec saasv1alpha1.APISpec
	Options config.APIOptions
}

// HPA returns a basereconciler.GeneratorFunction
func (gen *APIGenerator) HPA() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return hpa.New(key, gen.GetLabels(), *gen.APISpec.HPA)
}

// PDB returns a basereconciler.GeneratorFunction
func (gen *APIGenerator) PDB() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return pdb.New(key, gen.GetLabels(), gen.Selector().MatchLabels, *gen.APISpec.PDB)
}

// PodMonitor returns a basereconciler.GeneratorFunction
func (gen *APIGenerator) PodMonitor() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return podmonitor.New(key, gen.GetLabels(), gen.Selector().MatchLabels,
		podmonitor.PodMetricsEndpoint("/metrics", "metrics", 30),
	)
}

// QueGenerator has methods to generate resources for a
// Que environment
type QueGenerator struct {
	generators.BaseOptions
	Image   saasv1alpha1.ImageSpec
	QueSpec saasv1alpha1.QueSpec
	Options config.QueOptions
}

// HPA returns a basereconciler.GeneratorFunction
func (gen *QueGenerator) HPA() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return hpa.New(key, gen.GetLabels(), *gen.QueSpec.HPA)
}

// PDB returns a basereconciler.GeneratorFunction
func (gen *QueGenerator) PDB() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return pdb.New(key, gen.GetLabels(), gen.Selector().MatchLabels, *gen.QueSpec.PDB)
}

// PodMonitor returns a basereconciler.GeneratorFunction
func (gen *QueGenerator) PodMonitor() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return podmonitor.New(key, gen.GetLabels(), gen.Selector().MatchLabels,
		podmonitor.PodMetricsEndpoint("/metrics", "metrics", 30),
	)
}
