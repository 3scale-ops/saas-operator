/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"reflect"
	"testing"

	"github.com/3scale/saas-operator/pkg/util"
	"github.com/go-test/deep"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

func TestImageSpec_Default(t *testing.T) {
	type fields struct {
		Name           *string
		Tag            *string
		PullSecretName *string
		PullPolicy     *corev1.PullPolicy
	}
	type args struct {
		def defaultImageSpec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ImageSpec
	}{
		{
			name:   "Sets defaults",
			fields: fields{},
			args: args{def: defaultImageSpec{
				Name:           pointer.StringPtr("name"),
				Tag:            pointer.StringPtr("tag"),
				PullSecretName: pointer.StringPtr("pullSecret"),
				PullPolicy:     func() *corev1.PullPolicy { p := corev1.PullIfNotPresent; return &p }(),
			}},
			want: &ImageSpec{
				Name:           pointer.StringPtr("name"),
				Tag:            pointer.StringPtr("tag"),
				PullSecretName: pointer.StringPtr("pullSecret"),
				PullPolicy:     func() *corev1.PullPolicy { p := corev1.PullIfNotPresent; return &p }(),
			},
		},
		{
			name: "Combines explicitely set values with defaults",
			fields: fields{
				Name:       pointer.StringPtr("explicit"),
				PullPolicy: func() *corev1.PullPolicy { p := corev1.PullAlways; return &p }(),
			},
			args: args{def: defaultImageSpec{
				Name:           pointer.StringPtr("name"),
				Tag:            pointer.StringPtr("tag"),
				PullSecretName: pointer.StringPtr("pullSecret"),
				PullPolicy:     func() *corev1.PullPolicy { p := corev1.PullIfNotPresent; return &p }(),
			}},
			want: &ImageSpec{
				Name:           pointer.StringPtr("explicit"),
				Tag:            pointer.StringPtr("tag"),
				PullSecretName: pointer.StringPtr("pullSecret"),
				PullPolicy:     func() *corev1.PullPolicy { p := corev1.PullAlways; return &p }(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &ImageSpec{
				Name:           tt.fields.Name,
				Tag:            tt.fields.Tag,
				PullSecretName: tt.fields.PullSecretName,
				PullPolicy:     tt.fields.PullPolicy,
			}
			spec.Default(tt.args.def)
			if !reflect.DeepEqual(spec, tt.want) {
				t.Errorf("ImageSpec_Default() = %v, want %v", *spec, *tt.want)
			}
		})
	}
}

func TestImageSpec_IsDeactivated(t *testing.T) {
	tests := []struct {
		name string
		spec *ImageSpec
		want bool
	}{
		{"Wants false if empty", &ImageSpec{}, false},
		{"Wants false if nil", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.spec.IsDeactivated(); got != tt.want {
				t.Errorf("ImageSpec.IsDeactivated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitializeImageSpec(t *testing.T) {
	type args struct {
		spec *ImageSpec
		def  defaultImageSpec
	}
	tests := []struct {
		name string
		args args
		want *ImageSpec
	}{
		{
			name: "Initializes the struct with appropriate defaults if nil",
			args: args{nil, defaultImageSpec{
				Name:           pointer.StringPtr("name"),
				Tag:            pointer.StringPtr("tag"),
				PullSecretName: pointer.StringPtr("pullSecret"),
			}},
			want: &ImageSpec{
				Name:           pointer.StringPtr("name"),
				Tag:            pointer.StringPtr("tag"),
				PullSecretName: pointer.StringPtr("pullSecret"),
			},
		},
		{
			name: "Initializes the struct with appropriate defaults if empty",
			args: args{&ImageSpec{}, defaultImageSpec{
				Name:           pointer.StringPtr("name"),
				Tag:            pointer.StringPtr("tag"),
				PullSecretName: pointer.StringPtr("pullSecret"),
			}},
			want: &ImageSpec{
				Name:           pointer.StringPtr("name"),
				Tag:            pointer.StringPtr("tag"),
				PullSecretName: pointer.StringPtr("pullSecret"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitializeImageSpec(tt.args.spec, tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitializeImageSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProbeSpec_Default(t *testing.T) {
	type fields struct {
		InitialDelaySeconds *int32
		TimeoutSeconds      *int32
		PeriodSeconds       *int32
		SuccessThreshold    *int32
		FailureThreshold    *int32
	}
	type args struct {
		def defaultProbeSpec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ProbeSpec
	}{
		{
			name:   "Sets defaults",
			fields: fields{},
			args: args{def: defaultProbeSpec{
				InitialDelaySeconds: pointer.Int32Ptr(1),
				TimeoutSeconds:      pointer.Int32Ptr(2),
				PeriodSeconds:       pointer.Int32Ptr(3),
				SuccessThreshold:    pointer.Int32Ptr(4),
				FailureThreshold:    pointer.Int32Ptr(5),
			}},
			want: &ProbeSpec{
				InitialDelaySeconds: pointer.Int32Ptr(1),
				TimeoutSeconds:      pointer.Int32Ptr(2),
				PeriodSeconds:       pointer.Int32Ptr(3),
				SuccessThreshold:    pointer.Int32Ptr(4),
				FailureThreshold:    pointer.Int32Ptr(5),
			},
		},
		{
			name: "Combines explicitely set values with defaults",
			fields: fields{
				InitialDelaySeconds: pointer.Int32Ptr(9999),
			},
			args: args{def: defaultProbeSpec{
				InitialDelaySeconds: pointer.Int32Ptr(1),
				TimeoutSeconds:      pointer.Int32Ptr(2),
				PeriodSeconds:       pointer.Int32Ptr(3),
				SuccessThreshold:    pointer.Int32Ptr(4),
				FailureThreshold:    pointer.Int32Ptr(5),
			}},
			want: &ProbeSpec{
				InitialDelaySeconds: pointer.Int32Ptr(9999),
				TimeoutSeconds:      pointer.Int32Ptr(2),
				PeriodSeconds:       pointer.Int32Ptr(3),
				SuccessThreshold:    pointer.Int32Ptr(4),
				FailureThreshold:    pointer.Int32Ptr(5),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &ProbeSpec{
				InitialDelaySeconds: tt.fields.InitialDelaySeconds,
				TimeoutSeconds:      tt.fields.TimeoutSeconds,
				PeriodSeconds:       tt.fields.PeriodSeconds,
				SuccessThreshold:    tt.fields.SuccessThreshold,
				FailureThreshold:    tt.fields.FailureThreshold,
			}
			spec.Default(tt.args.def)
			if !reflect.DeepEqual(spec, tt.want) {
				t.Errorf("ProbeSpec_Default() = %v, want %v", *spec, *tt.want)
			}
		})
	}
}

func TestProbeSpec_IsDeactivated(t *testing.T) {
	tests := []struct {
		name string
		spec *ProbeSpec
		want bool
	}{
		{"Wants true if empty", &ProbeSpec{}, true},
		{"Wants false if nil", nil, false},
		{"Wants false if other", &ProbeSpec{InitialDelaySeconds: pointer.Int32Ptr(1)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.spec.IsDeactivated(); got != tt.want {
				t.Errorf("ProbeSpec.IsDeactivated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitializeProbeSpec(t *testing.T) {
	type args struct {
		spec *ProbeSpec
		def  defaultProbeSpec
	}
	tests := []struct {
		name string
		args args
		want *ProbeSpec
	}{
		{
			name: "Initializes the struct with appropriate defaults if nil",
			args: args{nil, defaultProbeSpec{
				InitialDelaySeconds: pointer.Int32Ptr(1),
				TimeoutSeconds:      pointer.Int32Ptr(2),
				PeriodSeconds:       pointer.Int32Ptr(3),
				SuccessThreshold:    pointer.Int32Ptr(4),
				FailureThreshold:    pointer.Int32Ptr(5),
			}},
			want: &ProbeSpec{
				InitialDelaySeconds: pointer.Int32Ptr(1),
				TimeoutSeconds:      pointer.Int32Ptr(2),
				PeriodSeconds:       pointer.Int32Ptr(3),
				SuccessThreshold:    pointer.Int32Ptr(4),
				FailureThreshold:    pointer.Int32Ptr(5),
			},
		},
		{
			name: "Deactivated",
			args: args{&ProbeSpec{}, defaultProbeSpec{}},
			want: &ProbeSpec{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitializeProbeSpec(tt.args.spec, tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitializeProbeSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadBalancerSpec_Default(t *testing.T) {
	type fields struct {
		ProxyProtocol                           *bool
		CrossZoneLoadBalancingEnabled           *bool
		ConnectionDrainingEnabled               *bool
		ConnectionDrainingTimeout               *int32
		ConnectionHealthcheckHealthyThreshold   *int32
		ConnectionHealthcheckUnhealthyThreshold *int32
		ConnectionHealthcheckInterval           *int32
		ConnectionHealthcheckTimeout            *int32
	}
	type args struct {
		def defaultLoadBalancerSpec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *LoadBalancerSpec
	}{
		{
			name:   "Sets defaults",
			fields: fields{},
			args: args{def: defaultLoadBalancerSpec{
				ProxyProtocol:                 pointer.BoolPtr(true),
				CrossZoneLoadBalancingEnabled: pointer.BoolPtr(true),
				ConnectionDrainingEnabled:     pointer.BoolPtr(true),
				ConnectionDrainingTimeout:     pointer.Int32Ptr(1),
				HealthcheckHealthyThreshold:   pointer.Int32Ptr(2),
				HealthcheckUnhealthyThreshold: pointer.Int32Ptr(3),
				HealthcheckInterval:           pointer.Int32Ptr(4),
				HealthcheckTimeout:            pointer.Int32Ptr(5),
			}},
			want: &LoadBalancerSpec{
				ProxyProtocol:                 pointer.BoolPtr(true),
				CrossZoneLoadBalancingEnabled: pointer.BoolPtr(true),
				ConnectionDrainingEnabled:     pointer.BoolPtr(true),
				ConnectionDrainingTimeout:     pointer.Int32Ptr(1),
				HealthcheckHealthyThreshold:   pointer.Int32Ptr(2),
				HealthcheckUnhealthyThreshold: pointer.Int32Ptr(3),
				HealthcheckInterval:           pointer.Int32Ptr(4),
				HealthcheckTimeout:            pointer.Int32Ptr(5),
			},
		},
		{
			name: "Combines explicitely set values with defaults",
			fields: fields{
				ProxyProtocol: pointer.BoolPtr(false),
			},
			args: args{def: defaultLoadBalancerSpec{
				ProxyProtocol:                 pointer.BoolPtr(true),
				CrossZoneLoadBalancingEnabled: pointer.BoolPtr(true),
				ConnectionDrainingEnabled:     pointer.BoolPtr(true),
				ConnectionDrainingTimeout:     pointer.Int32Ptr(1),
				HealthcheckHealthyThreshold:   pointer.Int32Ptr(2),
				HealthcheckUnhealthyThreshold: pointer.Int32Ptr(3),
				HealthcheckInterval:           pointer.Int32Ptr(4),
				HealthcheckTimeout:            pointer.Int32Ptr(5),
			}},
			want: &LoadBalancerSpec{
				ProxyProtocol:                 pointer.BoolPtr(false),
				CrossZoneLoadBalancingEnabled: pointer.BoolPtr(true),
				ConnectionDrainingEnabled:     pointer.BoolPtr(true),
				ConnectionDrainingTimeout:     pointer.Int32Ptr(1),
				HealthcheckHealthyThreshold:   pointer.Int32Ptr(2),
				HealthcheckUnhealthyThreshold: pointer.Int32Ptr(3),
				HealthcheckInterval:           pointer.Int32Ptr(4),
				HealthcheckTimeout:            pointer.Int32Ptr(5),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &LoadBalancerSpec{
				ProxyProtocol:                 tt.fields.ProxyProtocol,
				CrossZoneLoadBalancingEnabled: tt.fields.CrossZoneLoadBalancingEnabled,
				ConnectionDrainingEnabled:     tt.fields.ConnectionDrainingEnabled,
				ConnectionDrainingTimeout:     tt.fields.ConnectionDrainingTimeout,
				HealthcheckHealthyThreshold:   tt.fields.ConnectionHealthcheckHealthyThreshold,
				HealthcheckUnhealthyThreshold: tt.fields.ConnectionHealthcheckUnhealthyThreshold,
				HealthcheckInterval:           tt.fields.ConnectionHealthcheckInterval,
				HealthcheckTimeout:            tt.fields.ConnectionHealthcheckTimeout,
			}
			spec.Default(tt.args.def)
			if !reflect.DeepEqual(spec, tt.want) {
				t.Errorf("LoadBalancerSpec_Default() = %v, want %v", *spec, *tt.want)
			}
		})
	}
}

func TestLoadBalancerSpec_IsDeactivated(t *testing.T) {
	tests := []struct {
		name string
		spec *LoadBalancerSpec
		want bool
	}{
		{"Wants false if empty", &LoadBalancerSpec{}, false},
		{"Wants false if nil", nil, false},
		{"Wants false if other", &LoadBalancerSpec{ProxyProtocol: pointer.BoolPtr(false)}, false}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.spec.IsDeactivated(); got != tt.want {
				t.Errorf("LoadBalancerSpec.IsDeactivated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitializeLoadBalancerSpec(t *testing.T) {
	type args struct {
		spec *LoadBalancerSpec
		def  defaultLoadBalancerSpec
	}
	tests := []struct {
		name string
		args args
		want *LoadBalancerSpec
	}{
		{
			name: "Initializes the struct with appropriate defaults if nil",
			args: args{nil, defaultLoadBalancerSpec{
				ProxyProtocol:                 pointer.BoolPtr(true),
				CrossZoneLoadBalancingEnabled: pointer.BoolPtr(true),
				ConnectionDrainingEnabled:     pointer.BoolPtr(true),
				ConnectionDrainingTimeout:     pointer.Int32Ptr(1),
				HealthcheckHealthyThreshold:   pointer.Int32Ptr(2),
				HealthcheckUnhealthyThreshold: pointer.Int32Ptr(3),
				HealthcheckInterval:           pointer.Int32Ptr(4),
				HealthcheckTimeout:            pointer.Int32Ptr(5),
			}},
			want: &LoadBalancerSpec{
				ProxyProtocol:                 pointer.BoolPtr(true),
				CrossZoneLoadBalancingEnabled: pointer.BoolPtr(true),
				ConnectionDrainingEnabled:     pointer.BoolPtr(true),
				ConnectionDrainingTimeout:     pointer.Int32Ptr(1),
				HealthcheckHealthyThreshold:   pointer.Int32Ptr(2),
				HealthcheckUnhealthyThreshold: pointer.Int32Ptr(3),
				HealthcheckInterval:           pointer.Int32Ptr(4),
				HealthcheckTimeout:            pointer.Int32Ptr(5),
			},
		},
		{
			name: "Initializes the struct with appropriate defaults if empty",
			args: args{&LoadBalancerSpec{}, defaultLoadBalancerSpec{
				ProxyProtocol:                 pointer.BoolPtr(true),
				CrossZoneLoadBalancingEnabled: pointer.BoolPtr(true),
				ConnectionDrainingEnabled:     pointer.BoolPtr(true),
				ConnectionDrainingTimeout:     pointer.Int32Ptr(1),
				HealthcheckHealthyThreshold:   pointer.Int32Ptr(2),
				HealthcheckUnhealthyThreshold: pointer.Int32Ptr(3),
				HealthcheckInterval:           pointer.Int32Ptr(4),
				HealthcheckTimeout:            pointer.Int32Ptr(5),
			}},
			want: &LoadBalancerSpec{
				ProxyProtocol:                 pointer.BoolPtr(true),
				CrossZoneLoadBalancingEnabled: pointer.BoolPtr(true),
				ConnectionDrainingEnabled:     pointer.BoolPtr(true),
				ConnectionDrainingTimeout:     pointer.Int32Ptr(1),
				HealthcheckHealthyThreshold:   pointer.Int32Ptr(2),
				HealthcheckUnhealthyThreshold: pointer.Int32Ptr(3),
				HealthcheckInterval:           pointer.Int32Ptr(4),
				HealthcheckTimeout:            pointer.Int32Ptr(5),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitializeLoadBalancerSpec(tt.args.spec, tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitializeLoadBalancerSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNLBLoadBalancerSpec_Default(t *testing.T) {
	type fields struct {
		ProxyProtocol                 *bool
		CrossZoneLoadBalancingEnabled *bool
	}
	type args struct {
		def defaultNLBLoadBalancerSpec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *NLBLoadBalancerSpec
	}{
		{
			name:   "Sets defaults",
			fields: fields{},
			args: args{def: defaultNLBLoadBalancerSpec{
				ProxyProtocol:                 pointer.BoolPtr(true),
				CrossZoneLoadBalancingEnabled: pointer.BoolPtr(true),
			}},
			want: &NLBLoadBalancerSpec{
				ProxyProtocol:                 pointer.BoolPtr(true),
				CrossZoneLoadBalancingEnabled: pointer.BoolPtr(true),
			},
		},
		{
			name: "Combines explicitely set values with defaults",
			fields: fields{
				ProxyProtocol: pointer.BoolPtr(false),
			},
			args: args{def: defaultNLBLoadBalancerSpec{
				ProxyProtocol:                 pointer.BoolPtr(true),
				CrossZoneLoadBalancingEnabled: pointer.BoolPtr(true),
			}},
			want: &NLBLoadBalancerSpec{
				ProxyProtocol:                 pointer.BoolPtr(false),
				CrossZoneLoadBalancingEnabled: pointer.BoolPtr(true),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &NLBLoadBalancerSpec{
				ProxyProtocol:                 tt.fields.ProxyProtocol,
				CrossZoneLoadBalancingEnabled: tt.fields.CrossZoneLoadBalancingEnabled,
			}
			spec.Default(tt.args.def)
			if !reflect.DeepEqual(spec, tt.want) {
				t.Errorf("NLBLoadBalancerSpec_Default() = %v, want %v", *spec, *tt.want)
			}
		})
	}
}

func TestNLBLoadBalancerSpec_IsDeactivated(t *testing.T) {
	tests := []struct {
		name string
		spec *NLBLoadBalancerSpec
		want bool
	}{
		{"Wants false if empty", &NLBLoadBalancerSpec{}, false},
		{"Wants false if nil", nil, false},
		{"Wants false if other", &NLBLoadBalancerSpec{ProxyProtocol: pointer.BoolPtr(false)}, false}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.spec.IsDeactivated(); got != tt.want {
				t.Errorf("NLBLoadBalancerSpec.IsDeactivated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitializeNLBLoadBalancerSpec(t *testing.T) {
	type args struct {
		spec *NLBLoadBalancerSpec
		def  defaultNLBLoadBalancerSpec
	}
	tests := []struct {
		name string
		args args
		want *NLBLoadBalancerSpec
	}{
		{
			name: "Initializes the struct with appropriate defaults if nil",
			args: args{nil, defaultNLBLoadBalancerSpec{
				ProxyProtocol:                 pointer.BoolPtr(true),
				CrossZoneLoadBalancingEnabled: pointer.BoolPtr(true),
			}},
			want: &NLBLoadBalancerSpec{
				ProxyProtocol:                 pointer.BoolPtr(true),
				CrossZoneLoadBalancingEnabled: pointer.BoolPtr(true),
			},
		},
		{
			name: "Initializes the struct with appropriate defaults if empty",
			args: args{&NLBLoadBalancerSpec{}, defaultNLBLoadBalancerSpec{
				ProxyProtocol:                 pointer.BoolPtr(true),
				CrossZoneLoadBalancingEnabled: pointer.BoolPtr(true),
			}},
			want: &NLBLoadBalancerSpec{
				ProxyProtocol:                 pointer.BoolPtr(true),
				CrossZoneLoadBalancingEnabled: pointer.BoolPtr(true),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitializeNLBLoadBalancerSpec(tt.args.spec, tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitializeNLBLoadBalancerSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrafanaDashboardSpec_Default(t *testing.T) {
	type fields struct {
		SelectorKey   *string
		SelectorValue *string
	}
	type args struct {
		def defaultGrafanaDashboardSpec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *GrafanaDashboardSpec
	}{
		{
			name:   "Sets defaults",
			fields: fields{},
			args: args{def: defaultGrafanaDashboardSpec{
				SelectorKey:   pointer.StringPtr("key"),
				SelectorValue: pointer.StringPtr("label"),
			}},
			want: &GrafanaDashboardSpec{
				SelectorKey:   pointer.StringPtr("key"),
				SelectorValue: pointer.StringPtr("label"),
			},
		},
		{
			name: "Combines explicitely set values with defaults",
			fields: fields{
				SelectorKey: pointer.StringPtr("xxxx"),
			},
			args: args{def: defaultGrafanaDashboardSpec{
				SelectorKey:   pointer.StringPtr("key"),
				SelectorValue: pointer.StringPtr("label"),
			}},
			want: &GrafanaDashboardSpec{
				SelectorKey:   pointer.StringPtr("xxxx"),
				SelectorValue: pointer.StringPtr("label"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &GrafanaDashboardSpec{
				SelectorKey:   tt.fields.SelectorKey,
				SelectorValue: tt.fields.SelectorValue,
			}
			spec.Default(tt.args.def)
			if !reflect.DeepEqual(spec, tt.want) {
				t.Errorf("GrafanaDashboardSpec_Default() = %v, want %v", *spec, *tt.want)
			}
		})
	}
}

func TestGrafanaDashboardSpec_IsDeactivated(t *testing.T) {
	tests := []struct {
		name string
		spec *GrafanaDashboardSpec
		want bool
	}{
		{"Wants true if empty", &GrafanaDashboardSpec{}, true},
		{"Wants false if nil", nil, false},
		{"Wants false if other", &GrafanaDashboardSpec{SelectorKey: pointer.StringPtr("key")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.spec.IsDeactivated(); got != tt.want {
				t.Errorf("GrafanaDashboardSpec_IsDeactivated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitializeGrafanaDashboardSpec(t *testing.T) {
	type args struct {
		spec *GrafanaDashboardSpec
		def  defaultGrafanaDashboardSpec
	}
	tests := []struct {
		name string
		args args
		want *GrafanaDashboardSpec
	}{
		{
			name: "Initializes the struct with appropriate defaults if nil",
			args: args{nil, defaultGrafanaDashboardSpec{
				SelectorKey:   pointer.StringPtr("key"),
				SelectorValue: pointer.StringPtr("label"),
			}},
			want: &GrafanaDashboardSpec{
				SelectorKey:   pointer.StringPtr("key"),
				SelectorValue: pointer.StringPtr("label"),
			},
		},
		{
			name: "Deactivated",
			args: args{&GrafanaDashboardSpec{}, defaultGrafanaDashboardSpec{}},
			want: &GrafanaDashboardSpec{},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitializeGrafanaDashboardSpec(tt.args.spec, tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitializeGrafanaDashboardSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPodDisruptionBudgetSpec_Default(t *testing.T) {
	type fields struct {
		MinAvailable   *intstr.IntOrString
		MaxUnavailable *intstr.IntOrString
	}
	type args struct {
		def defaultPodDisruptionBudgetSpec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *PodDisruptionBudgetSpec
	}{
		{
			name:   "Sets defaults",
			fields: fields{},
			args: args{def: defaultPodDisruptionBudgetSpec{
				MinAvailable:   util.IntStrPtr(intstr.FromString("default")),
				MaxUnavailable: nil,
			}},
			want: &PodDisruptionBudgetSpec{
				MinAvailable:   util.IntStrPtr(intstr.FromString("default")),
				MaxUnavailable: nil,
			},
		},
		{
			name: "Combines explicitely set values with defaults",
			fields: fields{
				MinAvailable: util.IntStrPtr(intstr.FromString("explicit")),
			},
			args: args{def: defaultPodDisruptionBudgetSpec{
				MinAvailable:   util.IntStrPtr(intstr.FromString("default")),
				MaxUnavailable: nil,
			}},
			want: &PodDisruptionBudgetSpec{
				MinAvailable:   util.IntStrPtr(intstr.FromString("explicit")),
				MaxUnavailable: nil,
			},
		},
		{
			name: "Only one of MinAvailable or MaxUnavailable can be set",
			fields: fields{
				MinAvailable: util.IntStrPtr(intstr.FromString("explicit")),
			},
			args: args{def: defaultPodDisruptionBudgetSpec{
				MinAvailable:   nil,
				MaxUnavailable: util.IntStrPtr(intstr.FromString("default")),
			}},
			want: &PodDisruptionBudgetSpec{
				MinAvailable:   util.IntStrPtr(intstr.FromString("explicit")),
				MaxUnavailable: nil,
			},
		},
		{
			name:   "Only one of MinAvailable or MaxUnavailable can be set (II)",
			fields: fields{},
			args: args{def: defaultPodDisruptionBudgetSpec{
				MinAvailable:   util.IntStrPtr(intstr.IntOrString{Type: intstr.String, StrVal: "defaultMin"}),
				MaxUnavailable: util.IntStrPtr(intstr.IntOrString{Type: intstr.String, StrVal: "defaultMax"}),
			}},
			want: &PodDisruptionBudgetSpec{
				MinAvailable:   util.IntStrPtr(intstr.IntOrString{Type: intstr.String, StrVal: "defaultMin"}),
				MaxUnavailable: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &PodDisruptionBudgetSpec{
				MinAvailable:   tt.fields.MinAvailable,
				MaxUnavailable: tt.fields.MaxUnavailable,
			}
			spec.Default(tt.args.def)
			if !reflect.DeepEqual(spec, tt.want) {
				t.Errorf("PodDisruptionBudgetSpec_Default() = %v, want %v", *spec, *tt.want)
			}
		})
	}
}

func TestPodDisruptionBudgetSpec_IsDeactivated(t *testing.T) {
	tests := []struct {
		name string
		spec *PodDisruptionBudgetSpec
		want bool
	}{
		{"Wants true if empty", &PodDisruptionBudgetSpec{}, true},
		{"Wants false if nil", nil, false},
		{"Wants false if other", &PodDisruptionBudgetSpec{MinAvailable: util.IntStrPtr(intstr.FromInt(1))}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.spec.IsDeactivated(); got != tt.want {
				t.Errorf("PodDisruptionBudgetSpec.IsDeactivated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitializePodDisruptionBudgetSpec(t *testing.T) {
	type args struct {
		spec *PodDisruptionBudgetSpec
		def  defaultPodDisruptionBudgetSpec
	}
	tests := []struct {
		name string
		args args
		want *PodDisruptionBudgetSpec
	}{
		{
			name: "Initializes the struct with appropriate defaults if nil",
			args: args{nil, defaultPodDisruptionBudgetSpec{
				MinAvailable:   util.IntStrPtr(intstr.FromString("default")),
				MaxUnavailable: nil,
			}},
			want: &PodDisruptionBudgetSpec{
				MinAvailable:   util.IntStrPtr(intstr.FromString("default")),
				MaxUnavailable: nil,
			},
		},
		{
			name: "Deactivated",
			args: args{&PodDisruptionBudgetSpec{}, defaultPodDisruptionBudgetSpec{}},
			want: &PodDisruptionBudgetSpec{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitializePodDisruptionBudgetSpec(tt.args.spec, tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitializePodDisruptionBudgetSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHorizontalPodAutoscalerSpec_Default(t *testing.T) {
	type fields struct {
		MinReplicas         *int32
		MaxReplicas         *int32
		ResourceName        *string
		ResourceUtilization *int32
	}
	type args struct {
		def defaultHorizontalPodAutoscalerSpec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *HorizontalPodAutoscalerSpec
	}{
		{
			name:   "Sets defaults",
			fields: fields{},
			args: args{def: defaultHorizontalPodAutoscalerSpec{
				MinReplicas:         pointer.Int32Ptr(1),
				MaxReplicas:         pointer.Int32Ptr(2),
				ResourceUtilization: pointer.Int32Ptr(3),
				ResourceName:        pointer.StringPtr("xxxx"),
			}},
			want: &HorizontalPodAutoscalerSpec{
				MinReplicas:         pointer.Int32Ptr(1),
				MaxReplicas:         pointer.Int32Ptr(2),
				ResourceUtilization: pointer.Int32Ptr(3),
				ResourceName:        pointer.StringPtr("xxxx"),
			},
		},
		{
			name: "Combines explicitely set values with defaults",
			fields: fields{
				MinReplicas: pointer.Int32Ptr(9999),
			},
			args: args{def: defaultHorizontalPodAutoscalerSpec{
				MinReplicas:         pointer.Int32Ptr(1),
				MaxReplicas:         pointer.Int32Ptr(2),
				ResourceUtilization: pointer.Int32Ptr(3),
				ResourceName:        pointer.StringPtr("xxxx"),
			}},
			want: &HorizontalPodAutoscalerSpec{
				MinReplicas:         pointer.Int32Ptr(9999),
				MaxReplicas:         pointer.Int32Ptr(2),
				ResourceUtilization: pointer.Int32Ptr(3),
				ResourceName:        pointer.StringPtr("xxxx"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &HorizontalPodAutoscalerSpec{
				MinReplicas:         tt.fields.MinReplicas,
				MaxReplicas:         tt.fields.MaxReplicas,
				ResourceName:        tt.fields.ResourceName,
				ResourceUtilization: tt.fields.ResourceUtilization,
			}
			spec.Default(tt.args.def)
			if !reflect.DeepEqual(spec, tt.want) {
				t.Errorf("HorizontalPodAutoscalerSpec_Default() = %v, want %v", *spec, *tt.want)
			}
		})
	}
}

func TestHorizontalPodAutoscalerSpec_IsDeactivated(t *testing.T) {
	tests := []struct {
		name string
		spec *HorizontalPodAutoscalerSpec
		want bool
	}{
		{"Wants true if empty", &HorizontalPodAutoscalerSpec{}, true},
		{"Wants false if nil", nil, false},
		{"Wants false if other", &HorizontalPodAutoscalerSpec{MinReplicas: pointer.Int32Ptr(1)}, false}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.spec.IsDeactivated(); got != tt.want {
				t.Errorf("HorizontalPodAutoscalerSpec.IsDeactivated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitializeHorizontalPodAutoscalerSpec(t *testing.T) {
	type args struct {
		spec *HorizontalPodAutoscalerSpec
		def  defaultHorizontalPodAutoscalerSpec
	}
	tests := []struct {
		name string
		args args
		want *HorizontalPodAutoscalerSpec
	}{
		{
			name: "Initializes the struct with appropriate defaults if nil",
			args: args{nil, defaultHorizontalPodAutoscalerSpec{
				MinReplicas:         pointer.Int32Ptr(1),
				MaxReplicas:         pointer.Int32Ptr(2),
				ResourceUtilization: pointer.Int32Ptr(3),
				ResourceName:        pointer.StringPtr("xxxx"),
			}},
			want: &HorizontalPodAutoscalerSpec{
				MinReplicas:         pointer.Int32Ptr(1),
				MaxReplicas:         pointer.Int32Ptr(2),
				ResourceUtilization: pointer.Int32Ptr(3),
				ResourceName:        pointer.StringPtr("xxxx"),
			},
		},
		{
			name: "Deactivated",
			args: args{&HorizontalPodAutoscalerSpec{}, defaultHorizontalPodAutoscalerSpec{}},
			want: &HorizontalPodAutoscalerSpec{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitializeHorizontalPodAutoscalerSpec(tt.args.spec, tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitializeHorizontalPodAutoscalerSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceRequirementsSpec_Default(t *testing.T) {
	type fields struct {
		Limits   corev1.ResourceList
		Requests corev1.ResourceList
	}
	type args struct {
		def defaultResourceRequirementsSpec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ResourceRequirementsSpec
	}{
		{
			name:   "Sets defaults",
			fields: fields{},
			args: args{def: defaultResourceRequirementsSpec{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("200m"),
					corev1.ResourceMemory: resource.MustParse("200Mi"),
				},
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("100m"),
					corev1.ResourceMemory: resource.MustParse("100Mi"),
				},
			}},
			want: &ResourceRequirementsSpec{
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
		{
			name: "Combines explicitely set values with defaults",
			fields: fields{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("500m"),
					corev1.ResourceMemory: resource.MustParse("500Mi"),
				}},
			args: args{def: defaultResourceRequirementsSpec{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("200m"),
					corev1.ResourceMemory: resource.MustParse("200Mi"),
				},
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("100m"),
					corev1.ResourceMemory: resource.MustParse("100Mi"),
				},
			}},
			want: &ResourceRequirementsSpec{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("500m"),
					corev1.ResourceMemory: resource.MustParse("500Mi"),
				},
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("100m"),
					corev1.ResourceMemory: resource.MustParse("100Mi"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &ResourceRequirementsSpec{
				Limits:   tt.fields.Limits,
				Requests: tt.fields.Requests,
			}
			spec.Default(tt.args.def)
			if !reflect.DeepEqual(spec, tt.want) {
				t.Errorf("ResourceRequirementsSpec_Default() = %v, want %v", *spec, *tt.want)
			}
		})
	}
}

func TestResourceRequirementsSpec_IsDeactivated(t *testing.T) {
	tests := []struct {
		name string
		spec *ResourceRequirementsSpec
		want bool
	}{
		{"Wants true if empty", &ResourceRequirementsSpec{}, true},
		{"Wants false if nil", nil, false},
		{"Wants false if other",
			&ResourceRequirementsSpec{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("500m"),
					corev1.ResourceMemory: resource.MustParse("500Mi"),
				}},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.spec.IsDeactivated(); got != tt.want {
				t.Errorf("ResourceRequirementsSpec.IsDeactivated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitializeResourceRequirementsSpec(t *testing.T) {
	type args struct {
		spec *ResourceRequirementsSpec
		def  defaultResourceRequirementsSpec
	}
	tests := []struct {
		name string
		args args
		want *ResourceRequirementsSpec
	}{
		{
			name: "Initializes the struct with appropriate defaults if nil",
			args: args{nil, defaultResourceRequirementsSpec{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("500m"),
					corev1.ResourceMemory: resource.MustParse("500Mi"),
				},
			}},
			want: &ResourceRequirementsSpec{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("500m"),
					corev1.ResourceMemory: resource.MustParse("500Mi"),
				},
			},
		},
		{
			name: "Deactivated",
			args: args{&ResourceRequirementsSpec{}, defaultResourceRequirementsSpec{}},
			want: &ResourceRequirementsSpec{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitializeResourceRequirementsSpec(tt.args.spec, tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitializeResourceRequirementsSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func Test_stringOrDefault(t *testing.T) {
	type args struct {
		value    *string
		defValue *string
	}
	tests := []struct {
		name string
		args args
		want *string
	}{
		{
			name: "Value explicitely set",
			args: args{
				value:    pointer.StringPtr("value"),
				defValue: pointer.StringPtr("default"),
			},
			want: pointer.StringPtr("value"),
		},
		{
			name: "Value not set",
			args: args{
				value:    nil,
				defValue: pointer.StringPtr("default"),
			},
			want: pointer.StringPtr("default"),
		},
		{
			name: "Nor value not default set",
			args: args{
				value:    nil,
				defValue: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stringOrDefault(tt.args.value, tt.args.defValue)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stringOrDefault() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func Test_intOrDefault(t *testing.T) {
	type args struct {
		value    *int32
		defValue *int32
	}
	tests := []struct {
		name string
		args args
		want *int32
	}{
		{
			name: "Value explicitely set",
			args: args{
				value:    pointer.Int32Ptr(100),
				defValue: pointer.Int32Ptr(10),
			},
			want: pointer.Int32Ptr(100),
		},
		{
			name: "Value not set",
			args: args{
				value:    nil,
				defValue: pointer.Int32Ptr(10),
			},
			want: pointer.Int32Ptr(10),
		},
		{
			name: "Nor value not default set",
			args: args{
				value:    nil,
				defValue: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := intOrDefault(tt.args.value, tt.args.defValue)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("intOrDefault() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func Test_boolOrDefault(t *testing.T) {
	type args struct {
		value    *bool
		defValue *bool
	}
	tests := []struct {
		name string
		args args
		want *bool
	}{
		{
			name: "Value explicitely set",
			args: args{
				value:    pointer.BoolPtr(true),
				defValue: pointer.BoolPtr(false),
			},
			want: pointer.BoolPtr(true),
		},
		{
			name: "Value not set",
			args: args{
				value:    nil,
				defValue: pointer.BoolPtr(false),
			},
			want: pointer.BoolPtr(false),
		},
		{
			name: "Nor value not default set",
			args: args{
				value:    nil,
				defValue: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := boolOrDefault(tt.args.value, tt.args.defValue)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("boolOrDefault() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestCanary_CanarySpec(t *testing.T) {
	type fields struct {
		ImageName *string
		ImageTag  *string
		Replicas  *int32
		Patches   []string
	}
	type args struct {
		spec       interface{}
		canarySpec interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Returns a canary spec",
			fields: fields{
				Patches: []string{
					`[{"op": "replace", "path": "/image/name", "value": "new"}]`,
				},
			},
			args: args{
				spec: &BackendSpec{
					Image: &ImageSpec{
						Name: pointer.StringPtr("old"),
						Tag:  pointer.StringPtr("tag"),
					},
				},
				canarySpec: &BackendSpec{},
			},
			want: &BackendSpec{
				Image: &ImageSpec{
					Name: pointer.StringPtr("new"),
					Tag:  pointer.StringPtr("tag"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Canary{
				ImageName: tt.fields.ImageName,
				ImageTag:  tt.fields.ImageTag,
				Replicas:  tt.fields.Replicas,
				Patches:   tt.fields.Patches,
			}
			err := c.PatchSpec(tt.args.spec, tt.args.canarySpec)
			if (err != nil) != tt.wantErr {
				t.Errorf("Canary.CanarySpec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(tt.args.canarySpec, tt.want); len(diff) > 0 {
				t.Errorf("Canary.CanarySpec() = diff %v", diff)
			}
		})
	}
}

func TestVaultSecretStoreReferenceSpec_Default(t *testing.T) {
	type fields struct {
		Name *string
		Kind *string
	}
	type args struct {
		def defaultVaultSecretStoreReferenceSpec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *VaultSecretStoreReferenceSpec
	}{
		{
			name:   "Sets defaults",
			fields: fields{},
			args: args{def: defaultVaultSecretStoreReferenceSpec{
				Name: pointer.StringPtr("vault-mgmt"),
				Kind: pointer.StringPtr("ClusterSecretStore"),
			}},
			want: &VaultSecretStoreReferenceSpec{
				Name: pointer.StringPtr("vault-mgmt"),
				Kind: pointer.StringPtr("ClusterSecretStore"),
			},
		},
		{
			name: "Combines explicitely set values with defaults",
			fields: fields{
				Name: pointer.StringPtr("other-vault"),
			},
			args: args{def: defaultVaultSecretStoreReferenceSpec{
				Name: pointer.StringPtr("vault-mgmt"),
				Kind: pointer.StringPtr("ClusterSecretStore"),
			}},
			want: &VaultSecretStoreReferenceSpec{
				Name: pointer.StringPtr("other-vault"),
				Kind: pointer.StringPtr("ClusterSecretStore"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &VaultSecretStoreReferenceSpec{
				Name: tt.fields.Name,
				Kind: tt.fields.Kind,
			}
			spec.Default(tt.args.def)
			if !reflect.DeepEqual(spec, tt.want) {
				t.Errorf("VaultSecretStoreReferenceSpec_Default() = %v, want %v", *spec, *tt.want)
			}
		})
	}
}

func TestInitializeVaultSecretStoreReferenceSpec(t *testing.T) {
	type args struct {
		spec *VaultSecretStoreReferenceSpec
		def  defaultVaultSecretStoreReferenceSpec
	}
	tests := []struct {
		name string
		args args
		want *VaultSecretStoreReferenceSpec
	}{
		{
			name: "Initializes the struct with appropriate defaults if nil",
			args: args{nil, defaultVaultSecretStoreReferenceSpec{
				Name: pointer.StringPtr("vault-mgmt"),
				Kind: pointer.StringPtr("ClusterSecretStore"),
			}},
			want: &VaultSecretStoreReferenceSpec{
				Name: pointer.StringPtr("vault-mgmt"),
				Kind: pointer.StringPtr("ClusterSecretStore"),
			},
		},
		{
			name: "Deactivated",
			args: args{&VaultSecretStoreReferenceSpec{}, defaultVaultSecretStoreReferenceSpec{}},
			want: &VaultSecretStoreReferenceSpec{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitializeVaultSecretStoreReferenceSpec(tt.args.spec, tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitializeVaultSecretStoreReferenceSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}
