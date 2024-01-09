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

	"github.com/3scale-ops/basereconciler/util"
	"github.com/go-test/deep"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
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
				Name:           util.Pointer("name"),
				Tag:            util.Pointer("tag"),
				PullSecretName: util.Pointer("pullSecret"),
				PullPolicy:     func() *corev1.PullPolicy { p := corev1.PullIfNotPresent; return &p }(),
			}},
			want: &ImageSpec{
				Name:           util.Pointer("name"),
				Tag:            util.Pointer("tag"),
				PullSecretName: util.Pointer("pullSecret"),
				PullPolicy:     func() *corev1.PullPolicy { p := corev1.PullIfNotPresent; return &p }(),
			},
		},
		{
			name: "Combines explicitely set values with defaults",
			fields: fields{
				Name:       util.Pointer("explicit"),
				PullPolicy: func() *corev1.PullPolicy { p := corev1.PullAlways; return &p }(),
			},
			args: args{def: defaultImageSpec{
				Name:           util.Pointer("name"),
				Tag:            util.Pointer("tag"),
				PullSecretName: util.Pointer("pullSecret"),
				PullPolicy:     func() *corev1.PullPolicy { p := corev1.PullIfNotPresent; return &p }(),
			}},
			want: &ImageSpec{
				Name:           util.Pointer("explicit"),
				Tag:            util.Pointer("tag"),
				PullSecretName: util.Pointer("pullSecret"),
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
				Name:           util.Pointer("name"),
				Tag:            util.Pointer("tag"),
				PullSecretName: util.Pointer("pullSecret"),
			}},
			want: &ImageSpec{
				Name:           util.Pointer("name"),
				Tag:            util.Pointer("tag"),
				PullSecretName: util.Pointer("pullSecret"),
			},
		},
		{
			name: "Initializes the struct with appropriate defaults if empty",
			args: args{&ImageSpec{}, defaultImageSpec{
				Name:           util.Pointer("name"),
				Tag:            util.Pointer("tag"),
				PullSecretName: util.Pointer("pullSecret"),
			}},
			want: &ImageSpec{
				Name:           util.Pointer("name"),
				Tag:            util.Pointer("tag"),
				PullSecretName: util.Pointer("pullSecret"),
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
				InitialDelaySeconds: util.Pointer[int32](1),
				TimeoutSeconds:      util.Pointer[int32](2),
				PeriodSeconds:       util.Pointer[int32](3),
				SuccessThreshold:    util.Pointer[int32](4),
				FailureThreshold:    util.Pointer[int32](5),
			}},
			want: &ProbeSpec{
				InitialDelaySeconds: util.Pointer[int32](1),
				TimeoutSeconds:      util.Pointer[int32](2),
				PeriodSeconds:       util.Pointer[int32](3),
				SuccessThreshold:    util.Pointer[int32](4),
				FailureThreshold:    util.Pointer[int32](5),
			},
		},
		{
			name: "Combines explicitely set values with defaults",
			fields: fields{
				InitialDelaySeconds: util.Pointer[int32](9999),
			},
			args: args{def: defaultProbeSpec{
				InitialDelaySeconds: util.Pointer[int32](1),
				TimeoutSeconds:      util.Pointer[int32](2),
				PeriodSeconds:       util.Pointer[int32](3),
				SuccessThreshold:    util.Pointer[int32](4),
				FailureThreshold:    util.Pointer[int32](5),
			}},
			want: &ProbeSpec{
				InitialDelaySeconds: util.Pointer[int32](9999),
				TimeoutSeconds:      util.Pointer[int32](2),
				PeriodSeconds:       util.Pointer[int32](3),
				SuccessThreshold:    util.Pointer[int32](4),
				FailureThreshold:    util.Pointer[int32](5),
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
		{"Wants false if other", &ProbeSpec{InitialDelaySeconds: util.Pointer[int32](1)}, false},
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
				InitialDelaySeconds: util.Pointer[int32](1),
				TimeoutSeconds:      util.Pointer[int32](2),
				PeriodSeconds:       util.Pointer[int32](3),
				SuccessThreshold:    util.Pointer[int32](4),
				FailureThreshold:    util.Pointer[int32](5),
			}},
			want: &ProbeSpec{
				InitialDelaySeconds: util.Pointer[int32](1),
				TimeoutSeconds:      util.Pointer[int32](2),
				PeriodSeconds:       util.Pointer[int32](3),
				SuccessThreshold:    util.Pointer[int32](4),
				FailureThreshold:    util.Pointer[int32](5),
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
				ProxyProtocol:                 util.Pointer(true),
				CrossZoneLoadBalancingEnabled: util.Pointer(true),
				ConnectionDrainingEnabled:     util.Pointer(true),
				ConnectionDrainingTimeout:     util.Pointer[int32](1),
				HealthcheckHealthyThreshold:   util.Pointer[int32](2),
				HealthcheckUnhealthyThreshold: util.Pointer[int32](3),
				HealthcheckInterval:           util.Pointer[int32](4),
				HealthcheckTimeout:            util.Pointer[int32](5),
			}},
			want: &LoadBalancerSpec{
				ProxyProtocol:                 util.Pointer(true),
				CrossZoneLoadBalancingEnabled: util.Pointer(true),
				ConnectionDrainingEnabled:     util.Pointer(true),
				ConnectionDrainingTimeout:     util.Pointer[int32](1),
				HealthcheckHealthyThreshold:   util.Pointer[int32](2),
				HealthcheckUnhealthyThreshold: util.Pointer[int32](3),
				HealthcheckInterval:           util.Pointer[int32](4),
				HealthcheckTimeout:            util.Pointer[int32](5),
			},
		},
		{
			name: "Combines explicitely set values with defaults",
			fields: fields{
				ProxyProtocol: util.Pointer(false),
			},
			args: args{def: defaultLoadBalancerSpec{
				ProxyProtocol:                 util.Pointer(true),
				CrossZoneLoadBalancingEnabled: util.Pointer(true),
				ConnectionDrainingEnabled:     util.Pointer(true),
				ConnectionDrainingTimeout:     util.Pointer[int32](1),
				HealthcheckHealthyThreshold:   util.Pointer[int32](2),
				HealthcheckUnhealthyThreshold: util.Pointer[int32](3),
				HealthcheckInterval:           util.Pointer[int32](4),
				HealthcheckTimeout:            util.Pointer[int32](5),
			}},
			want: &LoadBalancerSpec{
				ProxyProtocol:                 util.Pointer(false),
				CrossZoneLoadBalancingEnabled: util.Pointer(true),
				ConnectionDrainingEnabled:     util.Pointer(true),
				ConnectionDrainingTimeout:     util.Pointer[int32](1),
				HealthcheckHealthyThreshold:   util.Pointer[int32](2),
				HealthcheckUnhealthyThreshold: util.Pointer[int32](3),
				HealthcheckInterval:           util.Pointer[int32](4),
				HealthcheckTimeout:            util.Pointer[int32](5),
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
		{"Wants false if other", &LoadBalancerSpec{ProxyProtocol: util.Pointer(false)}, false}}
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
				ProxyProtocol:                 util.Pointer(true),
				CrossZoneLoadBalancingEnabled: util.Pointer(true),
				ConnectionDrainingEnabled:     util.Pointer(true),
				ConnectionDrainingTimeout:     util.Pointer[int32](1),
				HealthcheckHealthyThreshold:   util.Pointer[int32](2),
				HealthcheckUnhealthyThreshold: util.Pointer[int32](3),
				HealthcheckInterval:           util.Pointer[int32](4),
				HealthcheckTimeout:            util.Pointer[int32](5),
			}},
			want: &LoadBalancerSpec{
				ProxyProtocol:                 util.Pointer(true),
				CrossZoneLoadBalancingEnabled: util.Pointer(true),
				ConnectionDrainingEnabled:     util.Pointer(true),
				ConnectionDrainingTimeout:     util.Pointer[int32](1),
				HealthcheckHealthyThreshold:   util.Pointer[int32](2),
				HealthcheckUnhealthyThreshold: util.Pointer[int32](3),
				HealthcheckInterval:           util.Pointer[int32](4),
				HealthcheckTimeout:            util.Pointer[int32](5),
			},
		},
		{
			name: "Initializes the struct with appropriate defaults if empty",
			args: args{&LoadBalancerSpec{}, defaultLoadBalancerSpec{
				ProxyProtocol:                 util.Pointer(true),
				CrossZoneLoadBalancingEnabled: util.Pointer(true),
				ConnectionDrainingEnabled:     util.Pointer(true),
				ConnectionDrainingTimeout:     util.Pointer[int32](1),
				HealthcheckHealthyThreshold:   util.Pointer[int32](2),
				HealthcheckUnhealthyThreshold: util.Pointer[int32](3),
				HealthcheckInterval:           util.Pointer[int32](4),
				HealthcheckTimeout:            util.Pointer[int32](5),
			}},
			want: &LoadBalancerSpec{
				ProxyProtocol:                 util.Pointer(true),
				CrossZoneLoadBalancingEnabled: util.Pointer(true),
				ConnectionDrainingEnabled:     util.Pointer(true),
				ConnectionDrainingTimeout:     util.Pointer[int32](1),
				HealthcheckHealthyThreshold:   util.Pointer[int32](2),
				HealthcheckUnhealthyThreshold: util.Pointer[int32](3),
				HealthcheckInterval:           util.Pointer[int32](4),
				HealthcheckTimeout:            util.Pointer[int32](5),
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
				ProxyProtocol:                 util.Pointer(true),
				CrossZoneLoadBalancingEnabled: util.Pointer(true),
			}},
			want: &NLBLoadBalancerSpec{
				ProxyProtocol:                 util.Pointer(true),
				CrossZoneLoadBalancingEnabled: util.Pointer(true),
			},
		},
		{
			name: "Combines explicitely set values with defaults",
			fields: fields{
				ProxyProtocol: util.Pointer(false),
			},
			args: args{def: defaultNLBLoadBalancerSpec{
				ProxyProtocol:                 util.Pointer(true),
				CrossZoneLoadBalancingEnabled: util.Pointer(true),
			}},
			want: &NLBLoadBalancerSpec{
				ProxyProtocol:                 util.Pointer(false),
				CrossZoneLoadBalancingEnabled: util.Pointer(true),
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
		{"Wants false if other", &NLBLoadBalancerSpec{ProxyProtocol: util.Pointer(false)}, false}}
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
				ProxyProtocol:                 util.Pointer(true),
				CrossZoneLoadBalancingEnabled: util.Pointer(true),
			}},
			want: &NLBLoadBalancerSpec{
				ProxyProtocol:                 util.Pointer(true),
				CrossZoneLoadBalancingEnabled: util.Pointer(true),
			},
		},
		{
			name: "Initializes the struct with appropriate defaults if empty",
			args: args{&NLBLoadBalancerSpec{}, defaultNLBLoadBalancerSpec{
				ProxyProtocol:                 util.Pointer(true),
				CrossZoneLoadBalancingEnabled: util.Pointer(true),
			}},
			want: &NLBLoadBalancerSpec{
				ProxyProtocol:                 util.Pointer(true),
				CrossZoneLoadBalancingEnabled: util.Pointer(true),
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
				SelectorKey:   util.Pointer("key"),
				SelectorValue: util.Pointer("label"),
			}},
			want: &GrafanaDashboardSpec{
				SelectorKey:   util.Pointer("key"),
				SelectorValue: util.Pointer("label"),
			},
		},
		{
			name: "Combines explicitely set values with defaults",
			fields: fields{
				SelectorKey: util.Pointer("xxxx"),
			},
			args: args{def: defaultGrafanaDashboardSpec{
				SelectorKey:   util.Pointer("key"),
				SelectorValue: util.Pointer("label"),
			}},
			want: &GrafanaDashboardSpec{
				SelectorKey:   util.Pointer("xxxx"),
				SelectorValue: util.Pointer("label"),
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
		{"Wants false if other", &GrafanaDashboardSpec{SelectorKey: util.Pointer("key")}, false},
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
				SelectorKey:   util.Pointer("key"),
				SelectorValue: util.Pointer("label"),
			}},
			want: &GrafanaDashboardSpec{
				SelectorKey:   util.Pointer("key"),
				SelectorValue: util.Pointer("label"),
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
				MinAvailable:   util.Pointer(intstr.FromString("default")),
				MaxUnavailable: nil,
			}},
			want: &PodDisruptionBudgetSpec{
				MinAvailable:   util.Pointer(intstr.FromString("default")),
				MaxUnavailable: nil,
			},
		},
		{
			name: "Combines explicitely set values with defaults",
			fields: fields{
				MinAvailable: util.Pointer(intstr.FromString("explicit")),
			},
			args: args{def: defaultPodDisruptionBudgetSpec{
				MinAvailable:   util.Pointer(intstr.FromString("default")),
				MaxUnavailable: nil,
			}},
			want: &PodDisruptionBudgetSpec{
				MinAvailable:   util.Pointer(intstr.FromString("explicit")),
				MaxUnavailable: nil,
			},
		},
		{
			name: "Only one of MinAvailable or MaxUnavailable can be set",
			fields: fields{
				MinAvailable: util.Pointer(intstr.FromString("explicit")),
			},
			args: args{def: defaultPodDisruptionBudgetSpec{
				MinAvailable:   nil,
				MaxUnavailable: util.Pointer(intstr.FromString("default")),
			}},
			want: &PodDisruptionBudgetSpec{
				MinAvailable:   util.Pointer(intstr.FromString("explicit")),
				MaxUnavailable: nil,
			},
		},
		{
			name:   "Only one of MinAvailable or MaxUnavailable can be set (II)",
			fields: fields{},
			args: args{def: defaultPodDisruptionBudgetSpec{
				MinAvailable:   util.Pointer(intstr.IntOrString{Type: intstr.String, StrVal: "defaultMin"}),
				MaxUnavailable: util.Pointer(intstr.IntOrString{Type: intstr.String, StrVal: "defaultMax"}),
			}},
			want: &PodDisruptionBudgetSpec{
				MinAvailable:   util.Pointer(intstr.IntOrString{Type: intstr.String, StrVal: "defaultMin"}),
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
		{"Wants false if other", &PodDisruptionBudgetSpec{MinAvailable: util.Pointer(intstr.FromInt(1))}, false},
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
				MinAvailable:   util.Pointer(intstr.FromString("default")),
				MaxUnavailable: nil,
			}},
			want: &PodDisruptionBudgetSpec{
				MinAvailable:   util.Pointer(intstr.FromString("default")),
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
				MinReplicas:         util.Pointer[int32](1),
				MaxReplicas:         util.Pointer[int32](2),
				ResourceUtilization: util.Pointer[int32](3),
				ResourceName:        util.Pointer("xxxx"),
			}},
			want: &HorizontalPodAutoscalerSpec{
				MinReplicas:         util.Pointer[int32](1),
				MaxReplicas:         util.Pointer[int32](2),
				ResourceUtilization: util.Pointer[int32](3),
				ResourceName:        util.Pointer("xxxx"),
			},
		},
		{
			name: "Combines explicitely set values with defaults",
			fields: fields{
				MinReplicas: util.Pointer[int32](9999),
			},
			args: args{def: defaultHorizontalPodAutoscalerSpec{
				MinReplicas:         util.Pointer[int32](1),
				MaxReplicas:         util.Pointer[int32](2),
				ResourceUtilization: util.Pointer[int32](3),
				ResourceName:        util.Pointer("xxxx"),
			}},
			want: &HorizontalPodAutoscalerSpec{
				MinReplicas:         util.Pointer[int32](9999),
				MaxReplicas:         util.Pointer[int32](2),
				ResourceUtilization: util.Pointer[int32](3),
				ResourceName:        util.Pointer("xxxx"),
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
		{"Wants false if other", &HorizontalPodAutoscalerSpec{MinReplicas: util.Pointer[int32](1)}, false}}
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
				MinReplicas:         util.Pointer[int32](1),
				MaxReplicas:         util.Pointer[int32](2),
				ResourceUtilization: util.Pointer[int32](3),
				ResourceName:        util.Pointer("xxxx"),
			}},
			want: &HorizontalPodAutoscalerSpec{
				MinReplicas:         util.Pointer[int32](1),
				MaxReplicas:         util.Pointer[int32](2),
				ResourceUtilization: util.Pointer[int32](3),
				ResourceName:        util.Pointer("xxxx"),
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
				value:    util.Pointer("value"),
				defValue: util.Pointer("default"),
			},
			want: util.Pointer("value"),
		},
		{
			name: "Value not set",
			args: args{
				value:    nil,
				defValue: util.Pointer("default"),
			},
			want: util.Pointer("default"),
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
				value:    util.Pointer[int32](100),
				defValue: util.Pointer[int32](10),
			},
			want: util.Pointer[int32](100),
		},
		{
			name: "Value not set",
			args: args{
				value:    nil,
				defValue: util.Pointer[int32](10),
			},
			want: util.Pointer[int32](10),
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
				value:    util.Pointer(true),
				defValue: util.Pointer(false),
			},
			want: util.Pointer(true),
		},
		{
			name: "Value not set",
			args: args{
				value:    nil,
				defValue: util.Pointer(false),
			},
			want: util.Pointer(false),
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
						Name: util.Pointer("old"),
						Tag:  util.Pointer("tag"),
					},
				},
				canarySpec: &BackendSpec{},
			},
			want: &BackendSpec{
				Image: &ImageSpec{
					Name: util.Pointer("new"),
					Tag:  util.Pointer("tag"),
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

func TestExternalSecretSecretStoreReferenceSpec_Default(t *testing.T) {
	type fields struct {
		Name *string
		Kind *string
	}
	type args struct {
		def defaultExternalSecretSecretStoreReferenceSpec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ExternalSecretSecretStoreReferenceSpec
	}{
		{
			name:   "Sets defaults",
			fields: fields{},
			args: args{def: defaultExternalSecretSecretStoreReferenceSpec{
				Name: util.Pointer("vault-mgmt"),
				Kind: util.Pointer("ClusterSecretStore"),
			}},
			want: &ExternalSecretSecretStoreReferenceSpec{
				Name: util.Pointer("vault-mgmt"),
				Kind: util.Pointer("ClusterSecretStore"),
			},
		},
		{
			name: "Combines explicitely set values with defaults",
			fields: fields{
				Name: util.Pointer("other-vault"),
			},
			args: args{def: defaultExternalSecretSecretStoreReferenceSpec{
				Name: util.Pointer("vault-mgmt"),
				Kind: util.Pointer("ClusterSecretStore"),
			}},
			want: &ExternalSecretSecretStoreReferenceSpec{
				Name: util.Pointer("other-vault"),
				Kind: util.Pointer("ClusterSecretStore"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &ExternalSecretSecretStoreReferenceSpec{
				Name: tt.fields.Name,
				Kind: tt.fields.Kind,
			}
			spec.Default(tt.args.def)
			if !reflect.DeepEqual(spec, tt.want) {
				t.Errorf("ExternalSecretSecretStoreReferenceSpec_Default() = %v, want %v", *spec, *tt.want)
			}
		})
	}
}

func TestInitializeVaultSecretStoreReferenceSpec(t *testing.T) {
	type args struct {
		spec *ExternalSecretSecretStoreReferenceSpec
		def  defaultExternalSecretSecretStoreReferenceSpec
	}
	tests := []struct {
		name string
		args args
		want *ExternalSecretSecretStoreReferenceSpec
	}{
		{
			name: "Initializes the struct with appropriate defaults if nil",
			args: args{nil, defaultExternalSecretSecretStoreReferenceSpec{
				Name: util.Pointer("vault-mgmt"),
				Kind: util.Pointer("ClusterSecretStore"),
			}},
			want: &ExternalSecretSecretStoreReferenceSpec{
				Name: util.Pointer("vault-mgmt"),
				Kind: util.Pointer("ClusterSecretStore"),
			},
		},
		{
			name: "Deactivated",
			args: args{&ExternalSecretSecretStoreReferenceSpec{}, defaultExternalSecretSecretStoreReferenceSpec{}},
			want: &ExternalSecretSecretStoreReferenceSpec{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitializeExternalSecretSecretStoreReferenceSpec(tt.args.spec, tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitializeExternalSecretSecretStoreReferenceSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}
