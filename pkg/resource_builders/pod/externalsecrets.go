package pod

import (
	"fmt"
	"reflect"

	externalsecretsv1alpha1 "github.com/3scale/saas-operator/pkg/apis/externalsecrets/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateExternalSecretFn generates a ExternalSecret
func GenerateExternalSecretFn(name, namespace, secretStoreName, secretStoreKind string, refreshInterval metav1.Duration, labels map[string]string,
	opts interface{}) func() *externalsecretsv1alpha1.ExternalSecret {

	return func() *externalsecretsv1alpha1.ExternalSecret {
		return &externalsecretsv1alpha1.ExternalSecret{
			TypeMeta: metav1.TypeMeta{
				Kind:       externalsecretsv1alpha1.ExtSecretKind,
				APIVersion: externalsecretsv1alpha1.ExtSecretGroupVersionKind.GroupVersion().String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
				Labels:    labels,
			},
			Spec: externalsecretsv1alpha1.ExternalSecretSpec{
				SecretStoreRef: externalsecretsv1alpha1.SecretStoreRef{
					Name: secretStoreName,
					Kind: secretStoreKind,
				},
				Target: externalsecretsv1alpha1.ExternalSecretTarget{
					Name: name,
				},
				RefreshInterval: &refreshInterval,
				Data:            keysSlice(name, opts),
			},
		}
	}
}

func keysSlice(name string, opts interface{}) []externalsecretsv1alpha1.ExternalSecretData {

	s := []externalsecretsv1alpha1.ExternalSecretData{}

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

		s = append(s, externalsecretsv1alpha1.ExternalSecretData{
			SecretKey: keyName,
			RemoteRef: externalsecretsv1alpha1.ExternalSecretDataRemoteRef{
				Key:      secretValue.Value.FromVault.Path,
				Property: secretValue.Value.FromVault.Key,
			},
		})
	}

	return s
}
