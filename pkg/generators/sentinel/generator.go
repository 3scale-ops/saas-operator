package sentinel

import (
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators"
	"github.com/3scale/saas-operator/pkg/generators/sentinel/config"
	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	basereconciler_resources "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2/resources"
	"github.com/3scale/saas-operator/pkg/resource_builders/grafanadashboard"
	"github.com/3scale/saas-operator/pkg/resource_builders/pdb"
)

const (
	component string = "redis-sentinel"
)

// Generator configures the generators for Sentinel
type Generator struct {
	generators.BaseOptionsV2
	Spec    saasv1alpha1.SentinelSpec
	Options config.Options
}

// NewGenerator returns a new Options struct
func NewGenerator(instance, namespace string, spec saasv1alpha1.SentinelSpec) Generator {
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
	}
}

// Returns all the resource templates that this generator manages
func (gen *Generator) Resources() []basereconciler.Resource {
	resources := []basereconciler.Resource{
		basereconciler_resources.StatefulSetTemplate{
			Template:        gen.statefulSet(),
			RolloutTriggers: []basereconciler_resources.RolloutTrigger{},
			IsEnabled:       true,
		},
		basereconciler_resources.ServiceTemplate{
			Template:  gen.statefulSetService(),
			IsEnabled: true,
		},
		basereconciler_resources.PodDisruptionBudgetTemplate{
			Template:  pdb.New(gen.GetKey(), gen.GetLabels(), gen.GetSelector(), *gen.Spec.PDB),
			IsEnabled: !gen.Spec.PDB.IsDeactivated(),
		},
		basereconciler_resources.ConfigMapTemplate{
			Template:  gen.configMap(),
			IsEnabled: true,
		},
	}

	for idx := 0; idx < int(*gen.Spec.Replicas); idx++ {
		resources = append(resources,
			basereconciler_resources.ServiceTemplate{Template: gen.podServices(idx), IsEnabled: true})
	}

	return resources
}

func (gen *Generator) GrafanaDashboard() basereconciler_resources.GrafanaDashboardTemplate {
	return basereconciler_resources.GrafanaDashboardTemplate{
		Template:  grafanadashboard.New(gen.GetKey(), gen.GetLabels(), *gen.Spec.GrafanaDashboard, "dashboards/redis-sentinel.json.gtpl"),
		IsEnabled: !gen.Spec.GrafanaDashboard.IsDeactivated(),
	}
}
