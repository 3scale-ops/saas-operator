package system

import (
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	"github.com/3scale/saas-operator/pkg/generators"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/grafanadashboard"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/hpa"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pdb"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/podmonitor"
	"k8s.io/apimachinery/pkg/types"
)

const (
	component string = "system"
	app       string = "app"
	sidekiq   string = "sidekiq"
	sphinx    string = "sphinx"
)

// Generator configures the generators for Apicast
type Generator struct {
	generators.BaseOptions
	Config               saasv1alpha1.SystemConfig
	App                  AppGenerator
	Sidekiq              SidekiqGenerator
	Sphinx               SphinxGenerator
	GrafanaDashboardSpec saasv1alpha1.GrafanaDashboardSpec
}

// Dashboard returns a basereconciler.GeneratorFunction
func (gen *Generator) Dashboard() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return grafanadashboard.New(key, gen.GetLabels(), gen.GrafanaDashboardSpec, "dashboards/system.json.tpl")
}

// NewGenerator returns a new Options struct
func NewGenerator(instance, namespace string, spec saasv1alpha1.SystemSpec) Generator {
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
		App: AppGenerator{
			BaseOptions: generators.BaseOptions{
				Component:    app,
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": app,
				},
			},
			Spec:   *spec.App,
			Config: spec.Config,
		},
		Sidekiq: SidekiqGenerator{
			BaseOptions: generators.BaseOptions{
				Component:    component,
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": sidekiq,
				},
			},
			Spec:   *spec.Sidekiq,
			Config: spec.Config,
		},
		Sphinx: SphinxGenerator{
			BaseOptions: generators.BaseOptions{
				Component:    component,
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": sphinx,
				},
			},
			Spec:   *spec.Sphinx,
			Config: spec.Config,
		},
		GrafanaDashboardSpec: *spec.GrafanaDashboard,
	}
}

// AppGenerator has methods to generate resources for system-app
type AppGenerator struct {
	generators.BaseOptions
	Spec   saasv1alpha1.SystemAppSpec
	Config saasv1alpha1.SystemConfig
}

// HPA returns a basereconciler.GeneratorFunction
func (gen *AppGenerator) HPA() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return hpa.New(key, gen.GetLabels(), *gen.Spec.HPA)
}

// PDB returns a basereconciler.GeneratorFunction
func (gen *AppGenerator) PDB() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return pdb.New(key, gen.GetLabels(), gen.Selector().MatchLabels, *gen.Spec.PDB)
}

// PodMonitor returns a basereconciler.GeneratorFunction
func (gen *AppGenerator) PodMonitor() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return podmonitor.New(key, gen.GetLabels(), gen.Selector().MatchLabels,
		podmonitor.PodMetricsEndpoint("/metrics", "metrics", 30),
		podmonitor.PodMetricsEndpoint("/yabeda-metrics", "metrics", 30),
	)
}

// SidekiqGenerator has methods to generate resources for system-sidekiq
type SidekiqGenerator struct {
	generators.BaseOptions
	Spec   saasv1alpha1.SystemSidekiqSpec
	Config saasv1alpha1.SystemConfig
}

// HPA returns a basereconciler.GeneratorFunction
func (gen *SidekiqGenerator) HPA() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return hpa.New(key, gen.GetLabels(), *gen.Spec.HPA)
}

// PDB returns a basereconciler.GeneratorFunction
func (gen *SidekiqGenerator) PDB() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return pdb.New(key, gen.GetLabels(), gen.Selector().MatchLabels, *gen.Spec.PDB)
}

// PodMonitor returns a basereconciler.GeneratorFunction
func (gen *SidekiqGenerator) PodMonitor() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return podmonitor.New(key, gen.GetLabels(), gen.Selector().MatchLabels,
		podmonitor.PodMetricsEndpoint("/metrics", "metrics", 30),
	)
}

// SphinxGenerator has methods to generate resources for system-sphinx
type SphinxGenerator struct {
	generators.BaseOptions
	Spec   saasv1alpha1.SystemSphinxSpec
	Config saasv1alpha1.SystemConfig
}
