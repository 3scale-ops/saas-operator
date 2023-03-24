package workloads

import (
	"testing"

	basereconciler_resources "github.com/3scale-ops/basereconciler/resources"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/util"
	"github.com/go-test/deep"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

// TEST GENERATORS
type TestWorkloadGenerator struct {
	TName            string
	TNamespace       string
	TTraffic         bool
	TLabels          map[string]string
	TSelector        map[string]string
	TTrafficSelector map[string]string
}

var _ DeploymentWorkload = &TestWorkloadGenerator{}
var _ WithTraffic = &TestWorkloadGenerator{}

func (gen *TestWorkloadGenerator) Deployment() basereconciler_resources.DeploymentTemplate {
	return basereconciler_resources.DeploymentTemplate{
		Template: func() *appsv1.Deployment {
			return &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Replicas: pointer.Int32Ptr(1),
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"orig-key": "orig-value"},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:      "container",
									Image:     "example.com:latest",
									Resources: corev1.ResourceRequirements{},
								},
							},
						},
					},
				},
			}
		},
		RolloutTriggers: []basereconciler_resources.RolloutTrigger{{
			Name:       "secret",
			SecretName: pointer.String("secret"),
		}},
		IsEnabled:       true,
		EnforceReplicas: true,
	}
}
func (gen *TestWorkloadGenerator) MonitoredEndpoints() []monitoringv1.PodMetricsEndpoint { return nil }
func (gen *TestWorkloadGenerator) GetKey() types.NamespacedName {
	return types.NamespacedName{Name: gen.TName, Namespace: gen.TNamespace}
}
func (gen *TestWorkloadGenerator) GetLabels() map[string]string { return gen.TLabels }
func (gen *TestWorkloadGenerator) GetSelector() map[string]string {
	return gen.TSelector
}
func (gen *TestWorkloadGenerator) HPASpec() *saasv1alpha1.HorizontalPodAutoscalerSpec {
	return &saasv1alpha1.HorizontalPodAutoscalerSpec{
		MinReplicas:         pointer.Int32Ptr(1),
		MaxReplicas:         pointer.Int32Ptr(2),
		ResourceUtilization: pointer.Int32Ptr(90),
		ResourceName:        pointer.StringPtr("cpu"),
	}
}
func (gen *TestWorkloadGenerator) PDBSpec() *saasv1alpha1.PodDisruptionBudgetSpec {
	return &saasv1alpha1.PodDisruptionBudgetSpec{
		MaxUnavailable: util.IntStrPtr(intstr.FromInt(1)),
	}
}
func (gen *TestWorkloadGenerator) SendTraffic() bool { return gen.TTraffic }

func (gen *TestWorkloadGenerator) Services() []basereconciler_resources.ServiceTemplate {
	return []basereconciler_resources.ServiceTemplate{{
		Template: func() *corev1.Service {
			return &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "service",
					Namespace: gen.TNamespace,
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{{
						Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP}},
				},
			}
		},
		IsEnabled: true,
	}}
}

func (gen *TestWorkloadGenerator) TrafficSelector() map[string]string {
	return gen.TTrafficSelector
}

