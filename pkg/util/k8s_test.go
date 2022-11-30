package util

import (
	"reflect"
	"testing"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestGetItems(t *testing.T) {
	type args struct {
		list client.ObjectList
	}
	tests := []struct {
		name string
		args args
		want []client.Object
	}{
		{
			name: "Returns items of a corev1.ServiceList as []client.Object",
			args: args{
				list: &corev1.ServiceList{
					Items: []corev1.Service{
						{ObjectMeta: metav1.ObjectMeta{Name: "one"}},
						{ObjectMeta: metav1.ObjectMeta{Name: "two"}},
					},
				},
			},
			want: []client.Object{
				&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "one"}},
				&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "two"}},
			},
		},
		{
			name: "Returns items of a monitoringv1.PodMonitorList as []client.Object",
			args: args{
				list: &monitoringv1.PodMonitorList{
					Items: []*monitoringv1.PodMonitor{
						{ObjectMeta: metav1.ObjectMeta{Name: "one"}},
						{ObjectMeta: metav1.ObjectMeta{Name: "two"}},
					},
				},
			},
			want: []client.Object{
				&monitoringv1.PodMonitor{ObjectMeta: metav1.ObjectMeta{Name: "one"}},
				&monitoringv1.PodMonitor{ObjectMeta: metav1.ObjectMeta{Name: "two"}},
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetItems(tt.args.list); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetItems() = %v, want %v", got, tt.want)
			}
		})
	}
}
