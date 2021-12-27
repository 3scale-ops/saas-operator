package system

import (
	"strings"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	basereconciler_types "github.com/3scale/saas-operator/pkg/basereconciler/types"
	"github.com/3scale/saas-operator/pkg/generators"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/grafanadashboard"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/hpa"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pdb"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pod"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/podmonitor"
	"github.com/3scale/saas-operator/pkg/generators/system/config"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
)

const (
	component      string = "system"
	app            string = "app"
	sidekiqDefault string = "sidekiq-default"
	sidekiqBilling string = "sidekiq-billing"
	sidekiqLow     string = "sidekiq-low"
	sphinx         string = "sphinx"
)

// Generator configures the generators for System
type Generator struct {
	generators.BaseOptions
	App                  AppGenerator
	SidekiqDefault       SidekiqGenerator
	SidekiqBilling       SidekiqGenerator
	SidekiqLow           SidekiqGenerator
	Sphinx               SphinxGenerator
	GrafanaDashboardSpec saasv1alpha1.GrafanaDashboardSpec
	ConfigFilesSecret    string
	Options              config.Options
}

// GrafanaDashboard returns a basereconciler_types.GeneratorFunction
func (gen *Generator) GrafanaDashboard() basereconciler_types.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return grafanadashboard.New(key, gen.GetLabels(), gen.GrafanaDashboardSpec, "dashboards/system.json.gtpl")
}

// DatabaseSecretDefinition returns a basereconciler_types.GeneratorFunction
func (gen *Generator) DatabaseSecretDefinition() basereconciler_types.GeneratorFunction {
	return pod.GenerateSecretDefinitionFn("system-database", gen.GetNamespace(), gen.GetLabels(), gen.Options)
}

// RecaptchaSecretDefinition returns a basereconciler_types.GeneratorFunction
func (gen *Generator) RecaptchaSecretDefinition() basereconciler_types.GeneratorFunction {
	return pod.GenerateSecretDefinitionFn("system-recaptcha", gen.GetNamespace(), gen.GetLabels(), gen.Options)
}

// EventsHookSecretDefinition returns a basereconciler_types.GeneratorFunction
func (gen *Generator) EventsHookSecretDefinition() basereconciler_types.GeneratorFunction {
	return pod.GenerateSecretDefinitionFn("system-events-hook", gen.GetNamespace(), gen.GetLabels(), gen.Options)
}

// SMTPSecretDefinition returns a basereconciler_types.GeneratorFunction
func (gen *Generator) SMTPSecretDefinition() basereconciler_types.GeneratorFunction {
	return pod.GenerateSecretDefinitionFn("system-smtp", gen.GetNamespace(), gen.GetLabels(), gen.Options)
}

// MasterApicastSecretDefinition returns a basereconciler_types.GeneratorFunction
func (gen *Generator) MasterApicastSecretDefinition() basereconciler_types.GeneratorFunction {
	return pod.GenerateSecretDefinitionFn("system-master-apicast", gen.GetNamespace(), gen.GetLabels(), gen.Options)
}

// ZyncSecretDefinition returns a basereconciler_types.GeneratorFunction
func (gen *Generator) ZyncSecretDefinition() basereconciler_types.GeneratorFunction {
	return pod.GenerateSecretDefinitionFn("system-zync", gen.GetNamespace(), gen.GetLabels(), gen.Options)
}

// BackendSecretDefinition returns a basereconciler_types.GeneratorFunction
func (gen *Generator) BackendSecretDefinition() basereconciler_types.GeneratorFunction {
	return pod.GenerateSecretDefinitionFn("system-backend", gen.GetNamespace(), gen.GetLabels(), gen.Options)
}

// MultitenantAssetsSecretDefinition returns a basereconciler_types.GeneratorFunction
func (gen *Generator) MultitenantAssetsSecretDefinition() basereconciler_types.GeneratorFunction {
	return pod.GenerateSecretDefinitionFn("system-multitenant-assets-s3", gen.GetNamespace(), gen.GetLabels(), gen.Options)
}

