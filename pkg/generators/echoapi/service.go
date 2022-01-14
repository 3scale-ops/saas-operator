package echoapi

import (
	"github.com/3scale/saas-operator/pkg/resource_builders/service"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// service returns a function that will return the corev1.Service for echo-api
func (gen *Generator) service() func() *corev1.Service {

	return func() *corev1.Service {

		return &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:        gen.GetComponent(),
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
			},
		}
	}
}
