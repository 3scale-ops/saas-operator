package resources

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-test/deep"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func Test_populateServiceSpecRuntimeValues(t *testing.T) {
	type args struct {
		ctx context.Context
		cl  client.Client
		svc *corev1.Service
	}
	tests := []struct {
		name    string
		args    args
		want    *corev1.Service
		wantErr bool
	}{
		{
			name: "Populates the runtime fields",
			args: args{
				ctx: context.TODO(),
				cl: fake.NewClientBuilder().WithScheme(scheme.Scheme).WithObjects(
					&corev1.Service{
						ObjectMeta: metav1.ObjectMeta{
							Name: "service", Namespace: "ns",
						},
						Spec: corev1.ServiceSpec{
							Type:       corev1.ServiceTypeLoadBalancer,
							IPFamilies: []corev1.IPFamily{corev1.IPv4Protocol},
							IPFamilyPolicy: func() *corev1.IPFamilyPolicyType {
								f := corev1.IPFamilyPolicySingleStack
								return &f
							}(),
							ClusterIP:  "1.1.1.1",
							ClusterIPs: []string{"1.1.1.1"},
							Ports: []corev1.ServicePort{{
								Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP, NodePort: 3333}},
						},
					}).Build(),
				svc: &corev1.Service{
					ObjectMeta: metav1.ObjectMeta{Name: "service", Namespace: "ns"},
					Spec: corev1.ServiceSpec{
						Type: corev1.ServiceTypeLoadBalancer,
						Ports: []corev1.ServicePort{{
							Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP}},
					},
				},
			},
			want: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "service", Namespace: "ns",
				},
				Spec: corev1.ServiceSpec{
					Type:       corev1.ServiceTypeLoadBalancer,
					IPFamilies: []corev1.IPFamily{corev1.IPv4Protocol},
					IPFamilyPolicy: func() *corev1.IPFamilyPolicyType {
						f := corev1.IPFamilyPolicySingleStack
						return &f
					}(),
					ClusterIP:  "1.1.1.1",
					ClusterIPs: []string{"1.1.1.1"},
					Ports: []corev1.ServicePort{{
						Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP, NodePort: 3333}},
				},
			},
			wantErr: false,
		},
		{
			name: "Populates the runtime fields (template adds a new port)",
			args: args{
				ctx: context.TODO(),
				cl: fake.NewClientBuilder().WithScheme(scheme.Scheme).WithObjects(
					&corev1.Service{
						ObjectMeta: metav1.ObjectMeta{
							Name: "service", Namespace: "ns",
						},
						Spec: corev1.ServiceSpec{
							Type:       corev1.ServiceTypeLoadBalancer,
							IPFamilies: []corev1.IPFamily{corev1.IPv4Protocol},
							IPFamilyPolicy: func() *corev1.IPFamilyPolicyType {
								f := corev1.IPFamilyPolicySingleStack
								return &f
							}(),
							ClusterIP:  "1.1.1.1",
							ClusterIPs: []string{"1.1.1.1"},
							Ports: []corev1.ServicePort{{
								Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP, NodePort: 3333}},
						},
					}).Build(),
				svc: &corev1.Service{
					ObjectMeta: metav1.ObjectMeta{Name: "service", Namespace: "ns"},
					Spec: corev1.ServiceSpec{
						Type: corev1.ServiceTypeLoadBalancer,
						Ports: []corev1.ServicePort{
							{Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP},
							{Name: "port", Port: 8080, TargetPort: intstr.FromInt(8080), Protocol: corev1.ProtocolTCP}},
					},
				},
			},
			want: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "service", Namespace: "ns",
				},
				Spec: corev1.ServiceSpec{
					Type:       corev1.ServiceTypeLoadBalancer,
					IPFamilies: []corev1.IPFamily{corev1.IPv4Protocol},
					IPFamilyPolicy: func() *corev1.IPFamilyPolicyType {
						f := corev1.IPFamilyPolicySingleStack
						return &f
					}(),
					ClusterIP:  "1.1.1.1",
					ClusterIPs: []string{"1.1.1.1"},
					Ports: []corev1.ServicePort{
						{Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP, NodePort: 3333},
						{Name: "port", Port: 8080, TargetPort: intstr.FromInt(8080), Protocol: corev1.ProtocolTCP},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Populates the runtime fields (template removes a port)",
			args: args{
				ctx: context.TODO(),
				cl: fake.NewClientBuilder().WithScheme(scheme.Scheme).WithObjects(
					&corev1.Service{
						ObjectMeta: metav1.ObjectMeta{
							Name: "service", Namespace: "ns",
						},
						Spec: corev1.ServiceSpec{
							Type:       corev1.ServiceTypeLoadBalancer,
							IPFamilies: []corev1.IPFamily{corev1.IPv4Protocol},
							IPFamilyPolicy: func() *corev1.IPFamilyPolicyType {
								f := corev1.IPFamilyPolicySingleStack
								return &f
							}(),
							ClusterIP:  "1.1.1.1",
							ClusterIPs: []string{"1.1.1.1"},
							Ports: []corev1.ServicePort{
								{Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP, NodePort: 3333},
								{Name: "port", Port: 8080, TargetPort: intstr.FromInt(8080), Protocol: corev1.ProtocolTCP, NodePort: 3334},
							},
						},
					}).Build(),
				svc: &corev1.Service{
					ObjectMeta: metav1.ObjectMeta{Name: "service", Namespace: "ns"},
					Spec: corev1.ServiceSpec{
						Type: corev1.ServiceTypeLoadBalancer,
						Ports: []corev1.ServicePort{{
							Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP}},
					},
				},
			},
			want: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "service", Namespace: "ns",
				},
				Spec: corev1.ServiceSpec{
					Type:       corev1.ServiceTypeLoadBalancer,
					IPFamilies: []corev1.IPFamily{corev1.IPv4Protocol},
					IPFamilyPolicy: func() *corev1.IPFamilyPolicyType {
						f := corev1.IPFamilyPolicySingleStack
						return &f
					}(),
					ClusterIP:  "1.1.1.1",
					ClusterIPs: []string{"1.1.1.1"},
					Ports: []corev1.ServicePort{
						{Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP, NodePort: 3333},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Populates the runtime fields (does not fail with ClusterIP service)",
			args: args{
				ctx: context.TODO(),
				cl: fake.NewClientBuilder().WithScheme(scheme.Scheme).WithObjects(
					&corev1.Service{
						ObjectMeta: metav1.ObjectMeta{
							Name: "service", Namespace: "ns",
						},
						Spec: corev1.ServiceSpec{
							Type:       corev1.ServiceTypeClusterIP,
							IPFamilies: []corev1.IPFamily{corev1.IPv4Protocol},
							IPFamilyPolicy: func() *corev1.IPFamilyPolicyType {
								f := corev1.IPFamilyPolicySingleStack
								return &f
							}(),
							ClusterIP:  "1.1.1.1",
							ClusterIPs: []string{"1.1.1.1"},
							Ports: []corev1.ServicePort{{
								Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP}},
						},
					}).Build(),
				svc: &corev1.Service{
					ObjectMeta: metav1.ObjectMeta{Name: "service", Namespace: "ns"},
					Spec: corev1.ServiceSpec{
						Type: corev1.ServiceTypeClusterIP,
						Ports: []corev1.ServicePort{{
							Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP}},
					},
				},
			},
			want: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "service", Namespace: "ns",
				},
				Spec: corev1.ServiceSpec{
					Type:       corev1.ServiceTypeClusterIP,
					IPFamilies: []corev1.IPFamily{corev1.IPv4Protocol},
					IPFamilyPolicy: func() *corev1.IPFamilyPolicyType {
						f := corev1.IPFamilyPolicySingleStack
						return &f
					}(),
					ClusterIP:  "1.1.1.1",
					ClusterIPs: []string{"1.1.1.1"},
					Ports: []corev1.ServicePort{{
						Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP}},
				},
			},
			wantErr: false,
		},
		{
			name: "Populates the runtime fields (does not fail if Service not found)",
			args: args{
				ctx: context.TODO(),
				cl:  fake.NewClientBuilder().WithScheme(scheme.Scheme).Build(),
				svc: &corev1.Service{
					ObjectMeta: metav1.ObjectMeta{Name: "service", Namespace: "ns"},
					Spec: corev1.ServiceSpec{
						Type: corev1.ServiceTypeClusterIP,
						Ports: []corev1.ServicePort{{
							Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP}},
					},
				},
			},
			want: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "service", Namespace: "ns",
				},
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeClusterIP,
					Ports: []corev1.ServicePort{{
						Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP}},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := populateServiceSpecRuntimeValues(tt.args.ctx, tt.args.cl, tt.args.svc); (err != nil) != tt.wantErr {
				t.Errorf("populateServiceSpecRuntimeValues() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := deep.Equal(tt.args.svc, tt.want); len(diff) > 0 {
				t.Errorf("populateServiceSpecRuntimeValues() = diff %s", diff)

			}
		})
	}
}

func Test_findPort(t *testing.T) {
	type args struct {
		pNumber   int32
		pProtocol corev1.Protocol
		ports     []corev1.ServicePort
	}
	tests := []struct {
		name string
		args args
		want *corev1.ServicePort
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findPort(tt.args.pNumber, tt.args.pProtocol, tt.args.ports); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findPort() = %v, want %v", got, tt.want)
			}
		})
	}
}
