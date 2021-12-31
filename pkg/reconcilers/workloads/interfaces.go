package workloads

import (
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2/resources"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/types"
)

/* Each of the workload types can be composed
of several of the features, each one of them
described by one of the following interfaces */

type WithKey interface {
	GetKey() types.NamespacedName
}

type WithLabels interface {
	GetLabels() map[string]string
}

type WithSelector interface {
	GetSelector() map[string]string
}

type WithWorkloadMeta interface {
	WithKey
	WithLabels
	WithSelector
}

type WithMonitoring interface {
	MonitoredEndpoints() []monitoringv1.PodMetricsEndpoint
}

type WithPodDisruptionBadget interface {
	PDBSpec() *saasv1alpha1.PodDisruptionBudgetSpec
}

type WithHorizontalPodAutoscaler interface {
	HPASpec() *saasv1alpha1.HorizontalPodAutoscalerSpec
}

type WithTraffic interface {
	WithSelector
	SendTraffic() bool
}

/* --------- */
/* WORKLOADS */
/* --------- */

type TrafficManager interface {
	WithKey
	WithLabels
	TrafficSelector() map[string]string
	Services() []resources.ServiceTemplate
}

type DeploymentWorkload interface {
	WithWorkloadMeta
	WithMonitoring
	WithHorizontalPodAutoscaler
	WithPodDisruptionBadget
	Deployment() resources.DeploymentTemplate
}

type DeploymentWorkloadWithTraffic interface {
	DeploymentWorkload
	WithTraffic
}
