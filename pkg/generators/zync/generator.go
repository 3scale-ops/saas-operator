package zync

import (
	"fmt"
	"strings"

	"github.com/3scale-ops/basereconciler/mutators"
	"github.com/3scale-ops/basereconciler/resource"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators"
	"github.com/3scale/saas-operator/pkg/generators/zync/config"
	"github.com/3scale/saas-operator/pkg/reconcilers/workloads"
	"github.com/3scale/saas-operator/pkg/resource_builders/grafanadashboard"
	"github.com/3scale/saas-operator/pkg/resource_builders/pod"
	"github.com/3scale/saas-operator/pkg/resource_builders/podmonitor"
	"github.com/3scale/saas-operator/pkg/util"
	externalsecretsv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	grafanav1alpha1 "github.com/grafana-operator/grafana-operator/v4/api/integreatly/v1alpha1"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

const (
	component string = "zync"
	api       string = "zync"
	que       string = "que"
	console   string = "console"
)

// Generator configures the generators for Zync
type Generator struct {
	generators.BaseOptionsV2
	API                  APIGenerator
	Que                  QueGenerator
	Console              ConsoleGenerator
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
			Spec:    *spec.Console,
			Options: config.NewAPIOptions(spec),
			Enabled: *spec.Console.Enabled,
		},
		GrafanaDashboardSpec: *spec.GrafanaDashboard,
		Config:               spec.Config,
	}
}

// Resources returns functions to generate all Zync's shared resources
func (gen *Generator) Resources() []resource.TemplateInterface {
	return []resource.TemplateInterface{
		// GrafanaDashboard
		resource.NewTemplate[*grafanav1alpha1.GrafanaDashboard](
			grafanadashboard.New(gen.GetKey(), gen.GetLabels(), gen.GrafanaDashboardSpec, "dashboards/zync.json.gtpl")).
			WithEnabled(!gen.GrafanaDashboardSpec.IsDeactivated()),
		// ExternalSecret
		resource.NewTemplate[*externalsecretsv1beta1.ExternalSecret](
			pod.GenerateExternalSecretFn("zync", gen.GetNamespace(),
				*gen.Config.ExternalSecret.SecretStoreRef.Name, *gen.Config.ExternalSecret.SecretStoreRef.Kind,
				*gen.Config.ExternalSecret.RefreshInterval, gen.GetLabels(), gen.API.Options),
		),
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

// Validate that APIGenerator implements workloads.DeploymentWorkload interface
var _ workloads.DeploymentWorkload = &APIGenerator{}

// Validate that APIGenerator implements workloads.WithTraffic interface
var _ workloads.WithTraffic = &APIGenerator{}

func (gen *APIGenerator) Labels() map[string]string {
	return gen.GetLabels()
}
func (gen *APIGenerator) Deployment() *resource.Template[*appsv1.Deployment] {
	return resource.NewTemplateFromObjectFunction(gen.deployment).
		WithMutation(mutators.SetDeploymentReplicas(gen.APISpec.HPA.IsDeactivated())).
		WithMutation(mutators.RolloutTrigger{Name: "zync", SecretName: util.Pointer("zync")}.Add())
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
func (gen *APIGenerator) Services() []*resource.Template[*corev1.Service] {
	return []*resource.Template[*corev1.Service]{
		resource.NewTemplateFromObjectFunction(gen.service).WithMutation(mutators.SetServiceLiveValues()),
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

func (gen *QueGenerator) Deployment() *resource.Template[*appsv1.Deployment] {
	return resource.NewTemplateFromObjectFunction(gen.deployment).
		WithMutation(mutators.SetDeploymentReplicas(gen.QueSpec.HPA.IsDeactivated())).
		WithMutation(mutators.RolloutTrigger{Name: "zync", SecretName: pointer.String("zync")}.Add())
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

// ConsoleGenerator has methods to generate resources for zync-console
type ConsoleGenerator struct {
	generators.BaseOptionsV2
	Image   saasv1alpha1.ImageSpec
	Spec    saasv1alpha1.ZyncRailsConsoleSpec
	Options config.APIOptions
	Enabled bool
}

func (gen *ConsoleGenerator) StatefulSet() *resource.Template[*appsv1.StatefulSet] {
	return resource.NewTemplateFromObjectFunction(gen.statefulset).
		WithEnabled(gen.Enabled).
		WithMutation(mutators.SetDeploymentReplicas(true)).
		WithMutation(mutators.RolloutTrigger{Name: "zync", SecretName: pointer.String("zync")}.Add())
}
