package system

import (
	"fmt"
	"strings"

	"github.com/3scale-ops/basereconciler/mutators"
	"github.com/3scale-ops/basereconciler/resource"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators"
	"github.com/3scale/saas-operator/pkg/generators/system/config"
	"github.com/3scale/saas-operator/pkg/reconcilers/workloads"
	"github.com/3scale/saas-operator/pkg/resource_builders/grafanadashboard"
	"github.com/3scale/saas-operator/pkg/resource_builders/pod"
	"github.com/3scale/saas-operator/pkg/resource_builders/podmonitor"
	"github.com/3scale/saas-operator/pkg/util"
	grafanav1alpha1 "github.com/grafana-operator/grafana-operator/v4/api/integreatly/v1alpha1"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	pipelinev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	res "k8s.io/apimachinery/pkg/api/resource"
)

const (
	component      string = "system"
	app            string = "app"
	console        string = "console"
	sidekiq        string = "sidekiq"
	sidekiqDefault string = "sidekiq-default"
	sidekiqBilling string = "sidekiq-billing"
	sidekiqLow     string = "sidekiq-low"
	searchd        string = "searchd"
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
	Searchd              SearchdGenerator
	Console              ConsoleGenerator
	Config               saasv1alpha1.SystemConfig
	GrafanaDashboardSpec saasv1alpha1.GrafanaDashboardSpec
	ConfigFilesSecret    string
	Options              config.Options
	Tekton               []SystemTektonGenerator
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
		Searchd: SearchdGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
				Component:    strings.Join([]string{component, searchd}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": searchd,
				},
			},
			Enabled:              *spec.Searchd.Enabled,
			Spec:                 *spec.Searchd,
			Image:                *spec.Searchd.Image,
			DatabasePort:         *spec.Searchd.Config.Port,
			DatabasePath:         *spec.Searchd.Config.DatabasePath,
			DatabaseStorageSize:  *spec.Searchd.Config.DatabaseStorageSize,
			DatabaseStorageClass: spec.Searchd.Config.DatabaseStorageClass,
		},
		Console: ConsoleGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
				Component:    strings.Join([]string{component, console}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": console,
				},
			},
			Spec:              *spec.Console,
			Options:           config.NewOptions(spec),
			Image:             *spec.Console.Image,
			ConfigFilesSecret: *spec.Config.ConfigFilesSecret,
			Enabled:           *spec.Config.Rails.Console,
			TwemproxySpec:     spec.Twemproxy,
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

	for _, task := range spec.Tasks {
		generator.Tekton = append(generator.Tekton,
			SystemTektonGenerator{
				BaseOptionsV2: generators.BaseOptionsV2{
					Component:    *task.Name,
					InstanceName: instance,
					Namespace:    namespace,
					Labels: map[string]string{
						"app":                          "3scale-api-management",
						"threescale_component":         component,
						"threescale_component_element": fmt.Sprintf("task-%s", *task.Name),
					},
				},
				Spec:              task,
				Image:             *task.Config.Image,
				Options:           config.NewOptions(spec),
				ConfigFilesSecret: *spec.Config.ConfigFilesSecret,
				TwemproxySpec:     spec.Twemproxy,
				Enabled:           *task.Enabled,
			})
	}

	return generator, nil
}

