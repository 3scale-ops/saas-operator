package workloads

import (
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	basereconciler_resources "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2/resources"
	"github.com/3scale/saas-operator/pkg/resource_builders/hpa"
	"github.com/3scale/saas-operator/pkg/resource_builders/pdb"
	"github.com/3scale/saas-operator/pkg/resource_builders/podmonitor"
	"github.com/3scale/saas-operator/pkg/util"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var (
	EmptyKey      types.NamespacedName = types.NamespacedName{}
	EmptyLabel    map[string]string    = map[string]string{}
	EmptySelector map[string]string    = map[string]string{}
)

// DeploymentTemplate specifies a Deployment resource and its rollout triggers
type DeploymentTemplate struct {
	basereconciler_resources.DeploymentTemplate
}

func NewDeploymentTemplate(t basereconciler_resources.DeploymentTemplate) DeploymentTemplate {
	return DeploymentTemplate{DeploymentTemplate: t}
}

func (dt DeploymentTemplate) ApplyMeta(gen DeploymentWorkload) DeploymentTemplate {

	fn := dt.Template
	dt.Template = func() *appsv1.Deployment {

		dep := fn()

		// Set basic resource metadata
		applyKey(dep, gen)
		applyLabels(dep, gen)

		// Set the Pod selector
		dep.Spec.Selector = &metav1.LabelSelector{MatchLabels: gen.GetSelector()}
		// Set the Pod labels
		dep.Spec.Template.ObjectMeta.Labels = util.MergeMaps(map[string]string{}, dep.Spec.Template.ObjectMeta.Labels, gen.GetLabels(), gen.GetSelector())
		// Set the Pod annotations
		dep.Spec.Template.ObjectMeta.Labels = util.MergeMaps(map[string]string{}, dep.Spec.Template.ObjectMeta.Labels, gen.GetLabels(), gen.GetSelector())

		return dep
	}

	return dt
}

func (dt DeploymentTemplate) ApplyTrafficSelector(wt WithTraffic) DeploymentTemplate {

	fn := dt.Template
	dt.Template = func() *appsv1.Deployment {
		dep := fn()
		dep.Spec.Template.ObjectMeta.Labels = util.MergeMaps(map[string]string{}, dep.Spec.Template.ObjectMeta.Labels, wt.TrafficSelector())
		return dep
	}

	return dt
}

// ServicesTemplate specifies a Services resource
type ServiceTemplate struct {
	basereconciler_resources.ServiceTemplate
}

func NewServiceTemplate(t basereconciler_resources.ServiceTemplate) ServiceTemplate {
	return ServiceTemplate{ServiceTemplate: t}
}

func (st ServiceTemplate) ApplyMeta(wt WithTraffic) ServiceTemplate {
	fn := st.Template
	st.Template = func() *corev1.Service {
		svc := fn()
		// Do not enforce metadata.name:
		//   Services are special because there can be more than one of them, so the Name
		//   is relevant and must be provided by the service template
		svc.SetNamespace(wt.GetKey().Namespace)
		applyLabels(svc, wt)
		return svc
	}

	return st
}

func (st ServiceTemplate) ApplyTrafficSelector(main WithTraffic, canary WithTraffic) ServiceTemplate {
	fn := st.Template
	st.Template = func() *corev1.Service {
		svc := fn()
		svc.Spec.Selector = trafficSwitcher(main, canary)
		return svc
	}

	return st
}

func trafficSwitcher(main WithTraffic, canary WithTraffic) map[string]string {

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
	for _, workload := range []WithTraffic{main, canary} {
		if workload != nil && workload.SendTraffic() {
			enabledSelectors = append(enabledSelectors, workload.GetSelector())
		}
	}

	switch c := len(enabledSelectors); c {
	case 0:
		return map[string]string{}
	case 1:
		// If there is only one Deployment with SendTraffic() active
		// return its selector together with the shared traffic selector
		return util.MergeMaps(enabledSelectors[0], main.TrafficSelector())
	default:
		// If there is more than one Deployment with SendTraffic() active
		// send traffic to all Deployments by using the shared traffic selector
		return main.TrafficSelector()
	}
}

// PodDisruptionBudget specifies a PodDisruptionBudget resource
type PodDisruptionBudgetTemplate struct {
	basereconciler_resources.PodDisruptionBudgetTemplate
}

func (pdbt PodDisruptionBudgetTemplate) ApplyMeta(w WithWorkloadMeta) PodDisruptionBudgetTemplate {
	fn := pdbt.Template
	pdbt.Template = func() *policyv1.PodDisruptionBudget {
		pdb := fn()
		applyKey(pdb, w)
		applyLabels(pdb, w)
		pdb.Spec.Selector = &metav1.LabelSelector{MatchLabels: w.GetSelector()}
		return pdb
	}
	return pdbt
}

func NewPodDisruptionBudgetTemplate(t basereconciler_resources.PodDisruptionBudgetTemplate) PodDisruptionBudgetTemplate {
	return PodDisruptionBudgetTemplate{PodDisruptionBudgetTemplate: t}
}

func NewPodDisruptionBudgetTemplateFromSpec(cfg saasv1alpha1.PodDisruptionBudgetSpec) PodDisruptionBudgetTemplate {
	return NewPodDisruptionBudgetTemplate(basereconciler_resources.PodDisruptionBudgetTemplate{
		Template:  pdb.New(EmptyKey, EmptyLabel, EmptySelector, cfg),
		IsEnabled: !cfg.IsDeactivated(),
	})
}

// HorizontalPodAutoscaler specifies a HorizontalPodAutoscaler resource
type HorizontalPodAutoscalerTemplate struct {
	basereconciler_resources.HorizontalPodAutoscalerTemplate
}

func (hpat HorizontalPodAutoscalerTemplate) ApplyMeta(w WithWorkloadMeta) HorizontalPodAutoscalerTemplate {
	fn := hpat.Template
	hpat.Template = func() *autoscalingv2beta2.HorizontalPodAutoscaler {
		hpa := fn()
		applyKey(hpa, w)
		applyLabels(hpa, w)
		hpa.Spec.ScaleTargetRef = autoscalingv2beta2.CrossVersionObjectReference{
			Kind:       "Deployment",
			Name:       w.GetKey().Name,
			APIVersion: appsv1.SchemeGroupVersion.String(),
		}
		return hpa
	}
	return hpat
}

func NewHorizontalPodAutoscalerTemplate(t basereconciler_resources.HorizontalPodAutoscalerTemplate) HorizontalPodAutoscalerTemplate {
	return HorizontalPodAutoscalerTemplate{HorizontalPodAutoscalerTemplate: t}
}

func NewHorizontalPodAutoscalerTemplateFromSpec(cfg saasv1alpha1.HorizontalPodAutoscalerSpec) HorizontalPodAutoscalerTemplate {
	return NewHorizontalPodAutoscalerTemplate(basereconciler_resources.HorizontalPodAutoscalerTemplate{
		Template:  hpa.New(EmptyKey, EmptyLabel, cfg),
		IsEnabled: !cfg.IsDeactivated(),
	})
}

// PodMonitor specifies a PodMonitor resource
type PodMonitorTemplate struct {
	basereconciler_resources.PodMonitorTemplate
}

func (pmt PodMonitorTemplate) ApplyMeta(w WithWorkloadMeta) PodMonitorTemplate {
	fn := pmt.Template
	pmt.Template = func() *monitoringv1.PodMonitor {
		pm := fn()
		applyKey(pm, w)
		applyLabels(pm, w)
		pm.Spec.Selector = metav1.LabelSelector{MatchLabels: w.GetSelector()}
		return pm
	}
	return pmt
}

func NewPodMonitorTemplate(t basereconciler_resources.PodMonitorTemplate) PodMonitorTemplate {
	return PodMonitorTemplate{PodMonitorTemplate: t}
}

func NewPodMonitorTemplateFromEndpoints(endpoints ...monitoringv1.PodMetricsEndpoint) PodMonitorTemplate {
	return NewPodMonitorTemplate(basereconciler_resources.PodMonitorTemplate{
		Template:  podmonitor.New(EmptyKey, EmptyLabel, EmptySelector, endpoints...),
		IsEnabled: len(endpoints) > 0,
	})
}
