package basereconciler

import (
	"context"

	basereconciler_types "github.com/3scale/saas-operator/pkg/basereconciler/types"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/hpa"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pdb"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/podmonitor"
)

func (r *Reconciler) NewControlledResourcesFromDeploymentGenerator(ctx context.Context, gen basereconciler_types.DeploymentWorkloadGenerator) (*ControlledResources, error) {

	// Calculate rollout triggers
	triggers, err := r.TriggersFromSecretDefs(ctx, gen.RolloutTriggers()...)
	if err != nil {
		return nil, err
	}

	return &ControlledResources{
		Deployments: []Deployment{{
			Template:        gen.Deployment(),
			HasHPA:          gen.HPASpec() != nil,
			RolloutTriggers: triggers,
		}},
		PodDisruptionBudgets: []PodDisruptionBudget{{
			Template: pdb.New(gen.Key(), gen.GetLabels(), gen.Selector().MatchLabels, *gen.PDBSpec()),
			Enabled:  !gen.PDBSpec().IsDeactivated(),
		}},
		HorizontalPodAutoscalers: []HorizontalPodAutoscaler{{
			Template: hpa.New(gen.Key(), gen.GetLabels(), *gen.HPASpec()),
			Enabled:  !gen.HPASpec().IsDeactivated(),
		}},
		PodMonitors: []PodMonitor{{
			Template: podmonitor.New(gen.Key(), gen.GetLabels(), gen.Selector().MatchLabels, gen.MonitoredEndpoints()...),
			Enabled:  len(gen.MonitoredEndpoints()) > 0,
		}},
		Services: func() []Service {
			svcs := make([]Service, 0, len(gen.Services()))
			for _, fn := range gen.Services() {
				svcs = append(svcs, Service{Template: fn, Enabled: fn != nil})
			}
			return svcs
		}(),
	}, nil
}
