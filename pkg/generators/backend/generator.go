package backend

import (
	"fmt"
	"strings"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	basereconciler_types "github.com/3scale/saas-operator/pkg/basereconciler/types"
	"github.com/3scale/saas-operator/pkg/generators"
	"github.com/3scale/saas-operator/pkg/generators/backend/config"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/grafanadashboard"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pod"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/podmonitor"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
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
	CanaryListener       *ListenerGenerator
	Worker               WorkerGenerator
	CanaryWorker         *WorkerGenerator
	Cron                 CronGenerator
	GrafanaDashboardSpec saasv1alpha1.GrafanaDashboardSpec
	Config               saasv1alpha1.BackendConfig
}

// NewGenerator returns a new Options struct
func NewGenerator(instance, namespace string, spec saasv1alpha1.BackendSpec) (Generator, error) {

	generator := Generator{
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
			Options:      config.NewListenerOptions(spec),
			Traffic:      true,
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
			Options:    config.NewWorkerOptions(spec),
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
			Options:  config.NewCronOptions(spec),
		},
		GrafanaDashboardSpec: *spec.GrafanaDashboard,
		Config:               spec.Config,
	}

	if spec.Listener.Canary != nil {
		canarySpec, err := spec.ResolveCanarySpec(spec.Listener.Canary)
		if err != nil {
			return Generator{}, err
		}
		generator.CanaryListener = &ListenerGenerator{
			BaseOptions: generators.BaseOptions{
				Component:    strings.Join([]string{component, "canary", listener}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": "canary-" + listener,
				},
			},
			ListenerSpec: canarySpec.Listener,
			Image:        *canarySpec.Image,
			Options:      config.NewListenerOptions(*canarySpec),
			Traffic:      spec.Listener.Canary.SendTraffic,
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
			BaseOptions: generators.BaseOptions{
				Component:    strings.Join([]string{component, "canary", worker}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": "canary-" + worker,
				},
			},
			WorkerSpec: *canarySpec.Worker,
			Image:      *canarySpec.Image,
			Options:    config.NewWorkerOptions(*canarySpec),
		}
		// Disable PDB and HPA for the canary Deployment
		generator.CanaryWorker.WorkerSpec.HPA = &saasv1alpha1.HorizontalPodAutoscalerSpec{}
		generator.CanaryWorker.WorkerSpec.PDB = &saasv1alpha1.PodDisruptionBudgetSpec{}
	}

	return generator, nil
}

// GrafanaDashboard returns a basereconciler_types.GeneratorFunction
func (gen *Generator) GrafanaDashboard() basereconciler_types.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return grafanadashboard.New(key, gen.GetLabels(), gen.GrafanaDashboardSpec, "dashboards/backend.json.gtpl")
}

// SystemEventsHookSecretDefinition returns a basereconciler_types.GeneratorFunction
func (gen *Generator) SystemEventsHookSecretDefinition() basereconciler_types.GeneratorFunction {
	return pod.GenerateSecretDefinitionFn("backend-system-events-hook", gen.GetNamespace(), gen.GetLabels(), gen.Worker.Options)
}

// InternalAPISecretDefinition returns a basereconciler_types.GeneratorFunction
func (gen *Generator) InternalAPISecretDefinition() basereconciler_types.GeneratorFunction {
	return pod.GenerateSecretDefinitionFn("backend-internal-api", gen.GetNamespace(), gen.GetLabels(), gen.Listener.Options)
}

// ErrorMonitoringSecretDefinition returns a basereconciler_types.GeneratorFunction
func (gen *Generator) ErrorMonitoringSecretDefinition() basereconciler_types.GeneratorFunction {
	return pod.GenerateSecretDefinitionFn("backend-error-monitoring", gen.GetNamespace(), gen.GetLabels(), gen.Listener.Options)
}

// ListenerGenerator has methods to generate resources for a
// Backend environment
type ListenerGenerator struct {
	generators.BaseOptions
	Image        saasv1alpha1.ImageSpec
	ListenerSpec saasv1alpha1.ListenerSpec
	Options      config.ListenerOptions
	Traffic      bool
}

// Validate that ListenerGenerator implements basereconciler_types.DeploymentWorkloadGenerator interface
var _ basereconciler_types.DeploymentWorkloadGenerator = &ListenerGenerator{}

