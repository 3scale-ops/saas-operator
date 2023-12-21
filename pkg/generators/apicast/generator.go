package apicast

import (
	"fmt"

	mutators "github.com/3scale-ops/basereconciler/mutators"
	"github.com/3scale-ops/basereconciler/resource"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/generators"
	"github.com/3scale-ops/saas-operator/pkg/generators/apicast/config"
	descriptor "github.com/3scale-ops/saas-operator/pkg/resource_builders/envoyconfig/descriptor"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/grafanadashboard"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/podmonitor"
	operatorutil "github.com/3scale-ops/saas-operator/pkg/util"
	deployment_workload "github.com/3scale-ops/saas-operator/pkg/workloads/deployment"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	apicastStaging          string = "apicast-staging"
	apicastCanaryStaging    string = "apicast-staging-canary"
	apicastProduction       string = "apicast-production"
	apicastCanaryProduction string = "apicast-production-canary"
	apicast                 string = "apicast"
)

// Generator configures the generators for Apicast
type Generator struct {
	generators.BaseOptionsV2
	Staging              EnvGenerator
	CanaryStaging        *EnvGenerator
	Production           EnvGenerator
	CanaryProduction     *EnvGenerator
	LoadBalancerSpec     saasv1alpha1.LoadBalancerSpec
	GrafanaDashboardSpec saasv1alpha1.GrafanaDashboardSpec
}

// NewGenerator returns a new Options struct
func NewGenerator(instance, namespace string, spec saasv1alpha1.ApicastSpec) (Generator, error) {
	generator := Generator{
		BaseOptionsV2: generators.BaseOptionsV2{
			Component:    apicast,
			InstanceName: instance,
			Namespace:    namespace,
			Labels: map[string]string{
				"app":                  "3scale-api-management",
				"threescale_component": apicast,
			},
		},
		Staging: EnvGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
				Component:    apicastStaging,
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         apicastStaging,
					"threescale_component_element": "gateway",
				},
			},
			Spec:    spec.Staging,
			Options: config.NewEnvOptions(spec.Staging, "staging"),
			Traffic: true,
		},
		Production: EnvGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
				Component:    apicastProduction,
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         apicastProduction,
					"threescale_component_element": "gateway",
				},
			},
			Spec:    spec.Production,
			Options: config.NewEnvOptions(spec.Production, "production"),
			Traffic: true,
		},
		GrafanaDashboardSpec: *spec.GrafanaDashboard,
	}

	if spec.Staging.Canary != nil {
		canarySpec, err := spec.ResolveCanarySpec(spec.Staging.Canary)
		if err != nil {
			return Generator{}, err
		}
		generator.CanaryStaging = &EnvGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
				Component:    apicastCanaryStaging,
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         apicastCanaryStaging,
					"threescale_component_element": "gateway",
				},
			},
			Spec:    canarySpec.Staging,
			Options: config.NewEnvOptions(canarySpec.Staging, "staging"),
			Traffic: canarySpec.Staging.Canary.SendTraffic,
		}
		// Disable PDB and HPA for the canary Deployment
		generator.CanaryStaging.Spec.HPA = &saasv1alpha1.HorizontalPodAutoscalerSpec{}
		generator.CanaryStaging.Spec.PDB = &saasv1alpha1.PodDisruptionBudgetSpec{}
	}

	if spec.Production.Canary != nil {
		canarySpec, err := spec.ResolveCanarySpec(spec.Production.Canary)
		if err != nil {
			return Generator{}, err
		}
		generator.CanaryProduction = &EnvGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
				Component:    apicastCanaryProduction,
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         apicastCanaryProduction,
					"threescale_component_element": "gateway",
				},
			},
			Spec:    canarySpec.Production,
			Options: config.NewEnvOptions(canarySpec.Production, "production"),
			Traffic: canarySpec.Production.Canary.SendTraffic,
		}
		// Disable PDB and HPA for the canary Deployment
		generator.CanaryProduction.Spec.HPA = &saasv1alpha1.HorizontalPodAutoscalerSpec{}
		generator.CanaryProduction.Spec.PDB = &saasv1alpha1.PodDisruptionBudgetSpec{}
	}

	return generator, nil
}

