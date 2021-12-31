package pod

import (
	"fmt"
	"reflect"

	secretsmanagerv1alpha1 "github.com/3scale/saas-operator/pkg/apis/secrets-manager/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateSecretDefinitionFn generates a SecretDefinition
func GenerateSecretDefinitionFn(name, namespace string, labels map[string]string,
	opts interface{}) func() *secretsmanagerv1alpha1.SecretDefinition {

	return func() *secretsmanagerv1alpha1.SecretDefinition {
		return &secretsmanagerv1alpha1.SecretDefinition{
			TypeMeta: metav1.TypeMeta{
				Kind:       "SecretDefinition",
				APIVersion: secretsmanagerv1alpha1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
				Labels:    labels,
			},
			Spec: secretsmanagerv1alpha1.SecretDefinitionSpec{
				Name:    name,
				Type:    "opaque",
				KeysMap: keysMap(name, opts),
			},
		}
	}
}

func keysMap(name string, opts interface{}) map[string]secretsmanagerv1alpha1.DataSource {

	m := map[string]secretsmanagerv1alpha1.DataSource{}

	t := reflect.TypeOf(opts)

	for i := 0; i < t.NumField(); i++ {

		field := t.Field(i)

		// Ensure field is of EnvVarValue type
		if field.Type.String() != "pod.EnvVarValue" {
			panic(fmt.Errorf("field in '%s/%s' is not a 'pod.EnvVarValue'", t.Name(), field.Name))
		}

		secretName, hasSecretTag := field.Tag.Lookup("secret")
		if !hasSecretTag || secretName != name {
			continue
		}

		keyName, hasEnvTag := field.Tag.Lookup("env")
		if !hasEnvTag {
			panic(fmt.Errorf("missing 'env' tag from field '%s/%s'", t.Name(), field.Name))
		}

		value := reflect.ValueOf(opts).FieldByName(field.Name)
		// Skip field if its value is not set
		if value.IsZero() {
			continue
		}

		// Value should be of SecretValue type
		valueType := value.Elem().Elem().Type().String()
		if valueType != "pod.SecretValue" {
			panic(fmt.Errorf("wrong type '%s' for field %s/%s", valueType, t.Name(), field.Name))
		}

		secretValue := value.Elem().Elem().Interface().(SecretValue)
		if secretValue.Value.Override != nil {
			continue
		}
		m[keyName] = secretsmanagerv1alpha1.DataSource{
			Path: secretValue.Value.FromVault.Path,
			Key:  secretValue.Value.FromVault.Key,
		}
	}

	return m
}
