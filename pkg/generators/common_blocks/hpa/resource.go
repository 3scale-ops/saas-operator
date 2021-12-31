package hpa

import (
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// New returns a basereconciler_types.GeneratorFunction function that will return an HorizontalPodAutoscaler
// resource when called
func New(key types.NamespacedName, labels map[string]string, cfg saasv1alpha1.HorizontalPodAutoscalerSpec) basereconciler.GeneratorFunction {

	return func() client.Object {

		return &autoscalingv2beta2.HorizontalPodAutoscaler{
			TypeMeta: metav1.TypeMeta{
				Kind:       "HorizontalPodAutoscaler",
				APIVersion: autoscalingv2beta2.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
				Labels:    labels,
			},
			Spec: autoscalingv2beta2.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: autoscalingv2beta2.CrossVersionObjectReference{
					APIVersion: appsv1.SchemeGroupVersion.String(),
					Kind:       "Deployment",
					Name:       key.Name,
				},
				MinReplicas: cfg.MinReplicas,
				MaxReplicas: *cfg.MaxReplicas,
				Metrics: []autoscalingv2beta2.MetricSpec{
					{
						Type: autoscalingv2beta2.ResourceMetricSourceType,
						Resource: &autoscalingv2beta2.ResourceMetricSource{
							Name: corev1.ResourceName(*cfg.ResourceName),
							Target: autoscalingv2beta2.MetricTarget{
								Type:               autoscalingv2beta2.UtilizationMetricType,
								AverageUtilization: cfg.ResourceUtilization,
							},
						},
					},
				},
			},
			Status: autoscalingv2beta2.HorizontalPodAutoscalerStatus{
				Conditions: []autoscalingv2beta2.HorizontalPodAutoscalerCondition{},
			},
		}
	}
}
