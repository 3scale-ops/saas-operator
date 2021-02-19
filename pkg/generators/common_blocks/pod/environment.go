package pod

import (
	"fmt"
	"reflect"
	"strings"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

type EnvVarValue interface {
	ToEnvVar(key string) corev1.EnvVar
}

type ClearTextValue struct {
	Value string
}

func (ctv *ClearTextValue) ToEnvVar(key string) corev1.EnvVar {
	return corev1.EnvVar{
		Name:  key,
		Value: ctv.Value,
	}
}

type SecretValue struct {
	Value saasv1alpha1.SecretReference
}

func (sv *SecretValue) ToEnvVar(key string) corev1.EnvVar {
	s := strings.Split(key, ":")
	envvar := s[0]
	secret := s[1]

	if sv.Value.Override != nil {
		return corev1.EnvVar{
			Name:  envvar,
			Value: *sv.Value.Override,
		}
	}

	return corev1.EnvVar{
		Name: envvar,
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				Key: envvar,
				LocalObjectReference: corev1.LocalObjectReference{
					Name: secret,
				},
			},
		},
	}
}

func BuildEnvironment(opts interface{}) []corev1.EnvVar {
	env := []corev1.EnvVar{}

	t := reflect.TypeOf(opts)

	for i := 0; i < t.NumField(); i++ {

		field := t.Field(i)

		// Ensure field is of EnvVarValue type
		if field.Type.String() != "pod.EnvVarValue" {
			panic(fmt.Errorf("Field in '%s/%s' is not a 'pod.EnvVarValue'", t.Name(), field.Name))
		}

		value := reflect.ValueOf(opts).FieldByName(field.Name)
		// Skip field if its value is not set
		if value.IsZero() {
			continue
		}

		// Parse the field "env" tag
		envVarName, ok := field.Tag.Lookup("env")
		if !ok {
			panic(fmt.Errorf("missing 'env' tag in  %s/%s", t.Name(), field.Name))
		}

		secretName, hasSecretTag := field.Tag.Lookup("secret")

		valueType := value.Elem().Elem().Type().String()
		// If value is of ClearTextValue type it shoud not have the 'secret' tag
		if valueType == "pod.ClearTextValue" && hasSecretTag {
			panic(fmt.Errorf("unexpected 'secret' tag in field  %s/%s", t.Name(), field.Name))
		}

		// If value is of SecretValue type it shoud have the 'secret' tag
		if valueType == "pod.SecretValue" && !hasSecretTag {
			panic(fmt.Errorf("missing 'secret' tag in field  %s/%s", t.Name(), field.Name))
		}

		var arg reflect.Value
		if valueType == "pod.ClearTextValue" {
			arg = reflect.ValueOf(envVarName)
		} else {
			arg = reflect.ValueOf(strings.Join([]string{envVarName, secretName}, ":"))
		}

		envvar := value.MethodByName("ToEnvVar").Call([]reflect.Value{arg})
		env = append(env, envvar[0].Interface().(corev1.EnvVar))
	}

	return env
}
