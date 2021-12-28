package deployment

import (
	"testing"

	basereconciler_types "github.com/3scale/saas-operator/pkg/basereconciler/types"
	"github.com/go-test/deep"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestNew(t *testing.T) {
	type args struct {
		key             types.NamespacedName
		labels          map[string]string
		selector        map[string]string
		trafficSelector map[string]string
		fn              basereconciler_types.GeneratorFunction
	}
	tests := []struct {
		name string
		args args
		want *appsv1.Deployment
	}{
		{
			name: "Generates a Deployment with the proper labels",
			args: args{
				key: types.NamespacedName{
					Name:      "test",
					Namespace: "ns",
				},
				labels:          map[string]string{"lkey1": "lvalue1", "lkey2": "lvalue2"},
				selector:        map[string]string{"skey": "svalue"},
				trafficSelector: map[string]string{"tkey": "tvalue"},
				fn: func() client.Object {
					return &appsv1.Deployment{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{"akey": "avalue"},
						},
						Spec: appsv1.DeploymentSpec{
							Replicas: pointer.Int32(3),
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{ /*...*/ },
							},
						},
					}
				},
			},
			want: &appsv1.Deployment{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Deployment",
					APIVersion: appsv1.SchemeGroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test",
					Namespace:   "ns",
					Labels:      map[string]string{"lkey1": "lvalue1", "lkey2": "lvalue2"},
					Annotations: map[string]string{"akey": "avalue"},
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: pointer.Int32(3),
					Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"skey": "svalue"}},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"lkey1": "lvalue1",
								"lkey2": "lvalue2",
								"skey":  "svalue",
								"tkey":  "tvalue",
							},
						},
						Spec: corev1.PodSpec{ /*...*/ },
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.key, tt.args.labels, tt.args.selector, tt.args.trafficSelector, tt.args.fn)().(*appsv1.Deployment)
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("New() got diff: %v", diff)
			}
		})
	}
}
