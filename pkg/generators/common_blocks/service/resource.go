package service

import (
	basereconciler_types "github.com/3scale/saas-operator/pkg/basereconciler/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// New returns a basereconciler_types.GeneratorFunction function that will return a Service
// resource when called
func New(labels map[string]string, selector map[string]string,
	fn basereconciler_types.GeneratorFunction) basereconciler_types.GeneratorFunction {

	return func() client.Object {

		svc := fn().(*corev1.Service)

		return &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: corev1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:        svc.GetName(),
				Namespace:   svc.GetNamespace(),
				Labels:      labels,
				Annotations: svc.GetAnnotations(),
			},
			Spec: func() corev1.ServiceSpec {
				svc.Spec.Selector = selector
				return svc.Spec
			}(),
		}
	}
}
