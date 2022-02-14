package backend

import (
	"fmt"
	"strings"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators"
	"github.com/3scale/saas-operator/pkg/generators/backend/config"
	"github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	basereconciler_resources "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2/resources"
	"github.com/3scale/saas-operator/pkg/reconcilers/workloads"
	"github.com/3scale/saas-operator/pkg/resource_builders/grafanadashboard"
	"github.com/3scale/saas-operator/pkg/resource_builders/pod"
	"github.com/3scale/saas-operator/pkg/resource_builders/podmonitor"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/utils/pointer"
)

const (
	component string = "backend"
	listener  string = "listener"
	worker    string = "worker"
	cron      string = "cron"
	twemproxy string = "twemproxy"
)

// Generator configures the generators for Backend
type Generator struct {
	generators.BaseOptionsV2
	Listener             ListenerGenerator
	CanaryListener       *ListenerGenerator
	Worker               WorkerGenerator
	CanaryWorker         *WorkerGenerator
	Cron                 CronGenerator
	grafanaDashboardSpec saasv1alpha1.GrafanaDashboardSpec
	config               saasv1alpha1.BackendConfig
}

// NewGenerator returns a new Options struct
func NewGenerator(instance, namespace string, spec saasv1alpha1.BackendSpec) (Generator, error) {

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
		Listener: ListenerGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
				Component:    strings.Join([]string{component, listener}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": listener,
				},
			},
			ListenerSpec:  spec.Listener,
			Image:         *spec.Image,
			Options:       config.NewListenerOptions(spec),
			Traffic:       true,
			TwemproxySpec: spec.Twemproxy,
		},
		Worker: WorkerGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
				Component:    strings.Join([]string{component, worker}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": worker,
				},
			},
			WorkerSpec:    *spec.Worker,
			Image:         *spec.Image,
			Options:       config.NewWorkerOptions(spec),
			TwemproxySpec: spec.Twemproxy,
		},
		Cron: CronGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
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
			Options:  config.NewCronOptions(spec),
		},
		grafanaDashboardSpec: *spec.GrafanaDashboard,
		config:               spec.Config,
	}

	if spec.Listener.Canary != nil {
		canarySpec, err := spec.ResolveCanarySpec(spec.Listener.Canary)
		if err != nil {
			return Generator{}, err
		}
		generator.CanaryListener = &ListenerGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
				Component:    strings.Join([]string{component, listener, "canary"}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": "canary-" + listener,
				},
			},
			ListenerSpec:  canarySpec.Listener,
			Image:         *canarySpec.Image,
			Options:       config.NewListenerOptions(*canarySpec),
			Traffic:       spec.Listener.Canary.SendTraffic,
			TwemproxySpec: canarySpec.Twemproxy,
		}
		// Disable PDB and HPA for the canary Deployment
		generator.CanaryListener.ListenerSpec.HPA = &saasv1alpha1.HorizontalPodAutoscalerSpec{}
		generator.CanaryListener.ListenerSpec.PDB = &saasv1alpha1.PodDisruptionBudgetSpec{}
	}

	if spec.Worker.Canary != nil {
		canarySpec, err := spec.ResolveCanarySpec(spec.Worker.Canary)
		if err != nil {
			return Generator{}, err
		}
		generator.CanaryWorker = &WorkerGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
				Component:    strings.Join([]string{component, worker, "canary"}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": "canary-" + worker,
				},
			},
			WorkerSpec:    *canarySpec.Worker,
			Image:         *canarySpec.Image,
			Options:       config.NewWorkerOptions(*canarySpec),
			TwemproxySpec: canarySpec.Twemproxy,
		}
		// Disable PDB and HPA for the canary Deployment
		generator.CanaryWorker.WorkerSpec.HPA = &saasv1alpha1.HorizontalPodAutoscalerSpec{}
		generator.CanaryWorker.WorkerSpec.PDB = &saasv1alpha1.PodDisruptionBudgetSpec{}
	}

	return generator, nil
}

// Resources returns functions to generate all Backend's shared resources
func (gen *Generator) Resources() []basereconciler.Resource {
	return []basereconciler.Resource{
		// GrafanaDashboard
		basereconciler_resources.GrafanaDashboardTemplate{
			Template:  grafanadashboard.New(gen.GetKey(), gen.GetLabels(), gen.grafanaDashboardSpec, "dashboards/backend.json.gtpl"),
			IsEnabled: !gen.grafanaDashboardSpec.IsDeactivated(),
		},
		// SecretDefinitions
		basereconciler_resources.SecretDefinitionTemplate{
			Template:  pod.GenerateSecretDefinitionFn("backend-system-events-hook", gen.GetNamespace(), gen.GetLabels(), gen.Worker.Options),
			IsEnabled: true,
		},
		basereconciler_resources.SecretDefinitionTemplate{
			Template:  pod.GenerateSecretDefinitionFn("backend-internal-api", gen.GetNamespace(), gen.GetLabels(), gen.Listener.Options),
			IsEnabled: true,
		},
		basereconciler_resources.SecretDefinitionTemplate{
			Template:  pod.GenerateSecretDefinitionFn("backend-error-monitoring", gen.GetNamespace(), gen.GetLabels(), gen.Listener.Options),
			IsEnabled: gen.config.ErrorMonitoringKey != nil,
		},
	}
}

// ListenerGenerator has methods to generate resources for a
// Backend environment
type ListenerGenerator struct {
	generators.BaseOptionsV2
	Image         saasv1alpha1.ImageSpec
	ListenerSpec  saasv1alpha1.ListenerSpec
	Options       config.ListenerOptions
	Traffic       bool
	TwemproxySpec *saasv1alpha1.TwemproxySpec
}

