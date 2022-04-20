package pod

import (
	"fmt"
	"reflect"

	externalsecretsv1beta1 "github.com/3scale/saas-operator/pkg/apis/externalsecrets/v1beta1"
	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GenerateExternalSecretFn generates a ExternalSecret
func GenerateExternalSecretFn(name, namespace, secretStoreName, secretStoreKind string, refreshInterval metav1.Duration, labels map[string]string,
	opts interface{}) basereconciler.GeneratorFunction {

	return func() client.Object {
		return &externalsecretsv1beta1.ExternalSecret{
			TypeMeta: metav1.TypeMeta{
				Kind:       externalsecretsv1beta1.ExtSecretKind,
				APIVersion: externalsecretsv1beta1.ExtSecretGroupVersionKind.GroupVersion().String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
				Labels:    labels,
			},
			Spec: externalsecretsv1beta1.ExternalSecretSpec{
				SecretStoreRef: externalsecretsv1beta1.SecretStoreRef{
					Name: secretStoreName,
					Kind: secretStoreKind,
				},
				Target: externalsecretsv1beta1.ExternalSecretTarget{
					Name: name,
				},
				RefreshInterval: &refreshInterval,
				Data:            keysSlice(name, opts),
			},
		}
	}
}

func keysSlice(name string, opts interface{}) []externalsecretsv1beta1.ExternalSecretData {

	s := []externalsecretsv1beta1.ExternalSecretData{}

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

		s = append(s, externalsecretsv1beta1.ExternalSecretData{
			SecretKey: keyName,
			RemoteRef: externalsecretsv1beta1.ExternalSecretDataRemoteRef{
				Key:      secretValue.Value.FromVault.Path,
				Property: secretValue.Value.FromVault.Key,
			},
		})
	}

	return s
}
