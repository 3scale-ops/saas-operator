package backend

import (
	"encoding/json"
	"strings"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	"github.com/3scale/saas-operator/pkg/generators"
	"github.com/3scale/saas-operator/pkg/generators/backend/config"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/grafanadashboard"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/hpa"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pdb"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/podmonitor"
	"k8s.io/apimachinery/pkg/types"
)

const (
	component string = "backend"
	listener  string = "listener"
	worker    string = "worker"
	cron      string = "cron"
)

// Generator configures the generators for Backend
type Generator struct {
	generators.BaseOptions
	Listener             ListenerGenerator
	Worker               WorkerGenerator
	Cron                 CronGenerator
	GrafanaDashboardSpec saasv1alpha1.GrafanaDashboardSpec
	Config               saasv1alpha1.BackendConfig
}

// NewGenerator returns a new Options struct
func NewGenerator(instance, namespace string, spec saasv1alpha1.BackendSpec) Generator {
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
		Listener: ListenerGenerator{
			BaseOptions: generators.BaseOptions{
				Component:    strings.Join([]string{component, listener}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": listener,
				},
			},
			ListenerSpec: spec.Listener,
			Image:        *spec.Image,
			Config:       spec.Config,
		},
		Worker: WorkerGenerator{
			BaseOptions: generators.BaseOptions{
				Component:    strings.Join([]string{component, worker}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": worker,
				},
			},
			WorkerSpec: *spec.Worker,
			Image:      *spec.Image,
			Config:     spec.Config,
		},
		Cron: CronGenerator{
			BaseOptions: generators.BaseOptions{
				Component:    strings.Join([]string{component, cron}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": cron,
				},
			},
			CronSpec: *spec.Cron,
			Image:    *spec.Image,
			Config:   spec.Config,
		},
		GrafanaDashboardSpec: *spec.GrafanaDashboard,
	}
}

// GrafanaDashboard returns a basereconciler.GeneratorFunction
func (gen *Generator) GrafanaDashboard() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return grafanadashboard.New(key, gen.GetLabels(), gen.GrafanaDashboardSpec, "dashboards/backend.json.tpl")
}

// SystemEventsHookSecretDefinition returns a basereconciler.GeneratorFunction
func (gen *Generator) SystemEventsHookSecretDefinition() basereconciler.GeneratorFunction {
	serializedConfig, _ := json.Marshal(gen.Config)
	sc := config.SecretDefinitions.LookupSecretConfiguration(config.SystemEventsHookSecretName)
	return sc.GenerateSecretDefinitionFn(gen.GetNamespace(), gen.GetLabels(), "/spec/config", serializedConfig)
}

// InternalAPISecretDefinition returns a basereconciler.GeneratorFunction
func (gen *Generator) InternalAPISecretDefinition() basereconciler.GeneratorFunction {
	serializedConfig, _ := json.Marshal(gen.Config)
	sc := config.SecretDefinitions.LookupSecretConfiguration(config.InternalAPISecretName)
	return sc.GenerateSecretDefinitionFn(gen.GetNamespace(), gen.GetLabels(), "/spec/config", serializedConfig)
}

// ErrorMonitoringSecretDefinition returns a basereconciler.GeneratorFunction
func (gen *Generator) ErrorMonitoringSecretDefinition() basereconciler.GeneratorFunction {
	serializedConfig, _ := json.Marshal(gen.Config)
	sc := config.SecretDefinitions.LookupSecretConfiguration(config.ErrorMonitoringSecretName)
	return sc.GenerateSecretDefinitionFn(gen.GetNamespace(), gen.GetLabels(), "/spec/config", serializedConfig)
}

// ListenerGenerator has methods to generate resources for a
// Backend environment
type ListenerGenerator struct {
	generators.BaseOptions
	Image        saasv1alpha1.ImageSpec
	ListenerSpec saasv1alpha1.ListenerSpec
	Config       saasv1alpha1.BackendConfig
}

// HPA returns a basereconciler.GeneratorFunction
func (gen *ListenerGenerator) HPA() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return hpa.New(key, gen.GetLabels(), *gen.ListenerSpec.HPA)
}

// PDB returns a basereconciler.GeneratorFunction
func (gen *ListenerGenerator) PDB() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return pdb.New(key, gen.GetLabels(), gen.Selector().MatchLabels, *gen.ListenerSpec.PDB)
}

// PodMonitor returns a basereconciler.GeneratorFunction
func (gen *ListenerGenerator) PodMonitor() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return podmonitor.New(key, gen.GetLabels(), gen.Selector().MatchLabels,
		podmonitor.PodMetricsEndpoint("/metrics", "metrics", 30),
		podmonitor.PodMetricsEndpoint("/stats/prometheus", "envoy-metrics", 60),
	)
}

// WorkerGenerator has methods to generate resources for a
// Backend environment
type WorkerGenerator struct {
	generators.BaseOptions
	Image      saasv1alpha1.ImageSpec
	WorkerSpec saasv1alpha1.WorkerSpec
	Config     saasv1alpha1.BackendConfig
}

// HPA returns a basereconciler.GeneratorFunction
func (gen *WorkerGenerator) HPA() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return hpa.New(key, gen.GetLabels(), *gen.WorkerSpec.HPA)
}

// PDB returns a basereconciler.GeneratorFunction
func (gen *WorkerGenerator) PDB() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return pdb.New(key, gen.GetLabels(), gen.Selector().MatchLabels, *gen.WorkerSpec.PDB)
}

// PodMonitor returns a basereconciler.GeneratorFunction
func (gen *WorkerGenerator) PodMonitor() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return podmonitor.New(key, gen.GetLabels(), gen.Selector().MatchLabels,
		podmonitor.PodMetricsEndpoint("/metrics", "metrics", 30),
	)
}

// CronGenerator has methods to generate resources for a
// Backend environment
type CronGenerator struct {
	generators.BaseOptions
	Image    saasv1alpha1.ImageSpec
	CronSpec saasv1alpha1.CronSpec
	Config   saasv1alpha1.BackendConfig
}
