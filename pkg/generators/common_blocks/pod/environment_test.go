package pod

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestGenerateEnvironment(t *testing.T) {
	type args struct {
		base   map[string]string
		config map[string]EnvVarValue
	}
	tests := []struct {
		name string
		args args
		want []corev1.EnvVar
	}{
		{
			"Fuck me if I get it",
			args{
				base: map[string]string{
					"env1": "baseValue1",
					"env2": "baseValue2",
				},
				config: map[string]EnvVarValue{
					"env2": &DirectValue{Value: "configValue2"},
					"env3": &DirectValue{Value: "configValue3"},
					"env4": &SecretRef{SecretName: "secret"},
				},
			},
			[]corev1.EnvVar{
				{
					Name:  "env1",
					Value: "baseValue1",
				},
				{
					Name:  "env2",
					Value: "configValue2",
				},
				{
					Name:  "env3",
					Value: "configValue3",
				},
				{
					Name: "env4",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							Key: "env4",
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "secret",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateEnvironment(tt.args.base, tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}
}
