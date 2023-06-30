package system

import (
	"github.com/3scale/saas-operator/pkg/resource_builders/service"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// service returns a function that will return the corev1.Service for searchd
func (gen *SearchdGenerator) service() func() *corev1.Service {

	return func() *corev1.Service {

		return &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      *gen.Spec.Config.ServiceName,
				Namespace: gen.GetNamespace(),
			},
			Spec: corev1.ServiceSpec{
				Type:            corev1.ServiceTypeClusterIP,
				SessionAffinity: corev1.ServiceAffinityNone,
				Ports: service.Ports(
					service.TCPPort("searchd", gen.DatabasePort, intstr.FromString("searchd")),
				),
				Selector: gen.GetSelector(),
			},
		}
	}
}