// Validate that ListenerGenerator implements basereconciler_types.DeploymentWorkloadGenerator interface
var _ basereconciler_types.DeploymentIngressGenerator = &ListenerGenerator{}

func (gen *ListenerGenerator) HPASpec() *saasv1alpha1.HorizontalPodAutoscalerSpec {
	return gen.ListenerSpec.HPA
}
func (gen *ListenerGenerator) PDBSpec() *saasv1alpha1.PodDisruptionBudgetSpec {
	return gen.ListenerSpec.PDB
}
func (gen *ListenerGenerator) RolloutTriggers() []basereconciler_types.GeneratorFunction {
	return []basereconciler_types.GeneratorFunction{
		pod.GenerateSecretDefinitionFn("backend-internal-api", gen.GetNamespace(), gen.GetLabels(), gen.Options),
		pod.GenerateSecretDefinitionFn("backend-error-monitoring", gen.GetNamespace(), gen.GetLabels(), gen.Options),
	}
}
func (gen *ListenerGenerator) MonitoredEndpoints() []monitoringv1.PodMetricsEndpoint {
	return []monitoringv1.PodMetricsEndpoint{
		podmonitor.PodMetricsEndpoint("/metrics", "metrics", 30),
		podmonitor.PodMetricsEndpoint("/stats/prometheus", "envoy-metrics", 60),
	}
}
func (gen *ListenerGenerator) Services() []basereconciler_types.GeneratorFunction {
	return []basereconciler_types.GeneratorFunction{gen.Service(), gen.InternalService()}
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
	generators.BaseOptions
	Image      saasv1alpha1.ImageSpec
	WorkerSpec saasv1alpha1.WorkerSpec
	Options    config.WorkerOptions
}

// Validate that WorkerGenerator implements basereconciler_types.DeploymentWorkloadGenerator interface
var _ basereconciler_types.DeploymentWorkloadGenerator = &WorkerGenerator{}

func (gen *WorkerGenerator) HPASpec() *saasv1alpha1.HorizontalPodAutoscalerSpec {
	return gen.WorkerSpec.HPA
}
func (gen *WorkerGenerator) PDBSpec() *saasv1alpha1.PodDisruptionBudgetSpec {
	return gen.WorkerSpec.PDB
}
func (gen *WorkerGenerator) RolloutTriggers() []basereconciler_types.GeneratorFunction {
	return []basereconciler_types.GeneratorFunction{
		pod.GenerateSecretDefinitionFn("backend-system-events-hook", gen.GetNamespace(), gen.GetLabels(), gen.Options),
		pod.GenerateSecretDefinitionFn("backend-error-monitoring", gen.GetNamespace(), gen.GetLabels(), gen.Options),
	}
}
func (gen *WorkerGenerator) MonitoredEndpoints() []monitoringv1.PodMetricsEndpoint {
	return []monitoringv1.PodMetricsEndpoint{
		podmonitor.PodMetricsEndpoint("/metrics", "metrics", 30),
	}
}
func (gen *WorkerGenerator) SendTraffic() bool { return false }

// CronGenerator has methods to generate resources for a
// Backend environment
type CronGenerator struct {
	generators.BaseOptions
	Image    saasv1alpha1.ImageSpec
	CronSpec saasv1alpha1.CronSpec
	Options  config.CronOptions
}

// Validate that CronGenerator implements basereconciler_types.DeploymentWorkloadGenerator interface
var _ basereconciler_types.DeploymentWorkloadGenerator = &CronGenerator{}

func (gen *CronGenerator) HPASpec() *saasv1alpha1.HorizontalPodAutoscalerSpec {
	return &saasv1alpha1.HorizontalPodAutoscalerSpec{}
}
func (gen *CronGenerator) PDBSpec() *saasv1alpha1.PodDisruptionBudgetSpec {
	return &saasv1alpha1.PodDisruptionBudgetSpec{}
}
func (gen *CronGenerator) MonitoredEndpoints() []monitoringv1.PodMetricsEndpoint     { return nil }
func (gen *CronGenerator) RolloutTriggers() []basereconciler_types.GeneratorFunction { return nil }
func (gen *CronGenerator) SendTraffic() bool                                         { return false }
