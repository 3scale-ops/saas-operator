package workloads

import (
	"context"

	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type WorkloadReconciler struct {
	basereconciler.Reconciler
}

// NewFromManager constructs a new Reconciler from the given manager
func NewFromManager(mgr manager.Manager, recorder record.EventRecorder, clusterWatchers bool) WorkloadReconciler {
	return WorkloadReconciler{
		Reconciler: basereconciler.NewFromManager(mgr, recorder, clusterWatchers),
	}
}

func (r *WorkloadReconciler) NewDeploymentWorkload(ctx context.Context, owner client.Object,
	scheme *runtime.Scheme, workloads ...DeploymentWorkload) ([]basereconciler.Resource, error) {

	resources := []basereconciler.Resource{}

	for _, workload := range workloads {
		resources = append(resources,
			NewDeploymentTemplate(workload.Deployment()).ApplyMeta(workload),
			NewHorizontalPodAutoscalerTemplateFromSpec(*workload.HPASpec()).ApplyMeta(workload),
			NewPodDisruptionBudgetTemplateFromSpec(*workload.PDBSpec()).ApplyMeta(workload),
			NewPodMonitorTemplateFromEndpoints(workload.MonitoredEndpoints()...).ApplyMeta(workload),
		)
	}

	return resources, nil
}

func (r *WorkloadReconciler) NewDeploymentWorkloadWithTraffic(ctx context.Context, owner client.Object,
	scheme *runtime.Scheme, trafficManager TrafficManager, workloads ...DeploymentWorkloadWithTraffic) ([]basereconciler.Resource, error) {

	resources := []basereconciler.Resource{}

	for _, workload := range workloads {
		resources = append(resources,
			NewDeploymentTemplate(workload.Deployment()).ApplyMeta(workload).ApplyTrafficSelector(trafficManager),
			NewHorizontalPodAutoscalerTemplateFromSpec(*workload.HPASpec()).ApplyMeta(workload),
			NewPodDisruptionBudgetTemplateFromSpec(*workload.PDBSpec()).ApplyMeta(workload),
			NewPodMonitorTemplateFromEndpoints(workload.MonitoredEndpoints()...).ApplyMeta(workload),
		)
	}

	for _, svct := range trafficManager.Services() {
		resources = append(resources, NewServiceTemplate(svct).ApplyMeta(trafficManager).
			ApplyTrafficSelector(trafficManager, sliceDeploymentWorkloadWithTraffic_to_sliceWithTraffic(workloads)...))
	}

	return resources, nil
}

func sliceDeploymentWorkloadWithTraffic_to_sliceWithTraffic(a []DeploymentWorkloadWithTraffic) []WithTraffic {
	// Go does not automatically convert []DeploymentWorkloadWithTraffic to []WithTraffic
	// even though all the elements implement both interfaces. Conversion must be manually performed.
	b := make([]WithTraffic, 0, len(a))
	for _, item := range a {
		b = append(b, WithTraffic(item))
	}
	return b
}
