package apicast

import (
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/service"
	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GatewayService returns a basereconciler.GeneratorFunction function that will return the
// gateway Service resource when called
func (gen *EnvGenerator) GatewayService() basereconciler.GeneratorFunction {

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
				Annotations: service.ELBServiceAnnotations(*gen.Spec.LoadBalancer, gen.Spec.Endpoint.DNS),
			},
			Spec: corev1.ServiceSpec{
				Type:                  corev1.ServiceTypeLoadBalancer,
				ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyTypeCluster,
				SessionAffinity:       corev1.ServiceAffinityNone,
				Ports: func() []corev1.ServicePort {
					if gen.Spec.Marin3r.IsDeactivated() {
						return service.Ports(
							service.TCPPort("http", 80, intstr.FromString("gateway")),
						)
					}
					return service.Ports(
						service.TCPPort("gateway-http", 80, intstr.FromString("gateway-http")),
						service.TCPPort("gateway-https", 443, intstr.FromString("gateway-https")),
					)
				}(),
				Selector: gen.Selector().MatchLabels,
			},
		}
	}
}

// MgmtService returns a basereconciler.GeneratorFunction function that will return the
// management Service resource when called
func (gen *EnvGenerator) MgmtService() basereconciler.GeneratorFunction {

	return func() client.Object {

		return &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: corev1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      gen.GetComponent() + "-management",
				Namespace: gen.GetNamespace(),
				Labels:    gen.GetLabels(),
			},
			Spec: corev1.ServiceSpec{
				Type:            corev1.ServiceTypeClusterIP,
				SessionAffinity: corev1.ServiceAffinityNone,
				Ports: service.Ports(
					service.TCPPort("management", 8090, intstr.FromString("management")),
				),
				Selector: gen.Selector().MatchLabels,
			},
		}
	}
}
