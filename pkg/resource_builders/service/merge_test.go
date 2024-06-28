package service

import (
	"testing"

	"github.com/3scale-ops/basereconciler/util"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/google/go-cmp/cmp"
	"github.com/imdario/mergo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestMergeWithDefaultPublishingStrategy(t *testing.T) {
	type args struct {
		def []ServiceDescriptor
		in  *saasv1alpha1.PublishingStrategies
	}
	tests := []struct {
		name    string
		args    args
		want    []ServiceDescriptor
		wantErr bool
	}{
		{
			name: "Merge: changes the publishing strategy",
			args: args{
				def: []ServiceDescriptor{
					{
						PublishingStrategy: saasv1alpha1.PublishingStrategy{
							Strategy:     saasv1alpha1.SimpleStrategy,
							EndpointName: "Gateway",
							Simple: &saasv1alpha1.Simple{
								ServiceType: util.Pointer(saasv1alpha1.ServiceTypeELB),
								ElasticLoadBalancerConfig: &saasv1alpha1.ElasticLoadBalancerSpec{
									ProxyProtocol:                 util.Pointer(true),
									CrossZoneLoadBalancingEnabled: util.Pointer(true),
									ConnectionDrainingEnabled:     util.Pointer(true),
									ConnectionDrainingTimeout:     util.Pointer[int32](60),
									HealthcheckHealthyThreshold:   util.Pointer[int32](2),
									HealthcheckUnhealthyThreshold: util.Pointer[int32](2),
									HealthcheckInterval:           util.Pointer[int32](5),
									HealthcheckTimeout:            util.Pointer[int32](3),
								},
							},
						},
						PortDefinitions: []corev1.ServicePort{{
							Name:       "gateway-http",
							Protocol:   corev1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromString("gateway-http"),
						}},
					},
				},
				in: &saasv1alpha1.PublishingStrategies{
					Mode: util.Pointer(saasv1alpha1.PublishingStrategiesReconcileModeMerge),
					Endpoints: []saasv1alpha1.PublishingStrategy{{
						Strategy:     saasv1alpha1.Marin3rSidecarStrategy,
						EndpointName: "Gateway",
						Marin3rSidecar: &saasv1alpha1.Marin3rSidecarSpec{
							Simple: &saasv1alpha1.Simple{
								ServiceType: util.Pointer(saasv1alpha1.ServiceTypeELB),
								ElasticLoadBalancerConfig: &saasv1alpha1.ElasticLoadBalancerSpec{
									ProxyProtocol:                 util.Pointer(true),
									CrossZoneLoadBalancingEnabled: util.Pointer(true),
									ConnectionDrainingEnabled:     util.Pointer(true),
									ConnectionDrainingTimeout:     util.Pointer[int32](60),
									HealthcheckHealthyThreshold:   util.Pointer[int32](2),
									HealthcheckUnhealthyThreshold: util.Pointer[int32](2),
									HealthcheckInterval:           util.Pointer[int32](5),
									HealthcheckTimeout:            util.Pointer[int32](3),
								},
							},
							NodeID: util.Pointer("test"),
						},
					}},
				},
			},
			want: []ServiceDescriptor{
				{
					PublishingStrategy: saasv1alpha1.PublishingStrategy{
						Strategy:     saasv1alpha1.Marin3rSidecarStrategy,
						EndpointName: "Gateway",
						Marin3rSidecar: &saasv1alpha1.Marin3rSidecarSpec{
							Simple: &saasv1alpha1.Simple{
								ServiceType: util.Pointer(saasv1alpha1.ServiceTypeELB),
								ElasticLoadBalancerConfig: &saasv1alpha1.ElasticLoadBalancerSpec{
									ProxyProtocol:                 util.Pointer(true),
									CrossZoneLoadBalancingEnabled: util.Pointer(true),
									ConnectionDrainingEnabled:     util.Pointer(true),
									ConnectionDrainingTimeout:     util.Pointer[int32](60),
									HealthcheckHealthyThreshold:   util.Pointer[int32](2),
									HealthcheckUnhealthyThreshold: util.Pointer[int32](2),
									HealthcheckInterval:           util.Pointer[int32](5),
									HealthcheckTimeout:            util.Pointer[int32](3),
								},
							},
							NodeID: util.Pointer("test"),
						},
					},
					PortDefinitions: []corev1.ServicePort{{
						Name:       "gateway-http",
						Protocol:   corev1.ProtocolTCP,
						Port:       80,
						TargetPort: intstr.FromString("gateway-http"),
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "Merge: modifies some parameters of the publishing strategy",
			args: args{
				def: []ServiceDescriptor{
					{
						PublishingStrategy: saasv1alpha1.PublishingStrategy{
							Strategy:     saasv1alpha1.SimpleStrategy,
							EndpointName: "Gateway",
							Simple: &saasv1alpha1.Simple{
								ServiceType: util.Pointer(saasv1alpha1.ServiceTypeELB),
								ElasticLoadBalancerConfig: &saasv1alpha1.ElasticLoadBalancerSpec{
									ProxyProtocol:                 util.Pointer(true),
									CrossZoneLoadBalancingEnabled: util.Pointer(true),
									ConnectionDrainingEnabled:     util.Pointer(true),
									ConnectionDrainingTimeout:     util.Pointer[int32](60),
									HealthcheckHealthyThreshold:   util.Pointer[int32](2),
									HealthcheckUnhealthyThreshold: util.Pointer[int32](2),
									HealthcheckInterval:           util.Pointer[int32](5),
									HealthcheckTimeout:            util.Pointer[int32](3),
								},
							},
						},
						PortDefinitions: []corev1.ServicePort{{
							Name:       "gateway-http",
							Protocol:   corev1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromString("gateway-http"),
						}},
					},
				},
				in: &saasv1alpha1.PublishingStrategies{
					Mode: util.Pointer(saasv1alpha1.PublishingStrategiesReconcileModeMerge),
					Endpoints: []saasv1alpha1.PublishingStrategy{{
						Strategy:     saasv1alpha1.SimpleStrategy,
						EndpointName: "Gateway",
						Simple: &saasv1alpha1.Simple{
							ServiceType: util.Pointer(saasv1alpha1.ServiceTypeELB),
							ElasticLoadBalancerConfig: &saasv1alpha1.ElasticLoadBalancerSpec{
								ProxyProtocol:      util.Pointer(false),
								HealthcheckTimeout: util.Pointer[int32](10),
							},
						},
					}},
				},
			},
			want: []ServiceDescriptor{
				{
					PublishingStrategy: saasv1alpha1.PublishingStrategy{
						Strategy:     saasv1alpha1.SimpleStrategy,
						EndpointName: "Gateway",
						Simple: &saasv1alpha1.Simple{
							ServiceType: util.Pointer(saasv1alpha1.ServiceTypeELB),
							ElasticLoadBalancerConfig: &saasv1alpha1.ElasticLoadBalancerSpec{
								ProxyProtocol:                 util.Pointer(false),
								CrossZoneLoadBalancingEnabled: util.Pointer(true),
								ConnectionDrainingEnabled:     util.Pointer(true),
								ConnectionDrainingTimeout:     util.Pointer[int32](60),
								HealthcheckHealthyThreshold:   util.Pointer[int32](2),
								HealthcheckUnhealthyThreshold: util.Pointer[int32](2),
								HealthcheckInterval:           util.Pointer[int32](5),
								HealthcheckTimeout:            util.Pointer[int32](10),
							},
						},
					},
					PortDefinitions: []corev1.ServicePort{{
						Name:       "gateway-http",
						Protocol:   corev1.ProtocolTCP,
						Port:       80,
						TargetPort: intstr.FromString("gateway-http"),
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "Replace: replaces the whole list of endpoints",
			args: args{
				def: []ServiceDescriptor{
					{
						PublishingStrategy: saasv1alpha1.PublishingStrategy{
							Strategy:     saasv1alpha1.SimpleStrategy,
							EndpointName: "Gateway",
							Simple: &saasv1alpha1.Simple{
								ServiceType: util.Pointer(saasv1alpha1.ServiceTypeELB),
								ElasticLoadBalancerConfig: &saasv1alpha1.ElasticLoadBalancerSpec{
									ProxyProtocol:                 util.Pointer(true),
									CrossZoneLoadBalancingEnabled: util.Pointer(true),
									ConnectionDrainingEnabled:     util.Pointer(true),
									ConnectionDrainingTimeout:     util.Pointer[int32](60),
									HealthcheckHealthyThreshold:   util.Pointer[int32](2),
									HealthcheckUnhealthyThreshold: util.Pointer[int32](2),
									HealthcheckInterval:           util.Pointer[int32](5),
									HealthcheckTimeout:            util.Pointer[int32](3),
								},
							},
						},
						PortDefinitions: []corev1.ServicePort{{
							Name:       "gateway",
							Protocol:   corev1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromString("gateway"),
						}},
					},
				},
				in: &saasv1alpha1.PublishingStrategies{
					Mode: util.Pointer(saasv1alpha1.PublishingStrategiesReconcileModeReplace),
					Endpoints: []saasv1alpha1.PublishingStrategy{
						{
							Strategy:     saasv1alpha1.SimpleStrategy,
							EndpointName: "Gateway",
							Simple: &saasv1alpha1.Simple{
								ServiceType: util.Pointer(saasv1alpha1.ServiceTypeELB),
							},
						},
						{
							Strategy:       saasv1alpha1.Marin3rSidecarStrategy,
							EndpointName:   "Gateway",
							Marin3rSidecar: &saasv1alpha1.Marin3rSidecarSpec{},
						},
					},
				},
			},
			want: []ServiceDescriptor{
				{
					PublishingStrategy: saasv1alpha1.PublishingStrategy{
						Strategy:     saasv1alpha1.SimpleStrategy,
						EndpointName: "Gateway",
						Simple: &saasv1alpha1.Simple{
							ServiceType: util.Pointer(saasv1alpha1.ServiceTypeELB),
						},
					},
					PortDefinitions: []corev1.ServicePort{{
						Name:       "gateway",
						Protocol:   corev1.ProtocolTCP,
						Port:       80,
						TargetPort: intstr.FromString("gateway"),
					}},
				},
				{
					PublishingStrategy: saasv1alpha1.PublishingStrategy{
						Strategy:       saasv1alpha1.Marin3rSidecarStrategy,
						EndpointName:   "Gateway",
						Marin3rSidecar: &saasv1alpha1.Marin3rSidecarSpec{},
					},
					PortDefinitions: []corev1.ServicePort{{
						Name:       "gateway",
						Protocol:   corev1.ProtocolTCP,
						Port:       80,
						TargetPort: intstr.FromString("gateway"),
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "Merge: undefined endpoint error",
			args: args{
				def: []ServiceDescriptor{
					{PublishingStrategy: saasv1alpha1.PublishingStrategy{EndpointName: "Gateway"}},
				},
				in: &saasv1alpha1.PublishingStrategies{
					Mode:      util.Pointer(saasv1alpha1.PublishingStrategiesReconcileModeMerge),
					Endpoints: []saasv1alpha1.PublishingStrategy{{EndpointName: "Other"}},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Replace: undefined endpoint error",
			args: args{
				def: []ServiceDescriptor{
					{PublishingStrategy: saasv1alpha1.PublishingStrategy{EndpointName: "Gateway"}},
				},
				in: &saasv1alpha1.PublishingStrategies{
					Mode:      util.Pointer(saasv1alpha1.PublishingStrategiesReconcileModeReplace),
					Endpoints: []saasv1alpha1.PublishingStrategy{{EndpointName: "Other"}},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Merge: create new endpoint",
			args: args{
				def: []ServiceDescriptor{
					{PublishingStrategy: saasv1alpha1.PublishingStrategy{EndpointName: "Gateway"}},
				},
				in: &saasv1alpha1.PublishingStrategies{
					Mode: util.Pointer(saasv1alpha1.PublishingStrategiesReconcileModeMerge),
					Endpoints: []saasv1alpha1.PublishingStrategy{{
						Strategy:     saasv1alpha1.SimpleStrategy,
						EndpointName: "Other",
						Create:       util.Pointer(true),
					}},
				},
			},
			want: []ServiceDescriptor{
				{PublishingStrategy: saasv1alpha1.PublishingStrategy{EndpointName: "Gateway"}},
				{PublishingStrategy: saasv1alpha1.PublishingStrategy{
					Strategy:     saasv1alpha1.SimpleStrategy,
					EndpointName: "Other",
					Create:       util.Pointer(true),
				}},
			},
			wantErr: false,
		},
		{
			name: "Replace: create new endpoint",
			args: args{
				def: []ServiceDescriptor{
					{PublishingStrategy: saasv1alpha1.PublishingStrategy{EndpointName: "Gateway"}},
				},
				in: &saasv1alpha1.PublishingStrategies{
					Mode: util.Pointer(saasv1alpha1.PublishingStrategiesReconcileModeReplace),
					Endpoints: []saasv1alpha1.PublishingStrategy{{
						Strategy:     saasv1alpha1.SimpleStrategy,
						EndpointName: "Other",
						Create:       util.Pointer(true),
					}},
				},
			},
			want: []ServiceDescriptor{{
				PublishingStrategy: saasv1alpha1.PublishingStrategy{
					Strategy:     saasv1alpha1.SimpleStrategy,
					EndpointName: "Other",
					Create:       util.Pointer(true),
				},
			}},
			wantErr: false,
		},
		{
			name: "Merge: no enpoint definintions in the API returns defaults",
			args: args{
				def: []ServiceDescriptor{
					{PublishingStrategy: saasv1alpha1.PublishingStrategy{EndpointName: "Gateway"}},
				},
				in: &saasv1alpha1.PublishingStrategies{
					Mode: util.Pointer(saasv1alpha1.PublishingStrategiesReconcileModeMerge),
				},
			},
			want: []ServiceDescriptor{{
				PublishingStrategy: saasv1alpha1.PublishingStrategy{EndpointName: "Gateway"},
			}},
			wantErr: false,
		},
		{
			name: "Replace: no enpoint definintions in the API returns empty list",
			args: args{
				def: []ServiceDescriptor{
					{PublishingStrategy: saasv1alpha1.PublishingStrategy{EndpointName: "Gateway"}},
				},
				in: &saasv1alpha1.PublishingStrategies{
					Mode: util.Pointer(saasv1alpha1.PublishingStrategiesReconcileModeReplace),
				},
			},
			want:    []ServiceDescriptor{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MergeWithDefaultPublishingStrategy(tt.args.def, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("MergeWithDefaultPublishingStrategy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); len(diff) > 0 {
				t.Errorf("MergeWithDefaultPublishingStrategy() got diff %v", diff)
			}
		})
	}
}

func TestNullTransformer(t *testing.T) {

	type P struct {
		D *bool
		E *int
	}

	type Foo struct {
		A *bool
		B *int
		C *P
	}

	in := Foo{
		A: util.Pointer(false),
		B: util.Pointer(3),
		C: &P{
			D: util.Pointer(false),
		},
	}
	def := Foo{
		A: util.Pointer(true),
		B: util.Pointer(10),
		C: &P{
			D: util.Pointer(true),
			E: util.Pointer(3),
		},
	}
	want := Foo{
		A: util.Pointer(false),
		B: util.Pointer(3),
		C: &P{
			D: util.Pointer(false),
			E: util.Pointer(3),
		},
	}

	mergo.Merge(&def, in, mergo.WithOverride, mergo.WithTransformers(&nullTransformer{}))

	if diff := cmp.Diff(def, want); len(diff) > 0 {
		t.Errorf("TestNullTransformer() got diff %v", diff)
	}
}
