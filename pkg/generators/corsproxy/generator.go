package corsproxy

import (
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	"github.com/3scale/saas-operator/pkg/generators"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/grafanadashboard"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/hpa"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pdb"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/podmonitor"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/secrets"

	"k8s.io/apimachinery/pkg/types"
)

const (
	component             string = "cors-proxy"
	databaseURLSecretName string = "cors-proxy-system-database"
)

// Generator configures the generators for CORSProxy
type Generator struct {
	generators.BaseOptions
	Spec saasv1alpha1.CORSProxySpec
}

// NewGenerator returns a new Options struct
func NewGenerator(instance, namespace string, spec saasv1alpha1.CORSProxySpec) Generator {
	return Generator{
		BaseOptions: generators.BaseOptions{
			Component:    component,
			InstanceName: instance,
			Namespace:    namespace,
			Labels: map[string]string{
				"app":     component,
				"part-of": "3scale-saas",
			},
		},
		Spec: spec,
	}
}

// HPA returns a basereconciler.GeneratorFunction
func (gen *Generator) HPA() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return hpa.New(key, gen.GetLabels(), *gen.Spec.HPA)
}

// PDB returns a basereconciler.GeneratorFunction
func (gen *Generator) PDB() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return pdb.New(key, gen.GetLabels(), gen.Selector().MatchLabels, *gen.Spec.PDB)
}

// PodMonitor returns a basereconciler.GeneratorFunction
func (gen *Generator) PodMonitor() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return podmonitor.New(key, gen.GetLabels(), gen.Selector().MatchLabels, podmonitor.PodMetricsEndpoint("/metrics", "metrics", 30))
}

// GrafanaDashboard returns a basereconciler.GeneratorFunction
func (gen *Generator) GrafanaDashboard() basereconciler.GeneratorFunction {
	key := types.NamespacedName{Name: gen.Component, Namespace: gen.Namespace}
	return grafanadashboard.New(key, gen.GetLabels(), *gen.Spec.GrafanaDashboard, "dashboards/cors-proxy.json.tpl")
}

// SecretDefinition returns a basereconciler.GeneratorFunction
func (gen *Generator) SecretDefinition() basereconciler.GeneratorFunction {
	key := types.NamespacedName{
		Name:      gen.Component + "-system-database",
		Namespace: gen.Namespace,
	}
	return secrets.NewSecretDefinition(key, gen.GetLabels(), key.Name, gen.Spec.Config.SystemDatabaseDSN.FromVault)
}