// TESTS START HERE
func TestDeploymentTemplate_ApplyMeta(t *testing.T) {
	type args struct {
		gen DeploymentWorkload
	}
	tests := []struct {
		name string
		dt   DeploymentTemplate
		args args
		want *appsv1.Deployment
	}{
		{
			name: "Applies meta to an empty Deployment",
			dt: DeploymentTemplate{
				DeploymentTemplate: basereconciler_resources.DeploymentTemplate{
					Template: func() *appsv1.Deployment {
						return &appsv1.Deployment{}
					},
					RolloutTriggers: []basereconciler_resources.RolloutTrigger{},
					EnforceReplicas: false,
					IsEnabled:       false,
				},
			},
			args: args{
				gen: &TestWorkloadGenerator{
					TName:      "test",
					TNamespace: "test",
					TTraffic:   false,
					TLabels:    map[string]string{"key": "value"},
					TSelector:  map[string]string{"skey": "svalue"},
				},
			},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "test",
					Labels:    map[string]string{"key": "value"}},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"skey": "svalue"},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"key": "value", "skey": "svalue"},
						},
					},
				},
			},
		},
		{
			name: "Applies meta keeping original tempplate meta",
			dt: DeploymentTemplate{
				DeploymentTemplate: basereconciler_resources.DeploymentTemplate{
					Template: func() *appsv1.Deployment {
						return &appsv1.Deployment{
							ObjectMeta: metav1.ObjectMeta{
								Labels: map[string]string{"okey": "ovalue"},
							},
						}
					},
					RolloutTriggers: []basereconciler_resources.RolloutTrigger{},
					EnforceReplicas: false,
					IsEnabled:       false,
				},
			},
			args: args{
				gen: &TestWorkloadGenerator{
					TName:      "test",
					TNamespace: "test",
					TTraffic:   false,
					TLabels:    map[string]string{"key": "value"},
					TSelector:  map[string]string{"skey": "svalue"},
				},
			},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "test",
					Labels:    map[string]string{"okey": "ovalue", "key": "value"}},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"skey": "svalue"},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"key": "value", "skey": "svalue"},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := deep.Equal(tt.dt.ApplyMeta(tt.args.gen).Template(), tt.want); len(diff) > 0 {
				t.Errorf("DeploymentTemplate.ApplyMeta() = diff %v", diff)
			}
		})
	}
}

func TestDeploymentTemplate_ApplyTrafficSelector(t *testing.T) {
	type args struct {
		wt WithTraffic
	}
	tests := []struct {
		name string
		dt   DeploymentTemplate
		args args
		want *appsv1.Deployment
	}{
		{
			name: "Applies the traffic selector to an empty Deployment",
			dt: DeploymentTemplate{
				DeploymentTemplate: basereconciler_resources.DeploymentTemplate{
					Template: func() *appsv1.Deployment {
						return &appsv1.Deployment{}
					},
					RolloutTriggers: []basereconciler_resources.RolloutTrigger{},
					EnforceReplicas: false,
					IsEnabled:       false,
				},
			},
			args: args{
				wt: &TestWorkloadGenerator{
					TTrafficSelector: map[string]string{"tskey": "tsvalue"},
				},
			},
			want: &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"tskey": "tsvalue"},
						},
					},
				},
			},
		},
		{
			name: "Adds the selector if the Deployment already has pod labels",
			dt: DeploymentTemplate{
				DeploymentTemplate: basereconciler_resources.DeploymentTemplate{
					Template: func() *appsv1.Deployment {
						return &appsv1.Deployment{
							Spec: appsv1.DeploymentSpec{
								Template: corev1.PodTemplateSpec{
									ObjectMeta: metav1.ObjectMeta{
										Labels: map[string]string{"xxx": "xxx"},
									},
								},
							},
						}
					},
					RolloutTriggers: []basereconciler_resources.RolloutTrigger{},
					EnforceReplicas: false,
					IsEnabled:       false,
				},
			},
			args: args{
				wt: &TestWorkloadGenerator{
					TTrafficSelector: map[string]string{"tskey": "tsvalue"},
				},
			},
			want: &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"xxx": "xxx", "tskey": "tsvalue"},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := deep.Equal(tt.dt.ApplyTrafficSelector(tt.args.wt).Template(), tt.want); len(diff) > 0 {
				t.Errorf("DeploymentTemplate.ApplyTrafficSelector() = diff %v", diff)
			}
		})
	}
}

func TestServiceTemplate_ApplyMeta(t *testing.T) {
	type args struct {
		wt WithTraffic
	}
	tests := []struct {
		name string
		st   ServiceTemplate
		args args
		want *corev1.Service
	}{
		{
			name: "Adds meta to an empty Service",
			st: ServiceTemplate{
				ServiceTemplate: basereconciler_resources.ServiceTemplate{
					Template: func() *corev1.Service {
						return &corev1.Service{}
					},
					IsEnabled: false,
				},
			},
			args: args{
				wt: &TestWorkloadGenerator{
					TNamespace: "ns",
					TLabels:    map[string]string{"key": "value"},
				},
			},
			want: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns",
					Labels:    map[string]string{"key": "value"},
				},
			},
		},
		{
			name: "Keeps the original Service labels and adds the new ones",
			st: ServiceTemplate{
				ServiceTemplate: basereconciler_resources.ServiceTemplate{
					Template: func() *corev1.Service {
						return &corev1.Service{ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"okey": "ovalue"},
						}}
					},
					IsEnabled: false,
				},
			},
			args: args{
				wt: &TestWorkloadGenerator{
					TNamespace: "ns",
					TLabels:    map[string]string{"key": "value"},
				},
			},
			want: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns",
					Labels:    map[string]string{"okey": "ovalue", "key": "value"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := deep.Equal(tt.st.ApplyMeta(tt.args.wt).Template(), tt.want); len(diff) > 0 {
				t.Errorf("ServiceTemplate.ApplyMeta() = diff %v", diff)
			}
		})
	}
}

