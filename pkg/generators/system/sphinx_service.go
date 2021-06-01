package system

import (
	"github.com/3scale/saas-operator/pkg/basereconciler"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/service"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Service returns a basereconciler.GeneratorFunction function that will return the
// Service resource when called
func (gen *SphinxGenerator) Service() basereconciler.GeneratorFunction {

	return func() client.Object {

		return &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: corev1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      gen.DatabaseService,
				Namespace: gen.GetNamespace(),
				Labels:    gen.GetLabels(),
			},
			Spec: corev1.ServiceSpec{
				Type:            corev1.ServiceTypeClusterIP,
				SessionAffinity: corev1.ServiceAffinityNone,
				Ports: service.Ports(
					service.TCPPort("sphinx", gen.DatabasePort, intstr.FromString("sphinx")),
				),
				Selector: gen.Selector().MatchLabels,
			},
		}
	}
}
