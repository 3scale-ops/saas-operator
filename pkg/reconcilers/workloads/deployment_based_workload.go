package workloads

import (
	"reflect"

	basereconciler "github.com/3scale-ops/basereconciler/reconciler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type WorkloadReconciler struct {
	basereconciler.Reconciler
}

// NewFromManager constructs a new Reconciler from the given manager
func NewFromManager(mgr manager.Manager) WorkloadReconciler {
	return WorkloadReconciler{
		Reconciler: basereconciler.NewFromManager(mgr),
	}
}

func (r *WorkloadReconciler) NewDeploymentWorkload(
	main DeploymentWorkload, canary DeploymentWorkload) ([]basereconciler.Resource, error) {

	resources := workloadResources(main)

	if unwrapNil(canary) != nil {
		resources = append(resources, workloadResources(canary)...)
	}

	// Generate services if the workload implements WithTraffic interface
	if _, ok := main.(WithTraffic); ok {
		for _, svct := range main.(WithTraffic).Services() {
			resources = append(resources, NewServiceTemplate(svct).ApplyMeta(main.(WithTraffic)).
				ApplyTrafficSelector(main.(WithTraffic), toWithTraffic(canary)))
		}
	}

	return resources, nil
}

func workloadResources(workload DeploymentWorkload) []basereconciler.Resource {

	dep := NewDeploymentTemplate(workload.Deployment()).ApplyMeta(workload)

	// if workload implements WithTraffic add the TrafficSelector
	if w, ok := workload.(WithTraffic); ok {
		dep = dep.ApplyTrafficSelector(w)
	}

	hpa := NewHorizontalPodAutoscalerTemplateFromSpec(*workload.HPASpec()).ApplyMeta(workload)
	pdb := NewPodDisruptionBudgetTemplateFromSpec(*workload.PDBSpec()).ApplyMeta(workload)
	pm := NewPodMonitorTemplateFromEndpoints(workload.MonitoredEndpoints()...).ApplyMeta(workload)

	resources := []basereconciler.Resource{dep, hpa, pdb, pm}

	// if workload implements WithEnvoySidecar add the EnvoyConfig
	if w, ok := workload.(WithEnvoySidecar); ok {
		resources = append(resources,
			NewEnvoyConfigTemplateFromEnvoyResources(w.EnvoyDynamicConfigurations()).ApplyMeta(w).SetNodeID(w))
	}

	return resources
}

func unwrapNil(w DeploymentWorkload) DeploymentWorkload {
	if w == nil || reflect.ValueOf(w).IsNil() {
		return nil
	}
	return w
}

func toWithTraffic(w DeploymentWorkload) WithTraffic {
	if w == nil || reflect.ValueOf(w).IsNil() {
		return nil
	}
	return w.(WithTraffic)
}
