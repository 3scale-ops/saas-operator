package zync

import (
	"fmt"
	"strings"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators"

	"github.com/3scale/saas-operator/pkg/generators/zync/config"
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
	component string = "zync"
	api       string = "zync"
	que       string = "que"
)

// Generator configures the generators for Zync
type Generator struct {
	generators.BaseOptionsV2
	API                  APIGenerator
	Que                  QueGenerator
	GrafanaDashboardSpec saasv1alpha1.GrafanaDashboardSpec
	Config               saasv1alpha1.ZyncConfig
}

// NewGenerator returns a new Options struct
func NewGenerator(instance, namespace string, spec saasv1alpha1.ZyncSpec) Generator {
	return Generator{
		BaseOptionsV2: generators.BaseOptionsV2{
			Component:    component,
			InstanceName: instance,
			Namespace:    namespace,
			Labels: map[string]string{
				"app":                  "3scale-api-management",
				"threescale_component": component,
			},
		},
		API: APIGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
				Component:    api,
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": api,
				},
			},
			APISpec: *spec.API,
			Image:   *spec.Image,
			Options: config.NewAPIOptions(spec),
			Traffic: true,
		},
		Que: QueGenerator{
			BaseOptionsV2: generators.BaseOptionsV2{
				Component:    strings.Join([]string{component, que}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": que,
				},
			},
			QueSpec: *spec.Que,
			Image:   *spec.Image,
			Options: config.NewQueOptions(spec),
		},
		GrafanaDashboardSpec: *spec.GrafanaDashboard,
		Config:               spec.Config,
	}
}

// Resources returns functions to generate all Zync's shared resources
func (gen *Generator) Resources() []basereconciler.Resource {
	return []basereconciler.Resource{
		// GrafanaDashboard
		basereconciler_resources.GrafanaDashboardTemplate{
			Template:  grafanadashboard.New(gen.GetKey(), gen.GetLabels(), gen.GrafanaDashboardSpec, "dashboards/zync.json.gtpl"),
			IsEnabled: !gen.GrafanaDashboardSpec.IsDeactivated(),
		},
		// ExternalSecret
		basereconciler_resources.ExternalSecretTemplate{
			Template:  pod.GenerateExternalSecretFn("zync", gen.GetNamespace(), *gen.Config.DatabaseDSN.FromVault.SecretStoreRef.Name, *gen.Config.DatabaseDSN.FromVault.SecretStoreRef.Kind, *gen.Config.DatabaseDSN.FromVault.RefreshInterval, gen.GetLabels(), gen.API.Options),
			IsEnabled: true,
		},
	}
}

// APIGenerator has methods to generate resources for a
// Zync environment
type APIGenerator struct {
	generators.BaseOptionsV2
	Image   saasv1alpha1.ImageSpec
	APISpec saasv1alpha1.APISpec
	Options config.APIOptions
	Traffic bool
}

// Validate that APIGenerator implements workloads.DeploymentWorkloadWithTraffic interface
var _ workloads.DeploymentWorkloadWithTraffic = &APIGenerator{}

// Validate that APIGenerator implements workloads.TrafficManager interface
var _ workloads.TrafficManager = &APIGenerator{}

func (gen *APIGenerator) Labels() map[string]string {
	return gen.GetLabels()
}
func (gen *APIGenerator) Deployment() basereconciler_resources.DeploymentTemplate {
	return basereconciler_resources.DeploymentTemplate{
		Template: gen.deployment(),
		RolloutTriggers: func() []basereconciler_resources.RolloutTrigger {
			return []basereconciler_resources.RolloutTrigger{
				{Name: "zync", SecretName: pointer.String("zync")},
			}
		}(),
		EnforceReplicas: gen.APISpec.HPA.IsDeactivated(),
		IsEnabled:       true,
	}
}

func (gen *APIGenerator) HPASpec() *saasv1alpha1.HorizontalPodAutoscalerSpec {
	return gen.APISpec.HPA
}
func (gen *APIGenerator) PDBSpec() *saasv1alpha1.PodDisruptionBudgetSpec {
	return gen.APISpec.PDB
}
func (gen *APIGenerator) MonitoredEndpoints() []monitoringv1.PodMetricsEndpoint {
	return []monitoringv1.PodMetricsEndpoint{
		podmonitor.PodMetricsEndpoint("/metrics", "metrics", 30),
	}
}
func (gen *APIGenerator) Services() []basereconciler_resources.ServiceTemplate {
	return []basereconciler_resources.ServiceTemplate{
		{Template: gen.service(), IsEnabled: true},
	}
}
func (gen *APIGenerator) SendTraffic() bool { return gen.Traffic }
func (gen *APIGenerator) TrafficSelector() map[string]string {
	return map[string]string{
		fmt.Sprintf("%s/traffic", saasv1alpha1.GroupVersion.Group): component,
	}
}

// QueGenerator has methods to generate resources for a
// Que environment
type QueGenerator struct {
	generators.BaseOptionsV2
	Image   saasv1alpha1.ImageSpec
	QueSpec saasv1alpha1.QueSpec
	Options config.QueOptions
}

// Validate that QueGenerator implements workloads.DeploymentWorkload interface
var _ workloads.DeploymentWorkload = &QueGenerator{}

func (gen *QueGenerator) Deployment() basereconciler_resources.DeploymentTemplate {
	return basereconciler_resources.DeploymentTemplate{
		Template: gen.deployment(),
		RolloutTriggers: func() []basereconciler_resources.RolloutTrigger {
			return []basereconciler_resources.RolloutTrigger{
				{Name: "zync", SecretName: pointer.String("zync")},
			}
		}(),
		EnforceReplicas: gen.QueSpec.HPA.IsDeactivated(),
		IsEnabled:       true,
	}
}
func (gen *QueGenerator) HPASpec() *saasv1alpha1.HorizontalPodAutoscalerSpec {
	return gen.QueSpec.HPA
}
func (gen *QueGenerator) PDBSpec() *saasv1alpha1.PodDisruptionBudgetSpec {
	return gen.QueSpec.PDB
}
func (gen *QueGenerator) MonitoredEndpoints() []monitoringv1.PodMetricsEndpoint {
	return []monitoringv1.PodMetricsEndpoint{
		podmonitor.PodMetricsEndpoint("/metrics", "metrics", 30),
	}
}
