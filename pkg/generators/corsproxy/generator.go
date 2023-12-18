package corsproxy

import (
	"fmt"

	"github.com/3scale-ops/basereconciler/mutators"
	"github.com/3scale-ops/basereconciler/resource"
	"github.com/3scale-ops/basereconciler/util"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/generators"
	"github.com/3scale-ops/saas-operator/pkg/generators/corsproxy/config"
	"github.com/3scale-ops/saas-operator/pkg/reconcilers/workloads"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/grafanadashboard"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/podmonitor"
	externalsecretsv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	grafanav1alpha1 "github.com/grafana-operator/grafana-operator/v4/api/integreatly/v1alpha1"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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

// Validate that Generator implements workloads.DeploymentWorkload interface
var _ workloads.DeploymentWorkload = &Generator{}

// Validate that Generator implements workloads.WithTraffic interface
var _ workloads.WithTraffic = &Generator{}

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

// Validate that Generator implements workloads.DeploymentWorkload interface
var _ workloads.DeploymentWorkload = &Generator{}

func (gen *Generator) Deployment() *resource.Template[*appsv1.Deployment] {
	return resource.NewTemplateFromObjectFunction(gen.deployment).
		WithMutation(mutators.SetDeploymentReplicas(gen.Spec.HPA.IsDeactivated())).
		WithMutation(mutators.RolloutTrigger{Name: "cors-proxy-system-database", SecretName: util.Pointer("cors-proxy-system-database")}.Add())
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

func (gen *Generator) GrafanaDashboard() *resource.Template[*grafanav1alpha1.GrafanaDashboard] {
	return resource.NewTemplate[*grafanav1alpha1.GrafanaDashboard](
		grafanadashboard.New(gen.GetKey(), gen.GetLabels(), *gen.Spec.GrafanaDashboard, "dashboards/cors-proxy.json.gtpl")).
		WithEnabled(!gen.Spec.GrafanaDashboard.IsDeactivated())
}

func (gen *Generator) ExternalSecret() *resource.Template[*externalsecretsv1beta1.ExternalSecret] {
	return resource.NewTemplate[*externalsecretsv1beta1.ExternalSecret](
		pod.GenerateExternalSecretFn("cors-proxy-system-database", gen.GetNamespace(),
			*gen.Spec.Config.ExternalSecret.SecretStoreRef.Name, *gen.Spec.Config.ExternalSecret.SecretStoreRef.Kind,
			*gen.Spec.Config.ExternalSecret.RefreshInterval, gen.GetLabels(), gen.Options))
}
