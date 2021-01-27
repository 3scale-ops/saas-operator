package secrets

import (
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	secretsmanagerv1alpha1 "github.com/3scale/saas-operator/pkg/apis/secrets-manager/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewSecretDefinition returns a basereconciler.GeneratorFunction function that will return a SecretDefinition
// resource when called
func NewSecretDefinition(key types.NamespacedName, labels map[string]string, secretName string,
	cfg saasv1alpha1.VaultSecretReference) basereconciler.GeneratorFunction {

	return func() client.Object {

		return &secretsmanagerv1alpha1.SecretDefinition{
			TypeMeta: metav1.TypeMeta{
				Kind:       "SecretDefinition",
				APIVersion: secretsmanagerv1alpha1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
				Labels:    labels,
			},
			Spec: secretsmanagerv1alpha1.SecretDefinitionSpec{
				Name: secretName,
				Type: "opaque",
				KeysMap: map[string]secretsmanagerv1alpha1.DataSource{
					cfg.Key: {
						Key:  cfg.Key,
						Path: cfg.Path,
					},
				},
			},
		}
	}
}
