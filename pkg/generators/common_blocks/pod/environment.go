package pod

import (
	"sort"

	corev1 "k8s.io/api/core/v1"
)

type EnvVarValue interface {
	ToEnvVar(key string) corev1.EnvVar
}

type DirectValue struct {
	Value string
}

func (dv *DirectValue) ToEnvVar(key string) corev1.EnvVar {
	return corev1.EnvVar{
		Name:  key,
		Value: dv.Value,
	}
}

type SecretRef struct {
	SecretName string
}

func (dv *SecretRef) ToEnvVar(key string) corev1.EnvVar {
	return corev1.EnvVar{
		Name: key,
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				Key: key,
				LocalObjectReference: corev1.LocalObjectReference{
					Name: dv.SecretName,
				},
			},
		},
	}
}

func GenerateEnvironment(base map[string]string, config map[string]EnvVarValue) []corev1.EnvVar {

	envmap := map[string]EnvVarValue{}

	for k, v := range base {
		envmap[k] = &DirectValue{Value: v}
	}

	for k, v := range config {
		envmap[k] = v
	}

	env := []corev1.EnvVar{}
	for k, v := range envmap {
		env = append(env, v.ToEnvVar(k))
	}

	// Sort to return always the same result
	sort.Slice(env, func(i, j int) bool {
		return env[i].Name < env[j].Name
	})

	return env
}
