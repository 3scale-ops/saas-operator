package autossl

import (
	"github.com/3scale/saas-operator/pkg/resource_builders/service"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// service returns a function that will return the corev1.Service for autossl
func (gen *Generator) service() func() *corev1.Service {

	return func() *corev1.Service {

		return &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:        gen.GetComponent(),
				Annotations: service.ELBServiceAnnotations(*gen.Spec.LoadBalancer, gen.Spec.Endpoint.DNS),
			},
			Spec: corev1.ServiceSpec{
				Type:                  corev1.ServiceTypeLoadBalancer,
				ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyTypeCluster,
				SessionAffinity:       corev1.ServiceAffinityNone,
				Ports: service.Ports(
					service.TCPPort("http", 80, intstr.FromString("http")),
					service.TCPPort("https", 443, intstr.FromString("https")),
				),
			},
		}
	}
}
