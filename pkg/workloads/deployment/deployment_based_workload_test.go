package deployment

import (
	"context"
	"reflect"
	"testing"

	"github.com/3scale-ops/basereconciler/mutators"
	"github.com/3scale-ops/basereconciler/resource"
	"github.com/3scale-ops/basereconciler/util"
	marin3rv1alpha1 "github.com/3scale-ops/marin3r/apis/marin3r/v1alpha1"
	"github.com/3scale-ops/marin3r/pkg/envoy"
	envoy_serializer "github.com/3scale-ops/marin3r/pkg/envoy/serializer"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	descriptor "github.com/3scale-ops/saas-operator/pkg/resource_builders/envoyconfig/descriptor"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/service"
	"github.com/google/go-cmp/cmp"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
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
var _ WithMarin3rSidecar = &TestWorkloadGenerator{}

func (gen *TestWorkloadGenerator) Deployment() *resource.Template[*appsv1.Deployment] {
	return &resource.Template[*appsv1.Deployment]{
		TemplateBuilder: func(client.Object) (*appsv1.Deployment, error) {
			return &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Replicas: util.Pointer[int32](1),
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
			}, nil
		},
		TemplateMutations: []resource.TemplateMutationFunction{
			mutators.RolloutTrigger{
				Name:       "secret",
				SecretName: util.Pointer("secret"),
			}.Add(),
			mutators.SetDeploymentReplicas(true),
		},
		IsEnabled: true,
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
		MinReplicas:         util.Pointer[int32](1),
		MaxReplicas:         util.Pointer[int32](2),
		ResourceUtilization: util.Pointer[int32](90),
		ResourceName:        util.Pointer("cpu"),
	}
}
func (gen *TestWorkloadGenerator) PDBSpec() *saasv1alpha1.PodDisruptionBudgetSpec {
	return &saasv1alpha1.PodDisruptionBudgetSpec{
		MaxUnavailable: util.Pointer(intstr.FromInt(1)),
	}
}
func (gen *TestWorkloadGenerator) SendTraffic() bool { return gen.TTraffic }

func (gen *TestWorkloadGenerator) Services() []*resource.Template[*corev1.Service] {
	return []*resource.Template[*corev1.Service]{{
		TemplateBuilder: func(client.Object) (*corev1.Service, error) {
			return &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "service",
					Namespace: gen.TNamespace,
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{{
						Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP}},
				},
			}, nil
		},
		IsEnabled: true,
	}}
}

func (gen *TestWorkloadGenerator) TrafficSelector() map[string]string {
	return gen.TTrafficSelector
}

func (gen *TestWorkloadGenerator) EnvoyDynamicConfigurations() []descriptor.EnvoyDynamicConfigDescriptor {
	return nil
}

// // TESTS START HERE

