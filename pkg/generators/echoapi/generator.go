package echoapi

import (
	"fmt"

	"github.com/3scale-ops/basereconciler/mutators"
	"github.com/3scale-ops/basereconciler/resource"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/generators"
	"github.com/3scale-ops/saas-operator/pkg/generators/echoapi/config"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/podmonitor"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/service"
	deployment_workload "github.com/3scale-ops/saas-operator/pkg/workloads/deployment"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
)

const (
	component string = "echo-api"
)

// Generator configures the generators for EchoAPI
type Generator struct {
	generators.BaseOptionsV2
	Spec    saasv1alpha1.EchoAPISpec
	Traffic bool
}

// Validate that Generator implements deployment_workload.DeploymentWorkload interface
var _ deployment_workload.DeploymentWorkload = &Generator{}

// Validate that Generator implements deployment_workload.WithTWithPublishingStrategiesraffic interface
var _ deployment_workload.WithPublishingStrategies = &Generator{}

// NewGenerator returns a new Options struct
func NewGenerator(instance, namespace string, spec saasv1alpha1.EchoAPISpec) Generator {
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
		Traffic: true,
	}
}

// Resources returns the list of resource templates
func (gen *Generator) Resources() ([]resource.TemplateInterface, error) {
	workload, err := deployment_workload.New(gen, nil)
	if err != nil {
		return nil, err
	}

	return workload, nil
}

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
		podmonitor.PodMetricsEndpoint("/stats/prometheus", "envoy-metrics", 60),
	}
}

func (gen *Generator) SendTraffic() bool { return gen.Traffic }
func (gen *Generator) TrafficSelector() map[string]string {
	return map[string]string{
		fmt.Sprintf("%s/traffic", saasv1alpha1.GroupVersion.Group): component,
	}
}

func (gen *Generator) PublishingStrategies() ([]service.ServiceDescriptor, error) {
	if pss, err := service.MergeWithDefaultPublishingStrategy(config.DefaultPublishingStrategy(), gen.Spec.PublishingStrategies); err != nil {
		return nil, err
	} else {
		return pss, nil
	}
}
