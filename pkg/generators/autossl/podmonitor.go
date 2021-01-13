package autossl

import (
	"github.com/3scale/saas-operator/pkg/basereconciler"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PodMonitor returns a basereconciler.GeneratorFunction funtion that will return a PodMonitor
// resource when called
func (opts *Options) PodMonitor() basereconciler.GeneratorFunction {

	return func() client.Object {

		return &monitoringv1.PodMonitor{
			TypeMeta: metav1.TypeMeta{
				Kind:       "PodMonitor",
				APIVersion: monitoringv1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      Component,
				Namespace: opts.Namespace,
				Labels:    opts.labels(),
			},
			Spec: monitoringv1.PodMonitorSpec{
				PodMetricsEndpoints: []monitoringv1.PodMetricsEndpoint{
					{
						Interval: "30s",
						Path:     "/metrics",
						Port:     "metrics",
						RelabelConfigs: []*monitoringv1.RelabelConfig{
							{
								SourceLabels: []string{" __meta_kubernetes_service_label_app"},
								TargetLabel:  "app",
							},
						},
					},
				},
				Selector: *opts.selector(),
			},
		}
	}
}