func TestWorkloadReconciler_NewDeploymentWorkload(t *testing.T) {
	type args struct {
		main   DeploymentWorkload
		canary DeploymentWorkload
	}
	tests := []struct {
		name    string
		args    args
		want    []client.Object
		wantErr bool
	}{
		{
			name: "Generates the workload resources",
			args: args{
				main: &TestWorkloadGenerator{
					TName:            "my-workload",
					TNamespace:       "ns",
					TTraffic:         true,
					TLabels:          map[string]string{"l-key": "l-value"},
					TSelector:        map[string]string{"sel-key": "sel-value"},
					TTrafficSelector: map[string]string{"traffic": "yes"},
				},
				canary: nil,
			},
			want: []client.Object{
				&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name: "my-workload", Namespace: "ns",
						Labels: map[string]string{"l-key": "l-value"}},
					Spec: appsv1.DeploymentSpec{
						Replicas: util.Pointer[int32](1),
						Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"sel-key": "sel-value"}},
						Template: corev1.PodTemplateSpec{
							ObjectMeta: metav1.ObjectMeta{
								Labels: map[string]string{
									"orig-key": "orig-value",
									"l-key":    "l-value",
									"sel-key":  "sel-value",
									"traffic":  "yes",
								},
								Annotations: map[string]string{"basereconciler.3cale.net/secret.secret-hash": ""},
							},
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{{
									Name:      "container",
									Image:     "example.com:latest",
									Resources: corev1.ResourceRequirements{},
								}}}}}},
				&policyv1.PodDisruptionBudget{
					ObjectMeta: metav1.ObjectMeta{
						Name: "my-workload", Namespace: "ns",
						Labels: map[string]string{"l-key": "l-value"},
					},
					Spec: policyv1.PodDisruptionBudgetSpec{
						Selector:       &metav1.LabelSelector{MatchLabels: map[string]string{"sel-key": "sel-value"}},
						MaxUnavailable: util.Pointer(intstr.FromInt(1)),
					}},
				&autoscalingv2.HorizontalPodAutoscaler{
					ObjectMeta: metav1.ObjectMeta{
						Name: "my-workload", Namespace: "ns",
						Labels: map[string]string{"l-key": "l-value"},
					},
					Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
						ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
							APIVersion: appsv1.SchemeGroupVersion.String(),
							Kind:       "Deployment",
							Name:       "my-workload",
						},
						MinReplicas: util.Pointer[int32](1),
						MaxReplicas: 2,
						Metrics: []autoscalingv2.MetricSpec{{
							Type: autoscalingv2.ResourceMetricSourceType,
							Resource: &autoscalingv2.ResourceMetricSource{
								Name: corev1.ResourceName("cpu"),
								Target: autoscalingv2.MetricTarget{
									Type:               autoscalingv2.UtilizationMetricType,
									AverageUtilization: util.Pointer[int32](90),
								}}}}}},
				&monitoringv1.PodMonitor{
					ObjectMeta: metav1.ObjectMeta{
						Name: "my-workload", Namespace: "ns",
						Labels: map[string]string{"l-key": "l-value"},
					},
					Spec: monitoringv1.PodMonitorSpec{
						PodMetricsEndpoints: nil,
						Selector: metav1.LabelSelector{
							MatchLabels: map[string]string{"sel-key": "sel-value"},
						},
					}},
				&marin3rv1alpha1.EnvoyConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name: "my-workload", Namespace: "ns",
						Labels: map[string]string{"l-key": "l-value"},
					},
					Spec: marin3rv1alpha1.EnvoyConfigSpec{
						NodeID:        "my-workload",
						Serialization: util.Pointer[envoy_serializer.Serialization]("yaml"),
						EnvoyAPI:      util.Pointer[envoy.APIVersion]("v3"),
						EnvoyResources: &marin3rv1alpha1.EnvoyResources{
							Clusters:  []marin3rv1alpha1.EnvoyResource{},
							Routes:    []marin3rv1alpha1.EnvoyResource{},
							Listeners: []marin3rv1alpha1.EnvoyResource{},
							Runtimes:  []marin3rv1alpha1.EnvoyResource{},
							Secrets:   []marin3rv1alpha1.EnvoySecretResource{},
						}}},
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name: "service", Namespace: "ns",
						Labels: map[string]string{"l-key": "l-value"},
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{"sel-key": "sel-value", "traffic": "yes"},
						Ports: []corev1.ServicePort{{
							Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP}},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templates, err := New(tt.args.main, tt.args.canary)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := make([]client.Object, 0, len(templates))
			for _, tpl := range templates {
				o, _ := tpl.Build(context.TODO(), fake.NewClientBuilder().Build(), nil)
				got = append(got, o)
			}
			if diff := cmp.Diff(got, tt.want, util.IgnoreProperty("Status")); len(diff) > 0 {
				t.Errorf("New() diff %v", diff)
			}
		})
	}
}

func Test_applyTrafficSelectorToDeployment(t *testing.T) {
	type args struct {
		w DeploymentWorkload
	}
	tests := []struct {
		name     string
		template *resource.Template[*appsv1.Deployment]
		args     args
		want     *appsv1.Deployment
	}{
		{
			name: "Applies the traffic selector to a Deployment",
			template: resource.NewTemplate[*appsv1.Deployment](
				func(client.Object) (*appsv1.Deployment, error) {
					return &appsv1.Deployment{}, nil
				}),
			args: args{
				w: &TestWorkloadGenerator{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.template.Apply(trafficSelectorToDeployment(tt.args.w)).Build(context.TODO(), nil, nil)
			if diff := cmp.Diff(got, tt.want); len(diff) > 0 {
				t.Errorf("applyTrafficSelectorToDeployment() got diff %v", diff)
			}
		})
	}
}

func Test_applyHPAScaleTargetRef(t *testing.T) {
	type args struct {
		w WithWorkloadMeta
	}
	tests := []struct {
		name     string
		template *resource.Template[*autoscalingv2.HorizontalPodAutoscaler]
		args     args
		want     *autoscalingv2.HorizontalPodAutoscaler
	}{
		{
			name: "Adds ScaleTargetRef to HPA",
			template: resource.NewTemplate[*autoscalingv2.HorizontalPodAutoscaler](
				func(client.Object) (*autoscalingv2.HorizontalPodAutoscaler, error) {
					return &autoscalingv2.HorizontalPodAutoscaler{}, nil
				}),
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
			got, _ := tt.template.Apply(scaleTargetRefToHPA(tt.args.w)).Build(context.TODO(), nil, nil)
			if diff := cmp.Diff(got, tt.want); len(diff) > 0 {
				t.Errorf("applyHPAScaleTargetRef() got diff %v", diff)
			}
		})
	}
}

