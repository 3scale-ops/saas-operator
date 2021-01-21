package podmonitor

import (
	"fmt"

	"github.com/3scale/saas-operator/pkg/basereconciler"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// New returns a basereconciler.GeneratorFunction funtion that will return a PodMonitor
// resource when called
func New(key types.NamespacedName, labels map[string]string, selector map[string]string,
	endpoints ...monitoringv1.PodMetricsEndpoint) basereconciler.GeneratorFunction {

	return func() client.Object {

		return &monitoringv1.PodMonitor{
			TypeMeta: metav1.TypeMeta{
				Kind:       "PodMonitor",
				APIVersion: monitoringv1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
				Labels:    labels,
			},
			Spec: monitoringv1.PodMonitorSpec{
				PodMetricsEndpoints: endpoints,
				Selector: metav1.LabelSelector{
					MatchLabels: selector,
				},
			},
		}
	}
}

// PodMetricsEndpoint returns a monitoringv1.PodMetricsEndpoint
func PodMetricsEndpoint(path, port string, interval int32) monitoringv1.PodMetricsEndpoint {
	return monitoringv1.PodMetricsEndpoint{
		Interval: fmt.Sprintf("%ds", interval),
		Path:     path,
		Port:     port,
		Scheme:   "http",
	}
}
