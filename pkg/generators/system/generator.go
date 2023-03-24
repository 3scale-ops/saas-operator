package system

import (
	"fmt"
	"strings"

	basereconciler "github.com/3scale-ops/basereconciler/reconciler"
	basereconciler_resources "github.com/3scale-ops/basereconciler/resources"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators"
	"github.com/3scale/saas-operator/pkg/generators/system/config"
	"github.com/3scale/saas-operator/pkg/reconcilers/workloads"
	"github.com/3scale/saas-operator/pkg/resource_builders/grafanadashboard"
	"github.com/3scale/saas-operator/pkg/resource_builders/pod"
	"github.com/3scale/saas-operator/pkg/resource_builders/podmonitor"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/pointer"
)

const (
	component      string = "system"
	app            string = "app"
	console        string = "console"
	sidekiq        string = "sidekiq"
	sidekiqDefault string = "sidekiq-default"
	sidekiqBilling string = "sidekiq-billing"
	sidekiqLow     string = "sidekiq-low"
	sphinx         string = "sphinx"
)

// Generator configures the generators for System
type Generator struct {
	generators.BaseOptionsV2
	App                  AppGenerator
	CanaryApp            *AppGenerator
	SidekiqDefault       SidekiqGenerator
	CanarySidekiqDefault *SidekiqGenerator
	SidekiqBilling       SidekiqGenerator
	CanarySidekiqBilling *SidekiqGenerator
	SidekiqLow           SidekiqGenerator
	CanarySidekiqLow     *SidekiqGenerator
	Sphinx               SphinxGenerator
	Console              ConsoleGenerator
	Config               saasv1alpha1.SystemConfig
	GrafanaDashboardSpec saasv1alpha1.GrafanaDashboardSpec
	ConfigFilesSecret    string
	Options              config.Options
}