func Test_applyMeta(t *testing.T) {
	type args struct {
		w WithWorkloadMeta
	}
	tests1 := []struct {
		name     string
		template *resource.Template[*corev1.ConfigMap]
		args     args
		want     *corev1.ConfigMap
	}{
		{
			name: "Adds meta to a ConfigMap",
			template: resource.NewTemplate[*corev1.ConfigMap](
				func(client.Object) (*corev1.ConfigMap, error) { return &corev1.ConfigMap{}, nil },
			),
			args: args{
				w: &TestWorkloadGenerator{
					TName:      "cm",
					TNamespace: "ns",
					TLabels:    map[string]string{"key": "value"},
				},
			},
			want: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cm",
					Namespace: "ns",
					Labels:    map[string]string{"key": "value"},
				},
			},
		},
		{
			name: "Keeps the original ConfigMap labels and adds the new ones",
			template: resource.NewTemplate[*corev1.ConfigMap](
				func(client.Object) (*corev1.ConfigMap, error) {
					return &corev1.ConfigMap{
						ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"okey": "ovalue"}},
					}, nil
				},
			),
			args: args{
				w: &TestWorkloadGenerator{
					TName:      "cm",
					TNamespace: "ns",
					TLabels:    map[string]string{"key": "value"},
				},
			},
			want: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cm",
					Namespace: "ns",
					Labels:    map[string]string{"okey": "ovalue", "key": "value"},
				},
			},
		},
	}
	for _, tt := range tests1 {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.template.Apply(meta[*corev1.ConfigMap](tt.args.w)).Build(context.TODO(), nil, nil)
			if diff := cmp.Diff(got, tt.want); len(diff) > 0 {
				t.Errorf("applyMeta() got diff %v", diff)
			}
		})
	}

	tests2 := []struct {
		name     string
		template *resource.Template[*corev1.Service]
		args     args
		want     *corev1.Service
	}{
		{
			name: "Adds meta to a Service preserving the name",
			template: resource.NewTemplate[*corev1.Service](
				func(client.Object) (*corev1.Service, error) {
					return &corev1.Service{
						ObjectMeta: metav1.ObjectMeta{Name: "my-service"},
					}, nil
				},
			),
			args: args{
				w: &TestWorkloadGenerator{
					TName:      "name",
					TNamespace: "ns",
					TLabels:    map[string]string{"key": "value"},
				},
			},
			want: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-service",
					Namespace: "ns",
					Labels:    map[string]string{"key": "value"},
				},
			},
		},
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.template.Apply(meta[*corev1.Service](tt.args.w)).Build(context.TODO(), nil, nil)
			if diff := cmp.Diff(got, tt.want); len(diff) > 0 {
				t.Errorf("applyMeta() got diff %v", diff)
			}
		})
	}

}

func Test_applyTrafficSelectorToService(t *testing.T) {
	type args struct {
		main   WithTraffic
		canary WithTraffic
	}
	tests := []struct {
		name     string
		template *resource.Template[*corev1.Service]
		args     args
		want     *corev1.Service
	}{
		{
			name: "Applies pod selector to Service (traffic to w1)",
			template: resource.NewTemplate[*corev1.Service](
				func(client.Object) (*corev1.Service, error) { return &corev1.Service{}, nil },
			),
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
			got, _ := tt.template.Apply(trafficSelectorToService(tt.args.main, tt.args.canary)).Build(context.TODO(), nil, nil)
			if diff := cmp.Diff(got, tt.want); len(diff) > 0 {
				t.Errorf("applyTrafficSelectorToService() got diff %v", diff)
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
			if diff := cmp.Diff(trafficSwitcher(&tt.args.main, &tt.args.canary), tt.want); len(diff) > 0 {
				t.Errorf("trafficSwitcher() = diff %v", diff)
			}
		})
	}
}