// Resources returns the list of resource templates
func (gen *Generator) Resources() ([]resource.TemplateInterface, error) {
	staging, err := deployment_workload.New(&gen.Staging, gen.CanaryStaging)
	if err != nil {
		return nil, err
	}
	production, err := deployment_workload.New(&gen.Production, gen.CanaryProduction)
	if err != nil {
		return nil, err
	}
	misc := []resource.TemplateInterface{
		resource.NewTemplate(
			grafanadashboard.New(
				types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace},
				gen.GetLabels(), gen.GrafanaDashboardSpec, "dashboards/apicast.json.gtpl")).
			WithEnabled(!gen.GrafanaDashboardSpec.IsDeactivated()),
		resource.NewTemplate(
			grafanadashboard.New(
				types.NamespacedName{Name: gen.Component + "-services", Namespace: gen.Namespace},
				gen.GetLabels(), gen.GrafanaDashboardSpec, "dashboards/apicast-services.json.gtpl")).
			WithEnabled(!gen.GrafanaDashboardSpec.IsDeactivated()),
	}

	return operatorutil.ConcatSlices(staging, production, misc), nil
}

// EnvGenerator has methods to generate resources for an
// Apicast environment
type EnvGenerator struct {
	generators.BaseOptionsV2
	Spec    saasv1alpha1.ApicastEnvironmentSpec
	Options pod.Options
	Traffic bool
}

// Validate that EnvGenerator implements deployment_workload.DeploymentWorkload interface
var _ deployment_workload.DeploymentWorkload = &EnvGenerator{}

// Validate that EnvGenerator implements deployment_workload.WithTraffic interface
var _ deployment_workload.WithTraffic = &EnvGenerator{}

// Validate that EnvGenerator implements deployment_workload.WithEnvoySidecar interface
var _ deployment_workload.WithEnvoySidecar = &EnvGenerator{}

func (gen *EnvGenerator) Labels() map[string]string {
	return gen.GetLabels()
}
func (gen *EnvGenerator) Deployment() *resource.Template[*appsv1.Deployment] {
	return resource.NewTemplateFromObjectFunction(gen.deployment).
		WithMutation(mutators.SetDeploymentReplicas(gen.Spec.HPA.IsDeactivated()))
}

func (gen *EnvGenerator) HPASpec() *saasv1alpha1.HorizontalPodAutoscalerSpec {
	return gen.Spec.HPA
}
func (gen *EnvGenerator) PDBSpec() *saasv1alpha1.PodDisruptionBudgetSpec {
	return gen.Spec.PDB
}
func (gen *EnvGenerator) MonitoredEndpoints() []monitoringv1.PodMetricsEndpoint {
	return []monitoringv1.PodMetricsEndpoint{
		podmonitor.PodMetricsEndpoint("/metrics", "metrics", 30),
		podmonitor.PodMetricsEndpoint("/stats/prometheus", "envoy-metrics", 60),
	}
}
func (gen *EnvGenerator) Services() []*resource.Template[*corev1.Service] {
	return []*resource.Template[*corev1.Service]{
		resource.NewTemplateFromObjectFunction(gen.gatewayService).WithMutation(mutators.SetServiceLiveValues()),
		resource.NewTemplateFromObjectFunction(gen.mgmtService).WithMutation(mutators.SetServiceLiveValues()),
	}
}
func (gen *EnvGenerator) SendTraffic() bool { return gen.Traffic }
func (gen *EnvGenerator) TrafficSelector() map[string]string {
	return map[string]string{
		// This is purposedly hardcoded as the TrafficSelector needs to be the same for all workloads produced
		// by the same generator so traffic can be sent to all of them at the same time
		fmt.Sprintf("%s/traffic", saasv1alpha1.GroupVersion.Group): gen.GetComponent(),
	}
}
func (gen *EnvGenerator) EnvoyDynamicConfigurations() []descriptor.EnvoyDynamicConfigDescriptor {
	return gen.Spec.Marin3r.EnvoyDynamicConfig.AsList()
}
