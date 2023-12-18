package autossl

import (
	"fmt"
	"strings"

	"github.com/3scale-ops/basereconciler/mutators"
	"github.com/3scale-ops/basereconciler/resource"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/generators"
	"github.com/3scale-ops/saas-operator/pkg/generators/autossl/config"
	"github.com/3scale-ops/saas-operator/pkg/reconcilers/workloads"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/grafanadashboard"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/podmonitor"
	grafanav1alpha1 "github.com/grafana-operator/grafana-operator/v4/api/integreatly/v1alpha1"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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

// Validate that Generator implements workloads.DeploymentWorkload interface
var _ workloads.DeploymentWorkload = &Generator{}

// Validate that Generator implements workloads.WithTraffic interface
var _ workloads.WithTraffic = &Generator{}

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
		Traffic: true,
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

func (gen *Generator) Services() []*resource.Template[*corev1.Service] {
	return []*resource.Template[*corev1.Service]{
		resource.NewTemplateFromObjectFunction(gen.service).
			WithMutation(mutators.SetServiceLiveValues()),
	}
}
func (gen *Generator) SendTraffic() bool { return gen.Traffic }
func (gen *Generator) TrafficSelector() map[string]string {
	return map[string]string{
		fmt.Sprintf("%s/traffic", saasv1alpha1.GroupVersion.Group): component,
	}
}

// Validate that Generator implements workloads.DeploymentWorkload interface
var _ workloads.DeploymentWorkload = &Generator{}

func (gen *Generator) Deployment() *resource.Template[*appsv1.Deployment] {
	return resource.NewTemplateFromObjectFunction(gen.deployment).
		WithMutation(mutators.SetDeploymentReplicas(gen.Spec.HPA.IsDeactivated()))
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
func (gen *Generator) GrafanaDashboard() *resource.Template[*grafanav1alpha1.GrafanaDashboard] {
	return resource.NewTemplate[*grafanav1alpha1.GrafanaDashboard](
		grafanadashboard.New(gen.GetKey(), gen.GetLabels(), *gen.Spec.GrafanaDashboard, "dashboards/autossl.json.gtpl")).
		WithEnabled(!gen.Spec.GrafanaDashboard.IsDeactivated())
}
