package system

import (
	basereconciler_types "github.com/3scale/saas-operator/pkg/basereconciler/types"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/service"
	"github.com/3scale/saas-operator/pkg/generators/system/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Service returns a basereconciler_types.GeneratorFunction function that will return the
// Service resource when called
func (gen *SphinxGenerator) Service() basereconciler_types.GeneratorFunction {

	return func() client.Object {

		return &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: corev1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      config.SystemSphinxServiceName,
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
