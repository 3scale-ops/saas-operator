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
			}},
			want: &ImageSpec{
				Name:           pointer.StringPtr("name"),
				Tag:            pointer.StringPtr("tag"),
				PullSecretName: pointer.StringPtr("pullSecret"),
			},
		},
		{
			name: "Combines explicitely set values with defaults",
			fields: fields{
				Name: pointer.StringPtr("explicit"),
			},
			args: args{def: defaultImageSpec{
				Name:           pointer.StringPtr("name"),
				Tag:            pointer.StringPtr("tag"),
				PullSecretName: pointer.StringPtr("pullSecret"),
			}},
			want: &ImageSpec{
				Name:           pointer.StringPtr("explicit"),
				Tag:            pointer.StringPtr("tag"),
				PullSecretName: pointer.StringPtr("pullSecret"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &ImageSpec{
				Name:           tt.fields.Name,
				Tag:            tt.fields.Tag,
				PullSecretName: tt.fields.PullSecretName,
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

func TestHTTPProbeSpec_Default(t *testing.T) {
	type fields struct {
		InitialDelaySeconds *int32
		TimeoutSeconds      *int32
		PeriodSeconds       *int32
		SuccessThreshold    *int32
		FailureThreshold    *int32
	}
	type args struct {
		def defaultHTTPProbeSpec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *HTTPProbeSpec
	}{
		{
			name:   "Sets defaults",
			fields: fields{},
			args: args{def: defaultHTTPProbeSpec{
				InitialDelaySeconds: pointer.Int32Ptr(1),
				TimeoutSeconds:      pointer.Int32Ptr(2),
				PeriodSeconds:       pointer.Int32Ptr(3),
				SuccessThreshold:    pointer.Int32Ptr(4),
				FailureThreshold:    pointer.Int32Ptr(5),
			}},
			want: &HTTPProbeSpec{
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
			args: args{def: defaultHTTPProbeSpec{
				InitialDelaySeconds: pointer.Int32Ptr(1),
				TimeoutSeconds:      pointer.Int32Ptr(2),
				PeriodSeconds:       pointer.Int32Ptr(3),
				SuccessThreshold:    pointer.Int32Ptr(4),
				FailureThreshold:    pointer.Int32Ptr(5),
			}},
			want: &HTTPProbeSpec{
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
			spec := &HTTPProbeSpec{
				InitialDelaySeconds: tt.fields.InitialDelaySeconds,
				TimeoutSeconds:      tt.fields.TimeoutSeconds,
				PeriodSeconds:       tt.fields.PeriodSeconds,
				SuccessThreshold:    tt.fields.SuccessThreshold,
				FailureThreshold:    tt.fields.FailureThreshold,
			}
			spec.Default(tt.args.def)
			if !reflect.DeepEqual(spec, tt.want) {
				t.Errorf("HTTPProbeSpec_Default() = %v, want %v", *spec, *tt.want)
			}
		})
	}
}

func TestHTTPProbeSpec_IsDeactivated(t *testing.T) {
	tests := []struct {
		name string
		spec *HTTPProbeSpec
		want bool
	}{
		{"Wants true if empty", &HTTPProbeSpec{}, true},
		{"Wants false if nil", nil, false},
		{"Wants false if other", &HTTPProbeSpec{InitialDelaySeconds: pointer.Int32Ptr(1)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.spec.IsDeactivated(); got != tt.want {
				t.Errorf("HTTPProbeSpec.IsDeactivated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitializeHTTPProbeSpec(t *testing.T) {
	type args struct {
		spec *HTTPProbeSpec
		def  defaultHTTPProbeSpec
	}
	tests := []struct {
		name string
		args args
		want *HTTPProbeSpec
	}{
		{
			name: "Initializes the struct with appropriate defaults if nil",
			args: args{nil, defaultHTTPProbeSpec{
				InitialDelaySeconds: pointer.Int32Ptr(1),
				TimeoutSeconds:      pointer.Int32Ptr(2),
				PeriodSeconds:       pointer.Int32Ptr(3),
				SuccessThreshold:    pointer.Int32Ptr(4),
				FailureThreshold:    pointer.Int32Ptr(5),
			}},
			want: &HTTPProbeSpec{
				InitialDelaySeconds: pointer.Int32Ptr(1),
				TimeoutSeconds:      pointer.Int32Ptr(2),
				PeriodSeconds:       pointer.Int32Ptr(3),
				SuccessThreshold:    pointer.Int32Ptr(4),
				FailureThreshold:    pointer.Int32Ptr(5),
			},
		},
		{
			name: "Deactivated",
			args: args{&HTTPProbeSpec{}, defaultHTTPProbeSpec{}},
			want: &HTTPProbeSpec{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitializeHTTPProbeSpec(tt.args.spec, tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitializeHTTPProbeSpec() = %v, want %v", got, tt.want)
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
				MinAvailable:   intstrPtr(intstr.IntOrString{Type: intstr.String, StrVal: "default"}),
				MaxUnavailable: nil,
			}},
			want: &PodDisruptionBudgetSpec{
				MinAvailable:   intstrPtr(intstr.IntOrString{Type: intstr.String, StrVal: "default"}),
				MaxUnavailable: nil,
			},
		},
		{
			name: "Combines explicitely set values with defaults",
			fields: fields{
				MinAvailable: intstrPtr(intstr.IntOrString{Type: intstr.String, StrVal: "explicit"}),
			},
			args: args{def: defaultPodDisruptionBudgetSpec{
				MinAvailable:   intstrPtr(intstr.IntOrString{Type: intstr.String, StrVal: "default"}),
				MaxUnavailable: nil,
			}},
			want: &PodDisruptionBudgetSpec{
				MinAvailable:   intstrPtr(intstr.IntOrString{Type: intstr.String, StrVal: "explicit"}),
				MaxUnavailable: nil,
			},
		},
		{
			name: "Only one of MinAvailable or MaxUnavailable can be set",
			fields: fields{
				MinAvailable: intstrPtr(intstr.IntOrString{Type: intstr.String, StrVal: "explicit"}),
			},
			args: args{def: defaultPodDisruptionBudgetSpec{
				MinAvailable:   nil,
				MaxUnavailable: intstrPtr(intstr.IntOrString{Type: intstr.String, StrVal: "default"}),
			}},
			want: &PodDisruptionBudgetSpec{
				MinAvailable:   intstrPtr(intstr.IntOrString{Type: intstr.String, StrVal: "explicit"}),
				MaxUnavailable: nil,
			},
		},
		{
			name:   "Only one of MinAvailable or MaxUnavailable can be set (II)",
			fields: fields{},
			args: args{def: defaultPodDisruptionBudgetSpec{
				MinAvailable:   intstrPtr(intstr.IntOrString{Type: intstr.String, StrVal: "defaultMin"}),
				MaxUnavailable: intstrPtr(intstr.IntOrString{Type: intstr.String, StrVal: "defaultMax"}),
			}},
			want: &PodDisruptionBudgetSpec{
				MinAvailable:   intstrPtr(intstr.IntOrString{Type: intstr.String, StrVal: "defaultMin"}),
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
		{"Wants false if other", &PodDisruptionBudgetSpec{MinAvailable: intstrPtr(intstr.FromInt(1))}, false},
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
				MinAvailable:   intstrPtr(intstr.IntOrString{Type: intstr.String, StrVal: "default"}),
				MaxUnavailable: nil,
			}},
			want: &PodDisruptionBudgetSpec{
				MinAvailable:   intstrPtr(intstr.IntOrString{Type: intstr.String, StrVal: "default"}),
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

func Test_intstrPtr(t *testing.T) {
	type args struct {
		value intstr.IntOrString
	}
	tests := []struct {
		name string
		args args
		want *intstr.IntOrString
	}{
		{"Returns a pointer", args{value: intstr.FromInt(1)}, &intstr.IntOrString{Type: intstr.Int, IntVal: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intstrPtr(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("intstrPtr() = %v, want %v", got, tt.want)
			}
		})
	}
}