func TestServiceTemplate_ApplyTrafficSelector(t *testing.T) {
	type args struct {
		main   WithTraffic
		canary WithTraffic
	}
	tests := []struct {
		name string
		st   ServiceTemplate
		args args
		want *corev1.Service
	}{
		{
			name: "Applies pod selector to Service (traffic to w1)",
			st: ServiceTemplate{
				ServiceTemplate: basereconciler_resources.ServiceTemplate{
					Template: func() *corev1.Service {
						return &corev1.Service{}
					},
					IsEnabled: false,
				},
			},
			args: args{
				main: &TestWorkloadGenerator{
					TName:            "w1",
					TNamespace:       "ns",
					TTraffic:         true,
					TLabels:          nil,
					TSelector:        map[string]string{"name": "w1"},
					TTrafficSelector: map[string]string{"aaa": "aaa"},
				},
				canary: &TestWorkloadGenerator{
					TName:            "w2",
					TNamespace:       "ns",
					TTraffic:         false,
					TLabels:          nil,
					TSelector:        map[string]string{"name": "w2"},
					TTrafficSelector: map[string]string{"aaa": "aaa"},
				},
			},
			want: &corev1.Service{
				Spec: corev1.ServiceSpec{
					Selector: map[string]string{
						"aaa":  "aaa",
						"name": "w1",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := deep.Equal(tt.st.ApplyTrafficSelector(tt.args.main, tt.args.canary).Template(), tt.want); len(diff) > 0 {
				t.Errorf("ServiceTemplate.ApplyTrafficSelector() = diff %v", diff)
			}
		})
	}
}

func Test_trafficSwitcher(t *testing.T) {
	type args struct {
		main   TestWorkloadGenerator
		canary TestWorkloadGenerator
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Returns selector for a single Deployment",
			args: args{
				main: TestWorkloadGenerator{
					TTraffic:         true,
					TSelector:        map[string]string{"selector": "main"},
					TTrafficSelector: map[string]string{"traffic": "yes"},
				},
				canary: TestWorkloadGenerator{
					TTraffic:         false,
					TSelector:        map[string]string{"selector": "canary"},
					TTrafficSelector: map[string]string{"traffic": "yes"},
				},
			},
			want: map[string]string{"selector": "main", "traffic": "yes"},
		},
		{
			name: "Returns selector for all Deployments",
			args: args{
				main: TestWorkloadGenerator{
					TTraffic:         true,
					TSelector:        map[string]string{"selector": "main"},
					TTrafficSelector: map[string]string{"traffic": "yes"},
				},
				canary: TestWorkloadGenerator{
					TTraffic:         true,
					TSelector:        map[string]string{"selector": "canary"},
					TTrafficSelector: map[string]string{"traffic": "yes"},
				},
			},
			want: map[string]string{"traffic": "yes"},
		},
		{
			name: "Returns an empty map",
			args: args{
				main: TestWorkloadGenerator{
					TTraffic:         false,
					TSelector:        map[string]string{"selector": "main"},
					TTrafficSelector: map[string]string{"traffic": "yes"},
				},
				canary: TestWorkloadGenerator{
					TTraffic:         false,
					TSelector:        map[string]string{"selector": "canary"},
					TTrafficSelector: map[string]string{"traffic": "yes"},
				},
			},
			want: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := deep.Equal(trafficSwitcher(&tt.args.main, &tt.args.canary), tt.want); len(diff) > 0 {
				t.Errorf("trafficSwitcher() = diff %v", diff)
			}
		})
	}
}