// Validate that ListenerGenerator implements workloads.DeploymentWorkloadWithTraffic interface
var _ workloads.DeploymentWorkloadWithTraffic = &ListenerGenerator{}

// Validate that ListenerGenerator implements workloads.TrafficManager interface
var _ workloads.TrafficManager = &ListenerGenerator{}

func (gen *ListenerGenerator) Labels() map[string]string {
	return gen.GetLabels()
}
func (gen *ListenerGenerator) Deployment() basereconciler_resources.DeploymentTemplate {
	return basereconciler_resources.DeploymentTemplate{
		Template: gen.deployment(),
		RolloutTriggers: func() []basereconciler_resources.RolloutTrigger {
			return []basereconciler_resources.RolloutTrigger{
				{Name: "backend-internal-api", SecretName: pointer.String("backend-internal-api")},
				{Name: "backend-error-monitoring", SecretName: pointer.String("backend-error-monitoring")},
			}
		}(),
		EnforceReplicas: gen.ListenerSpec.HPA.IsDeactivated(),
		IsEnabled:       true,
	}
}

func (gen *ListenerGenerator) HPASpec() *saasv1alpha1.HorizontalPodAutoscalerSpec {
	return gen.ListenerSpec.HPA
}
func (gen *ListenerGenerator) PDBSpec() *saasv1alpha1.PodDisruptionBudgetSpec {
	return gen.ListenerSpec.PDB
}
func (gen *ListenerGenerator) MonitoredEndpoints() []monitoringv1.PodMetricsEndpoint {
	return []monitoringv1.PodMetricsEndpoint{
		podmonitor.PodMetricsEndpoint("/metrics", "metrics", 30),
		podmonitor.PodMetricsEndpoint("/stats/prometheus", "envoy-metrics", 60),
	}
}
func (gen *ListenerGenerator) Services() []basereconciler_resources.ServiceTemplate {
	return []basereconciler_resources.ServiceTemplate{
		{Template: gen.service(), IsEnabled: true},
		{Template: gen.internalService(), IsEnabled: true},
	}
}
func (gen *ListenerGenerator) SendTraffic() bool { return gen.Traffic }
func (gen *ListenerGenerator) TrafficSelector() map[string]string {
	return map[string]string{
		// This is purposedly hardcoded as the TrafficSelector needs to be the same for all workloads produced
		// by the same generator so traffic can be sent to all of them at the same time
		fmt.Sprintf("%s/traffic", saasv1alpha1.GroupVersion.Group): fmt.Sprintf("%s-%s", component, listener),
	}
}

// WorkerGenerator has methods to generate resources for a
// Backend environment
type WorkerGenerator struct {
	generators.BaseOptionsV2
	Image         saasv1alpha1.ImageSpec
	WorkerSpec    saasv1alpha1.WorkerSpec
	Options       config.WorkerOptions
	TwemproxySpec *saasv1alpha1.TwemproxySpec
}

// Validate that WorkerGenerator implements workloads.DeploymentWorkload interface
var _ workloads.DeploymentWorkload = &WorkerGenerator{}

func (gen *WorkerGenerator) Deployment() basereconciler_resources.DeploymentTemplate {
	return basereconciler_resources.DeploymentTemplate{
		Template: gen.deployment(),
		RolloutTriggers: func() []basereconciler_resources.RolloutTrigger {
			return []basereconciler_resources.RolloutTrigger{
				{Name: "backend-system-events-hook", SecretName: pointer.String("backend-system-events-hook")},
				{Name: "backend-error-monitoring", SecretName: pointer.String("backend-error-monitoring")},
			}
		}(),
		EnforceReplicas: gen.WorkerSpec.HPA.IsDeactivated(),
		IsEnabled:       true,
	}
}
func (gen *WorkerGenerator) HPASpec() *saasv1alpha1.HorizontalPodAutoscalerSpec {
	return gen.WorkerSpec.HPA
}
func (gen *WorkerGenerator) PDBSpec() *saasv1alpha1.PodDisruptionBudgetSpec {
	return gen.WorkerSpec.PDB
}
func (gen *WorkerGenerator) MonitoredEndpoints() []monitoringv1.PodMetricsEndpoint {
	return []monitoringv1.PodMetricsEndpoint{
		podmonitor.PodMetricsEndpoint("/metrics", "metrics", 30),
	}
}

// CronGenerator has methods to generate resources for a
// Backend environment
type CronGenerator struct {
	generators.BaseOptionsV2
	Image    saasv1alpha1.ImageSpec
	CronSpec saasv1alpha1.CronSpec
	Options  config.CronOptions
}

// Validate that CronGenerator implements workloads.DeploymentWorkload interface
var _ workloads.DeploymentWorkload = &CronGenerator{}

func (gen *CronGenerator) Deployment() basereconciler_resources.DeploymentTemplate {
	return basereconciler_resources.DeploymentTemplate{
		Template: gen.deployment(),
		RolloutTriggers: []basereconciler_resources.RolloutTrigger{
			{Name: "backend-error-monitoring", SecretName: pointer.String("backend-error-monitoring")},
		},
		EnforceReplicas: true,
		IsEnabled:       true,
	}
}
func (gen *CronGenerator) HPASpec() *saasv1alpha1.HorizontalPodAutoscalerSpec {
	return &saasv1alpha1.HorizontalPodAutoscalerSpec{}
}
func (gen *CronGenerator) PDBSpec() *saasv1alpha1.PodDisruptionBudgetSpec {
	return &saasv1alpha1.PodDisruptionBudgetSpec{}
}
func (gen *CronGenerator) MonitoredEndpoints() []monitoringv1.PodMetricsEndpoint { return nil }
