package basereconciler

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func Test_isManaged(t *testing.T) {
	type args struct {
		key     types.NamespacedName
		kind    string
		managed []corev1.ObjectReference
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Returns true",
			args: args{
				key:  types.NamespacedName{Name: "system-recaptcha", Namespace: "ns"},
				kind: "Secret",
				managed: []corev1.ObjectReference{
					{Namespace: "ns", Name: "system-recaptcha", Kind: "Secret"},
					{Namespace: "ns", Name: "system-smtp", Kind: "Secret"},
					{Namespace: "ns", Name: "system-zync", Kind: "Secret"},
					{Namespace: "ns", Name: "system", Kind: "Secret"},
					{Namespace: "ns", Name: "system-app", Kind: "Deployment"},
					{Namespace: "ns", Name: "system-app", Kind: "ServiceAccount"},
					{Namespace: "ns", Name: "system-app", Kind: "HorizontalPodAutoscaler"},
					{Namespace: "ns", Name: "system-app", Kind: "PodDisruptionBudget"},
					{Namespace: "ns", Name: "system-app", Kind: "PodMonitor"},
				},
			},
			want: true,
		},
		{
			name: "Returns false",
			args: args{
				key:  types.NamespacedName{Name: "system-recaptcha", Namespace: "ns"},
				kind: "Secret",
				managed: []corev1.ObjectReference{
					{Namespace: "ns", Name: "system-smtp", Kind: "Secret"},
					{Namespace: "ns", Name: "system-zync", Kind: "Secret"},
					{Namespace: "ns", Name: "system", Kind: "Secret"},
					{Namespace: "ns", Name: "system-app", Kind: "Deployment"},
					{Namespace: "ns", Name: "system-app", Kind: "ServiceAccount"},
					{Namespace: "ns", Name: "system-app", Kind: "HorizontalPodAutoscaler"},
					{Namespace: "ns", Name: "system-app", Kind: "PodDisruptionBudget"},
					{Namespace: "ns", Name: "system-app", Kind: "PodMonitor"},
				},
			},
			want: false,
		},
		{
			name: "Returns false",
			args: args{
				key:  types.NamespacedName{Name: "system-app", Namespace: "ns"},
				kind: "Role",
				managed: []corev1.ObjectReference{
					{Namespace: "ns", Name: "system-smtp", Kind: "Secret"},
					{Namespace: "ns", Name: "system-zync", Kind: "Secret"},
					{Namespace: "ns", Name: "system", Kind: "Secret"},
					{Namespace: "ns", Name: "system-app", Kind: "Deployment"},
					{Namespace: "ns", Name: "system-app", Kind: "ServiceAccount"},
					{Namespace: "ns", Name: "system-app", Kind: "HorizontalPodAutoscaler"},
					{Namespace: "ns", Name: "system-app", Kind: "PodDisruptionBudget"},
					{Namespace: "ns", Name: "system-app", Kind: "PodMonitor"},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isManaged(tt.args.key, tt.args.kind, tt.args.managed); got != tt.want {
				t.Errorf("isManaged() = %v, want %v", got, tt.want)
			}
		})
	}
}
