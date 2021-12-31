package echoapi

import (
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/service"
	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Service returns a basereconciler.GeneratorFunction function that will return a Service
// resource when called
func (gen *Generator) Service() basereconciler.GeneratorFunction {

	return func() client.Object {

		return &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: corev1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:        gen.GetComponent(),
				Namespace:   gen.GetNamespace(),
				Labels:      gen.GetLabels(),
				Annotations: service.NLBServiceAnnotations(*gen.Spec.LoadBalancer, gen.Spec.Endpoint.DNS),
			},
			Spec: corev1.ServiceSpec{
				Type:                  corev1.ServiceTypeLoadBalancer,
				ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyTypeCluster,
				SessionAffinity:       corev1.ServiceAffinityNone,
				Ports: func() []corev1.ServicePort {
					if gen.Spec.Marin3r.IsDeactivated() {
						return service.Ports(
							service.TCPPort("http", 80, intstr.FromString("http")),
						)
					}
					return service.Ports(
						service.TCPPort("http", 80, intstr.FromString("echo-api-http")),
						service.TCPPort("https", 443, intstr.FromString("echo-api-https")),
					)
				}(),
				Selector: gen.Selector().MatchLabels,
			},
		}
	}
}
