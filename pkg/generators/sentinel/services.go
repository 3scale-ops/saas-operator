package sentinel

import (
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	statefulsetPodSelectorLabelKey string = "statefulset.kubernetes.io/pod-name"
)

// statefulSetService returns a function function that returns a Service
// resource when called
func (gen *Generator) service() func() *corev1.Service {

	return func() *corev1.Service {

		return &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: corev1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      gen.GetComponent(),
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
				Selector: gen.GetSelector(),
			},
		}
	}
}

// SentinelServiceEndpoint returns the URI of the ClusterIP Service
func (gen *Generator) SentinelServiceEndpoint() string {
	return fmt.Sprintf("redis://%s.%s.svc.cluster.local:%d", gen.GetComponent(), gen.GetNamespace(), saasv1alpha1.SentinelPort)
}

// statefulSetService returns a function function that returns a Service
// resource when called
func (gen *Generator) statefulSetService() func() *corev1.Service {

	return func() *corev1.Service {

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
				Selector:        gen.GetSelector(),
			},
		}
	}
}

// podServices returns a function that returns a Service that points
// ot a specific StatefulSet Pod when called
// resource when called
func (gen *Generator) podServices(index int) func() *corev1.Service {

	return func() *corev1.Service {
		return &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: corev1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      gen.PodServiceName(index),
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
					statefulsetPodSelectorLabelKey: fmt.Sprintf("%s-%d", gen.GetComponent(), index),
				},
			},
		}
	}
}

// PodServiceName generates the name of the pod specific Service
func (gen *Generator) PodServiceName(index int) string {
	return fmt.Sprintf("%s-%d", gen.GetComponent(), index)
}

// SentinelEndpoints returns the list of redis URLs of all the sentinels
// These URLs point to the Pod specific Service of each sentinel Pod
func (gen *Generator) SentinelEndpoints(replicas int) []string {
	urls := make([]string, 0, replicas)
	for idx := 0; idx < int(replicas); idx++ {
		urls = append(urls,
			fmt.Sprintf("redis://%s.%s.svc.cluster.local:%d", gen.PodServiceName(idx), gen.GetNamespace(), saasv1alpha1.SentinelPort))
	}
	return urls
}
