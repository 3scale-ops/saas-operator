package autossl

import (
	"github.com/3scale/saas-operator/pkg/basereconciler"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Service returns a basereconciler.GeneratorFunction funtion that will return a Service
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
				Labels:      gen.Labels(),
				Annotations: gen.ELBServiceAnnotations(*gen.Spec.LoadBalancer, gen.Spec.Endpoint.DNS),
			},
			Spec: corev1.ServiceSpec{
				Type:                  corev1.ServiceTypeLoadBalancer,
				ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyTypeCluster,
				SessionAffinity:       corev1.ServiceAffinityNone,
				Ports: gen.ServicePorts(
					gen.ServicePortTCP("http", 80, intstr.FromString("http")),
					gen.ServicePortTCP("https", 443, intstr.FromString("https")),
				),
				Selector: gen.Selector().MatchLabels,
			},
		}
	}
}
