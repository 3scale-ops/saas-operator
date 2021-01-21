package apicast

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
	apicastStaging    string = "apicast-staging"
	apicastProduction string = "apicast-production"
	apicast           string = "apicast"
)

// Generator configures the generators for Apicast
type Generator struct {
	generators.BaseOptions
	Staging              EnvGenerator
	Production           EnvGenerator
	LoadBalancerSpec     saasv1alpha1.LoadBalancerSpec
	GrafanaDashboardSpec saasv1alpha1.GrafanaDashboardSpec
}

// ApicastDashboard returns a basereconciler.GeneratorFunction
func (gen *Generator) ApicastDashboard() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return grafanadashboard.New(key, gen.GetLabels(), gen.GrafanaDashboardSpec, "dashboards/apicast.json.tpl")
}

// ApicastServicesDashboard returns a basereconciler.GeneratorFunction
func (gen *Generator) ApicastServicesDashboard() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component + "-services", Namespace: gen.Namespace}
	return grafanadashboard.New(key, gen.GetLabels(), gen.GrafanaDashboardSpec, "dashboards/apicast-services.json.tpl")
}

// NewGenerator returns a new Options struct
func NewGenerator(instance, namespace string, spec saasv1alpha1.ApicastSpec) Generator {
	return Generator{
		BaseOptions: generators.BaseOptions{
			Component:    apicast,
			InstanceName: instance,
			Namespace:    namespace,
			Labels: map[string]string{
				"app":                  "3scale-api-management",
				"threescale_component": apicast,
			},
		},
		Staging: EnvGenerator{
			BaseOptions: generators.BaseOptions{
				Component:    apicastStaging,
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         apicastStaging,
					"threescale_component_element": "gateway",
				},
			},
			EnvName: "staging",
			Spec:    spec.Staging,
		},
		Production: EnvGenerator{
			BaseOptions: generators.BaseOptions{
				Component:    apicastProduction,
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         apicastProduction,
					"threescale_component_element": "gateway",
				},
			},
			EnvName: "production",
			Spec:    spec.Production,
		},
		GrafanaDashboardSpec: *spec.GrafanaDashboard,
	}
}

// EnvGenerator has methods to generate resources for an
// Apicast environment
type EnvGenerator struct {
	generators.BaseOptions
	EnvName string
	Spec    saasv1alpha1.ApicastEnvironmentSpec
}

// HPA returns a basereconciler.GeneratorFunction
func (gen *EnvGenerator) HPA() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return hpa.New(key, gen.GetLabels(), *gen.Spec.HPA)
}

// PDB returns a basereconciler.GeneratorFunction
func (gen *EnvGenerator) PDB() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return pdb.New(key, gen.GetLabels(), gen.Selector().MatchLabels, *gen.Spec.PDB)
}

// PodMonitor returns a basereconciler.GeneratorFunction
func (gen *EnvGenerator) PodMonitor() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return podmonitor.New(key, gen.GetLabels(), gen.Selector().MatchLabels,
		podmonitor.PodMetricsEndpoint("/metrics", "metrics", 30),
		podmonitor.PodMetricsEndpoint("/stats/prometheus", "envoy-metrics", 60),
	)
}
