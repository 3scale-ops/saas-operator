package v1alpha1

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
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

func TestEnvoyDynamicConfigRaw_GetRawConfig(t *testing.T) {
	type fields struct {
		RawConfig *runtime.RawExtension
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "returns the raw config",
			fields: fields{
				RawConfig: &runtime.RawExtension{
					Raw:    []byte("whatever"),
					Object: nil,
				},
			},
			want: []byte("whatever"),
		},
		{
			name: "returns nil",
			fields: fields{
				RawConfig: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := &EnvoyDynamicConfigRaw{
				RawConfig: tt.fields.RawConfig,
			}
			if got := raw.GetRawConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnvoyDynamicConfigRaw.GetRawConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}