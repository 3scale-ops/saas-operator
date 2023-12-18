package hpa

import (
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// New returns a basereconciler_types.GeneratorFunction function that will return an HorizontalPodAutoscaler
// resource when called
func New(key types.NamespacedName, labels map[string]string, cfg saasv1alpha1.HorizontalPodAutoscalerSpec) func(client.Object) (*autoscalingv2.HorizontalPodAutoscaler, error) {

	return func(client.Object) (*autoscalingv2.HorizontalPodAutoscaler, error) {
		hpa := autoscalingv2.HorizontalPodAutoscaler{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
				Labels:    labels,
			},
			Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
					APIVersion: appsv1.SchemeGroupVersion.String(),
					Kind:       "Deployment",
					Name:       key.Name,
				},
				MinReplicas: cfg.MinReplicas,
				Behavior:    cfg.Behavior,
			},
			Status: autoscalingv2.HorizontalPodAutoscalerStatus{
				Conditions: []autoscalingv2.HorizontalPodAutoscalerCondition{},
			},
		}

		if cfg.MaxReplicas != nil {
			hpa.Spec.MaxReplicas = *cfg.MaxReplicas
		}

		if cfg.ResourceName != nil {
			hpa.Spec.Metrics = []autoscalingv2.MetricSpec{
				{
					Type: autoscalingv2.ResourceMetricSourceType,
					Resource: &autoscalingv2.ResourceMetricSource{
						Name: corev1.ResourceName(*cfg.ResourceName),
						Target: autoscalingv2.MetricTarget{
							Type:               autoscalingv2.UtilizationMetricType,
							AverageUtilization: cfg.ResourceUtilization,
						},
					},
				},
			}
		}

		return &hpa, nil
	}
}
