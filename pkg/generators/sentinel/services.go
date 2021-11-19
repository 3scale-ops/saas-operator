package sentinel

import (
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	statefulsetPodSelectorLabelKey string = "statefulset.kubernetes.io/pod-name"
)

// Service returns a basereconciler.GeneratorFunction function that will return a Service
// resource when called
func (gen *Generator) StatefulSetService() basereconciler.GeneratorFunction {

	return func() client.Object {

		return &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: corev1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      gen.GetComponent() + "-headless",
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

// Service returns a basereconciler.GeneratorFunction function that will return a Service
// resource when called
func (gen *Generator) PodServices(replicas int) []basereconciler.GeneratorFunction {

	fn := func(i int) func() client.Object {

		return func() client.Object {
			return &corev1.Service{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Service",
					APIVersion: corev1.SchemeGroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-%d", gen.GetComponent(), i),
					Namespace: gen.GetNamespace(),
					Labels:    gen.GetLabels(),
				},
				Spec: corev1.ServiceSpec{
					Type:            corev1.ServiceTypeClusterIP,
					SessionAffinity: corev1.ServiceAffinityNone,
					Ports: []corev1.ServicePort{{
						Name:       gen.GetComponent(),
						Protocol:   corev1.ProtocolTCP,
						Port:       int32(saasv1alpha1.SentinelPort),
						TargetPort: intstr.FromString(gen.GetComponent()),
					}},
					Selector: map[string]string{
						statefulsetPodSelectorLabelKey: fmt.Sprintf("%s-%d", gen.GetComponent(), i),
					},
				},
			}
		}
	}

	svcFns := make([]basereconciler.GeneratorFunction, replicas)
	for idx := 0; idx < replicas; idx++ {
		// log.Printf("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ %d", idx)
		svcFns[idx] = fn(idx)
	}

	return svcFns
}
