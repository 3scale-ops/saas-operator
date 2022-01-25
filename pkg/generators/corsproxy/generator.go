package corsproxy

import (
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators"
	"github.com/3scale/saas-operator/pkg/generators/corsproxy/config"
	basereconciler_resources "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2/resources"
	"github.com/3scale/saas-operator/pkg/reconcilers/workloads"
	"github.com/3scale/saas-operator/pkg/resource_builders/grafanadashboard"
	"github.com/3scale/saas-operator/pkg/resource_builders/pod"
	"github.com/3scale/saas-operator/pkg/resource_builders/podmonitor"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/utils/pointer"
)

const (
	component string = "cors-proxy"
)

// Generator configures the generators for CORSProxy
type Generator struct {
	generators.BaseOptionsV2
	Spec    saasv1alpha1.CORSProxySpec
	Options config.Options
	Traffic bool
}

// NewGenerator returns a new Options struct
func NewGenerator(instance, namespace string, spec saasv1alpha1.CORSProxySpec) Generator {
	return Generator{
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
		Template: gen.deployment(),
		RolloutTriggers: []basereconciler_resources.RolloutTrigger{
			{Name: "cors-proxy-system-database", SecretName: pointer.String("cors-proxy-system-database")},
		},
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

func (gen *Generator) GrafanaDashboard() basereconciler_resources.GrafanaDashboardTemplate {
	return basereconciler_resources.GrafanaDashboardTemplate{
		Template:  grafanadashboard.New(gen.GetKey(), gen.GetLabels(), *gen.Spec.GrafanaDashboard, "dashboards/cors-proxy.json.gtpl"),
		IsEnabled: !gen.Spec.GrafanaDashboard.IsDeactivated(),
	}
}

func (gen *Generator) SecretDefinition() basereconciler_resources.SecretDefinitionTemplate {
	return basereconciler_resources.SecretDefinitionTemplate{
		Template:  pod.GenerateSecretDefinitionFn("cors-proxy-system-database", gen.GetNamespace(), gen.GetLabels(), gen.Options),
		IsEnabled: true,
	}
}