func Test_applyNodeIdToEnvoyConfig(t *testing.T) {
	type args struct {
		sd service.ServiceDescriptor
	}
	tests := []struct {
		name     string
		template *resource.Template[*marin3rv1alpha1.EnvoyConfig]
		args     args
		want     *marin3rv1alpha1.EnvoyConfig
	}{
		{
			name: "Adds the nodeID",
			template: resource.NewTemplate(
				func(client.Object) (*marin3rv1alpha1.EnvoyConfig, error) {
					return &marin3rv1alpha1.EnvoyConfig{}, nil
				},
			),
			args: args{
				sd: service.ServiceDescriptor{
					PortDefinition: corev1.ServicePort{},
					PublishingStrategy: saasv1alpha1.PublishingStrategy{
						Marin3rSidecar: &saasv1alpha1.Marin3rSidecarSpec{
							NodeID: util.Pointer("aaaa"),
						},
					},
				},
			},
			want: &marin3rv1alpha1.EnvoyConfig{
				Spec: marin3rv1alpha1.EnvoyConfigSpec{NodeID: "aaaa"},
			},
		},
		{
			name: "Adds the resource name as the nodeID when unset",
			template: resource.NewTemplate(
				func(client.Object) (*marin3rv1alpha1.EnvoyConfig, error) {
					return &marin3rv1alpha1.EnvoyConfig{ObjectMeta: metav1.ObjectMeta{Name: "bbbb"}}, nil
				},
			),
			args: args{
				sd: service.ServiceDescriptor{
					PortDefinition: corev1.ServicePort{},
					PublishingStrategy: saasv1alpha1.PublishingStrategy{
						Marin3rSidecar: &saasv1alpha1.Marin3rSidecarSpec{},
					},
				},
			},
			want: &marin3rv1alpha1.EnvoyConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "bbbb"},
				Spec:       marin3rv1alpha1.EnvoyConfigSpec{NodeID: "bbbb"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.template.Apply(nodeIdToEnvoyConfig(tt.args.sd)).Build(context.TODO(), nil, nil)
			if diff := cmp.Diff(got, tt.want); len(diff) > 0 {
				t.Errorf("applyNodeIdToEnvoyConfig() got diff %v", diff)
			}
		})
	}
}

func Test_toWithTraffic(t *testing.T) {
	type args struct {
		w DeploymentWorkload
	}
	tests := []struct {
		name string
		args args
		want WithTraffic
	}{
		{
			name: "Detects a nil value",
			args: args{
				w: nil,
			},
			want: nil,
		},
		{
			name: "Detects interface containing nil value",
			args: args{
				w: func() DeploymentWorkload {
					val := (*TestWorkloadGenerator)(nil)
					return val
				}(),
			},
			want: nil,
		},
		{
			name: "Converts DeploymentWorkload to WithTraffic",
			args: args{
				w: &TestWorkloadGenerator{},
			},
			want: &TestWorkloadGenerator{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toWithTraffic(tt.args.w); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toWithTraffic() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_marin3rSidecarToDeployment(t *testing.T) {
	type args struct {
		sd service.ServiceDescriptor
	}
	tests := []struct {
		name     string
		template *resource.Template[*appsv1.Deployment]
		args     args
		want     *appsv1.Deployment
	}{
		{
			name: "Adds the marin3r annotations for sidecar injection",
			template: resource.NewTemplateFromObjectFunction(
				func() *appsv1.Deployment {
					return &appsv1.Deployment{
						ObjectMeta: metav1.ObjectMeta{Name: "test"},
						Spec: appsv1.DeploymentSpec{
							Template: corev1.PodTemplateSpec{
								ObjectMeta: metav1.ObjectMeta{},
							},
						},
					}
				},
			),
			args: args{
				sd: service.ServiceDescriptor{
					PublishingStrategy: saasv1alpha1.PublishingStrategy{
						Strategy:     saasv1alpha1.Marin3rSidecarStrategy,
						EndpointName: "Endpoint",
						Marin3rSidecar: &saasv1alpha1.Marin3rSidecarSpec{
							Ports:                              []saasv1alpha1.SidecarPort{{Name: "port", Port: 8888}},
							ShutdownManagerExtraLifecycleHooks: []string{"container"},
						},
					},
				},
			},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"marin3r.3scale.net/status": "enabled",
							},
							Annotations: map[string]string{
								"marin3r.3scale.net/node-id":                                "test",
								"marin3r.3scale.net/ports":                                  "port:8888",
								"marin3r.3scale.net/shutdown-manager.enabled":               "true",
								"marin3r.3scale.net/shutdown-manager.extra-lifecycle-hooks": "container",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.template.Apply(marin3rSidecarToDeployment(tt.args.sd)).Build(context.TODO(), nil, nil)
			if diff := cmp.Diff(got, tt.want); len(diff) > 0 {
				t.Errorf("marin3rSidecarToDeployment() got diff %v", diff)
			}
		})
	}
}
