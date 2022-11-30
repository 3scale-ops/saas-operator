package mappingservice

import (
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators"
	"github.com/3scale/saas-operator/pkg/generators/mappingservice/config"
	basereconciler_resources "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2/resources"
	"github.com/3scale/saas-operator/pkg/reconcilers/workloads"
	"github.com/3scale/saas-operator/pkg/resource_builders/grafanadashboard"
	"github.com/3scale/saas-operator/pkg/resource_builders/pod"
	"github.com/3scale/saas-operator/pkg/resource_builders/podmonitor"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/utils/pointer"
)

const (
	component string = "mapping-service"
)

// Generator configures the generators for MappingService
type Generator struct {
	generators.BaseOptionsV2
	Spec    saasv1alpha1.MappingServiceSpec
	Options config.Options
	Traffic bool
}

// Validate that Generator implements workloads.DeploymentWorkload interface
var _ workloads.DeploymentWorkload = &Generator{}

// Validate that Generator implements workloads.WithTraffic interface
var _ workloads.WithTraffic = &Generator{}

// NewGenerator returns a new Options struct
func NewGenerator(instance, namespace string, spec saasv1alpha1.MappingServiceSpec) Generator {
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
var _ workloads.DeploymentWorkload = &Generator{}

func (gen *Generator) Deployment() basereconciler_resources.DeploymentTemplate {
	return basereconciler_resources.DeploymentTemplate{
		Template: gen.deployment(),
		RolloutTriggers: []basereconciler_resources.RolloutTrigger{
			{
				Name:       "mapping-service-system-master-access-token",
				SecretName: pointer.String("mapping-service-system-master-access-token"),
			},
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
		Template:  grafanadashboard.New(gen.GetKey(), gen.GetLabels(), *gen.Spec.GrafanaDashboard, "dashboards/mapping-service.json.gtpl"),
		IsEnabled: !gen.Spec.GrafanaDashboard.IsDeactivated(),
	}
}

func (gen *Generator) ExternalSecret() basereconciler_resources.ExternalSecretTemplate {
	return basereconciler_resources.ExternalSecretTemplate{
		Template:  pod.GenerateExternalSecretFn("mapping-service-system-master-access-token", gen.GetNamespace(), *gen.Spec.Config.ExternalSecret.SecretStoreRef.Name, *gen.Spec.Config.ExternalSecret.SecretStoreRef.Kind, *gen.Spec.Config.ExternalSecret.RefreshInterval, gen.GetLabels(), gen.Options),
		IsEnabled: true,
	}
}
