package marin3r

import (
	"reflect"
	"testing"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestEnableSidecar(t *testing.T) {
	type args struct {
		dep  appsv1.Deployment
		spec saasv1alpha1.Marin3rSidecarSpec
	}
	tests := []struct {
		name string
		args args
		want *appsv1.Deployment
	}{
		{
			name: "Adds marin3r labels and annotations to a Deployment",
			args: args{
				dep: appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{
						Template: corev1.PodTemplateSpec{
							ObjectMeta: metav1.ObjectMeta{},
						},
					},
				},
				spec: saasv1alpha1.Marin3rSidecarSpec{
					Ports: []saasv1alpha1.SidecarPort{
						{
							Name: "test",
							Port: 9999,
						},
					},
					Resources: &saasv1alpha1.ResourceRequirementsSpec{
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("200m"),
							corev1.ResourceMemory: resource.MustParse("200Mi"),
						},
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("100Mi"),
						},
					},
					ExtraPodAnnotations: map[string]string{
						"extra-key": "extra-value",
					},
				},
			},
			want: &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								"marin3r.3scale.net/resources.limits.cpu":      "200m",
								"marin3r.3scale.net/resources.limits.memory":   "200Mi",
								"marin3r.3scale.net/resources.requests.cpu":    "100m",
								"marin3r.3scale.net/resources.requests.memory": "100Mi",
								"marin3r.3scale.net/ports":                     "test:9999",
								"extra-key":                                    "extra-value",
							},
							Labels: map[string]string{
								"marin3r.3scale.net/status": "enabled",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EnableSidecar(tt.args.dep, tt.args.spec); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnableSidecar() = %v, want %v", got, tt.want)
			}
		})
	}
}
