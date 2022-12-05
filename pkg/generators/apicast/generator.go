package apicast

import (
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators"
	"github.com/3scale/saas-operator/pkg/generators/apicast/config"
	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	basereconciler_resources "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2/resources"
	"github.com/3scale/saas-operator/pkg/reconcilers/workloads"
	"github.com/3scale/saas-operator/pkg/resource_builders/grafanadashboard"
	"github.com/3scale/saas-operator/pkg/resource_builders/podmonitor"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
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
			Traffic: *&canarySpec.Production.Canary.SendTraffic,
		}
		// Disable PDB and HPA for the canary Deployment
		generator.CanaryProduction.Spec.HPA = &saasv1alpha1.HorizontalPodAutoscalerSpec{}
		generator.CanaryProduction.Spec.PDB = &saasv1alpha1.PodDisruptionBudgetSpec{}
	}

	return generator, nil
}

// Resources returns a list of basereconciler_v2.Resource
func (gen *Generator) Resources() []basereconciler.Resource {
	return []basereconciler.Resource{
		basereconciler_resources.GrafanaDashboardTemplate{
			Template: grafanadashboard.New(
				types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace},
				gen.GetLabels(), gen.GrafanaDashboardSpec, "dashboards/apicast.json.gtpl"),
			IsEnabled: !gen.GrafanaDashboardSpec.IsDeactivated(),
		},
		basereconciler_resources.GrafanaDashboardTemplate{
			Template: grafanadashboard.New(
				types.NamespacedName{Name: gen.Component + "-services", Namespace: gen.Namespace},
				gen.GetLabels(), gen.GrafanaDashboardSpec, "dashboards/apicast-services.json.gtpl"),
			IsEnabled: !gen.GrafanaDashboardSpec.IsDeactivated(),
		},
	}
}

// EnvGenerator has methods to generate resources for an
// Apicast environment
type EnvGenerator struct {
	generators.BaseOptionsV2
	Spec    saasv1alpha1.ApicastEnvironmentSpec
	Options config.EnvOptions
	Traffic bool
}

// Validate that EnvGenerator implements workloads.DeploymentWorkload interface
var _ workloads.DeploymentWorkload = &EnvGenerator{}

// Validate that EnvGenerator implements workloads.WithTraffic interface
var _ workloads.WithTraffic = &EnvGenerator{}

func (gen *EnvGenerator) Labels() map[string]string {
	return gen.GetLabels()
}
func (gen *EnvGenerator) Deployment() basereconciler_resources.DeploymentTemplate {
	return basereconciler_resources.DeploymentTemplate{
		Template:        gen.deployment(),
		RolloutTriggers: nil,
		EnforceReplicas: gen.Spec.HPA.IsDeactivated(),
		IsEnabled:       true,
	}
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
func (gen *EnvGenerator) Services() []basereconciler_resources.ServiceTemplate {
	return []basereconciler_resources.ServiceTemplate{
		{Template: gen.gatewayService(), IsEnabled: true},
		{Template: gen.mgmtService(), IsEnabled: true},
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
