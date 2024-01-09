package redisshard

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (gen *Generator) service() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gen.ServiceName(),
			Namespace: gen.GetNamespace(),
			Labels:    gen.GetLabels(),
		},
		Spec: corev1.ServiceSpec{
			Type:       corev1.ServiceTypeClusterIP,
			ClusterIP:  corev1.ClusterIPNone,
			ClusterIPs: []string{corev1.ClusterIPNone},
			Ports:      []corev1.ServicePort{},
			Selector:   gen.GetSelector(),
		},
	}
}