// NewGenerator returns a new Options struct
func NewGenerator(instance, namespace string, spec saasv1alpha1.SystemSpec) (Generator, error) {

	generator := Generator{
		BaseOptionsV2: generators.BaseOptionsV2{
			Component:    component,
			InstanceName: instance,
			Namespace:    namespace,
			Labels: map[string]string{
				"app":                  "3scale-api-management",
				"threescale_component": component,
			},
		},
		App: AppGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
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
			Image:             *spec.Image,
			ConfigFilesSecret: *spec.Config.ConfigFilesSecret,
			Traffic:           true,
			TwemproxySpec:     spec.Twemproxy,
		},
		SidekiqDefault: SidekiqGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
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
			Image:             *spec.Image,
			ConfigFilesSecret: *spec.Config.ConfigFilesSecret,
			TwemproxySpec:     spec.Twemproxy,
		},
		SidekiqBilling: SidekiqGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
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
			Image:             *spec.Image,
			ConfigFilesSecret: *spec.Config.ConfigFilesSecret,
			TwemproxySpec:     spec.Twemproxy,
		},
		SidekiqLow: SidekiqGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
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
			Image:             *spec.Image,
			ConfigFilesSecret: *spec.Config.ConfigFilesSecret,
			TwemproxySpec:     spec.Twemproxy,
		},
		Sphinx: SphinxGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
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
			Image:                *spec.Sphinx.Image,
			DatabasePort:         *spec.Sphinx.Config.Thinking.Port,
			DatabasePath:         *spec.Sphinx.Config.Thinking.DatabasePath,
			DatabaseStorageSize:  *spec.Sphinx.Config.Thinking.DatabaseStorageSize,
			DatabaseStorageClass: spec.Sphinx.Config.Thinking.DatabaseStorageClass,
		},
		Console: ConsoleGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
				Component:    strings.Join([]string{component, console}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": sphinx,
				},
			},
			Spec:              *spec.Console,
			Options:           config.NewOptions(spec),
			Image:             *spec.Image,
			ConfigFilesSecret: *spec.Config.ConfigFilesSecret,
			Enabled:           *spec.Config.Rails.Console,
		},
		GrafanaDashboardSpec: *spec.GrafanaDashboard,
		Config:               spec.Config,
		ConfigFilesSecret:    *spec.Config.ConfigFilesSecret,
		Options:              config.NewOptions(spec),
	}

	if spec.App.Canary != nil {
		canarySpec, err := spec.ResolveCanarySpec(spec.App.Canary)
		if err != nil {
			return Generator{}, err
		}
		generator.CanaryApp = &AppGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
				Component:    strings.Join([]string{component, app, "canary"}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": app + "-canary",
				},
			},
			Spec:              *canarySpec.App,
			Image:             *canarySpec.Image,
			Options:           config.NewOptions(*canarySpec),
			ConfigFilesSecret: *canarySpec.Config.ConfigFilesSecret,
			Traffic:           spec.App.Canary.SendTraffic,
			TwemproxySpec:     canarySpec.Twemproxy,
		}
		// Disable PDB and HPA for the canary Deployment
		generator.CanaryApp.Spec.HPA = &saasv1alpha1.HorizontalPodAutoscalerSpec{}
		generator.CanaryApp.Spec.PDB = &saasv1alpha1.PodDisruptionBudgetSpec{}
	}

	if spec.SidekiqDefault.Canary != nil {
		canarySpec, err := spec.ResolveCanarySpec(spec.SidekiqDefault.Canary)
		if err != nil {
			return Generator{}, err
		}
		generator.CanarySidekiqDefault = &SidekiqGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
				Component:    strings.Join([]string{component, sidekiqDefault, "canary"}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": sidekiqDefault + "-canary",
				},
			},
			Spec:              *canarySpec.SidekiqDefault,
			Image:             *canarySpec.Image,
			Options:           config.NewOptions(*canarySpec),
			ConfigFilesSecret: *canarySpec.Config.ConfigFilesSecret,
			TwemproxySpec:     canarySpec.Twemproxy,
		}
		// Disable PDB and HPA for the canary Deployment
		generator.CanarySidekiqDefault.Spec.HPA = &saasv1alpha1.HorizontalPodAutoscalerSpec{}
		generator.CanarySidekiqDefault.Spec.PDB = &saasv1alpha1.PodDisruptionBudgetSpec{}
	}

	if spec.SidekiqLow.Canary != nil {
		canarySpec, err := spec.ResolveCanarySpec(spec.SidekiqLow.Canary)
		if err != nil {
			return Generator{}, err
		}
		generator.CanarySidekiqLow = &SidekiqGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
				Component:    strings.Join([]string{component, sidekiqLow, "canary"}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": sidekiqLow + "-canary",
				},
			},
			Spec:              *canarySpec.SidekiqLow,
			Image:             *canarySpec.Image,
			Options:           config.NewOptions(*canarySpec),
			ConfigFilesSecret: *canarySpec.Config.ConfigFilesSecret,
			TwemproxySpec:     canarySpec.Twemproxy,
		}
		// Disable PDB and HPA for the canary Deployment
		generator.CanarySidekiqLow.Spec.HPA = &saasv1alpha1.HorizontalPodAutoscalerSpec{}
		generator.CanarySidekiqLow.Spec.PDB = &saasv1alpha1.PodDisruptionBudgetSpec{}
	}

	if spec.SidekiqBilling.Canary != nil {
		canarySpec, err := spec.ResolveCanarySpec(spec.SidekiqBilling.Canary)
		if err != nil {
			return Generator{}, err
		}
		generator.CanarySidekiqBilling = &SidekiqGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
				Component:    strings.Join([]string{component, sidekiqBilling, "canary"}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": sidekiqBilling + "-canary",
				},
			},
			Spec:              *canarySpec.SidekiqBilling,
			Image:             *canarySpec.Image,
			Options:           config.NewOptions(*canarySpec),
			ConfigFilesSecret: *canarySpec.Config.ConfigFilesSecret,
			TwemproxySpec:     canarySpec.Twemproxy,
		}
		// Disable PDB and HPA for the canary Deployment
		generator.CanarySidekiqBilling.Spec.HPA = &saasv1alpha1.HorizontalPodAutoscalerSpec{}
		generator.CanarySidekiqBilling.Spec.PDB = &saasv1alpha1.PodDisruptionBudgetSpec{}
	}

	return generator, nil
}

