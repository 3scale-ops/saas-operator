package system

import (
	secretsmanagerv1alpha1 "github.com/3scale/saas-operator/pkg/apis/secrets-manager/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ConfigFilesSecretDefinition generates a SecretDefinition
func (gen *Generator) ConfigFilesSecretDefinition() basereconciler.GeneratorFunction {

	return func() client.Object {
		return &secretsmanagerv1alpha1.SecretDefinition{
			TypeMeta: metav1.TypeMeta{
				Kind:       "SecretDefinition",
				APIVersion: secretsmanagerv1alpha1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "system-config",
				Namespace: gen.GetNamespace(),
				Labels:    gen.GetLabels(),
			},
			Spec: secretsmanagerv1alpha1.SecretDefinitionSpec{
				Name: "system-config",
				Type: "opaque",
				KeysMap: func() map[string]secretsmanagerv1alpha1.DataSource {
					m := map[string]secretsmanagerv1alpha1.DataSource{}
					for _, file := range gen.ConfigFilesSpec.Files {
						m[file] = secretsmanagerv1alpha1.DataSource{
							Path: gen.ConfigFilesSpec.VaultPath,
							Key:  file,
						}
					}
					return m
				}(),
			},
		}
	}
}
