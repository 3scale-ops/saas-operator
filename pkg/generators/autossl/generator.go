package autossl

import (
	"fmt"
	"strings"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators"
	"github.com/3scale/saas-operator/pkg/generators/autossl/config"
	basereconciler_resources "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2/resources"
	"github.com/3scale/saas-operator/pkg/reconcilers/workloads"
	"github.com/3scale/saas-operator/pkg/resource_builders/grafanadashboard"
	"github.com/3scale/saas-operator/pkg/resource_builders/podmonitor"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
)

const (
	component string = "autossl"
)

// Generator configures the generators for AutoSSL
type Generator struct {
	generators.BaseOptionsV2
	Spec    saasv1alpha1.AutoSSLSpec
	Options config.Options
	Canary  *Generator
	Traffic bool
}

// NewGenerator returns a new Options struct
func NewGenerator(instance, namespace string, spec saasv1alpha1.AutoSSLSpec) (Generator, error) {

	generator := Generator{
		BaseOptionsV2: generators.BaseOptionsV2{
			Component:    component,
			InstanceName: instance,
			Namespace:    namespace,
			Labels: map[string]string{
				"app":     component,
				"part-of": "3scale-saas",
			},
		},
		Spec:    spec,
		Options: config.NewOptions(spec),
	}

	if spec.Canary != nil {
		canarySpec, err := spec.ResolveCanarySpec(spec.Canary)
		if err != nil {
			return Generator{}, err
		}
		generator.Canary = &Generator{
			BaseOptionsV2: generators.BaseOptionsV2{
				Component:    strings.Join([]string{component, "canary"}, "-"),
				InstanceName: instance,
				Namespace:    namespace,
				Labels: map[string]string{
					"app":                          "3scale-api-management",
					"threescale_component":         component,
					"threescale_component_element": "canary",
				},
			},
			Spec:    *canarySpec,
			Options: config.NewOptions(*canarySpec),
			Traffic: spec.Canary.SendTraffic,
		}
		// Disable PDB and HPA for the canary Deployment
		generator.Canary.Spec.HPA = &saasv1alpha1.HorizontalPodAutoscalerSpec{}
		generator.Canary.Spec.PDB = &saasv1alpha1.PodDisruptionBudgetSpec{}
	}

	return generator, nil
}

// Validate that Generator implements workloads.TrafficManager interface
var _ workloads.TrafficManager = &Generator{}

func (gen *Generator) Services() []basereconciler_resources.ServiceTemplate {
	return []basereconciler_resources.ServiceTemplate{
		{Template: gen.service(), IsEnabled: true},
	}
}
func (gen *Generator) SendTraffic() bool { return gen.Traffic }
func (gen *Generator) TrafficSelector() map[string]string {
	return map[string]string{
		fmt.Sprintf("%s/traffic", saasv1alpha1.GroupVersion.Group): component,
	}
}

// Validate that Generator implements workloads.DeploymentWorkload interface
var _ workloads.DeploymentWorkloadWithTraffic = &Generator{}

func (gen *Generator) Deployment() basereconciler_resources.DeploymentTemplate {
	return basereconciler_resources.DeploymentTemplate{
		Template:        gen.deployment(),
		EnforceReplicas: gen.Spec.HPA.IsDeactivated(),
		IsEnabled:       true,
	}
}

func (gen *Generator) HPASpec() *saasv1alpha1.HorizontalPodAutoscalerSpec {
	return gen.Spec.HPA
}

func (gen *Generator) PDBSpec() *saasv1alpha1.PodDisruptionBudgetSpec {
	return gen.Spec.PDB
}

func (gen *Generator) MonitoredEndpoints() []monitoringv1.PodMetricsEndpoint {
	return []monitoringv1.PodMetricsEndpoint{
		podmonitor.PodMetricsEndpoint("/metrics", "metrics", 30),
	}
}

// GrafanaDashboard returns a basereconciler_resources.GrafanaDashboardTemplate
func (gen *Generator) GrafanaDashboard() basereconciler_resources.GrafanaDashboardTemplate {
	return basereconciler_resources.GrafanaDashboardTemplate{
		Template:  grafanadashboard.New(gen.GetKey(), gen.GetLabels(), *gen.Spec.GrafanaDashboard, "dashboards/autossl.json.gtpl"),
		IsEnabled: !gen.Spec.GrafanaDashboard.IsDeactivated(),
	}
}