// GrafanaDashboard returns a basereconciler.GeneratorFunction
func (gen *Generator) GrafanaDashboard() basereconciler_resources.GrafanaDashboardTemplate {
	return basereconciler_resources.GrafanaDashboardTemplate{
		Template: grafanadashboard.New(
			gen.GetKey(), gen.GetLabels(), gen.GrafanaDashboardSpec, "dashboards/system.json.gtpl",
		),
		IsEnabled: !gen.GrafanaDashboardSpec.IsDeactivated(),
	}
}

func getSystemSecrets() []string {
	return []string{
		"system-app",
		"system-backend",
		"system-database",
		"system-events-hook",
		"system-master-apicast",
		"system-multitenant-assets-s3",
		"system-recaptcha",
		"system-smtp",
		"system-zync",
	}
}

// Resources returns functions to generate all System's external secrets resources
func (gen *Generator) ExternalSecrets() []basereconciler.Resource {

	resources := []basereconciler.Resource{}
	for _, es := range getSystemSecrets() {
		resources = append(
			resources,
			basereconciler_resources.ExternalSecretTemplate{
				Template: pod.GenerateExternalSecretFn(
					es, gen.GetNamespace(), *gen.Config.ExternalSecret.SecretStoreRef.Name, *gen.Config.ExternalSecret.SecretStoreRef.Kind, *gen.Config.ExternalSecret.RefreshInterval, gen.GetLabels(), gen.Options,
				), IsEnabled: true,
			},
		)
	}

	return resources
}

func getSystemSecretsRolloutTriggers(additionalSecrets ...string) []basereconciler_resources.RolloutTrigger {

	triggers := []basereconciler_resources.RolloutTrigger{}

	secrets := append(getSystemSecrets(), additionalSecrets...)

	for _, secret := range secrets {
		triggers = append(
			triggers,
			basereconciler_resources.RolloutTrigger{
				Name:       secret,
				SecretName: pointer.String(secret),
			},
		)
	}

	return triggers
}

// AppGenerator has methods to generate resources for system-app
type AppGenerator struct {
	generators.BaseOptionsV2
	Spec              saasv1alpha1.SystemAppSpec
	Options           config.Options
	Image             saasv1alpha1.ImageSpec
	ConfigFilesSecret string
	Traffic           bool
	TwemproxySpec     *saasv1alpha1.TwemproxySpec
}

// Validate that AppGenerator implements workloads.DeploymentWorkload interface
var _ workloads.DeploymentWorkload = &AppGenerator{}

// Validate that AppGenerator implements workloads.WithTraffic interface
var _ workloads.WithTraffic = &AppGenerator{}

func (gen *AppGenerator) Services() []basereconciler_resources.ServiceTemplate {
	return []basereconciler_resources.ServiceTemplate{
		{Template: gen.service(), IsEnabled: true},
	}
}
func (gen *AppGenerator) SendTraffic() bool { return gen.Traffic }
func (gen *AppGenerator) TrafficSelector() map[string]string {
	return map[string]string{
		fmt.Sprintf("%s/traffic", saasv1alpha1.GroupVersion.Group): fmt.Sprintf("%s-%s", component, app),
	}
}

func (gen *AppGenerator) Deployment() basereconciler_resources.DeploymentTemplate {
	return basereconciler_resources.DeploymentTemplate{
		Template:        gen.deployment(),
		RolloutTriggers: getSystemSecretsRolloutTriggers(gen.ConfigFilesSecret),
		EnforceReplicas: gen.Spec.HPA.IsDeactivated(),
		IsEnabled:       true,
	}
}

func (gen *AppGenerator) HPASpec() *saasv1alpha1.HorizontalPodAutoscalerSpec {
	return gen.Spec.HPA
}

func (gen *AppGenerator) PDBSpec() *saasv1alpha1.PodDisruptionBudgetSpec {
	return gen.Spec.PDB
}