// AppSecretDefinition returns a basereconciler_types.GeneratorFunction
func (gen *Generator) AppSecretDefinition() basereconciler_types.GeneratorFunction {
	return pod.GenerateSecretDefinitionFn("system-app", gen.GetNamespace(), gen.GetLabels(), gen.Options)
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
				Component:    strings.Join([]string{component, app}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": app,
				},
			},
			Spec:              *spec.App,
			Options:           config.NewOptions(spec),
			ImageSpec:         *spec.Image,
			ConfigFilesSecret: *spec.Config.ConfigFilesSecret,
		},
		SidekiqDefault: SidekiqGenerator{
			BaseOptions: generators.BaseOptions{
				Component:    strings.Join([]string{component, sidekiqDefault}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": sidekiqDefault,
				},
			},
			Spec:              *spec.SidekiqDefault,
			Options:           config.NewOptions(spec),
			ImageSpec:         *spec.Image,
			ConfigFilesSecret: *spec.Config.ConfigFilesSecret,
		},
		SidekiqBilling: SidekiqGenerator{
			BaseOptions: generators.BaseOptions{
				Component:    strings.Join([]string{component, sidekiqBilling}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": sidekiqBilling,
				},
			},
			Spec:              *spec.SidekiqBilling,
			Options:           config.NewOptions(spec),
			ImageSpec:         *spec.Image,
			ConfigFilesSecret: *spec.Config.ConfigFilesSecret,
		},
		SidekiqLow: SidekiqGenerator{
			BaseOptions: generators.BaseOptions{
				Component:    strings.Join([]string{component, sidekiqLow}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": sidekiqLow,
				},
			},
			Spec:              *spec.SidekiqLow,
			Options:           config.NewOptions(spec),
			ImageSpec:         *spec.Image,
			ConfigFilesSecret: *spec.Config.ConfigFilesSecret,
		},
		Sphinx: SphinxGenerator{
			BaseOptions: generators.BaseOptions{
				Component:    strings.Join([]string{component, sphinx}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": sphinx,
				},
			},
			Spec:                 *spec.Sphinx,
			Options:              config.NewSphinxOptions(spec),
			ImageSpec:            *spec.Sphinx.Image,
			DatabasePort:         *spec.Sphinx.Config.Thinking.Port,
			DatabasePath:         *spec.Sphinx.Config.Thinking.DatabasePath,
			DatabaseStorageSize:  *spec.Sphinx.Config.Thinking.DatabaseStorageSize,
			DatabaseStorageClass: spec.Sphinx.Config.Thinking.DatabaseStorageClass,
		},
		GrafanaDashboardSpec: *spec.GrafanaDashboard,
		ConfigFilesSecret:    *spec.Config.ConfigFilesSecret,
		Options:              config.NewOptions(spec),
	}
}

// AppGenerator has methods to generate resources for system-app
type AppGenerator struct {
	generators.BaseOptions
	Spec              saasv1alpha1.SystemAppSpec
	Options           config.Options
	ImageSpec         saasv1alpha1.ImageSpec
	ConfigFilesSecret string
}

// HPA returns a basereconciler_types.GeneratorFunction
func (gen *AppGenerator) HPA() basereconciler_types.GeneratorFunction {
	return hpa.New(gen.Key(), gen.GetLabels(), *gen.Spec.HPA)
}

// PDB returns a basereconciler_types.GeneratorFunction
func (gen *AppGenerator) PDB() basereconciler_types.GeneratorFunction {
	return pdb.New(gen.Key(), gen.GetLabels(), gen.Selector().MatchLabels, *gen.Spec.PDB)
}

// PodMonitor returns a basereconciler_types.GeneratorFunction
func (gen *AppGenerator) PodMonitor() basereconciler_types.GeneratorFunction {
	return podmonitor.New(gen.Key(), gen.GetLabels(), gen.Selector().MatchLabels,
		podmonitor.PodMetricsEndpoint("/metrics", "metrics", 30),
		podmonitor.PodMetricsEndpoint("/yabeda-metrics", "metrics", 30),
	)
}

// SidekiqGenerator has methods to generate resources for system-sidekiq
type SidekiqGenerator struct {
	generators.BaseOptions
	Spec              saasv1alpha1.SystemSidekiqSpec
	Options           config.Options
	ImageSpec         saasv1alpha1.ImageSpec
	ConfigFilesSecret string
}

// HPA returns a basereconciler_types.GeneratorFunction
func (gen *SidekiqGenerator) HPA() basereconciler_types.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return hpa.New(key, gen.GetLabels(), *gen.Spec.HPA)
}

// PDB returns a basereconciler_types.GeneratorFunction
func (gen *SidekiqGenerator) PDB() basereconciler_types.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return pdb.New(key, gen.GetLabels(), gen.Selector().MatchLabels, *gen.Spec.PDB)
}

// PodMonitor returns a basereconciler_types.GeneratorFunction
func (gen *SidekiqGenerator) PodMonitor() basereconciler_types.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return podmonitor.New(key, gen.GetLabels(), gen.Selector().MatchLabels,
		podmonitor.PodMetricsEndpoint("/metrics", "metrics", 30),
	)
}

// SphinxGenerator has methods to generate resources for system-sphinx
type SphinxGenerator struct {
	generators.BaseOptions
	Spec                 saasv1alpha1.SystemSphinxSpec
	Options              config.SphinxOptions
	ImageSpec            saasv1alpha1.ImageSpec
	DatabasePort         int32
	DatabasePath         string
	DatabaseStorageSize  resource.Quantity
	DatabaseStorageClass *string
}
