package basereconciler

import (
	"reflect"
	"testing"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	basereconciler_types "github.com/3scale/saas-operator/pkg/basereconciler/types"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type TestWorkloadGenerator struct {
	TSendTraffic bool
	TSelector    *metav1.LabelSelector
}

var _ basereconciler_types.DeploymentWorkloadGenerator = &TestWorkloadGenerator{}

func (gen *TestWorkloadGenerator) Key() types.NamespacedName                          { return types.NamespacedName{} }
func (gen *TestWorkloadGenerator) GetLabels() map[string]string                       { return nil }
func (gen *TestWorkloadGenerator) Selector() *metav1.LabelSelector                    { return gen.TSelector }
func (gen *TestWorkloadGenerator) Deployment() basereconciler_types.GeneratorFunction { return nil }
func (gen *TestWorkloadGenerator) HPASpec() *saasv1alpha1.HorizontalPodAutoscalerSpec {
	return &saasv1alpha1.HorizontalPodAutoscalerSpec{}
}
func (gen *TestWorkloadGenerator) PDBSpec() *saasv1alpha1.PodDisruptionBudgetSpec {
	return &saasv1alpha1.PodDisruptionBudgetSpec{}
}
func (gen *TestWorkloadGenerator) MonitoredEndpoints() []monitoringv1.PodMetricsEndpoint { return nil }
func (gen *TestWorkloadGenerator) RolloutTriggers() []basereconciler_types.GeneratorFunction {
	return nil
}
func (gen *TestWorkloadGenerator) SendTraffic() bool { return gen.TSendTraffic }

type TestIngressGenerator struct {
	TTrafficSelector map[string]string
}

var _ basereconciler_types.DeploymentIngressGenerator = &TestIngressGenerator{}

func (gen *TestIngressGenerator) GetLabels() map[string]string                       { return nil }
func (gen *TestIngressGenerator) Services() []basereconciler_types.GeneratorFunction { return nil }
func (gen *TestIngressGenerator) TrafficSelector() map[string]string                 { return gen.TTrafficSelector }

func TestTrafficSwitcher(t *testing.T) {
	type args struct {
		ingress   basereconciler_types.DeploymentIngressGenerator
		workloads []basereconciler_types.DeploymentWorkloadGenerator
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Returns selector for a single Deployment",
			args: args{
				ingress: &TestIngressGenerator{
					TTrafficSelector: map[string]string{"traffic": "yes"},
				},
				workloads: []basereconciler_types.DeploymentWorkloadGenerator{
					&TestWorkloadGenerator{
						TSendTraffic: true,
						TSelector:    &metav1.LabelSelector{MatchLabels: map[string]string{"selector": "dep1"}},
					},
					&TestWorkloadGenerator{
						TSendTraffic: false,
						TSelector:    &metav1.LabelSelector{MatchLabels: map[string]string{"selector": "dep2"}},
					},
				},
			},
			want: map[string]string{"selector": "dep1", "traffic": "yes"},
		},
		{
			name: "Returns selector for all Deployments",
			args: args{
				ingress: &TestIngressGenerator{
					TTrafficSelector: map[string]string{"traffic": "yes"},
				},
				workloads: []basereconciler_types.DeploymentWorkloadGenerator{
					&TestWorkloadGenerator{
						TSendTraffic: true,
						TSelector:    &metav1.LabelSelector{MatchLabels: map[string]string{"selector": "dep1"}},
					},
					&TestWorkloadGenerator{
						TSendTraffic: true,
						TSelector:    &metav1.LabelSelector{MatchLabels: map[string]string{"selector": "dep2"}},
					},
				},
			},
			want: map[string]string{"traffic": "yes"},
		},
		{
			name: "Returns an empty map",
			args: args{
				ingress: &TestIngressGenerator{
					TTrafficSelector: map[string]string{"traffic": "yes"},
				},
				workloads: []basereconciler_types.DeploymentWorkloadGenerator{
					&TestWorkloadGenerator{
						TSendTraffic: false,
						TSelector:    &metav1.LabelSelector{MatchLabels: map[string]string{"selector": "dep1"}},
					},
					&TestWorkloadGenerator{
						TSendTraffic: false,
						TSelector:    &metav1.LabelSelector{MatchLabels: map[string]string{"selector": "dep2"}},
					},
				},
			},
			want: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrafficSwitcher(tt.args.ingress, tt.args.workloads...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TrafficSwitcher() = %v, want %v", got, tt.want)
			}
		})
	}
}