func TestPodDisruptionBudgetTemplate_ApplyMeta(t *testing.T) {
	type args struct {
		w WithWorkloadMeta
	}
	tests := []struct {
		name string
		pdbt PodDisruptionBudgetTemplate
		args args
		want *policyv1.PodDisruptionBudget
	}{
		{
			name: "Applies meta to PDB",
			pdbt: PodDisruptionBudgetTemplate{
				PodDisruptionBudgetTemplate: basereconciler_resources.PodDisruptionBudgetTemplate{
					Template: func() *policyv1.PodDisruptionBudget {
						return &policyv1.PodDisruptionBudget{}
					},
					IsEnabled: false,
				},
			},
			args: args{
				w: &TestWorkloadGenerator{
					TName:      "test",
					TNamespace: "ns",
					TTraffic:   false,
					TLabels:    map[string]string{"key": "value"},
					TSelector:  map[string]string{"skey": "svalue"},
				},
			},
			want: &policyv1.PodDisruptionBudget{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "ns",
					Labels:    map[string]string{"key": "value"},
				},
				Spec: policyv1.PodDisruptionBudgetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"skey": "svalue"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := deep.Equal(tt.pdbt.ApplyMeta(tt.args.w).Template(), tt.want); len(diff) > 0 {
				t.Errorf("PodDisruptionBudgetTemplate.ApplyMeta() = diff %v", diff)
			}
		})
	}
}

func TestHorizontalPodAutoscalerTemplate_ApplyMeta(t *testing.T) {
	type args struct {
		w WithWorkloadMeta
	}
	tests := []struct {
		name string
		hpat HorizontalPodAutoscalerTemplate
		args args
		want *autoscalingv2.HorizontalPodAutoscaler
	}{
		{
			name: "Adds meta to HPA",
			hpat: HorizontalPodAutoscalerTemplate{
				HorizontalPodAutoscalerTemplate: basereconciler_resources.HorizontalPodAutoscalerTemplate{
					Template: func() *autoscalingv2.HorizontalPodAutoscaler {
						return &autoscalingv2.HorizontalPodAutoscaler{}
					},
					IsEnabled: false,
				},
			},
			args: args{
				w: &TestWorkloadGenerator{
					TName:      "test",
					TNamespace: "ns",
					TTraffic:   false,
					TLabels:    map[string]string{"key": "value"},
					TSelector:  nil,
				},
			},
			want: &autoscalingv2.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "ns",
					Labels:    map[string]string{"key": "value"},
				},
				Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
						Kind:       "Deployment",
						Name:       "test",
						APIVersion: "apps/v1",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := deep.Equal(tt.hpat.ApplyMeta(tt.args.w).Template(), tt.want); len(diff) > 0 {
				t.Errorf("HorizontalPodAutoscalerTemplate.ApplyMeta() = diff %v", diff)
			}
		})
	}
}

func TestPodMonitorTemplate_ApplyMeta(t *testing.T) {
	type args struct {
		w WithWorkloadMeta
	}
	tests := []struct {
		name string
		pmt  PodMonitorTemplate
		args args
		want *monitoringv1.PodMonitor
	}{
		{
			name: "Apply meta to PodMonitor",
			pmt: PodMonitorTemplate{
				PodMonitorTemplate: basereconciler_resources.PodMonitorTemplate{
					Template: func() *monitoringv1.PodMonitor {
						return &monitoringv1.PodMonitor{}
					},
					IsEnabled: false,
				},
			},
			args: args{
				w: &TestWorkloadGenerator{
					TName:      "test",
					TNamespace: "ns",
					TTraffic:   false,
					TLabels:    map[string]string{"key": "value"},
					TSelector:  nil,
				},
			},
			want: &monitoringv1.PodMonitor{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "ns",
					Labels:    map[string]string{"key": "value"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := deep.Equal(tt.pmt.ApplyMeta(tt.args.w).Template(), tt.want); len(diff) > 0 {
				t.Errorf("PodMonitorTemplate.ApplyMeta() = diff %v", diff)
			}
		})
	}
}
