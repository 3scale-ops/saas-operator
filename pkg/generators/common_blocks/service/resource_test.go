package service

import (
	"testing"

	basereconciler_types "github.com/3scale/saas-operator/pkg/basereconciler/types"
	"github.com/go-test/deep"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestNew(t *testing.T) {
	type args struct {
		labels   map[string]string
		selector map[string]string
		fn       basereconciler_types.GeneratorFunction
	}
	tests := []struct {
		name string
		args args
		want *corev1.Service
	}{
		{
			name: "Returns a Service",
			args: args{
				labels:   map[string]string{"key": "value"},
				selector: map[string]string{"skey": "svalue"},
				fn: func() client.Object {
					return &corev1.Service{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "test",
							Namespace:   "ns",
							Annotations: map[string]string{"akey": "avalue"},
						},
						Spec: corev1.ServiceSpec{
							Ports: []corev1.ServicePort{{
								Name:       "port",
								Protocol:   "TCP",
								Port:       8888,
								TargetPort: intstr.FromInt(8888),
							}},
						},
					}
				},
			},
			want: &corev1.Service{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Service",
					APIVersion: corev1.SchemeGroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test",
					Namespace:   "ns",
					Labels:      map[string]string{"key": "value"},
					Annotations: map[string]string{"akey": "avalue"},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{{
						Name:       "port",
						Protocol:   "TCP",
						Port:       8888,
						TargetPort: intstr.FromInt(8888),
					}},
					Selector: map[string]string{"skey": "svalue"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.labels, tt.args.selector, tt.args.fn)().(*corev1.Service)
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("New() got diff: %v", diff)
			}
		})
	}
}
