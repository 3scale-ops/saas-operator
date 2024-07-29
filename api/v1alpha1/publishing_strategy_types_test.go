package v1alpha1

import (
	"context"
	"reflect"
	"testing"

	"github.com/3scale-ops/basereconciler/util"
	envoyconfig "github.com/3scale-ops/saas-operator/pkg/resource_builders/envoyconfig/descriptor"
	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestMarin3rSidecarSpec_Default(t *testing.T) {
	type fields struct {
		Ports               []SidecarPort
		Resources           *ResourceRequirementsSpec
		ExtraPodAnnotations map[string]string
	}
	type args struct {
		def defaultMarin3rSidecarSpec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Marin3rSidecarSpec
	}{
		{
			name:   "Sets defaults",
			fields: fields{},
			args: args{def: defaultMarin3rSidecarSpec{
				Ports: []SidecarPort{
					{
						Name: "test",
						Port: 9999,
					},
				},
				Resources: defaultResourceRequirementsSpec{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("200m"),
						corev1.ResourceMemory: resource.MustParse("200Mi"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("100m"),
						corev1.ResourceMemory: resource.MustParse("100Mi"),
					},
				},
			}},
			want: &Marin3rSidecarSpec{
				Ports: []SidecarPort{
					{
						Name: "test",
						Port: 9999,
					},
				},
				Resources: &ResourceRequirementsSpec{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("200m"),
						corev1.ResourceMemory: resource.MustParse("200Mi"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("100m"),
						corev1.ResourceMemory: resource.MustParse("100Mi"),
					},
				},
			},
		},
		{
			name: "Combines explicitely set values with defaults",
			fields: fields{
				Resources: &ResourceRequirementsSpec{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("99m"),
						corev1.ResourceMemory: resource.MustParse("99Mi"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("99m"),
						corev1.ResourceMemory: resource.MustParse("99Mi"),
					},
				},
			},
			args: args{def: defaultMarin3rSidecarSpec{
				Ports: []SidecarPort{
					{
						Name: "test",
						Port: 9999,
					},
				},
				Resources: defaultResourceRequirementsSpec{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("200m"),
						corev1.ResourceMemory: resource.MustParse("200Mi"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("100m"),
						corev1.ResourceMemory: resource.MustParse("100Mi"),
					},
				},
			}},
			want: &Marin3rSidecarSpec{
				Ports: []SidecarPort{
					{
						Name: "test",
						Port: 9999,
					},
				},
				Resources: &ResourceRequirementsSpec{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("99m"),
						corev1.ResourceMemory: resource.MustParse("99Mi"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("99m"),
						corev1.ResourceMemory: resource.MustParse("99Mi"),
					},
				},
			},
		},
		{
			name:   "Default is deactivated",
			fields: fields{},
			args:   args{def: defaultMarin3rSidecarSpec{}},
			want:   &Marin3rSidecarSpec{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &Marin3rSidecarSpec{
				Ports:               tt.fields.Ports,
				Resources:           tt.fields.Resources,
				ExtraPodAnnotations: tt.fields.ExtraPodAnnotations,
			}
			spec.Default(tt.args.def)
			if !reflect.DeepEqual(spec, tt.want) {
				t.Errorf("Marin3rSidecarSpec_Default() = %v, want %v", *spec, *tt.want)
			}
		})
	}
}

func TestMarin3rSidecarSpec_IsDeactivated(t *testing.T) {
	tests := []struct {
		name string
		spec *Marin3rSidecarSpec
		want bool
	}{
		{"Wants true if empty", &Marin3rSidecarSpec{}, true},
		{"Wants false if nil", nil, false},
		{"Wants false if other", &Marin3rSidecarSpec{
			Ports: []SidecarPort{{Port: 9999, Name: "test"}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.spec.IsDeactivated(); got != tt.want {
				t.Errorf("Marin3rSidecarSpec_IsDeactivated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitializeMarin3rSidecarSpec(t *testing.T) {
	type args struct {
		spec *Marin3rSidecarSpec
		def  defaultMarin3rSidecarSpec
	}
	tests := []struct {
		name string
		args args
		want *Marin3rSidecarSpec
	}{
		{
			name: "Initializes the struct with appropriate defaults if nil",
			args: args{nil, defaultMarin3rSidecarSpec{
				Ports: []SidecarPort{
					{
						Name: "test",
						Port: 9999,
					},
				},
				Resources: defaultResourceRequirementsSpec{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("200m"),
						corev1.ResourceMemory: resource.MustParse("200Mi"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("100m"),
						corev1.ResourceMemory: resource.MustParse("100Mi"),
					},
				},
			}},
			want: &Marin3rSidecarSpec{
				Ports: []SidecarPort{
					{
						Name: "test",
						Port: 9999,
					},
				},
				Resources: &ResourceRequirementsSpec{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("200m"),
						corev1.ResourceMemory: resource.MustParse("200Mi"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("100m"),
						corev1.ResourceMemory: resource.MustParse("100Mi"),
					},
				},
			},
		},
		{
			name: "Deactivated",
			args: args{&Marin3rSidecarSpec{}, defaultMarin3rSidecarSpec{}},
			want: &Marin3rSidecarSpec{},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitializeMarin3rSidecarSpec(tt.args.spec, tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitializeMarin3rSidecarSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapOfEnvoyDynamicConfig_AsList(t *testing.T) {
	tests := []struct {
		name       string
		mapofconfs MapOfEnvoyDynamicConfig
		want       []envoyconfig.EnvoyDynamicConfigDescriptor
	}{
		{
			name: "Returns the map as a list of EnvoyDynamicConfigDescriptor",
			mapofconfs: map[string]EnvoyDynamicConfig{
				"one": {
					Name:             "",
					GeneratorVersion: new(string),
					ListenerHttp:     &ListenerHttp{},
				},
				"two": {
					Name:             "",
					GeneratorVersion: new(string),
					Cluster:          &Cluster{},
				},
			},
			want: []envoyconfig.EnvoyDynamicConfigDescriptor{
				&EnvoyDynamicConfig{
					Name:             "one",
					GeneratorVersion: new(string),
					ListenerHttp:     &ListenerHttp{},
				},
				&EnvoyDynamicConfig{
					Name:             "two",
					GeneratorVersion: new(string),
					Cluster:          &Cluster{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mapofconfs.AsList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapOfEnvoyDynamicConfig.AsList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkloadPublishingStrategyUpgrader_Build(t *testing.T) {
	type fields struct {
		EndpointName         string
		ServiceName          string
		ServiceType          ServiceType
		Namespace            string
		Endpoint             *Endpoint
		Marin3r              *Marin3rSidecarSpec
		ELBSpec              *ElasticLoadBalancerSpec
		NLBSpec              *NetworkLoadBalancerSpec
		ServicePortsOverride []corev1.ServicePort
		Create               bool
	}
	type args struct {
		ctx context.Context
		cl  client.Client
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *PublishingStrategy
		wantErr bool
	}{
		{
			name: "Simple/No_Service: does not generate strategy",
			fields: fields{
				EndpointName: "HTTP",
				ServiceName:  "service",
				ServiceType:  "NLB",
				Namespace:    "ns",
			},
			args: args{
				ctx: context.TODO(),
				cl:  fake.NewClientBuilder().Build(),
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Simple/No_Service: upgrades api",
			fields: fields{
				EndpointName: "HTTP",
				ServiceName:  "service",
				ServiceType:  "NLB",
				Namespace:    "ns",
				Endpoint:     &Endpoint{DNS: []string{"hostname"}},
			},
			args: args{
				ctx: context.TODO(),
				cl:  fake.NewClientBuilder().Build(),
			},
			want: &PublishingStrategy{
				Strategy:     "Simple",
				EndpointName: "HTTP",
				Simple: &Simple{
					ExternalDnsHostnames: []string{"hostname"},
				},
			},
			wantErr: false,
		},
		{
			name: "Simple/With_Service: keeps old service",
			fields: fields{
				EndpointName: "HTTP",
				ServiceName:  "service",
				ServiceType:  "NLB",
				Namespace:    "ns",
			},
			args: args{
				ctx: context.TODO(),
				cl:  fake.NewClientBuilder().WithObjects(&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "service", Namespace: "ns"}}).Build(),
			},
			want: &PublishingStrategy{
				Strategy:     "Simple",
				EndpointName: "HTTP",
				Simple: &Simple{
					ServiceType:         util.Pointer(ServiceTypeNLB),
					ServiceNameOverride: util.Pointer("service"),
				},
			},
			wantErr: false,
		},
		{
			name: "Marin3r/No_Service: upgrades API",
			fields: fields{
				EndpointName: "HTTP",
				ServiceName:  "service",
				ServiceType:  "NLB",
				Namespace:    "ns",
				Marin3r: &Marin3rSidecarSpec{
					Ports:              []SidecarPort{{Name: "port1", Port: 1111}, {Name: "port2", Port: 2222}},
					EnvoyDynamicConfig: map[string]EnvoyDynamicConfig{},
				},
				ServicePortsOverride: []corev1.ServicePort{
					{Name: "port1", Port: 1111, TargetPort: intstr.FromString("port1")},
					{Name: "port1", Port: 2222, TargetPort: intstr.FromString("port2")},
				},
			},
			args: args{
				ctx: context.TODO(),
				cl:  fake.NewClientBuilder().Build(),
			},
			want: &PublishingStrategy{
				Strategy:     "Marin3rSidecar",
				EndpointName: "HTTP",
				Marin3rSidecar: &Marin3rSidecarSpec{
					Simple: &Simple{
						ServiceType: util.Pointer(ServiceTypeNLB),
					},
					Ports:              []SidecarPort{{Name: "port1", Port: 1111}, {Name: "port2", Port: 2222}},
					EnvoyDynamicConfig: map[string]EnvoyDynamicConfig{},
				},
			},
			wantErr: false,
		},
		{
			name: "Marin3r/With_Service: upgrades API and keeps old Service",
			fields: fields{
				EndpointName: "HTTP",
				ServiceName:  "service",
				ServiceType:  "NLB",
				Namespace:    "ns",
				Marin3r: &Marin3rSidecarSpec{
					Ports:              []SidecarPort{{Name: "port1", Port: 1111}, {Name: "port2", Port: 2222}},
					EnvoyDynamicConfig: map[string]EnvoyDynamicConfig{},
				},
				ServicePortsOverride: []corev1.ServicePort{
					{Name: "port1", Port: 1111, TargetPort: intstr.FromString("port1")},
					{Name: "port1", Port: 2222, TargetPort: intstr.FromString("port2")},
				},
			},
			args: args{
				ctx: context.TODO(),
				cl:  fake.NewClientBuilder().WithObjects(&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "service", Namespace: "ns"}}).Build(),
			},
			want: &PublishingStrategy{
				Strategy:     "Marin3rSidecar",
				EndpointName: "HTTP",
				Marin3rSidecar: &Marin3rSidecarSpec{
					Simple: &Simple{
						ServiceNameOverride: util.Pointer("service"),
						ServiceType:         util.Pointer(ServiceTypeNLB),
						ServicePortsOverride: []corev1.ServicePort{
							{Name: "port1", Port: 1111, TargetPort: intstr.FromString("port1")},
							{Name: "port1", Port: 2222, TargetPort: intstr.FromString("port2")},
						},
					},
					Ports:              []SidecarPort{{Name: "port1", Port: 1111}, {Name: "port2", Port: 2222}},
					EnvoyDynamicConfig: map[string]EnvoyDynamicConfig{},
				},
			},
			wantErr: false,
		},
		{
			name: "Appends new endpoint only if Service exists",
			fields: fields{
				EndpointName: "New",
				ServiceName:  "service",
				ServiceType:  "ClusterIP",
				Namespace:    "ns",
				Endpoint:     &Endpoint{DNS: []string{"hostname"}},
				Create:       true,
			},
			args: args{
				ctx: context.TODO(),
				cl:  fake.NewClientBuilder().WithObjects(&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "service", Namespace: "ns"}}).Build(),
			},
			want: &PublishingStrategy{
				Strategy:     "Simple",
				EndpointName: "New",
				Simple: &Simple{
					ServiceNameOverride:  util.Pointer("service"),
					ServiceType:          util.Pointer(ServiceTypeClusterIP),
					ExternalDnsHostnames: []string{"hostname"},
				},
				Create: util.Pointer(true),
			},
			wantErr: false,
		},
		{
			name: "Don't add new endpoint only if Service doesn't exist",
			fields: fields{
				EndpointName: "New",
				ServiceName:  "service",
				ServiceType:  "ClusterIP",
				Namespace:    "ns",
				Endpoint:     &Endpoint{DNS: []string{"hostname"}},
				Create:       true,
			},
			args: args{
				ctx: context.TODO(),
				cl:  fake.NewClientBuilder().Build(),
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := WorkloadPublishingStrategyUpgrader{
				EndpointName:         tt.fields.EndpointName,
				ServiceName:          tt.fields.ServiceName,
				ServiceType:          tt.fields.ServiceType,
				Namespace:            tt.fields.Namespace,
				Endpoint:             tt.fields.Endpoint,
				Marin3r:              tt.fields.Marin3r,
				ELBSpec:              tt.fields.ELBSpec,
				NLBSpec:              tt.fields.NLBSpec,
				ServicePortOverrides: tt.fields.ServicePortsOverride,
				Create:               tt.fields.Create,
			}
			got, err := gen.Build(tt.args.ctx, tt.args.cl)
			if (err != nil) != tt.wantErr {
				t.Errorf("WorkloadPublishingStrategyUpgrader.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); len(diff) > 0 {
				t.Errorf("WorkloadPublishingStrategyUpgrader.Build() = diff %v", diff)
			}
		})
	}
}
