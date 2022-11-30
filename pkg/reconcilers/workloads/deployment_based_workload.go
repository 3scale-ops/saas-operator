package workloads

import (
	"reflect"

	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
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

	// if workload implements TrafficManager add the TrafficSelector
	if workloadWithTraffic, ok := workload.(WithTraffic); ok {
		dep = dep.ApplyTrafficSelector(workloadWithTraffic)
	}

	hpa := NewHorizontalPodAutoscalerTemplateFromSpec(*workload.HPASpec()).ApplyMeta(workload)
	pdb := NewPodDisruptionBudgetTemplateFromSpec(*workload.PDBSpec()).ApplyMeta(workload)
	pm := NewPodMonitorTemplateFromEndpoints(workload.MonitoredEndpoints()...).ApplyMeta(workload)

	return []basereconciler.Resource{dep, hpa, pdb, pm}
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
