package pod

import (
	"reflect"
	"testing"

	"github.com/3scale/saas-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func TestBuildEnvironment(t *testing.T) {
	type args struct {
		opts interface{}
	}
	tests := []struct {
		name      string
		args      args
		want      []corev1.EnvVar
		wantPanic bool
	}{
		{
			name: "Returns a slice of EnvVar",
			args: args{
				opts: struct {
					Option1 EnvVarValue `env:"OPTION1"`
					Option2 EnvVarValue `env:"OPTION2"`
					Option3 EnvVarValue `env:"OPTION3" secret:"my-secret"`
				}{
					Option1: &ClearTextValue{Value: "value1"},
					Option2: &ClearTextValue{Value: "value2"},
					Option3: &SecretValue{Value: v1alpha1.SecretReference{}},
				},
			},
			want: []corev1.EnvVar{
				{
					Name:  "OPTION1",
					Value: "value1",
				},
				{
					Name:  "OPTION2",
					Value: "value2",
				},
				{
					Name: "OPTION3",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							Key: "OPTION3",
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "my-secret",
							},
						},
					},
				},
			},
		},
		{
			name: "Panics due to field not being an EnvVarValue",
			args: args{
				opts: struct {
					Option1 string `env:"OPTION1"`
				}{
					Option1: "value",
				},
			},
			want:      nil,
			wantPanic: true,
		},
		{
			name: "Panics due to field missing 'secret' tag",
			args: args{
				opts: struct {
					Option1 EnvVarValue `env:"OPTION1"`
				}{
					Option1: &SecretValue{Value: v1alpha1.SecretReference{}},
				},
			},
			want:      nil,
			wantPanic: true,
		},
		{
			name: "Panics due to unexpected 'secret' tag",
			args: args{
				opts: struct {
					Option1 EnvVarValue `env:"OPTION1" secret:"some-secret"`
				}{
					Option1: &ClearTextValue{Value: "xxxx"},
				},
			},
			want:      nil,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil && tt.wantPanic {
					t.Errorf("code did not panic")
				}
				if r != nil && !tt.wantPanic {
					t.Errorf("code caused a panic")
				}
			}()

			if got := BuildEnvironment(tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}
}