// GrafanaDashboard returns a basereconciler.GeneratorFunction
func (gen *Generator) GrafanaDashboard() *resource.Template[*grafanav1alpha1.GrafanaDashboard] {
	return resource.NewTemplate[*grafanav1alpha1.GrafanaDashboard](
		grafanadashboard.New(gen.GetKey(), gen.GetLabels(), gen.GrafanaDashboardSpec, "dashboards/system.json.gtpl")).
		WithEnabled(!gen.GrafanaDashboardSpec.IsDeactivated())
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
func (gen *Generator) ExternalSecrets() []resource.TemplateInterface {

	resources := []resource.TemplateInterface{}
	for _, es := range getSystemSecrets() {
		resources = append(
			resources,
			resource.NewTemplate(
				pod.GenerateExternalSecretFn(es, gen.GetNamespace(),
					*gen.Config.ExternalSecret.SecretStoreRef.Name, *gen.Config.ExternalSecret.SecretStoreRef.Kind,
					*gen.Config.ExternalSecret.RefreshInterval, gen.GetLabels(), gen.Options,
				),
			),
		)
	}

	return resources
}

func getSystemSecretsRolloutTriggers(additionalSecrets ...string) []resource.TemplateMutationFunction {
	triggers := []resource.TemplateMutationFunction{}
	secrets := append(getSystemSecrets(), additionalSecrets...)
	for _, secret := range secrets {
		triggers = append(
			triggers,
			mutators.RolloutTrigger{Name: secret, SecretName: util.Pointer(secret)}.Add(),
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

func (gen *AppGenerator) Services() []*resource.Template[*corev1.Service] {
	return []*resource.Template[*corev1.Service]{
		resource.NewTemplateFromObjectFunction(gen.service).
			WithMutation(mutators.SetServiceLiveValues()),
	}
}
func (gen *AppGenerator) SendTraffic() bool { return gen.Traffic }
func (gen *AppGenerator) TrafficSelector() map[string]string {
	return map[string]string{
		fmt.Sprintf("%s/traffic", saasv1alpha1.GroupVersion.Group): fmt.Sprintf("%s-%s", component, app),
	}
}

func (gen *AppGenerator) Deployment() *resource.Template[*appsv1.Deployment] {
	return resource.NewTemplateFromObjectFunction(gen.deployment).
		WithMutations(getSystemSecretsRolloutTriggers(gen.ConfigFilesSecret)).
		WithMutation(mutators.SetDeploymentReplicas(gen.Spec.HPA.IsDeactivated()))
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

func (gen *SidekiqGenerator) Deployment() *resource.Template[*appsv1.Deployment] {
	return resource.NewTemplateFromObjectFunction(gen.deployment).
		WithMutations(getSystemSecretsRolloutTriggers(gen.ConfigFilesSecret)).
		WithMutation(mutators.SetDeploymentReplicas(gen.Spec.HPA.IsDeactivated()))
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

// SearchdGenerator has methods to generate resources for system-Searchd
type SearchdGenerator struct {
	generators.BaseOptionsV2
	Spec                 saasv1alpha1.SystemSearchdSpec
	Image                saasv1alpha1.ImageSpec
	DatabasePort         int32
	DatabasePath         string
	DatabaseStorageSize  res.Quantity
	DatabaseStorageClass *string
	Enabled              bool
}

func (gen *SearchdGenerator) StatefulSetWithTraffic() []resource.TemplateInterface {
	return []resource.TemplateInterface{
		resource.NewTemplateFromObjectFunction[*appsv1.StatefulSet](gen.statefulset).WithEnabled(gen.Enabled),
		resource.NewTemplateFromObjectFunction[*corev1.Service](gen.service).WithEnabled(gen.Enabled).WithMutation(mutators.SetServiceLiveValues()),
	}
}

// ConsoleGenerator has methods to generate resources for system-console
type ConsoleGenerator struct {
	generators.BaseOptionsV2
	Spec              saasv1alpha1.SystemRailsConsoleSpec
	Options           config.Options
	Image             saasv1alpha1.ImageSpec
	ConfigFilesSecret string
	Enabled           bool
	TwemproxySpec     *saasv1alpha1.TwemproxySpec
}

func (gen *ConsoleGenerator) StatefulSet() *resource.Template[*appsv1.StatefulSet] {
	return resource.NewTemplateFromObjectFunction(gen.statefulset).
		WithEnabled(gen.Enabled).
		WithMutations(getSystemSecretsRolloutTriggers(gen.ConfigFilesSecret))
}

// SystemTektonGenerator has methods to generate resources for system tekton tasks
type SystemTektonGenerator struct {
	generators.BaseOptionsV2
	Spec              saasv1alpha1.SystemTektonTaskSpec
	Options           config.Options
	Image             saasv1alpha1.ImageSpec
	ConfigFilesSecret string
	TwemproxySpec     *saasv1alpha1.TwemproxySpec
	Enabled           bool
}

// Resources returns functions to generate all System's tekton pipeline resources
func (gen *Generator) Pipelines() []resource.TemplateInterface {

	resources := []resource.TemplateInterface{}
	// Tekton resources
	for _, tr := range gen.Tekton {
		copy := tr
		resources = append(resources,
			resource.NewTemplateFromObjectFunction[*pipelinev1beta1.Task](copy.task).WithEnabled(copy.Enabled),
			resource.NewTemplateFromObjectFunction[*pipelinev1beta1.Pipeline](copy.pipeline).WithEnabled(copy.Enabled),
		)
	}

	return resources
}
