package types

import (
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GeneratorFunction is a function that returns a client.Object
type GeneratorFunction func() client.Object

// LockedResource is a struct that instructs the reconciler how to
// generate and reconcile a resource
type LockedResource struct {
	GeneratorFn  GeneratorFunction
	ExcludePaths []string
}

// ExtendedObjectList is an extension of client.ObjectList with methods
// to manipulate generically the objects in the list
type ExtendedObjectList interface {
	client.ObjectList
	GetItem(int) client.Object
	CountItems() int
}

type DeploymentWorkloadGenerator interface {
	Deployment() GeneratorFunction
	RolloutTriggers() []GeneratorFunction
	MonitoredEndpoints() []monitoringv1.PodMetricsEndpoint
	Key() types.NamespacedName
	GetLabels() map[string]string
	Selector() *metav1.LabelSelector
	HPASpec() *saasv1alpha1.HorizontalPodAutoscalerSpec
	PDBSpec() *saasv1alpha1.PodDisruptionBudgetSpec
	SendTraffic() bool
}

type DeploymentIngressGenerator interface {
	GetLabels() map[string]string
	TrafficSelector() map[string]string
	Services() []GeneratorFunction
}
