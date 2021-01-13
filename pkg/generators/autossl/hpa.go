package autossl

import (
	"github.com/3scale/saas-operator/pkg/basereconciler"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HPA returns a basereconciler.GeneratorFunction funtion that will return an HPA
// resource when called
func (opts *Options) HPA() basereconciler.GeneratorFunction {

	return func() client.Object {

		return &autoscalingv2beta2.HorizontalPodAutoscaler{
			TypeMeta: metav1.TypeMeta{
				Kind:       "HorizontalPodAutoscaler",
				APIVersion: autoscalingv2beta2.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      Component,
				Namespace: opts.Namespace,
				Labels:    opts.labels(),
			},
			Spec: autoscalingv2beta2.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: autoscalingv2beta2.CrossVersionObjectReference{
					APIVersion: appsv1.SchemeGroupVersion.String(),
					Kind:       "Deployment",
					Name:       Component,
				},
				MinReplicas: opts.Spec.HPA.MinReplicas,
				MaxReplicas: *opts.Spec.HPA.MaxReplicas,
				Metrics: []autoscalingv2beta2.MetricSpec{
					{
						Type: autoscalingv2beta2.ResourceMetricSourceType,
						Resource: &autoscalingv2beta2.ResourceMetricSource{
							Name: corev1.ResourceName(*opts.Spec.HPA.ResourceName),
							Target: autoscalingv2beta2.MetricTarget{
								Type:               autoscalingv2beta2.UtilizationMetricType,
								AverageUtilization: opts.Spec.HPA.ResourceUtilization,
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
