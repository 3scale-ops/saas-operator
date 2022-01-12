package basereconciler

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestService_PopulateSpecRuntimeValues(t *testing.T) {
	type fields struct {
		Template GeneratorFunction
		Enabled  bool
	}
	type args struct {
		ctx context.Context
		cl  client.Client
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *corev1.Service
		wantErr bool
	}{
		{
			name: "Populates the runtime fields",
			fields: fields{
				Template: func() client.Object {
					return &corev1.Service{
						ObjectMeta: metav1.ObjectMeta{Name: "service", Namespace: "ns"},
						Spec: corev1.ServiceSpec{
							Type: corev1.ServiceTypeLoadBalancer,
							Ports: []corev1.ServicePort{{
								Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP}},
						},
					}
				},
				Enabled: true,
			},
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
			fields: fields{
				Template: func() client.Object {
					return &corev1.Service{
						ObjectMeta: metav1.ObjectMeta{Name: "service", Namespace: "ns"},
						Spec: corev1.ServiceSpec{
							Type: corev1.ServiceTypeLoadBalancer,
							Ports: []corev1.ServicePort{
								{Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP},
								{Name: "port", Port: 8080, TargetPort: intstr.FromInt(8080), Protocol: corev1.ProtocolTCP}},
						},
					}
				},
				Enabled: true,
			},
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
			fields: fields{
				Template: func() client.Object {
					return &corev1.Service{
						ObjectMeta: metav1.ObjectMeta{Name: "service", Namespace: "ns"},
						Spec: corev1.ServiceSpec{
							Type: corev1.ServiceTypeLoadBalancer,
							Ports: []corev1.ServicePort{{
								Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP}},
						},
					}
				},
				Enabled: true,
			},
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
			fields: fields{
				Template: func() client.Object {
					return &corev1.Service{
						ObjectMeta: metav1.ObjectMeta{Name: "service", Namespace: "ns"},
						Spec: corev1.ServiceSpec{
							Type: corev1.ServiceTypeClusterIP,
							Ports: []corev1.ServicePort{{
								Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP}},
						},
					}
				},
				Enabled: true,
			},
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
			fields: fields{
				Template: func() client.Object {
					return &corev1.Service{
						ObjectMeta: metav1.ObjectMeta{Name: "service", Namespace: "ns"},
						Spec: corev1.ServiceSpec{
							Type: corev1.ServiceTypeClusterIP,
							Ports: []corev1.ServicePort{{
								Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP}},
						},
					}
				},
				Enabled: true,
			},
			args: args{
				ctx: context.TODO(),
				cl:  fake.NewClientBuilder().WithScheme(scheme.Scheme).Build(),
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
			s := &Service{
				Template: tt.fields.Template,
				Enabled:  tt.fields.Enabled,
			}
			got, err := s.PopulateSpecRuntimeValues(tt.args.ctx, tt.args.cl)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.PopulateSpecRuntimeValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got().(*corev1.Service), tt.want); len(diff) > 0 {
				t.Errorf("Service.PopulateSpecRuntimeValues() = diff %s", diff)

			}
		})
	}
}