func (gen *AppGenerator) MonitoredEndpoints() []monitoringv1.PodMetricsEndpoint {
	pmes := []monitoringv1.PodMetricsEndpoint{
		podmonitor.PodMetricsEndpoint("/metrics", "metrics", 30),
	}
	if gen.TwemproxySpec != nil {
		pmes = append(pmes, podmonitor.PodMetricsEndpoint("/metrics", "twem-metrics", 30))
	}
	return pmes
}

// Validate that SidekiqGenerator implements workloads.DeploymentWorkloadWithTraffic interface
var _ workloads.DeploymentWorkload = &SidekiqGenerator{}

// SidekiqGenerator has methods to generate resources for system-sidekiq
type SidekiqGenerator struct {
	generators.BaseOptionsV2
	Spec              saasv1alpha1.SystemSidekiqSpec
	Options           config.Options
	Image             saasv1alpha1.ImageSpec
	ConfigFilesSecret string
	TwemproxySpec     *saasv1alpha1.TwemproxySpec
}

func (gen *SidekiqGenerator) Deployment() basereconciler_resources.DeploymentTemplate {
	return basereconciler_resources.DeploymentTemplate{
		Template:        gen.deployment(),
		RolloutTriggers: getSystemSecretsRolloutTriggers(gen.ConfigFilesSecret),
		EnforceReplicas: gen.Spec.HPA.IsDeactivated(),
		IsEnabled:       true,
	}
}

func (gen *SidekiqGenerator) HPASpec() *saasv1alpha1.HorizontalPodAutoscalerSpec {
	return gen.Spec.HPA
}

func (gen *SidekiqGenerator) PDBSpec() *saasv1alpha1.PodDisruptionBudgetSpec {
	return gen.Spec.PDB
}

func (gen *SidekiqGenerator) MonitoredEndpoints() []monitoringv1.PodMetricsEndpoint {
	pmes := []monitoringv1.PodMetricsEndpoint{
		podmonitor.PodMetricsEndpoint("/metrics", "metrics", 30),
	}
	if gen.TwemproxySpec != nil {
		pmes = append(pmes, podmonitor.PodMetricsEndpoint("/metrics", "twem-metrics", 30))
	}
	return pmes
}

// SphinxGenerator has methods to generate resources for system-sphinx
type SphinxGenerator struct {
	generators.BaseOptionsV2
	Spec                 saasv1alpha1.SystemSphinxSpec
	Options              config.SphinxOptions
	Image                saasv1alpha1.ImageSpec
	DatabasePort         int32
	DatabasePath         string
	DatabaseStorageSize  resource.Quantity
	DatabaseStorageClass *string
}

func (gen *SphinxGenerator) StatefulSetWithTraffic() []basereconciler.Resource {
	return []basereconciler.Resource{
		gen.StatefulSet(), gen.Service(),
	}
}

func (gen *SphinxGenerator) StatefulSet() basereconciler_resources.StatefulSetTemplate {
	return basereconciler_resources.StatefulSetTemplate{
		Template: gen.statefulset(),
		RolloutTriggers: []basereconciler_resources.RolloutTrigger{
			{Name: "system-database", SecretName: pointer.String("system-database")},
		},
		IsEnabled: true,
	}
}

func (gen *SphinxGenerator) Service() basereconciler_resources.ServiceTemplate {
	return basereconciler_resources.ServiceTemplate{
		Template:  gen.service(),
		IsEnabled: true,
	}
}

// ConsoleGenerator has methods to generate resources for system-sphinx
type ConsoleGenerator struct {
	generators.BaseOptionsV2
	Spec              saasv1alpha1.SystemRailsConsoleSpec
	Options           config.Options
	Image             saasv1alpha1.ImageSpec
	ConfigFilesSecret string
	Enabled           bool
}

func (gen *ConsoleGenerator) StatefulSet() basereconciler_resources.StatefulSetTemplate {
	return basereconciler_resources.StatefulSetTemplate{
		Template:        gen.statefulset(),
		RolloutTriggers: getSystemSecretsRolloutTriggers(gen.ConfigFilesSecret),
		IsEnabled:       gen.Enabled,
	}
}
