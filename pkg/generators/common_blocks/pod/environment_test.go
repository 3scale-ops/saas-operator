package pod

import (
	"reflect"
	"testing"

	"github.com/3scale/saas-operator/api/v1alpha1"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

func TestClearTextValue_ToEnvVar(t *testing.T) {
	type fields struct {
		Value string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   corev1.EnvVar
	}{
		{
			name:   "Returns EnvVar from clear text value",
			fields: fields{Value: "value"},
			args:   args{key: "key"},
			want:   corev1.EnvVar{Name: "key", Value: "value"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctv := &ClearTextValue{
				Value: tt.fields.Value,
			}
			if got := ctv.ToEnvVar(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClearTextValue.ToEnvVar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecretValue_ToEnvVar(t *testing.T) {
	type fields struct {
		Value saasv1alpha1.SecretReference
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   corev1.EnvVar
	}{
		{
			name:   "Returns EnvVar from a Secret",
			fields: fields{Value: saasv1alpha1.SecretReference{FromVault: &saasv1alpha1.VaultSecretReference{}}},
			args:   args{key: "key:my-secret"},
			want: corev1.EnvVar{
				Name: "key",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						Key: "key",
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "my-secret",
						},
					},
				},
			},
		},
		{
			name:   "Returns EnvVar from an overrided Secret",
			fields: fields{Value: saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")}},
			args:   args{key: "key:my-secret"},
			want:   corev1.EnvVar{Name: "key", Value: "override"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sv := &SecretValue{
				Value: tt.fields.Value,
			}
			if got := sv.ToEnvVar(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SecretValue.ToEnvVar() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
					Option4 EnvVarValue `env:"OPTION4"`
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
			name: "Returns a slice of EnvVar, with overrides",
			args: args{
				opts: struct {
					Option1 EnvVarValue `env:"OPTION1"`
					Option2 EnvVarValue `env:"OPTION2" secret:"my-secret"`
					Option3 EnvVarValue `env:"OPTION3" secret:"my-secret"`
				}{
					Option1: &ClearTextValue{Value: "value1"},
					Option2: &SecretValue{Value: v1alpha1.SecretReference{Override: pointer.StringPtr("override")}},
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
					Value: "override",
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
