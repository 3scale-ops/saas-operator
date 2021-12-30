package redisshard

import (
	"fmt"

	basereconciler_types "github.com/3scale/saas-operator/pkg/basereconciler/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Service returns a basereconciler_types.GeneratorFunction function that will return a Service
// resource when called
func (gen *Generator) Service() basereconciler_types.GeneratorFunction {

	return func() client.Object {

		return &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: corev1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s", gen.GetComponent(), gen.GetInstanceName()),
				Namespace: gen.GetNamespace(),
				Labels:    gen.GetLabels(),
			},
			Spec: corev1.ServiceSpec{
				Type:            corev1.ServiceTypeClusterIP,
				ClusterIP:       corev1.ClusterIPNone,
				SessionAffinity: corev1.ServiceAffinityNone,
				Ports:           []corev1.ServicePort{},
				Selector:        gen.Selector().MatchLabels,
			},
		}
	}
}