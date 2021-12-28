package basereconciler

import (
	"context"

	basereconciler_types "github.com/3scale/saas-operator/pkg/basereconciler/types"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/deployment"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/hpa"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pdb"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/podmonitor"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/service"
	"github.com/3scale/saas-operator/pkg/util"
)

func (r *Reconciler) NewControlledResourcesFromDeploymentGenerator(ctx context.Context,
	ingress basereconciler_types.DeploymentIngressGenerator, workloads ...basereconciler_types.DeploymentWorkloadGenerator) (*ControlledResources, error) {

	resources := &ControlledResources{}

	for _, workload := range workloads {

		// Calculate rollout triggers
		triggers, err := r.TriggersFromSecretDefs(ctx, workload.RolloutTriggers()...)
		if err != nil {
			return nil, err
		}

		resources.Add(&ControlledResources{
			Deployments: []Deployment{{
				Template: func() basereconciler_types.GeneratorFunction {
					if ingress != nil {
						return deployment.New(workload.Key(), workload.GetLabels(),
							workload.Selector().MatchLabels, ingress.TrafficSelector(), workload.Deployment())

					} else {
						return deployment.New(workload.Key(), workload.GetLabels(),
							workload.Selector().MatchLabels, map[string]string{}, workload.Deployment())
					}
				}(),
				HasHPA:          !workload.HPASpec().IsDeactivated(),
				RolloutTriggers: triggers,
			}},
			PodDisruptionBudgets: []PodDisruptionBudget{{
				Template: pdb.New(workload.Key(), workload.GetLabels(), workload.Selector().MatchLabels, *workload.PDBSpec()),
				Enabled:  !workload.PDBSpec().IsDeactivated(),
			}},
			HorizontalPodAutoscalers: []HorizontalPodAutoscaler{{
				Template: hpa.New(workload.Key(), workload.GetLabels(), *workload.HPASpec()),
				Enabled:  !workload.HPASpec().IsDeactivated(),
			}},
			PodMonitors: []PodMonitor{{
				Template: podmonitor.New(workload.Key(), workload.GetLabels(), workload.Selector().MatchLabels, workload.MonitoredEndpoints()...),
				Enabled:  len(workload.MonitoredEndpoints()) > 0,
			}},
		})

	}

	if ingress != nil {
		resources.Add(&ControlledResources{
			Services: func() []Service {
				resolvedTrafficSelector := TrafficSwitcher(ingress, workloads...)
				svcs := make([]Service, 0, len(ingress.Services()))
				for _, fn := range ingress.Services() {
					svcs = append(svcs, Service{
						Template: service.New(ingress.GetLabels(), resolvedTrafficSelector, fn),
						Enabled:  fn != nil})
				}
				return svcs
			}(),
		})
	}

	return resources, nil
}

func TrafficSwitcher(ingress basereconciler_types.DeploymentIngressGenerator,
	workloads ...basereconciler_types.DeploymentWorkloadGenerator) map[string]string {

	// NOTE: due to the fact that services do not yet support set-based selectors, only MatchLabels selectors
	// can be used. This limits a lot what can be done in terms of deciding where to send traffic, as all
	// Deployments that should receive traffic need to have the same labels. The only way of doing this
	// without modifying the Deployment labels (which would trigger a rollout) and acting on the Service
	// selector alone is to choose only between three options:
	//                   traffic to noone / traffic to a single Deployment / traffic to all
	//
	// There seems to be great demand for set-based selectors for Services but it is not yet implamented:
	// https://github.com/kubernetes/kubernetes/issues/48528
	enabledSelectors := []map[string]string{}
	for _, workload := range workloads {
		if workload.SendTraffic() {
			enabledSelectors = append(enabledSelectors, workload.Selector().MatchLabels)
		}
	}

	switch c := len(enabledSelectors); c {
	case 0:
		return map[string]string{}
	case 1:
		// If there is only one Deployment with SendTraffic() active
		// return its selector together with the shared traffic selector
		return util.MergeMaps(enabledSelectors[0], ingress.TrafficSelector())
	default:
		// If there is more than one Deployment with SendTraffic() active
		// send traffic to all Deployments by using the shared traffic selector
		return ingress.TrafficSelector()
	}
}
