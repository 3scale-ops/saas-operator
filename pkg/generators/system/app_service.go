package system

import (
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/service"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// service returns a function that will return the corev1.Service for system
func (gen *AppGenerator) service() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gen.GetComponent(),
			Namespace: gen.GetNamespace(),
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: service.Ports(
				service.TCPPort("http", 3000, intstr.FromString("ui-api")),
			),
		},
	}
}
