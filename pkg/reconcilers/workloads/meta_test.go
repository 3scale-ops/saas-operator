package workloads

import (
	"testing"

	"github.com/go-test/deep"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type TestWithMetaGenerator struct {
	TName      string
	TNamespace string
	TLabels    map[string]string
	TSelector  map[string]string
}

func (gen *TestWithMetaGenerator) GetKey() types.NamespacedName {
	return types.NamespacedName{Name: gen.TName, Namespace: gen.TNamespace}
}
func (gen *TestWithMetaGenerator) GetLabels() map[string]string { return gen.TLabels }
func (gen *TestWithMetaGenerator) GetSelector() map[string]string {
	return gen.TSelector
}

func Test_applyKey(t *testing.T) {
	type args struct {
		o client.Object
		w WithKey
	}
	tests := []struct {
		name string
		args args
		want client.Object
	}{
		{
			name: "Sets resource Name and Namespace",
			args: args{
				o: &corev1.Service{},
				w: &TestWithMetaGenerator{
					TName:      "name",
					TNamespace: "ns",
					TLabels:    map[string]string{"key": "value"},
					TSelector:  nil,
				},
			},
			want: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "name",
					Namespace: "ns",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applyKey(tt.args.o, tt.args.w)
			if diff := deep.Equal(tt.args.o, tt.want); len(diff) > 0 {
				t.Errorf("applyKey() = diff %v", diff)
			}
		})
	}
}

func Test_applyLabels(t *testing.T) {
	type args struct {
		o client.Object
		w WithLabels
	}
	tests := []struct {
		name string
		args args
		want client.Object
	}{
		{
			name: "Applies labels to a resource without labels",
			args: args{
				o: &corev1.Service{},
				w: &TestWithMetaGenerator{
					TName:      "name",
					TNamespace: "ns",
					TLabels:    map[string]string{"key": "value"},
					TSelector:  nil,
				},
			},
			want: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"key": "value"},
				},
			},
		},
		{
			name: "Merges labels with the resource's original labels",
			args: args{
				o: &corev1.Service{
					ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"okey": "ovalue"}},
				},
				w: &TestWithMetaGenerator{
					TName:      "name",
					TNamespace: "ns",
					TLabels:    map[string]string{"key": "value"},
					TSelector:  nil,
				},
			},
			want: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"okey": "ovalue", "key": "value"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applyLabels(tt.args.o, tt.args.w)
			if diff := deep.Equal(tt.args.o, tt.want); len(diff) > 0 {
				t.Errorf("applyLabels() = diff %v", diff)
			}
		})
	}
}
