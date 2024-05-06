package mappingservice

import (
	"fmt"

	"github.com/3scale-ops/basereconciler/mutators"
	"github.com/3scale-ops/basereconciler/resource"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/generators"
	"github.com/3scale-ops/saas-operator/pkg/generators/mappingservice/config"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/grafanadashboard"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/podmonitor"
	operatorutil "github.com/3scale-ops/saas-operator/pkg/util"
	deployment_workload "github.com/3scale-ops/saas-operator/pkg/workloads/deployment"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

const (
	component string = "mapping-service"
)

// Generator configures the generators for MappingService
type Generator struct {
	generators.BaseOptionsV2
	Spec    saasv1alpha1.MappingServiceSpec
	Options pod.Options
	Traffic bool
}

// Validate that Generator implements deployment_workload.DeploymentWorkload interface
var _ deployment_workload.DeploymentWorkload = &Generator{}

// Validate that Generator implements deployment_workload.WithTraffic interface
var _ deployment_workload.WithTraffic = &Generator{}

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

// Resources returns the list of resource templates
func (gen *Generator) Resources() ([]resource.TemplateInterface, error) {
	workload, err := deployment_workload.New(gen, nil)
	if err != nil {
		return nil, err
	}

	externalsecrets := gen.Options.GenerateExternalSecrets(gen.GetKey().Namespace, gen.GetLabels(),
		*gen.Spec.Config.ExternalSecret.SecretStoreRef.Name, *gen.Spec.Config.ExternalSecret.SecretStoreRef.Kind,
		*gen.Spec.Config.ExternalSecret.RefreshInterval)

	misc := []resource.TemplateInterface{
		resource.NewTemplate(
			grafanadashboard.New(gen.GetKey(), gen.GetLabels(), *gen.Spec.GrafanaDashboard, "dashboards/mapping-service.json.gtpl")).
			WithEnabled(!gen.Spec.GrafanaDashboard.IsDeactivated()),
	}
	return operatorutil.ConcatSlices(workload, externalsecrets, misc), nil
}

func (gen *Generator) Services() []*resource.Template[*corev1.Service] {
	return []*resource.Template[*corev1.Service]{
		resource.NewTemplateFromObjectFunction(gen.service).WithMutation(mutators.SetServiceLiveValues()),
	}
}
func (gen *Generator) SendTraffic() bool { return gen.Traffic }
func (gen *Generator) TrafficSelector() map[string]string {
	return map[string]string{
		fmt.Sprintf("%s/traffic", saasv1alpha1.GroupVersion.Group): component,
	}
}

// Validate that Generator implements deployment_workload.DeploymentWorkload interface
var _ deployment_workload.DeploymentWorkload = &Generator{}

func (gen *Generator) Deployment() *resource.Template[*appsv1.Deployment] {
	return resource.NewTemplateFromObjectFunction(gen.deployment).
		WithMutation(mutators.SetDeploymentReplicas(gen.Spec.HPA.IsDeactivated())).
		WithMutations(gen.Options.GenerateRolloutTriggers())
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
