package workloads

import (
	"github.com/3scale-ops/basereconciler/resources"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	descriptor "github.com/3scale/saas-operator/pkg/resource_builders/envoyconfig/descriptor"
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
	WithWorkloadMeta
	WithSelector
	SendTraffic() bool
	TrafficSelector() map[string]string
	Services() []resources.ServiceTemplate
}

type WithEnvoySidecar interface {
	WithWorkloadMeta
	EnvoyDynamicConfigurations() []descriptor.EnvoyDynamicConfigDescriptor
}

type DeploymentWorkload interface {
	WithWorkloadMeta
	WithMonitoring
	WithHorizontalPodAutoscaler
	WithPodDisruptionBadget
	Deployment() resources.DeploymentTemplate
}
