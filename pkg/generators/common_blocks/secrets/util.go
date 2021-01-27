package secrets

import (
	corev1 "k8s.io/api/core/v1"
)

// EnvVarFromSecret is a helper to take an environment variable from
// a Secret resource
func EnvVarFromSecret(envvar, name, key string) corev1.EnvVar {

	return corev1.EnvVar{
		Name: envvar,
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: name,
				},
				Key: key,
			},
		},
	}
}
